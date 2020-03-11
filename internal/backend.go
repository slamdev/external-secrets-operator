package internal

type Backend interface {
	Connect(properties map[string]interface{}) error
	GetValue(key string) (map[string]interface{}, error)
}
