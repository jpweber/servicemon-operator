package servicemonitor

import (
	"context"
	"io/ioutil"
	"os"

	monitoring "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_servicemonitor")

// Add creates a new ServiceMonitor Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileServiceMonitor{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("servicemonitor-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Service
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource servicemonitors and requeue the owner Service
	err = c.Watch(&source.Kind{Type: &monitoring.ServiceMonitor{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &corev1.Service{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileServiceMonitor{}

// ReconcileServiceMonitor reconciles a ServiceMonitor object
type ReconcileServiceMonitor struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Service object and makes changes based on the state read
// and what is in the Service.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileServiceMonitor) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling ServiceMonitor")

	// Fetch the Service instance
	instance := &corev1.Service{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Only do work on service that have the "prometheus.io/probe: true" annotation
	if instance.Annotations["prometheus.io/probe"] == "true" {
	} else {
		// we aren't doing anything with this service, just return without errors
		return reconcile.Result{}, nil
	}
	// Define a new SeviceMonitor object
	serviceMon := newServiceMon(instance)

	// Set Service instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, serviceMon, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this ServiceMonitor already exists
	found := &monitoring.ServiceMonitor{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: serviceMon.Name, Namespace: serviceMon.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ServiceMonitor", "ServiceMonitor.Namespace", serviceMon.Namespace, "ServiceMonitor.Name", serviceMon.Name)
		err = r.client.Create(context.TODO(), serviceMon)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Service Monitor created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Service Monitor already exists - don't requeue
	reqLogger.Info("Skip reconcile: ServiceMonitor already exists", "ServiceMonitor.Namespace", found.Namespace, "ServiceMonitor.Name", found.Name)
	return reconcile.Result{}, nil
}

// newServiceMon return a new service monitor object
func newServiceMon(svc *corev1.Service) *monitoring.ServiceMonitor {

	// read what namespace we are running in
	// service monitors must always be in the same namespace as prometheus-operator
	// and we must be running in that namespace as well.
	file, err := os.Open("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		log.Error(err, "Could not determine operator namespace")
	}
	defer file.Close()
	opNameSpace, err := ioutil.ReadAll(file)

	return &monitoring.ServiceMonitor{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svc.Name,
			Namespace: string(opNameSpace),
			Labels:    svc.Labels,
		},
		Spec: monitoring.ServiceMonitorSpec{
			JobLabel:        "k8s-app",  // TODO this probably shouldn't be hardcoded
			TargetLabels:    []string{}, // TODO need to figure out how to handle this
			PodTargetLabels: []string{}, // TODO need to figre out how to handle this
			Endpoints: []monitoring.Endpoint{
				monitoring.Endpoint{
					Port:            svc.Annotations["prometheus.io/port"],
					Path:            svc.Annotations["prometheus.io/path"],
					Scheme:          "http",
					Interval:        "30s",
					BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
				},
			},
			// TODO: thought/work is needed to make these more flexible
			Selector: metav1.LabelSelector{
				MatchLabels: svc.Labels,
			},
			NamespaceSelector: monitoring.NamespaceSelector{
				Any:        false,
				MatchNames: []string{svc.Namespace},
			},
			// SampleLimit: svc.Annotations["prometheus.io/sampleLimit"],
		},
	}
}
