package converters

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"text/template"
)

var k8sSecretTemplate = `apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: {{ .Name }}
  namespace: {{ .Namespace }}
data:{{range $key, $value := .KeyValue}}
  {{ $key }} : {{ $value }}{{end}}
`

var k8sSecretParsedTemplate = template.Must(template.New("k8sSecret").Parse(k8sSecretTemplate))

type K8SSecret struct {
	namespace  *string
	targetName *string
}

type k8sVars struct {
	Name      string
	Namespace string
	KeyValue  map[string]string
}

func (k *K8SSecret) FlagSet() *flag.FlagSet {
	set := flag.NewFlagSet("k8ssecret", flag.ExitOnError)
	k.namespace = set.String("namespace", "default", "k8s namespace")
	k.targetName = set.String("name", "", "k8s namespace")
	return set
}

func (k *K8SSecret) Template() *template.Template {
	return k8sSecretParsedTemplate
}

func (k *K8SSecret) Values(name string, values map[string]*json.RawMessage) (interface{}, error) {
	if *k.targetName != "" {
		name = *k.targetName
	}

	keyValue := make(map[string]string)
	for k, v := range values {
		keyValue[k] = base64.StdEncoding.EncodeToString([]byte(*v)[1 : len(*v)-1])
	}

	return k8sVars{
		Name:      name,
		Namespace: *k.namespace,
		KeyValue:  keyValue,
	}, nil
}

func (k *K8SSecret) Filename(name string) string {
	if *k.targetName != "" {
		name = *k.targetName
	}
	return name + "-k8s-secret.yaml"
}
