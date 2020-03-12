package internal

import "fmt"

type Backend interface {
	Connect(properties map[string]interface{}) error
	GetValue(key string) (map[string]interface{}, error)
}

type BackendFactory interface {
	Create(backendType string, backendName string, properties map[string]interface{}) error
	Get(backendName string) (Backend, error)
}

type backendFactory struct {
	backends map[string]Backend
}

func NewBackendFactory() BackendFactory {
	return &backendFactory{
		backends: make(map[string]Backend),
	}
}

func (b *backendFactory) Create(backendType string, backendName string, properties map[string]interface{}) error {
	var backend Backend
	switch backendType {
	case "consul":
		backend = NewConsulBackend()
	case "vault":
		backend = NewVaultBackend()
	default:
		return fmt.Errorf("%s backend type is not known", backendType)
	}
	if err := backend.Connect(properties); err != nil {
		return err
	}
	b.backends[backendType] = backend
	return nil
}

func (b *backendFactory) Get(backendName string) (Backend, error) {
	if backend, ok := b.backends[backendName]; ok {
		return backend, nil
	}
	return nil, fmt.Errorf("backend with %s name is not found", backendName)
}
