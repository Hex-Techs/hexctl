package templates

const (
	Types = `package {{.Version}}

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// {{ .Kind }}Spec defines the desired state of {{ .Kind }}
type {{ .Kind }}Spec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Foo is an example field of {{ .Kind }}. Edit {{ .Kind }}_types.go to remove/update
	Foo string ` + "`" + `json:"foo,omitempty"` + "`" + `
}

// {{ .Kind }}Status defines the observed state of {{ .Kind }}
type {{ .Kind }}Status struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +genclient
{{- if not .Status }}
// +genclient:noStatus
{{- end}}
{{- if not .Namespaced }}
// +genclient:nonNamespaced
{{- end }}
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true
{{- if .Status }}
// +kubebuilder:subresource:status
{{- end}}
{{- if not .Namespaced }}
// +kubebuilder:resource:scope=Cluster
{{- end }}

// {{ .Kind }} is the Schema for the {{ .Kind }} API
type {{ .Kind }} struct {
	metav1.TypeMeta   ` + "`" + `json:",inline"` + "`" + `
	metav1.ObjectMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `
	Spec              {{ .Kind }}Spec   ` + "`" + `json:"spec,omitempty"` + "`" + `
	Status            {{ .Kind }}Status ` + "`" + `json:"status,omitempty"` + "`" + `
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:object:root=true

// {{ .Kind }}List contains a list of {{ .Kind }}
type {{ .Kind }}List struct {
	metav1.TypeMeta ` + "`" + `json:",inline"` + "`" + `
	metav1.ListMeta ` + "`" + `json:"metadata,omitempty"` + "`" + `
	Items           []{{ .Kind }} ` + "`" + `json:"items"` + "`" + `
}

func init() {
	SchemeBuilder.Register(&{{ .Kind }}{}, &{{ .Kind }}List{})
}
`
	Groupversion = `// +kubebuilder:object:generate=true
// +k8s:deepcopy-gen=package
// +groupName={{.Group}}.{{.Domain}}
package {{.Version}}

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// GroupVersion is group version used to register these objects
	GroupVersion = schema.GroupVersion{Group: "{{.Group}}.{{.Domain}}", Version: "{{.Version}}"}

	SchemeGroupVersion = GroupVersion

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: GroupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

func Kind(kind string) schema.GroupKind {
	return GroupVersion.WithKind(kind).GroupKind()
}

func Resource(resource string) schema.GroupResource {
	return GroupVersion.WithResource(resource).GroupResource()
}
`
)
