package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	externalsecretsoperatorv1alpha1 "external-secrets-operator/api/v1alpha1"
)

// ExternalBackendReconciler reconciles a ExternalBackend object
type ExternalBackendReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=external-secrets-operator.slamdev.net,resources=externalbackends,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=external-secrets-operator.slamdev.net,resources=externalbackends/status,verbs=get;update;patch

func (r *ExternalBackendReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("externalbackend", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *ExternalBackendReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&externalsecretsoperatorv1alpha1.ExternalBackend{}).
		Complete(r)
}
