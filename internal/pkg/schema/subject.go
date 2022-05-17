package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
)

type subjects []string

func NewSubjectList(config *configuration.Configuration) ([]string, error) {
	resp := BuildGetRequestSchemaRegistry(config, "subjects")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
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

func NewSubjecVersion(config *configuration.Configuration, subsubjectName string) string {
	params := fmt.Sprintf("subjects/%v/versions", subsubjectName)

	resp := BuildGetRequestSchemaRegistry(config, params)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}
