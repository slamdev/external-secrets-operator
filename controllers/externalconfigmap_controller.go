package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	externalsecretsoperatorv1alpha1 "external-secrets-operator/api/v1alpha1"
)

// ExternalConfigMapReconciler reconciles a ExternalConfigMap object
type ExternalConfigMapReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=external-secrets-operator.slamdev.net,resources=externalconfigmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=external-secrets-operator.slamdev.net,resources=externalconfigmaps/status,verbs=get;update;patch

func (r *ExternalConfigMapReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("externalconfigmap", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *ExternalConfigMapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&externalsecretsoperatorv1alpha1.ExternalConfigMap{}).
		Complete(r)
}
