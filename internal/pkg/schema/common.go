package schema

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/pkg/cmd/util"
)

func BuildGetRequestSchemaRegistry(config *configuration.Configuration, path string,) *http.Response {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)

	defer cancel()

	client := &http.Client{}
	url := fmt.Sprintf("%v/%v", config.GetContext().SchemaRegistry, path)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil)
	if err != nil {
		util.CheckErr(fmt.Errorf("error creating request: %w", err))
	}

	if config.GetContext().SchemaRegistryAuth.Secret != "" &&
		config.GetContext().SchemaRegistryAuth.Key != "" {
		secret := config.GetContext().SchemaRegistryAuth.Secret
		key := config.GetContext().SchemaRegistryAuth.Key

		req.Header.Add(
			"Authorization",
			"Basic "+basicAuth(key, secret),
		)
	} else {
		fmt.Println("Missing Schema Registry authentication credentials")
	}

	resp, err := client.Do(req)
	if err != nil {
		util.CheckErr(err)
	}

	return resp
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func BuildPostRequestSchemaRegistry(
	config *configuration.Configuration,
	path string,
	body string,
) *http.Response {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)

	defer cancel()

	client := &http.Client{}
	url := fmt.Sprintf("%v/%v", config.GetContext().SchemaRegistry, path)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		strings.NewReader(body),
	)
	if err != nil {
		util.CheckErr(fmt.Errorf("error creating request: %w", err))
	}

	req.Header.Set("Content-Type", "application/json")

	if config.GetContext().SchemaRegistryAuth.Secret != "" &&
		config.GetContext().SchemaRegistryAuth.Key != "" {
		secret := config.GetContext().SchemaRegistryAuth.Secret
		key := config.GetContext().SchemaRegistryAuth.Key

		req.Header.Add(
			"Authorization",
			"Basic "+basicAuth(key, secret),
		)
	} else {
		fmt.Println("Missing Schema Registry authentication credentials")
	}

	resp, err := client.Do(req)
	if err != nil {
		util.CheckErr(err)
	}

	return resp
}
