package controllers

import (
	"context"
	"encoding/base64"
	"external-secrets-operator/internal"
	v1 "k8s.io/api/core/v1"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	externalsecretsoperatorv1alpha1 "external-secrets-operator/api/v1alpha1"
)

// ExternalBackendReconciler reconciles a ExternalBackend object
type ExternalBackendReconciler struct {
	client.Client
	Log            logr.Logger
	Scheme         *runtime.Scheme
	BackendFactory internal.BackendFactory
}

// +kubebuilder:rbac:groups=external-secrets-operator.slamdev.net,resources=externalbackends,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=external-secrets-operator.slamdev.net,resources=externalbackends/status,verbs=get;update;patch

func (r *ExternalBackendReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("externalbackend", req.NamespacedName)

	var externalBackend externalsecretsoperatorv1alpha1.ExternalBackend
	if err := r.Get(ctx, req.NamespacedName, &externalBackend); err != nil {
		log.Error(err, "unable to fetch ExternalBackend")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.V(0).Info("reconcile", "externalBackend", externalBackend)

	externalBackend.Status.Connected = new(bool)
	*externalBackend.Status.Connected = false

	if err := r.Status().Update(ctx, &externalBackend); err != nil {
		log.Error(err, "unable to update ExternalBackend status")
		return ctrl.Result{}, err
	}

	var secret v1.Secret
	if err := r.Get(ctx, client.ObjectKey{
		Namespace: externalBackend.Namespace,
		Name:      externalBackend.Spec.SecretName,
	}, &secret); err != nil {
		return ctrl.Result{}, err
	}

	log.V(0).Info("reconcile", "secret", secret)

	properties := make(map[string]string)
	for key, encodedValue := range secret.Data {
		var value []byte
		_, err := base64.StdEncoding.Decode(value, encodedValue)
		if err != nil {
			return ctrl.Result{}, err
		}
		properties[key] = string(value)
	}

	if err := r.BackendFactory.Create(string(externalBackend.Spec.Type), externalBackend.Name, secret.StringData); err != nil {
		return ctrl.Result{}, err
	}

	*externalBackend.Status.Connected = true

	if err := r.Status().Update(ctx, &externalBackend); err != nil {
		log.Error(err, "unable to update ExternalBackend status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ExternalBackendReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&externalsecretsoperatorv1alpha1.ExternalBackend{}).
		Complete(r)
}
