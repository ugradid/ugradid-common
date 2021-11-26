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
	"github.com/pkg/errors"
)

// JsonSchema is a JSON Schema draft-07 document
type JsonSchema struct {
	Comment              *string                      `json:"$comment,omitempty"`
	ID                   *string                      `json:"$id,omitempty"`
	Reference            *string                      `json:"$ref,omitempty"`
	SchemaRef            string                       `json:"$schema,omitempty"`
	AdditionalItems      *JsonSchema                  `json:"additionalItems,omitempty"`
	AdditionalProperties *AdditionalPropertyValue     `json:"additionalProperties,omitempty"`
	AllOf                []*JsonSchema                `json:"allOf,omitempty"`
	AnyOf                []*JsonSchema                `json:"anyOf,omitempty"`
	Const                *interface{}                 `json:"const,omitempty"`
	Contains             *JsonSchema                  `json:"contains,omitempty"`
	Default              *interface{}                 `json:"default,omitempty"`
	Definitions          *map[string]*JsonSchema      `json:"definitions,omitempty"`
	Dependencies         *map[string]*DependencyValue `json:"dependencies,omitempty"`
	Description          *string                      `json:"description,omitempty"`
	Else                 *JsonSchema                  `json:"else,omitempty"`
	Enum                 EnumList                     `json:"enum,omitempty"`
	Examples             []interface{}                `json:"examples,omitempty"`
	ExclusiveMaximum     *float64                     `json:"exclusiveMaximum,omitempty"`
	ExclusiveMinimum     *float64                     `json:"exclusiveMinimum,omitempty"`
	Format               *Format                      `json:"format,omitempty"`
	If                   *JsonSchema                  `json:"if,omitempty"`
	Items                *ItemSchemaList              `json:"items,omitempty"`
	MaxItems             *int64                       `json:"maxItems,omitempty"`
	MaxLength            *int64                       `json:"maxLength,omitempty"`
	MaxProperties        *int64                       `json:"maxProperties,omitempty"`
	Maximum              *float64                     `json:"maximum,omitempty"`
	MinItems             *int64                       `json:"minItems,omitempty"`
	MinLength            *int64                       `json:"minLength,omitempty"`
	MinProperties        *int64                       `json:"minProperties,omitempty"`
	Minimum              *float64                     `json:"minimum,omitempty"`
	MultipleOf           *float64                     `json:"multipleOf,omitempty"`
	Not                  *JsonSchema                  `json:"not,omitempty"`
	OneOf                []*JsonSchema                `json:"oneOf,omitempty"`
	Pattern              *string                      `json:"pattern,omitempty"`
	PatternProperties    *map[string]*JsonSchema      `json:"patternProperties,omitempty"`
	Properties           *map[string]*JsonSchema      `json:"properties,omitempty"`
	PropertyNames        *JsonSchema                  `json:"propertyNames,omitempty"`
	Required             []string                     `json:"required,omitempty"`
	Then                 *JsonSchema                  `json:"then,omitempty"`
	Title                *string                      `json:"title,omitempty"`
	Type                 *PrimitiveTypeList            `json:"type,omitempty"`
	UniqueItems          *bool                        `json:"uniqueItems,omitempty"`
}

type DependencyValue struct {
	Schema             *JsonSchema
	RequiredProperties []string
}

type AdditionalPropertyValue struct {
	Schema *JsonSchema
	Status bool
}

// MarshalJSON implements json.Marshaler.
func (s *AdditionalPropertyValue) MarshalJSON() ([]byte, error) {
	if s.Schema != nil {
		return json.Marshal(s.Schema)
	}
	return json.Marshal(s.Status)
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *AdditionalPropertyValue) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case "true":
		s.Status = true
	case "false":
		s.Status = false
	default:
		return json.Unmarshal(data, &s.Schema)
	}
	return nil
}

type ItemSchemaList struct {
	Schema  *JsonSchema
	Schemas []*JsonSchema
}

// MarshalJSON implements json.Marshaler.
func (s *ItemSchemaList) MarshalJSON() ([]byte, error) {
	if s.Schema != nil {
		return json.Marshal(s.Schema)
	}
	return json.Marshal(s.Schemas)
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *ItemSchemaList) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '[' {
		return json.Unmarshal(data, &s.Schemas)
	}
	return json.Unmarshal(data, &s.Schema)
}

// MarshalJSON implements json.Marshaler.
func (s *JsonSchema) MarshalJSON() ([]byte, error) {
	type schema2 JsonSchema
	return json.Marshal((*schema2)(s))
}

// UnmarshalJSON implements json.Unmarshaler.
func (s *JsonSchema) UnmarshalJSON(data []byte) error {
	type schema2 JsonSchema
	if err := json.Unmarshal(data, (*schema2)(s)); err != nil {
		return errors.New("failed to unmarshal JSON Schema")
	}
	return nil
}

// MarshalJSON implements json.Marshaler.
func (v *DependencyValue) MarshalJSON() ([]byte, error) {
	if v.Schema != nil {
		return json.Marshal(v.Schema)
	}
	return json.Marshal(v.RequiredProperties)
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *DependencyValue) UnmarshalJSON(data []byte) error {
	*v = DependencyValue{}
	if len(data) > 0 && data[0] == '[' {
		return json.Unmarshal(data, &v.RequiredProperties)
	}
	return json.Unmarshal(data, &v.Schema)
}
