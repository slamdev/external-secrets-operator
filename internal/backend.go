package internal

import (
	"fmt"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
)

type Backend interface {
	Connect(properties map[string]string) error
	GetValues(key string) (map[string]string, error)
}

type BackendFactory interface {
	Create(backendType string, backendName string, properties map[string]string) error
	Get(backendName string) (Backend, error)
}

type backendFactory struct {
	backends map[string]Backend
	log      logr.Logger
}

func NewBackendFactory() BackendFactory {
	return &backendFactory{
		backends: make(map[string]Backend),
		log:      ctrl.Log.WithName("backendFactory"),
	}
}

func (b *backendFactory) Create(backendType string, backendName string, properties map[string]string) error {
	var backend Backend
	switch backendType {
	case "Consul":
		backend = NewConsulBackend()
	case "Vault":
		backend = NewVaultBackend()
	default:
		return fmt.Errorf("%s backend type is not known", backendType)
	}
	if err := backend.Connect(properties); err != nil {
		return err
	}
	b.backends[backendName] = backend
	return nil
}

func (b *backendFactory) Get(backendName string) (Backend, error) {
	if backend, ok := b.backends[backendName]; ok {
		return backend, nil
	}
	return nil, fmt.Errorf("backend with %s name is not found", backendName)
}
