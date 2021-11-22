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
	"github.com/go-playground/validator/v10"
)

// This schema builder supports simple draft 7 JSON Schemas

const (
	Draft7Schema = "http://json-schema.org/draft-07/schema#"
	TypeSchema   = "object"
)

type Builder struct {
	Name                 string `validate:"required"`
	Description          string `validate:"required"`
	AdditionalProperties bool
	Attributes           []Attribute `validate:"required"`
}

type AttributeType string

const (
	String  AttributeType = "string"
	Number  AttributeType = "number"
	Object  AttributeType = "object"
	Array   AttributeType = "array"
	Boolean AttributeType = "boolean"
)

type FormatType string

const (
	Date     FormatType = "date"
	Time     FormatType = "time"
	DateTime FormatType = "date-time"
	Email    FormatType = "email"
)

type Attribute struct {
	Name     string        `json:"-" validate:"required"`
	Type     AttributeType `json:"type,omitempty" validate:"required"`
	Required bool          `json:"-" validate:"required"`
	*StringType
	*NumberType
	*ObjectType
	*ArrayType
}

type StringType struct {
	Format FormatType `json:"format,omitempty"`
}

type NumberType struct {
	Minimum          float32 `json:"minimum,omitempty"`
	Maximum          float32 `json:"maximum,omitempty"`
	ExclusiveMinimum float32 `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum float32 `json:"exclusiveMaximum,omitempty"`
}

type ObjectType struct {
	Properties           map[string]Attribute `json:"properties" validate:"required"`
	RequiredProperties   []string             `json:"requiredProperties" validate:"required"`
	AdditionalProperties bool                 `json:"additionalProperties" validate:"required"`
}

type ArrayType struct {
	AttributeType `json:"type" validate:"required"`
	*StringType
	*NumberType
}

func (a *ArrayType) Validate(name string) error {
	if a.AttributeType == "" {
		return fmt.Errorf("field %s type is not empty", name)
	}
	return nil
}

// Private structs for construction

type jsonSchema struct {
	Draft                string              `json:"$schema"`
	Description          string              `json:"description"`
	Type                 AttributeType       `json:"type"`
	Properties           map[string]jsonAttr `json:"properties"`
	Required             []string            `json:"required,omitempty"`
	AdditionalProperties bool                `json:"additionalProperties"`
}

type jsonAttr struct {
	Type AttributeType `json:"type"`
	*StringType
	*NumberType
	*ArrayType `json:"items,omitempty"`
	*ObjectType
}

// Build takes in a schema descriptor and an ordered set of attributes and builds
// a valid JSON schema. The output is run through a schema validator before being returned
func (b Builder) Build() (JsonSchemaMap, error) {

	// validate all required fields are set
	if err := validator.New().Struct(b); err != nil {
		return nil, err
	}

	properties, required, err := buildProperties(b.Attributes)
	if err != nil {
		return nil, err
	}

	schema := jsonSchema{
		Draft:                Draft7Schema,
		Description:          b.Description,
		Type:                 TypeSchema,
		Properties:           properties,
		Required:             required,
		AdditionalProperties: b.AdditionalProperties,
	}

	// Turn it into JSON
	bytes, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	// Put it into the expected return type
	var jsonMap JsonSchemaMap
	if err := json.Unmarshal(bytes, &jsonMap); err != nil {
		return nil, err
	}

	// Validate the schema
	if err := ValidateJSONSchema(jsonMap); err != nil {
		return nil, err
	}
	return jsonMap, nil
}

func buildProperties(attrs []Attribute) (properties map[string]jsonAttr, required []string, err error) {
	properties = make(map[string]jsonAttr)
	required = make([]string, 0)

	for _, attr := range attrs {
		if _, ok := properties[attr.Name]; ok {
			return nil, nil, fmt.Errorf("duplicate property: %s", attr.Name)
		}

		var toAdd jsonAttr

		switch attr.Type {
		case String:
			toAdd = jsonAttr{
				Type:       attr.Type,
				StringType: attr.StringType,
			}
		case Number:
			toAdd = jsonAttr{
				Type:       attr.Type,
				NumberType: attr.NumberType,
			}
		case Boolean:
			toAdd = jsonAttr{
				Type: attr.Type,
			}
		case Array:
			toAdd = jsonAttr{
				Type: attr.Type,
				ArrayType: &ArrayType{
					AttributeType: attr.ArrayType.AttributeType,
					StringType:    attr.ArrayType.StringType,
					NumberType:    attr.ArrayType.NumberType,
				},
			}
		case Object:
			toAdd = jsonAttr{
				Type: attr.Type,
				ObjectType: &ObjectType{
					Properties:           attr.ObjectType.Properties,
					RequiredProperties:   attr.ObjectType.RequiredProperties,
					AdditionalProperties: attr.ObjectType.AdditionalProperties,
				},
			}
		default:
			return nil, nil, fmt.Errorf("unknown attr type: %s", attr.Type)
		}

		properties[attr.Name] = toAdd
		if attr.Required {
			required = append(required, attr.Name)
		}
	}
	return
}
