// Copyright (c) OpenFaaS Author(s) 2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package commands

import (
	"testing"

	"github.com/openfaas/faas-cli/schema"

	"github.com/openfaas/faas-cli/stack"
)

var generateTestcases = []struct {
	Name       string
	Input      string
	Output     string
	Format     schema.BuildFormat
	APIVersion string
	Namespace  string
	Branch     string
	Version    string
}{
	{
		Name: "Default Namespace and API Version",
		Input: `
provider:
  name: faas
  gateway: http://127.0.0.1:8080
  network: "func_functions"      
functions:
 url-ping:
   lang: python
   handler: ./sample/url-ping
   image: alexellis/faas-url-ping:0.2`,
		Output: `---
apiVersion: openfaas.com/v1alpha2
kind: Function
metadata:
  name: url-ping
  namespace: openfaas-fn
spec:
  name: url-ping
  image: alexellis/faas-url-ping:0.2
`,
		Format:     schema.DefaultFormat,
		APIVersion: "openfaas.com/v1alpha2",
		Namespace:  "openfaas-fn",
		Branch:     "",
		Version:    "",
	},

	{
		Name: "Blank namespace",
		Input: `
provider:
  name: faas
  gateway: http://127.0.0.1:8080
  network: "func_functions"
functions:
 url-ping:
  lang: python
  handler: ./sample/url-ping
  image: alexellis/faas-url-ping:0.2`,
		Output: `---
apiVersion: openfaas.com/v1alpha2
kind: Function
metadata:
  name: url-ping
spec:
  name: url-ping
  image: alexellis/faas-url-ping:0.2
`,
		Format:     schema.DefaultFormat,
		APIVersion: "openfaas.com/v1alpha2",
		Namespace:  "",
		Branch:     "",
		Version:    "",
	},
	{
		Name: "BranchAndSHA Image format",
		Input: `
provider:
  name: faas
  gateway: http://127.0.0.1:8080
  network: "func_functions"
functions:
 url-ping:
  lang: python
  handler: ./sample/url-ping
  image: alexellis/faas-url-ping:0.2`,
		Output: `---
apiVersion: openfaas.com/v1alpha2
kind: Function
metadata:
  name: url-ping
  namespace: openfaas-function
spec:
  name: url-ping
  image: alexellis/faas-url-ping:0.2-master-6bgf36qd
`,
		Format:     schema.BranchAndSHAFormat,
		APIVersion: "openfaas.com/v1alpha2",
		Namespace:  "openfaas-function",
		Branch:     "master",
		Version:    "6bgf36qd",
	},
	{
		Name: "Multiple functions",
		Input: `
provider:
  name: faas
  gateway: http://127.0.0.1:8080  
  network: "func_functions"
functions:
 url-ping:
  lang: python
  handler: ./sample/url-ping
  image: alexellis/faas-url-ping:0.2
 astronaut-finder:
  lang: python3
  handler: ./astronaut-finder
  image: astronaut-finder
  environment:
   write_debug: true`,
		Output: `---
apiVersion: openfaas.com/v2alpha2
kind: Function
metadata:
  name: url-ping
  namespace: openfaas-fn
spec:
  name: url-ping
  image: alexellis/faas-url-ping:0.2
---
apiVersion: openfaas.com/v2alpha2
kind: Function
metadata:
  name: astronaut-finder
  namespace: openfaas-fn
spec:
  name: astronaut-finder
  image: astronaut-finder:latest
  environment:
    write_debug: "true"
`,
		Format:     schema.DefaultFormat,
		APIVersion: "openfaas.com/v2alpha2",
		Namespace:  "openfaas-fn",
		Branch:     "",
		Version:    "",
	},
}

func Test_generateCRDYAML(t *testing.T) {

	for _, testcase := range generateTestcases {
		parsedServices, err := stack.ParseYAMLData([]byte(testcase.Input), "", "")

		if err != nil {
			t.Fatalf("%s failed: error while parsing the input data.", testcase.Name)
		}

		if parsedServices == nil {
			t.Fatalf("%s failed: empty input file", testcase.Name)
		}
		services := *parsedServices

		generatedYAML, err := generateCRDYAML(services, testcase.Format, testcase.APIVersion, testcase.Namespace, testcase.Branch, testcase.Version)
		if err != nil {
			t.Fatalf("%s failed: error while generating CRD YAML.", testcase.Name)
		}

		if generatedYAML != testcase.Output {
			t.Fatalf("%s failed: ouput is not as expected: %s", testcase.Name, generatedYAML)
		}
	}

}
