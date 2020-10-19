package v1alpha1

import (
	runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/core/v1alpha1"
	"github.com/crossplane/oam-kubernetes-runtime/pkg/oam"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KubeWatchSpec defines the desired state of KubeWatch
type KubeWatchSpec struct {
	// WorkloadReference to the workload whose metrics needs to be exposed
	WorkloadReference TypedReference `json:"workloadRef"`

	Handler   Handler `json:"handler"`
	Namespace string  `json:"namespace,omitempty"`
}

type TypedReference struct {
	// APIVersion of the referenced object.
	APIVersion string `json:"apiVersion"`

	// Kind of the referenced object.
	Kind string `json:"kind"`

	// Name of the referenced object.
	Name string `json:"name,omitempty"`
}

// Handler contains handler configuration
type Handler struct {
	Slack      *Slack      `json:"slack,omitempty"`
	Hipchat    *Hipchat    `json:"hipchat,omitempty"`
	Mattermost *Mattermost `json:"mattermost,omitempty"`
	Flock      *Flock      `json:"flock,omitempty"`
	Webhook    *Webhook    `json:"webhook,omitempty"`
	MSTeams    *MSTeams    `json:"msteams,omitempty"`
	SMTP       *SMTP       `json:"smtp,omitempty"`
}

// Slack contains slack configuration
type Slack struct {
	// Slack "legacy" API token.
	Token string `json:"token"`
	// Slack channel.
	Channel string `json:"channel"`
	// Title of the message.
	Title string `json:"title"`
}

// Hipchat contains hipchat configuration
type Hipchat struct {
	// Hipchat token.
	Token string `json:"token"`
	// Room name.
	Room string `json:"room"`
	// URL of the hipchat server.
	Url string `json:"url"`
}

// Mattermost contains mattermost configuration
type Mattermost struct {
	Channel  string `json:"room"`
	Url      string `json:"url"`
	Username string `json:"username"`
}

// Flock contains flock configuration
type Flock struct {
	// URL of the flock API.
	Url string `json:"url"`
}

// Webhook contains webhook configuration
type Webhook struct {
	// Webhook URL.
	Url string `json:"url"`
}

// MSTeams contains MSTeams configuration
type MSTeams struct {
	// MSTeams API Webhook URL.
	WebhookURL string `json:"webhookurl"`
}

// SMTP contains SMTP configuration.
type SMTP struct {
	// Destination e-mail address.
	To string `json:"to" yaml:"to,omitempty"`
	// Sender e-mail address .
	From string `json:"from" yaml:"from,omitempty"`
	// Smarthost, aka "SMTP server"; address of server used to send email.
	Smarthost string `json:"smarthost" yaml:"smarthost,omitempty"`
	// Subject of the outgoing emails.
	Subject string `json:"subject" yaml:"subject,omitempty"`
	// Extra e-mail headers to be added to all outgoing messages.
	Headers map[string]string `json:"headers" yaml:"headers,omitempty"`
	// Authentication parameters.
	Auth SMTPAuth `json:"auth" yaml:"auth,omitempty"`
	// If "true" forces secure SMTP protocol (AKA StartTLS).
	RequireTLS bool `json:"requireTLS" yaml:"requireTLS"`
	// SMTP hello field (optional)
	Hello string `json:"hello" yaml:"hello,omitempty"`
}

type SMTPAuth struct {
	// Username for PLAN and LOGIN auth mechanisms.
	Username string `json:"username" yaml:"username,omitempty"`
	// Password for PLAIN and LOGIN auth mechanisms.
	Password string `json:"password" yaml:"password,omitempty"`
	// Identity for PLAIN auth mechanism
	Identity string `json:"identity" yaml:"identity,omitempty"`
	// Secret for CRAM-MD5 auth mechanism
	Secret string `json:"secret" yaml:"secret,omitempty"`
}

// WatchStatus defines the observed state of KubeWatch
type WatchStatus struct {
	runtimev1alpha1.ConditionedStatus `json:",inline"`
	Watching                          bool `json:"watching,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:categories={oam}
// +kubebuilder:subresource:status
// KubeWatch is the Schema for the routes API
type KubeWatch struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KubeWatchSpec `json:"spec,omitempty"`
	Status WatchStatus   `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// KubeWatchList contains a list of KubeWatch
type KubeWatchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KubeWatch `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KubeWatch{}, &KubeWatchList{})
}

var _ oam.Trait = &KubeWatch{}

func (r *KubeWatch) SetConditions(c ...runtimev1alpha1.Condition) {
	r.Status.SetConditions(c...)
}

func (r *KubeWatch) GetCondition(c runtimev1alpha1.ConditionType) runtimev1alpha1.Condition {
	return r.Status.GetCondition(c)
}

func (r *KubeWatch) GetWorkloadReference() runtimev1alpha1.TypedReference {
	return runtimev1alpha1.TypedReference{
		APIVersion: r.Spec.WorkloadReference.APIVersion,
		Kind:       r.Spec.WorkloadReference.Kind,
		Name:       r.Spec.WorkloadReference.Name,
	}
}

func (r *KubeWatch) SetWorkloadReference(rt runtimev1alpha1.TypedReference) {
	r.Spec.WorkloadReference = TypedReference{
		APIVersion: rt.APIVersion,
		Kind:       rt.Kind,
		Name:       rt.Name,
	}
}
