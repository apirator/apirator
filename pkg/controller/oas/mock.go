package oas

import (
	"encoding/json"
	"github.com/getkin/kin-openapi/openapi3"
	"strings"
)

// find mock api path
// https://swagger.io/docs/specification/openapi-extensions/
func Path(doc *openapi3.Swagger) string {
	i := doc.Info.Extensions["x-apirator-mock-path"]
	json := i.(json.RawMessage)
	path := strings.Trim(string(json), "\"")
	return path
}
