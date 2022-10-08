package templates

const ControllerTemp = `package {{ .ToLower .Kind }}

import (
	"context"
	"time"

	"golang.org/x/time/rate"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	{{.Group}}{{.Version}} "{{.Repo}}/api/{{.Version}}"
)

// return a new {{.Kind}} reconciler.
func New{{.Kind}}Reconciler(mgr manager.Manager) *{{.Kind}}Reconciler {
	return &{{.Kind}}Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}
}

// {{.Kind}}Reconciler reconciles a {{.Kind}} object
type {{.Kind}}Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups={{.Group}}.{{.Domain}},resources={{.Kind}},verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups={{.Group}}.{{.Domain}},resources={{.Kind}}/status,verbs=get;update;patch
// +kubebuilder:rbac:groups={{.Group}}.{{.Domain}},resources={{.Kind}}/finalizers,verbs=update

// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *{{.Kind}}Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	klog.V(3).Info("reconcile {{.Kind}}")
	// TODO: reconcile, change it
    return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *{{.Kind}}Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&{{.Group}}{{.Version}}.{{.Kind}}{}).
		WithOptions(controller.Options{RateLimiter: workqueue.NewMaxOfRateLimiter(
			workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 1000*time.Second),
			// 10 qps, 100 bucket size for default ratelimiter workqueue
			&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(10), 100)},
		)}).
		Complete(r)
}
`
