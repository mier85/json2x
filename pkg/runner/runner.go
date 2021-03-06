package runner

import (
	"github.com/mier85/json2x/pkg/converters"
	"io/ioutil"

	"bytes"
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var availableConverters = map[string]Converter{
	"k8ssecret":  &converters.K8SSecret{},
	"bashexport": &converters.Bash{},
}

type Converter interface {
	FlagSet() *flag.FlagSet
	Template() *template.Template
	Filename(string) string
	Values(name string, values map[string]*json.RawMessage) (interface{}, error)
}

func mustGetValues() (Converter, map[string]*json.RawMessage) {
	var format = flag.String("format", "k8ssecret", "format to convert to")
	flag.Parse()
	converter, ok := availableConverters[*format]
	if !ok {
		log.Fatalf("target format does not exist")
	}

	if flag.NArg() < 1 {
		log.Fatal("expected at least one parameter")
	}
	f, err := os.Open(flag.Arg(0))
	if nil != err {
		log.Fatalf("failed opening input file: %s", err.Error())
	}
	var target map[string]*json.RawMessage
	err = json.NewDecoder(f).Decode(&target)
	f.Close()
	if nil != err {
		log.Fatalf("failed parsing json object: %s", err.Error())
	}
	return converter, target
}

func convert(tpl *template.Template, val interface{}) ([]byte, error) {
	target := &bytes.Buffer{}
	err := tpl.Execute(target, val)
	return target.Bytes(), err
}

func execute(converter Converter, target map[string]*json.RawMessage) {
	set := converter.FlagSet()
	err := set.Parse(flag.Args()[1:])
	if nil != err {
		log.Fatalf("failed parsing flags for converter: %s", err.Error())
	}

	name := extractNameFromFilename(flag.Arg(0))
	values, err := converter.Values(name, target)
	if nil != err {
		log.Fatalf("failed getting values for template: %s", err.Error())
	}
	raw, err := convert(converter.Template(), values)
	if nil != err {
		log.Fatalf("failed executing template: %s", err.Error())
	}

	if err := ioutil.WriteFile(converter.Filename(name), raw, os.ModePerm); nil != err {
		log.Fatalf("failed writing file: %s", err.Error())
	}
}

func extractNameFromFilename(filename string) string {
	name := filepath.Base(filename)
	index := strings.Index(name, ".")
	if index != -1 {
		return name[0:index]
	}

	return name
}

func Run() {
	execute(mustGetValues())
}
