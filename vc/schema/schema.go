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
	"github.com/ugradid/ugradid-common/did"
	"regexp"
)

// Version around modified semantic version where only major and minor numbers are valid
type Version struct {
	Major int
	Minor int
}

// IDFormatErr is a formatting error in the Schema ID,
// which should be in the form <author_did>;id=<uuid>;version=<major.minor>.
type IDFormatErr struct {
	schemaID string
}

func (e IDFormatErr) Error() string {
	return fmt.Sprintf("'%s' schema id is in an unrecognized format", e.schemaID)
}

// UnRecognisedVersionError is a formatting error in the Schema version, which should be in the form "major.minor".
type UnRecognisedVersionError struct {
	submittedVersion string
}

func (e UnRecognisedVersionError) Error() string {
	return fmt.Sprintf("'%s' is an unrecognized version format", e.submittedVersion)
}

// Metadata for ledger objects and how they will be stored
// Type, Model ModelVersion, and ID should always be present
// Depending on the model object, the remainder of the fields may be optional.
type Metadata struct {
	Type     string        `json:"type"`
	Version  string        `json:"version"`
	ID       string        `json:"id"`
	Name     string        `json:"name,omitempty"`
	Author   did.DID       `json:"author,omitempty"`
	Authored string        `json:"authored,omitempty"`
	Proof    []interface{} `json:"proof,omitempty"`
}

type Schema struct {
	*Metadata
	*JsonSchema
}

type UpdateInput struct {
	PreviousSchema *Schema
	UpdatedSchema  *Schema
}

type UpdateResult struct {
	Valid          bool
	MajorChange    bool
	MinorChange    bool
	DerivedVersion string
	Message        string
}

// JsonSchema for a credential that has not been signed
type JsonSchema struct {
	Schema JsonSchemaMap `json:"schema"`
}

// JsonSchemaMap representation of json schema document
type JsonSchemaMap map[string]interface{}

type Properties map[string]interface{}

// Properties Assumes the json schema has a properties field
func (j JsonSchemaMap) Properties() Properties {
	if properties, ok := j["properties"]; ok {
		return properties.(map[string]interface{})
	}
	return map[string]interface{}{}
}

// Description Assumes the json schema has a description field
func (j JsonSchemaMap) Description() string {
	if description, ok := j["description"]; ok {
		return description.(string)
	}
	return ""
}

func (j JsonSchemaMap) AllowsAdditionalProperties() bool {
	if v, exists := j["additionalProperties"]; exists {
		if additionalProps, ok := v.(bool); ok {
			return additionalProps
		}
	}
	return false
}

func (j JsonSchemaMap) RequiredFields() []string {
	if v, exists := j["required"]; exists {
		if requiredFields, ok := v.([]interface{}); ok {
			required := make([]string, 0, len(requiredFields))
			for _, f := range requiredFields {
				required = append(required, f.(string))
			}
			return required
		}
	}
	return []string{}
}

func Type(field interface{}) string {
	if asMap, isMap := field.(map[string]interface{}); isMap {
		if v, exists := asMap["type"]; exists {
			if typeString, ok := v.(string); ok {
				return typeString
			}
		}
	}
	return ""
}

func Format(field interface{}) string {
	if asMap, isMap := field.(map[string]interface{}); isMap {
		if v, exists := asMap["format"]; exists {
			if formatString, ok := v.(string); ok {
				return formatString
			}
		}
	}
	return ""
}

func Contains(field string, required []string) bool {
	for _, f := range required {
		if f == field {
			return true
		}
	}
	return false
}

func (j JsonSchemaMap) ToJSON() (string, error) {
	bytes, err := json.Marshal(j)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func GenerateSchemaID(author did.DID, id string, version string) string {
	return fmt.Sprintf("%s;id=%s;version=%s", author, id, version)
}

// Validates a schema for a correctly composed Credential Schema
// Currently only validates the ID property. Add additional validation if required

// ValidateID ID validation is based on our public schema specification:
// This identifier is a method-specific DID parameter name based upon the author of the
// schema. For example, if the author had a did like did:ugra:abcdefghi a possible schema
// ID the author created would have an identifier such as:
// did:ugra:abcdefghi;schema=17de181feb67447da4e78259d92d0240;version=1.0
func (s Schema) ValidateID() error {
	regx := "^did:ugra:\\S+\\;id=\\S+;version=\\d+\\.\\d+$"
	r, err := regexp.Compile(regx)
	if err != nil {
		return fmt.Errorf("failed to compile regular expression: %s", regx)
	}

	result := r.Match([]byte(s.ID))
	if !result {
		return fmt.Errorf("ledger schema 'id': %s is not valid against pattern: %s", s.ID, regx)
	}

	return nil
}

// Version assumes the version property is the only version in the identifier separated by periods
func (s Schema) Version() (string, error) {
	regx := "\\d+\\.\\d+$"
	r, err := regexp.Compile(regx)
	if err != nil {
		return "", fmt.Errorf("failed to compile regular expression: %s", regx)
	}

	result := r.Find([]byte(s.ID))
	if result == nil {
		return "", fmt.Errorf("error returning version property with regular expression: %s", regx)
	}

	return string(result), nil
}
