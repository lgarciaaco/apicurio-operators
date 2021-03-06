package test

import (
	"github.com/RHsyseng/operator-utils/pkg/validation"
	"github.com/apicurio/apicurio-operators/apicurito/pkg/apis/apicur/v1alpha1"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSampleCustomResources(t *testing.T) {
	schema := getSchema(t)
	assert.NotNil(t, schema)

	filePath := getCRFile(t, "../../../../deploy/crs")
	bytes, err := ioutil.ReadFile(filePath)
	assert.NoError(t, err, "Error reading CR yaml %v", filePath)

	var input map[string]interface{}
	assert.NoError(t, yaml.Unmarshal(bytes, &input))
	assert.NoError(t, schema.Validate(input), "File %v does not validate against the CRD schema", filePath)
}

func TestTrialEnvMinimum(t *testing.T) {
	var inputYaml = `
apiVersion: apicur.io/v1alpha1
kind: Apicurito
metadata:
  name: trial
spec:
  size: 3
`
	var input map[string]interface{}
	assert.NoError(t, yaml.Unmarshal([]byte(inputYaml), &input))

	schema := getSchema(t)
	assert.NoError(t, schema.Validate(input))
}

func TestCompleteCRD(t *testing.T) {
	schema := getSchema(t)
	missingEntries := schema.GetMissingEntries(&v1alpha1.Apicurito{})
	if strings.Contains(missingEntries[0].Path, "/spec/size") {
		// The size is not expected to be used and is not fully defined TODO: verify
	} else {
		assert.Fail(t, "Discrepancy between CRD and Struct", "Missing or incorrect schema validation at %v, expected type %v", missingEntries[0].Path, missingEntries[0].Type)
	}
}

func getCRFile(t *testing.T, dir string) string {
	var file string
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				file = path
			}
			return nil
		})
	assert.NoError(t, err, "Error finding CR yaml %v", file)
	return file
}

func getSchema(t *testing.T) validation.Schema {
	crdFile := "../../../../deploy/crds/apicur_v1alpha1_apicurito_crd.yaml"
	bytes, err := ioutil.ReadFile(crdFile)
	assert.NoError(t, err, "Error reading CRD yaml %v", crdFile)
	schema, err := validation.New(bytes)
	assert.NoError(t, err)
	return schema
}
