package internal

import (
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

type consul struct {
	log logr.Logger
}

func NewConsulBackend() Backend {
	return &consul{
		log: ctrl.Log.WithName("backend").WithName("consul"),
	}
}

func (v *consul) Connect(properties map[string]interface{}) error {
	v.log.Info("connecting", "properties", properties)
	return nil
}

func (v *consul) GetValue(key string) (map[string]interface{}, error) {
	v.log.Info("getting", "key", key)
	return make(map[string]interface{}), nil
}
