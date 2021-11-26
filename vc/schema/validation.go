/*
 * Copyright (c) 2021 ugradid community
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <https://www.gnu.org/licenses/>.
 */

package schema

import (
	"encoding/json"
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"regexp"
	"strings"
)

const (
	idRxStr = `^(did:(?:ugra):\S+)\;id=(\S+);version=(\d+\.\d+)$`
)

var IDRx = regexp.MustCompile(idRxStr)


type InvalidSchemaError struct {
	Errors []string
}

func (err InvalidSchemaError) Error() string {
	return fmt.Sprintf("Invalid schema: %s", strings.Join(err.Errors, ", "))
}

// Validate exists to hide gojsonschema logic within this file
// it is the entry-point to validation logic, requiring the caller pass in valid json strings for each argument
func Validate(schema, document string) error {
	if !IsJSON(schema) {
		return fmt.Errorf("schema is not valid json: %s", schema)
	} else if !IsJSON(document) {
		return fmt.Errorf("document is not valid json: %s", document)
	}
	return ValidateWithJSONLoader(gojsonschema.NewStringLoader(schema), gojsonschema.NewStringLoader(document))
}

// ValidateWithJSONLoader takes schema and document loaders; the document from the loader is validated against
// the schema from the loader. Nil if good, error if bad
func ValidateWithJSONLoader(schemaLoader, documentLoader gojsonschema.JSONLoader) error {
	// Add custom validator(s) and then ValidateWithJSONLoader
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		// Accumulate errs
		var errs []string
		for _, err := range result.Errors() {
			errs = append(errs, err.String())
		}
		return InvalidSchemaError{Errors: errs}
	}
	return nil
}

// ValidateJSONSchema takes in a string that is purported to be a JSON schema (schema definition)
// An error is returned if it is not a valid JSON schema, and nil is returned on success
func ValidateJSONSchema(maybeSchema JsonSchema) error {
	jsonSchema, err := json.Marshal(maybeSchema)
	if err != nil {
		return err
	}
	return ValidateJSONSchemaBytes(jsonSchema)
}

// ValidateJSONSchemaBytes takes in a map that is purported to be a JSON schema (schema definition)
// An error is returned if it is not a valid JSON schema, and nil is returned on success
func ValidateJSONSchemaBytes(maybeSchema []byte) error {
	schemaLoader := gojsonschema.NewSchemaLoader()
	schemaLoader.Validate = true
	return schemaLoader.AddSchemas(gojsonschema.NewBytesLoader(maybeSchema))
}

// IsJSON True if string is valid JSON, false otherwise
func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}


