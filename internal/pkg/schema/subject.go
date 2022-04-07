package schema

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
)

type subjects []string

func NewSubjectList(config *configuration.Configuration) []string {
	resp := BuildGetRequestSchemaRegistry(config, "subjects")
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data subjects

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
