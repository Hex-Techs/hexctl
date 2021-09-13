package templates

const ControllerTemp = `package {{ .ToLower .Kind }}

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	{{.Group}}{{.Version}} "{{.Repo}}/api/{{.Version}}"
	clientset "{{.Repo}}/pkg/generated/clientset/versioned"
	{{.Group}}scheme "{{.Repo}}/pkg/generated/clientset/versioned/scheme"
	informers "{{.Repo}}/pkg/generated/informers/externalversions/{{.Group}}/v1alpha1"
	listers "{{.Repo}}/pkg/generated/listers/{{.Group}}/v1alpha1"
)

const controllerAgentName = "{{.Group}}-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a {{.Kind}} is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a {{.Kind}} fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by {{.Kind}}"
	// MessageResourceSynced is the message used for an Event fired when a {{.Kind}}
	// is synced successfully
	MessageResourceSynced = "{{.Kind}} synced successfully"
)

// Controller is the controller implementation for {{.Kind}} resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// {{.Group}}clientset is a clientset for our own API group
	{{.Group}}clientset clientset.Interface

	{{.Kind}}sLister listers.{{.Kind}}Lister
	{{.Kind}}sSynced cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	workqueue workqueue.RateLimitingInterface
	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new {{.Group}} controller
func NewController(
	kubeclientset kubernetes.Interface,
	{{.Group}}clientset clientset.Interface,
	{{.Kind}}Informer informers.{{.Kind}}Informer) *Controller {

	// Create event broadcaster
	// Add {{.Group}}-controller types to the default Kubernetes Scheme so Events can be
	// logged for {{.Group}}-controller types.
	utilruntime.Must({{.Group}}scheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:   kubeclientset,
		{{.Group}}clientset: {{.Group}}clientset,
		{{.Kind}}sLister:      {{.Kind}}Informer.Lister(),
		{{.Kind}}sSynced:      {{.Kind}}Informer.Informer().HasSynced,
		workqueue:       workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "{{.Kind}}s"),
		recorder:        recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when {{.Kind}} resources change
	{{.Kind}}Informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueue{{.Kind}},
		UpdateFunc: func(old, new interface{}) {
			controller.enqueue{{.Kind}}(new)
		},
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting {{.Kind}} controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.{{.Kind}}sSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process {{.Kind}} resources
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// {{.Kind}} resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the {{.Kind}} resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	{{ if .UseNamespace -}}namespace{{- else -}}_{{- end -}}, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the {{.Kind}} resource with this namespace/name
	{{.Kind}}, err := c.{{.Kind}}sLister.{{- if .UseNamespace -}}{{.Kind}}s(namespace).{{- end -}}Get(name)
	if err != nil {
		// The {{.Kind}} resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("{{.Kind}} '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}


	// Finally, we update the status block of the {{.Kind}} resource to reflect the
	// current state of the world
	err = c.update{{.Kind}}Status({{.Kind}})
	if err != nil {
		return err
	}

	c.recorder.Event({{.Kind}}, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) update{{.Kind}}Status({{.Kind}} *{{.Group}}v1alpha1.{{.Kind}}) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	{{.Kind}}Copy := {{.Kind}}.DeepCopy()
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the {{.Kind}} resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.{{.Group}}clientset.{{.ToTitle .Group}}{{.ToTitle .Version}}().{{.Kind}}s({{- if .UseNamespace -}}{{.Kind}}.Namespace{{- end -}}).Update(context.TODO(), {{.Kind}}Copy, metav1.UpdateOptions{})
	return err
}

// enqueue{{.Kind}} takes a {{.Kind}} resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than {{.Kind}}.
func (c *Controller) enqueue{{.Kind}}(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the {{.Kind}} resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that {{.Kind}} resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	klog.V(4).Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a {{.Kind}}, we should not do anything more
		// with it.
		if ownerRef.Kind != "{{.Kind}}" {
			return
		}

		{{.Kind}}, err := c.{{.Kind}}sLister.{{- if .UseNamespace -}}{{.Kind}}s(namespace).{{- end -}}Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned object '%s' of {{.Kind}} '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueue{{.Kind}}({{.Kind}})
		return
	}
}
`
