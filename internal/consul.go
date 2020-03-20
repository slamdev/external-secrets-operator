package internal

import (
	"github.com/go-logr/logr"
	"github.com/hashicorp/consul/api"
	"os"
	"regexp"
	ctrl "sigs.k8s.io/controller-runtime"
	"strings"
)

type consul struct {
	log         logr.Logger
	kv          *api.KV
	alphaNumReg *regexp.Regexp
}

func NewConsulBackend() Backend {
	alphaNumReg := regexp.MustCompile("[^a-zA-Z0-9]+")
	return &consul{
		log:         ctrl.Log.WithName("backend").WithName("consul"),
		alphaNumReg: alphaNumReg,
	}
}

func (v *consul) Connect(properties map[string]string) error {
	if _, ok := properties["httpAddr"]; ok {
		if err := os.Setenv(api.HTTPAddrEnvName, properties["httpAddr"]); err != nil {
			return err
		}
	}
	if _, ok := properties["httpToken"]; ok {
		if err := os.Setenv(api.HTTPTokenEnvName, properties["httpToken"]); err != nil {
			return err
		}
	}
	cfg := api.DefaultConfig()
	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	v.kv = client.KV()
	return nil
}

func (v *consul) GetValues(key string) (map[string]string, error) {
	p, _, err := v.kv.List(key, nil)
	if err != nil {
		return nil, err
	}
	values := make(map[string]string)
	singlePair := len(p) == 1
	for _, kv := range p {
		v.log.V(0).Info("processing", "key", kv.Key)
		if strings.HasSuffix(kv.Key, "/") {
			v.log.V(0).Info("skipped", "key", kv.Key)
			continue
		}
		var resKey string
		if singlePair {
			resKey = v.getLastKeyPart(kv.Key)
		} else {
			resKey = v.guessKey(kv.Key, key)
		}
		resKey = v.asEnvVar(resKey)
		v.log.V(0).Info("converted", "key", kv.Key, "env", resKey)
		values[resKey] = string(kv.Value)
	}
	return values, nil
}

func (v *consul) getLastKeyPart(key string) string {
	parts := strings.Split(key, "/")
	return parts[len(parts)-1]
}

func (v *consul) asEnvVar(key string) string {
	value := v.alphaNumReg.ReplaceAllString(key, "_")
	return strings.ToUpper(value)
}

func (v *consul) guessKey(key string, root string) string {
	// (simple/nested/key, simple/) => nested/key
	if strings.HasSuffix(root, "/") {
		return strings.Replace(key, root, "", 1)
	}
	// (simple/nested/key, simple) => nested/key
	noRoot := strings.Replace(key, root, "", 1)
	if strings.HasPrefix(noRoot, "/") {
		return strings.TrimPrefix(noRoot, "/")
	}
	// (simple/nested/key, simple/nes) => nested/key
	parentRoot := root[:strings.LastIndex(root, "/")+1]
	return strings.Replace(key, parentRoot, "", 1)
}
