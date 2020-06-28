package controllers

import (
	"context"
	"errors"
	externalsecretsoperatorv1alpha1 "external-secrets-operator/api/v1alpha1"
	"external-secrets-operator/internal"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

// ExternalConfigMapReconciler reconciles a ExternalConfigMap object
type ExternalConfigMapReconciler struct {
	client.Client
	Log            logr.Logger
	Scheme         *runtime.Scheme
	BackendFactory internal.BackendFactory
}

// +kubebuilder:rbac:groups=external-secrets-operator.slamdev.net,resources=externalconfigmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=external-secrets-operator.slamdev.net,resources=externalconfigmaps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

func (r *ExternalConfigMapReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("externalconfigmap", req.NamespacedName)

	var externalConfigMap externalsecretsoperatorv1alpha1.ExternalConfigMap
	if err := r.Get(ctx, req.NamespacedName, &externalConfigMap); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "unable to fetch ExternalConfigMap")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	backend, err := r.BackendFactory.Get(externalConfigMap.Spec.BackendName)
	if err != nil {
		return ctrl.Result{}, err
	}
	values := make(map[string]string)
	for i := range externalConfigMap.Spec.Keys {
		kv, err := backend.GetValues(externalConfigMap.Spec.Keys[i])
		if err != nil {
			log.Error(err, "skipping", "key", externalConfigMap.Spec.Keys[i])
			continue
		}
		for k, v := range kv {
			values[k] = v
		}
	}

	if len(values) == 0 && len(externalConfigMap.Spec.Keys) != 0 {
		err := errors.New("discovered 0 values for provided keys, assuming the backend is down")
		log.Error(err, "stopping")
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	configMap, err := r.constructConfigMap(&externalConfigMap, values)
	if err != nil {
		return ctrl.Result{}, err
	}

	var existingConfigMap corev1.ConfigMap
	if err := r.Get(ctx, client.ObjectKey{
		Namespace: configMap.Namespace,
		Name:      configMap.Name,
	}, &existingConfigMap); err != nil {
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, configMap); err != nil {
			log.Error(err, "unable to create ConfigMap for ExternalConfigMap", "configMap", configMap)
			return ctrl.Result{}, err
		}
	} else {
		if err := r.Update(ctx, configMap); err != nil {
			return ctrl.Result{}, err
		}
	}

	if externalConfigMap.Status.LastSyncedTime == nil {
		externalConfigMap.Status.LastSyncedTime = &metav1.Time{}
	}

	*externalConfigMap.Status.LastSyncedTime = metav1.Now()
	if err := r.Status().Update(ctx, &externalConfigMap); err != nil {
		log.Error(err, "unable to update ExternalConfigMap status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Minute}, nil
}

func (r *ExternalConfigMapReconciler) constructConfigMap(externalConfigMap *externalsecretsoperatorv1alpha1.ExternalConfigMap, values map[string]string) (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        externalConfigMap.Name,
			Namespace:   externalConfigMap.Namespace,
			Labels:      make(map[string]string),
			Annotations: make(map[string]string),
		},
		Data: values,
	}
	for k, v := range externalConfigMap.Annotations {
		configMap.Annotations[k] = v
	}
	for k, v := range externalConfigMap.Labels {
		configMap.Labels[k] = v
	}
	if err := ctrl.SetControllerReference(externalConfigMap, configMap, r.Scheme); err != nil {
		return nil, err
	}
	return configMap, nil
}

func (r *ExternalConfigMapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(&corev1.ConfigMap{}, ownerKey, func(rawObj runtime.Object) []string {
		configMap := rawObj.(*corev1.ConfigMap)
		owner := metav1.GetControllerOf(configMap)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != apiGVStr || owner.Kind != "ExternalConfigMap" {
			return nil
		}
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&externalsecretsoperatorv1alpha1.ExternalConfigMap{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
