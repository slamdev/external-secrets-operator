package internal

import (
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

type vault struct {
	log logr.Logger
}

func NewVaultBackend() Backend {
	return &vault{
		log: ctrl.Log.WithName("backend").WithName("vault"),
	}
}

func (v *vault) Connect(properties map[string]interface{}) error {
	v.log.Info("connecting", "properties", properties)
	return nil
}

func (v *vault) GetValue(key string) (map[string]interface{}, error) {
	v.log.Info("getting", "key", key)
	return make(map[string]interface{}), nil
}
