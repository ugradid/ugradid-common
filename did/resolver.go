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

package did

import (
	"time"
)

type constError string

func (err constError) Error() string {
	return string(err)
}

const (
	// InvalidDIDErr indicates: "The DID supplied to the DID resolution function does not conform to valid syntax. (See § 3.1 DID Syntax.)"
	InvalidDIDErr = constError("supplied DID is invalid")
	// NotFoundErr indicates: "The DID resolver was unable to find the DID document resulting from this resolution request."
	NotFoundErr = constError("supplied DID wasn't found")
	// DeactivatedErr indicates: The DID supplied to the DID resolution function has been deactivated. (See § 7.2.4 Deactivate .)
	DeactivatedErr = constError("supplied DID is deactivated")
)

// Resolver defines the interface for DID resolution as specified by the DID Core specification (https://www.w3.org/TR/did-core/#did-resolution).
type Resolver interface {
	// Resolve tries to resolve the given input DID to its DID Document and Metadata. In addition to errors specific
	// to this resolver it can return InvalidDIDErr, NotFoundErr and DeactivatedErr as specified by the DID Core specification.
	// If no error occurs the DID Document and Medata are returned.
	Resolve(inputDID string) (*Document, *DocumentMetadata, error)
}

// DocumentMedata represents DID Document Metadata as specified by the DID Core specification (https://www.w3.org/TR/did-core/#did-document-metadata-properties).
type DocumentMetadata struct {
	Created    *time.Time
	Updated    *time.Time
	Properties map[string]interface{}
}

