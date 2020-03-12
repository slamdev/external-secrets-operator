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

func (v *vault) Connect(properties map[string]string) error {
	return nil
}

func (v *vault) GetValues(key string) (map[string]string, error) {
	values := make(map[string]string)
	values["STR"] = "some-string-value"
	values["INT"] = "123456"
	values["BOOL"] = "true"
	return values, nil
}
