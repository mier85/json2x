package main

import (
	"github.com/mier85/json2x/pkg/runner"

	"encoding/json"
	"strings"
)

func toExportEnv(name string, args []string, values map[string]*json.RawMessage) ([]byte, string, error) {
	keyValue := make(map[string]string)
	for k, v := range values {
		value := string([]byte(*v))
		// escape quotes
		value = strings.Replace(value, `"`, `\"`, -1)
		// replace line endings
		value = strings.Replace(value, "\n", "\\", -1)
		keyValue[strings.ToUpper(k)] = value
	}
	return nil, "", nil
}

func main() {
	runner.Run()
}
