package controller

import (
	"context"
	"time"

	"github.com/bitnami-labs/kubewatch/pkg/controller"

	client2 "github.com/bitnami-labs/kubewatch/pkg/client"

	"github.com/bitnami-labs/kubewatch/config"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/bitnami-labs/kubewatch/api/v1alpha1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Reconciler reconciles a MetricsTrait object
type Reconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

var ReconcileWaitResult = reconcile.Result{RequeueAfter: 30 * time.Second}

var pool = make(map[v1alpha1.TypedReference]chan struct{})

// +kubebuilder:rbac:groups=labs.bitnami.com,resources=kubewatchs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=labs.bitnami.com,resources=kubewatchs/status,verbs=get;update;patch
func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	mLog := r.Log.WithValues("kubewatch", req.NamespacedName)
	mLog.Info("Reconcile start")
	// fetch the trait
	var kubewatch v1alpha1.KubeWatch
	if err := r.Get(ctx, req.NamespacedName, &kubewatch); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	wtype := kubewatch.Spec.WorkloadReference
	wtype.Name = ""
	oldstop, ok := pool[wtype]
	if ok {
		close(oldstop)
	}
	hander := config.Handler{}
	if kubewatch.Spec.Handler.Flock != nil {
		hander.Flock = *kubewatch.Spec.Handler.Flock
	}
	if kubewatch.Spec.Handler.Hipchat != nil {
		hander.Hipchat = *kubewatch.Spec.Handler.Hipchat
	}
	if kubewatch.Spec.Handler.Mattermost != nil {
		hander.Mattermost = *kubewatch.Spec.Handler.Mattermost
	}
	if kubewatch.Spec.Handler.MSTeams != nil {
		hander.MSTeams = *kubewatch.Spec.Handler.MSTeams
	}
	if kubewatch.Spec.Handler.Slack != nil {
		hander.Slack = *kubewatch.Spec.Handler.Slack
	}
	if kubewatch.Spec.Handler.SMTP != nil {
		hander.SMTP = *kubewatch.Spec.Handler.SMTP
	}
	if kubewatch.Spec.Handler.Webhook != nil {
		hander.Webhook = *kubewatch.Spec.Handler.Webhook
	}

	stop := make(chan struct{})
	pool[wtype] = stop
	conf := &config.Config{
		Handler: hander,
		CRD: &config.TypedReference{
			APIVersion: wtype.APIVersion,
			Kind:       wtype.Kind,
		},
		Namespace: kubewatch.Spec.Namespace,
	}
	var eventHandler = client2.ParseEventHandler(conf)
	go controller.Start(conf, eventHandler, stop)
	kubewatch.Status.Watching = true
	return ctrl.Result{}, r.Client.Status().Update(ctx, &kubewatch)
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.KubeWatch{}).
		Complete(r)
}

// Setup adds a controller that reconciles MetricsTrait.
func Setup(mgr ctrl.Manager) error {
	reconciler := Reconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("KubeWatch"),
		Scheme: mgr.GetScheme(),
	}
	return reconciler.SetupWithManager(mgr)
}
