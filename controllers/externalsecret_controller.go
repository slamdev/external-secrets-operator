package controllers

import (
	"context"
	"external-secrets-operator/internal"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	externalsecretsoperatorv1alpha1 "external-secrets-operator/api/v1alpha1"
)

// ExternalSecretReconciler reconciles a ExternalSecret object
type ExternalSecretReconciler struct {
	client.Client
	Log            logr.Logger
	Scheme         *runtime.Scheme
	BackendFactory internal.BackendFactory
}

// +kubebuilder:rbac:groups=external-secrets-operator.slamdev.net,resources=externalsecrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=external-secrets-operator.slamdev.net,resources=externalsecrets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=create;get;list;watch

func (r *ExternalSecretReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("externalsecret", req.NamespacedName)

	var externalSecret externalsecretsoperatorv1alpha1.ExternalSecret
	if err := r.Get(ctx, req.NamespacedName, &externalSecret); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "unable to fetch ExternalSecret")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	backend, err := r.BackendFactory.Get(externalSecret.Spec.BackendName)
	if err != nil {
		return ctrl.Result{}, err
	}
	values := make(map[string]string)
	for i := range externalSecret.Spec.Keys {
		kv, err := backend.GetValues(externalSecret.Spec.Keys[i])
		if err != nil {
			log.Error(err, "skipping", "key", externalSecret.Spec.Keys[i])
			continue
		}
		for k, v := range kv {
			values[k] = v
		}
	}

	secret, err := r.constructSecret(&externalSecret, values)
	if err != nil {
		return ctrl.Result{}, err
	}

	var existingSecret corev1.Secret
	if err := r.Get(ctx, client.ObjectKey{
		Namespace: secret.Namespace,
		Name:      secret.Name,
	}, &existingSecret); err != nil {
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, secret); err != nil {
			log.Error(err, "unable to create Secret for ExternalSecret", "secret", secret)
			return ctrl.Result{}, err
		}
	} else {
		if err := r.Update(ctx, secret); err != nil {
			return ctrl.Result{}, err
		}
	}

	if externalSecret.Status.LastSyncedTime == nil {
		externalSecret.Status.LastSyncedTime = &metav1.Time{}
	}

	*externalSecret.Status.LastSyncedTime = metav1.Now()
	if err := r.Status().Update(ctx, &externalSecret); err != nil {
		log.Error(err, "unable to update ExternalSecret status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

func (r *ExternalSecretReconciler) constructSecret(externalSecret *externalsecretsoperatorv1alpha1.ExternalSecret, values map[string]string) (*corev1.Secret, error) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        externalSecret.Name,
			Namespace:   externalSecret.Namespace,
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
		},
		StringData: values,
	}
	for k, v := range externalSecret.Annotations {
		secret.Annotations[k] = v
	}
	for k, v := range externalSecret.Labels {
		secret.Labels[k] = v
	}
	if err := ctrl.SetControllerReference(externalSecret, secret, r.Scheme); err != nil {
		return nil, err
	}
	return secret, nil
}

func (r *ExternalSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(&corev1.Secret{}, ownerKey, func(rawObj runtime.Object) []string {
		secret := rawObj.(*corev1.Secret)
		owner := metav1.GetControllerOf(secret)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != apiGVStr || owner.Kind != "ExternalSecret" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&externalsecretsoperatorv1alpha1.ExternalSecret{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
