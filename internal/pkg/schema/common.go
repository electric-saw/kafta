package schema

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
)

func BuildGetRequestSchemaRegistry(config *configuration.Configuration, path string) *http.Response {
	client := &http.Client{}
	url := fmt.Sprintf("%v/%v", config.GetContext().SchemaRegistry, path)
	req, _ := http.NewRequest("GET", url, nil)

	if config.GetContext().SchemaRegistryAuth.Secret != "" && config.GetContext().SchemaRegistryAuth.Key != "" {
		req.Header.Add("Authorization", "Basic "+basicAuth(config.GetContext().SchemaRegistryAuth.Secret, config.GetContext().SchemaRegistryAuth.Key))
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	return resp
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
