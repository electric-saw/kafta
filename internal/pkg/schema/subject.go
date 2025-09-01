package schema

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/electric-saw/kafta/internal/pkg/configuration"
	"github.com/electric-saw/kafta/pkg/cmd/util"
)

type subjects []string

type SubjectConfig struct {
	CompatibilityLevel string `json:"compatibilityLevel"`
}

type SubjectVersionInfo struct {
	Version      int    `json:"version"`
	ID           int    `json:"id"`
	Schema       string `json:"schema"`
	Subject      string `json:"subject"`
	References   []any  `json:"references"`
}

type SubjectWithCompatibility struct {
	Name               string `json:"name"`
	CompatibilityLevel string `json:"compatibilityLevel"`
	VersionCount       int    `json:"versionCount"`
}

type SubjectVersionWithCompatibility struct {
	Version            int    `json:"version"`
	ID                 int    `json:"id"`
	Schema             string `json:"schema"`
	Subject            string `json:"subject"`
	CompatibilityLevel string `json:"compatibilityLevel,omitempty"`
}

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

func NewSubjectCompatibility(config *configuration.Configuration, subjectName string) (string, error) {
	params := fmt.Sprintf("config/%v", subjectName)

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

	if resp.StatusCode == 404 {
		return NewGlobalCompatibility(config)
	}

	return string(body), err
}

func NewGlobalCompatibility(config *configuration.Configuration) (string, error) {
	resp := BuildGetRequestSchemaRegistry(config, "config")

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

func NewSubjectVersionsWithCompatibility(config *configuration.Configuration, subjectName string) ([]SubjectVersionWithCompatibility, error) {
	versionsJSON, err := NewSubjectVersion(config, subjectName)
	if err != nil {
		return nil, err
	}

	var versions []int
	if err := json.Unmarshal([]byte(versionsJSON), &versions); err != nil {
		return nil, err
	}

	var defaultCompatibility = "BACKWARD"
	if compatibilityJSON, err := NewSubjectCompatibility(config, subjectName); err == nil {
		var subjectConfig SubjectConfig
		if err := json.Unmarshal([]byte(compatibilityJSON), &subjectConfig); err == nil {
			defaultCompatibility = subjectConfig.CompatibilityLevel
		}
	}

	type versionResult struct {
		version            int
		id                 int
		schema             string
		subject            string
		compatibilityLevel string
	}

	resultChan := make(chan versionResult, len(versions))
	
	for _, version := range versions {
		go func(v int) {
			params := fmt.Sprintf("subjects/%v/versions/%d", subjectName, v)
			resp := BuildGetRequestSchemaRegistry(config, params)

			defer func() {
				if err := resp.Body.Close(); err != nil {
					util.CheckErr(err)
				}
			}()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				resultChan <- versionResult{
					version:            v,
					id:                 0,
					schema:             "",
					subject:            subjectName,
					compatibilityLevel: defaultCompatibility,
				}
				return
			}

			var versionInfo SubjectVersionInfo
			if err := json.Unmarshal(body, &versionInfo); err != nil {
				resultChan <- versionResult{
					version:            v,
					id:                 0,
					schema:             "",
					subject:            subjectName,
					compatibilityLevel: defaultCompatibility,
				}
				return
			}

			resultChan <- versionResult{
				version:            versionInfo.Version,
				id:                 versionInfo.ID,
				schema:             versionInfo.Schema,
				subject:            versionInfo.Subject,
				compatibilityLevel: defaultCompatibility,
			}
		}(version)
	}

	var result []SubjectVersionWithCompatibility
	for i := 0; i < len(versions); i++ {
		res := <-resultChan
		result = append(result, SubjectVersionWithCompatibility{
			Version:            res.version,
			ID:                 res.id,
			Schema:             res.schema,
			Subject:            res.subject,
			CompatibilityLevel: res.compatibilityLevel,
		})
	}

	return result, nil
}

func NewSubjectListWithCompatibility(config *configuration.Configuration) ([]SubjectWithCompatibility, error) {
	subjects, err := NewSubjectList(config)
	if err != nil {
		return nil, err
	}

	globalCompatibilityJSON, err := NewGlobalCompatibility(config)
	if err != nil {
		return nil, err
	}

	var globalConfig SubjectConfig
	var defaultCompatibility = "BACKWARD"
	if err := json.Unmarshal([]byte(globalCompatibilityJSON), &globalConfig); err == nil {
		defaultCompatibility = globalConfig.CompatibilityLevel
	}

	type subjectResult struct {
		subject            string
		compatibilityLevel string
		versionCount       int
	}

	resultChan := make(chan subjectResult, len(subjects))
	
	for _, subject := range subjects {
		go func(subjectName string) {
			compatibility := defaultCompatibility
			versionCount := 0

			if compatibilityJSON, err := NewSubjectCompatibility(config, subjectName); err == nil {
				var subjectConfig SubjectConfig
				if err := json.Unmarshal([]byte(compatibilityJSON), &subjectConfig); err == nil {
					compatibility = subjectConfig.CompatibilityLevel
				}
			}

			if versionsJSON, err := NewSubjectVersion(config, subjectName); err == nil {
				var versions []int
				if err := json.Unmarshal([]byte(versionsJSON), &versions); err == nil {
					versionCount = len(versions)
				}
			}

			resultChan <- subjectResult{
				subject:            subjectName,
				compatibilityLevel: compatibility,
				versionCount:       versionCount,
			}
		}(subject)
	}

	var result []SubjectWithCompatibility
	for i := 0; i < len(subjects); i++ {
		res := <-resultChan
		result = append(result, SubjectWithCompatibility{
			Name:               res.subject,
			CompatibilityLevel: res.compatibilityLevel,
			VersionCount:       res.versionCount,
		})
	}

	return result, nil
}
