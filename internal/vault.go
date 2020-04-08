package internal

import (
	"fmt"
	"github.com/go-logr/logr"
	"github.com/hashicorp/vault/api"
	"os"
	"regexp"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

type vault struct {
	log         logr.Logger
	logical     *api.Logical
	alphaNumReg *regexp.Regexp
}

func NewVaultBackend() Backend {
	alphaNumReg := regexp.MustCompile("[^a-zA-Z0-9]+")
	return &vault{
		log:         ctrl.Log.WithName("backend").WithName("vault"),
		alphaNumReg: alphaNumReg,
	}
}

func (v *vault) Connect(properties map[string]string) error {
	if _, ok := properties["addr"]; ok {
		if err := os.Setenv(api.EnvVaultAddress, properties["addr"]); err != nil {
			return err
		}
	}
	if _, ok := properties["token"]; ok {
		if err := os.Setenv(api.EnvVaultToken, properties["token"]); err != nil {
			return err
		}
	}
	cfg := api.DefaultConfig()
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	v.logical = client.Logical()
	return nil
}

func (v *vault) GetValues(key string) (map[string]string, error) {
	v.log.V(0).Info("processing", "key", key)
	data, err := v.getSingleSecret(key)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		data, err = v.getNestedSecrets(key)
		if err != nil {
			return nil, err
		}
	}
	var res = make(map[string]string)
	for k, n := range data {
		e := strings.Replace(k, key, "", 1)
		e = strings.TrimPrefix(e, "/")
		e = v.alphaNumReg.ReplaceAllString(e, "_")
		e = strings.ToUpper(e)
		v.log.V(0).Info("converted", "key", k, "env", e)
		res[e] = n
	}
	return res, nil
}

func (v *vault) asEnvVar(key string) string {
	value := v.alphaNumReg.ReplaceAllString(key, "_")
	return strings.ToUpper(value)
}

func (v *vault) getSingleSecret(key string) (map[string]string, error) {
	var res = make(map[string]string)
	s, err := v.logical.Read(key)
	if err != nil {
		return res, err
	}
	if s == nil {
		return res, nil
	}
	if s.Data["data"] == nil {
		v.log.Info("got nil in [data]. skipping", "key", key)
		return res, nil
	}
	data := s.Data["data"].(map[string]interface{})
	for k, v := range data {
		res[k] = fmt.Sprint(v)
	}
	return res, nil
}

func (v *vault) getNestedSecrets(key string) (map[string]string, error) {
	var res = make(map[string]string)
	metaKey := strings.Replace(key, "/data/", "/metadata/", 1)
	s, err := v.logical.List(metaKey)
	if err != nil {
		return res, err
	}
	if s == nil {
		return res, nil
	}
	keys := s.Data["keys"].([]interface{})
	for _, i := range keys {
		k := i.(string)
		nestedKey := key + "/" + strings.TrimSuffix(k, "/")
		if strings.HasSuffix(k, "/") {
			nested, err := v.getNestedSecrets(nestedKey)
			if err != nil {
				return nil, err
			}
			res = v.mergeMaps(res, nested)
		} else {
			nested, err := v.getSingleSecret(nestedKey)
			if err != nil {
				return nil, err
			}
			mod := make(map[string]string)
			for n := range nested {
				mod[nestedKey+"/"+n] = nested[n]
			}
			res = v.mergeMaps(res, mod)
		}
	}
	return res, nil
}

func (v *vault) mergeMaps(m1, m2 map[string]string) map[string]string {
	merged := m1
	for k, v := range m2 {
		merged[k] = v
	}
	return merged
}
