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
	defer util.CheckErr(resp.Body.Close())

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data subjects

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func NewSubjecVersion(config *configuration.Configuration, subsubjectName string) (string, error) {
	params := fmt.Sprintf("subjects/%v/versions", subsubjectName)

	resp := BuildGetRequestSchemaRegistry(config, params)
	defer util.CheckErr(resp.Body.Close())

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), err
}
