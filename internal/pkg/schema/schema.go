package schema

import (
	"fmt"
	"io"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/pkg/cmd/util"
)

func NewSchemaList(
	config *configuration.Configuration,
	subjectName string,
	version string,
) (string, error) {
	params := fmt.Sprintf("subjects/%v/versions/%v/schema", subjectName, version)

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
