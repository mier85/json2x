package converters

import (
	"encoding/json"
	"flag"
	"strings"
	"text/template"
)

type bashTemplateVars struct {
	KeyValue map[string]string
}

var bashExportEnvTemplate = `#! /usr/bin/sh
{{range $key, $value := .KeyValue }}
export {{ $key }}="{{ $value }}"{{end}}`

var bashExportEnvParsedTemplate = template.Must(template.New("bashEnv").Parse(bashExportEnvTemplate))

type Bash struct {
	filename *string
}

func (b *Bash) FlagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("bash", flag.ExitOnError)
	return fs
}

func (b *Bash) Filename(name string) string {
	return name + ".sh"
}

func (b *Bash) Template() *template.Template {
	return bashExportEnvParsedTemplate
}

func (b *Bash) Values(name string, values map[string]*json.RawMessage) (interface{}, error) {
	keyValue := make(map[string]string)
	for k, v := range values {
		value := string([]byte(*v)[1 : len(*v)-1])
		// escape quotes
		value = strings.Replace(value, `"`, `\"`, -1)
		// replace line endings
		value = strings.Replace(value, "\n", "\\", -1)
		keyValue[strings.ToUpper(k)] = value
	}
	return bashTemplateVars{
		KeyValue: keyValue,
	}, nil
}
