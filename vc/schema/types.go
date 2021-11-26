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
)

type Enum interface{}
type EnumList []Enum
type Format string
type PrimitiveType string
type PrimitiveTypeList []PrimitiveType

const (
	UnspecifiedType PrimitiveType = "unspecified"
	NullType                      = "null"
	BooleanType                   = "boolean"
	ObjectType                    = "object"
	ArrayType                     = "array"
	NumberType                    = "number"
	StringType                    = "string"
	IntegerType                   = "integer"
)

func (l PrimitiveTypeList) Len() int {
	return len(l)
}

func (l PrimitiveTypeList) Less(i, j int) bool {
	return l[i] < l[j]
}

func (l PrimitiveTypeList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l *PrimitiveTypeList) MarshalJSON() ([]byte, error) {
	if len(*l) > 1 {
		return json.Marshal([]PrimitiveType(*l))
	}
	return json.Marshal((*l)[0])
}

func (l *PrimitiveTypeList) UnmarshalJSON(buf []byte) error {
	var sl []string
	if len(buf) > 0 && buf[0] == '[' {
		if err := json.Unmarshal(buf, &sl); err != nil {
			return fmt.Errorf(`failed to parse primitive types list: %s`, err)
		}
	} else {
		var s string
		if err := json.Unmarshal(buf, &s); err != nil {
			return fmt.Errorf(`failed to parse primitive types list: %s`, err)
		}
		sl = []string{s}
	}

	ptl := make(PrimitiveTypeList, 0, len(sl))
	for _, s := range sl {
		var pt PrimitiveType
		switch s {
		case "null":
			pt = NullType
		case "boolean":
			pt = BooleanType
		case "object":
			pt = ObjectType
		case "array":
			pt = ArrayType
		case "number":
			pt = NumberType
		case "string":
			pt = StringType
		case "integer":
			pt = IntegerType
		default:
			return fmt.Errorf(`invalid primitive type: %s`, s)
		}
		ptl = append(ptl, pt)
	}

	*l = ptl
	return nil
}
