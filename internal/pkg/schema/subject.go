package schema

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/pkg/cmd/util"
)

type subjects []string

func NewSubjectList(config *configuration.Configuration) ([]string, error) {
	resp := BuildGetRequestSchemaRegistry(config, "subjects")

	defer func() {
		if err := resp.Body.Close(); err != nil {
			util.CheckErr(err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data subjects
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	return data, nil
}

func NewSubjectVersion(config *configuration.Configuration, subsubjectName string) (string, error) {
	params := fmt.Sprintf("subjects/%v/versions", subsubjectName)

	resp := BuildGetRequestSchemaRegistry(config, params)

	defer func() {
		if err := resp.Body.Close(); err != nil {
			util.CheckErr(err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}

func NewSubjectCreate(
	config *configuration.Configuration,
	subjectName string,
	schema string,
) (string, error) {
	params := fmt.Sprintf("subjects/%v/versions", subjectName)

	resp := BuildPostRequestSchemaRegistry(config, params, schema)

	defer func() {
		if err := resp.Body.Close(); err != nil {
			util.CheckErr(err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}
