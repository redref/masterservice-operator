package masterservice

import (
	"context"

	blablacarv1 "github.com/blablacar/masterservice-operator/pkg/apis/blablacar/v1"

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

var log = logf.Log.WithName("controller_masterservice")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new MasterService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMasterService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("masterservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource MasterService
	err = c.Watch(&source.Kind{Type: &blablacarv1.MasterService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Services and requeue the owner MasterService
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &blablacarv1.MasterService{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Endpoints and requeue the owner MasterService
	err = c.Watch(&source.Kind{Type: &corev1.Endpoints{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &blablacarv1.MasterService{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileMasterService{}

// ReconcileMasterService reconciles a MasterService object
type ReconcileMasterService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a MasterService object and makes changes based on the state read
// and what is in the MasterService.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileMasterService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling MasterService")

	// Fetch the MasterService instance
	instance := &blablacarv1.MasterService{}
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

	// Check or create catchall service => <MasterService.Name>-all
	allService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name + "-all", Namespace: instance.Namespace}, allService)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		dep := r.serviceForMasterservice(instance, false)
		reqLogger.Info("Creating a new Service", "Service.Namespace", dep.Namespace, "Service.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", dep.Namespace, "Service.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// Service created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Service")
		return reconcile.Result{}, err
	}
	err = controllerutil.SetControllerReference(instance, allService, r.scheme)
	if err != nil {
		reqLogger.Error(err, "We do not own service", "Service.Namespace", allService.Namespace, "Service.Name", allService.Name)
		return reconcile.Result{}, err
	}

	// Check or create empty service => <MasterService.Name>
	masterService := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, masterService)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		dep := r.serviceForMasterservice(instance, true)
		reqLogger.Info("Creating a new Service", "Service.Namespace", dep.Namespace, "Service.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Service", "Service.Namespace", dep.Namespace, "Service.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// Service created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Service")
		return reconcile.Result{}, err
	}
	err = controllerutil.SetControllerReference(instance, masterService, r.scheme)
	if err != nil {
		reqLogger.Error(err, "We do not own service", "Service.Namespace", masterService.Namespace, "Service.Name", masterService.Name)
		return reconcile.Result{}, err
	}

	// Get allService endpoint and set ownership
	allServiceEndpoint := &corev1.Endpoints{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: allService.Name, Namespace: allService.Namespace}, allServiceEndpoint)
	if err != nil {
		reqLogger.Error(err, "Failed to get Endpoint", "Endpoint.Namespace", allService.Namespace, "Endpoint.Name", allService.Name)
		return reconcile.Result{}, err
	}
	ownerLen := len(allServiceEndpoint.OwnerReferences)
	err = controllerutil.SetControllerReference(instance, allServiceEndpoint, r.scheme)
	if err == nil && ownerLen != len(allServiceEndpoint.OwnerReferences) {
		err = r.client.Update(context.TODO(), allServiceEndpoint)
		if err != nil {
			reqLogger.Error(err, "Failed set owner for Endpoint", "Endpoint.Namespace", allServiceEndpoint.Namespace, "Endpoint.Name", allServiceEndpoint.Name)
			return reconcile.Result{}, err
		}
	}

	// Sort and get oldest address (from pod StartTime TOFIX ?)
	var oldest *metav1.Time
	var oldestAddress corev1.EndpointAddress
	for _, subset := range allServiceEndpoint.Subsets {
		for _, address := range subset.Addresses {
			pod := &corev1.Pod{}
			err = r.client.Get(context.TODO(), types.NamespacedName{Name: address.TargetRef.Name, Namespace: address.TargetRef.Namespace}, pod)
			if err != nil {
				reqLogger.Error(err, "Failed to get Pod", "Pod.Namespace", address.TargetRef.Namespace, "pod.Name", address.TargetRef.Name)
				return reconcile.Result{}, err
			}
			if oldest == nil || oldest.UnixNano() > pod.Status.StartTime.UnixNano() {
				oldest = pod.Status.StartTime
				oldestAddress = address
			}
		}
	}
	// all endpoint is empty: exit success
	if oldest == nil {
		return reconcile.Result{}, nil
	}

	// Create or update our endpoint
	masterEndpoint := &corev1.Endpoints{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, masterEndpoint)
	if err != nil && errors.IsNotFound(err) {
		dep := r.endpointForMasterService(instance, masterService.Name, masterService.Namespace, allServiceEndpoint, oldestAddress)
		reqLogger.Info("Creating a new Endpoint", "Endpoint.Namespace", dep.Namespace, "Endpoint.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Endpoint", "Endpoint.Namespace", dep.Namespace, "Endpoint.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// Endpoint created successfully - return success
		return reconcile.Result{}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Service")
		return reconcile.Result{}, err
	}

	// We need to update the endpoint
	dep := r.endpointForMasterService(instance, masterService.Name, masterService.Namespace, allServiceEndpoint, oldestAddress)
	reqLogger.Info("Update an Endpoint", "Endpoint.Namespace", dep.Namespace, "Endpoint.Name", dep.Name, "Address", oldestAddress)
	err = r.client.Update(context.TODO(), dep)
	if err != nil {
		reqLogger.Error(err, "Failed to update Endpoint", "Endpoint.Namespace", dep.Namespace, "Endpoint.Name", dep.Name)
		return reconcile.Result{}, err
	}
	// Endpoint updated successfully - return success
	return reconcile.Result{}, nil
}

// serviceForMasterService returns a masterservice Service object
func (r *ReconcileMasterService) serviceForMasterservice(m *blablacarv1.MasterService, master bool) *corev1.Service {
	template := m.Spec
	var name string
	if master {
		name = m.Name
	} else {
		name = m.Name + "-all"
	}

	// Override selector in order to push enpoint in this controller
	if master {
		template.Selector = nil
	}

	dep := &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "core/v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: m.Namespace,
		},
		Spec: template,
	}
	// Set Memcached instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// endpointForMasterService returns a masterservice Endpoint object
func (r *ReconcileMasterService) endpointForMasterService(m *blablacarv1.MasterService, name string, namespace string, e *corev1.Endpoints, a corev1.EndpointAddress) *corev1.Endpoints {
	dep := &corev1.Endpoints{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "core/v1",
			Kind:       "Endpoint",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: e.Namespace,
		},
		Subsets: []corev1.EndpointSubset{{
			Ports:     e.Subsets[0].Ports,
			Addresses: []corev1.EndpointAddress{a},
		}},
	}
	// Set Memcached instance as the owner and controller
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}
