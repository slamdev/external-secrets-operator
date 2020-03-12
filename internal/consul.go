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

func (v *consul) Connect(properties map[string]string) error {
	return nil
}

func (v *consul) GetValues(key string) (map[string]string, error) {
	values := make(map[string]string)
	values["STR"] = "some-string-value"
	values["INT"] = "123456"
	values["BOOL"] = "true"
	return values, nil
}
