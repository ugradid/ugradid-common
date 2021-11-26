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
	"fmt"
	ssi "github.com/ugradid/ugradid-common"
	"github.com/ugradid/ugradid-common/vc"
	"regexp"
	"time"
)

// Draft7Schema This schema supports simple draft 7 JSON Schemas
const (
	Draft7Schema = "http://json-schema.org/draft-07/schema"
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

// Schema for ledger objects and how they will be stored
// Type, Model ModelVersion, Schema, Author and ID should always be present
type Schema struct {
	Type     ssi.URI                       `json:"type"`
	Version  string                        `json:"version"`
	ID       *ssi.URI                      `json:"id"`
	Name     string                        `json:"name,omitempty"`
	Author   ssi.URI                       `json:"author,omitempty"`
	Authored time.Time                     `json:"authored,omitempty"`
	Schema   JsonSchema                    `json:"schema"`
	Proof    *vc.JSONWebSignature2020Proof `json:"proof,omitempty"`
}

func GenerateSchemaID(author ssi.URI, id string, version string) string {
	return fmt.Sprintf("%s;id=%s;version=%s", author.String(), id, version)
}

func Contains(field string, required []string) bool {
	for _, f := range required {
		if f == field {
			return true
		}
	}
	return false
}

// IsType returns true when a credential contains the requested type
func (s Schema) IsType(sType ssi.URI) bool {
	return s.Type == sType
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

	result := r.Match([]byte(s.ID.String()))
	if !result {
		return fmt.Errorf("schema 'id': %s is not valid against pattern: %s", s.ID, regx)
	}

	return nil
}

// ValidateVersion assumes the version property is the only version in the identifier separated by periods
func (s Schema) ValidateVersion() error {
	regx := "\\d+\\.\\d+$"
	r, err := regexp.Compile(regx)
	if err != nil {
		return fmt.Errorf("failed to compile version regular expression: %s", regx)
	}

	result := r.Find([]byte(s.Version))
	if result == nil {
		return fmt.Errorf("error returning version property with regular expression: %s", regx)
	}

	return nil
}
