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

package vc

import (
	"encoding/json"
	"fmt"
	ssi "github.com/ugradid/ugradid-common"
	"github.com/ugradid/ugradid-common/marshal"
	"net/url"
	"time"
)

// VerifiableCredentialType is the default credential type required for every credential
const VerifiableCredentialType = "VerifiableCredential"

// SchemaCredentialType use when a subject data schema needs to be specified
const SchemaCredentialType = "SchemaCredential"

// VerifiableCredentialTypeV1URI returns VerifiableCredential as URI
func VerifiableCredentialTypeV1URI() ssi.URI {
	if pURI, err := ssi.ParseURI(VerifiableCredentialType); err != nil {
		panic(err)
	} else {
		return *pURI
	}
}

// SchemaCredentialTypeV1URI returns SchemaCredentialType as URI
func SchemaCredentialTypeV1URI() ssi.URI {
	if pURI, err := ssi.ParseURI(SchemaCredentialType); err != nil {
		panic(err)
	} else {
		return *pURI
	}
}

// VCContextV1 is the context required for every credential
const VCContextV1 = "https://www.w3.org/2018/credentials/v1"

// VCContextV1URI returns 'https://www.w3.org/2018/credentials/v1' as URI
func VCContextV1URI() ssi.URI {
	if pURI, err := ssi.ParseURI(VCContextV1); err != nil {
		panic(err)
	} else {
		return *pURI
	}
}

// VerifiableCredential represents a credential as defined by the Verifiable Credentials Data Model 1.0 specification (https://www.w3.org/TR/vc-data-model/).
type VerifiableCredential struct {
	// Context defines the json-ld context to dereference the URIs
	Context []ssi.URI `json:"@context"`
	// ID is an unique identifier for the credential. It is optional
	ID *ssi.URI `json:"id,omitempty"`
	// Type holds multiplte types for a credential. A credential must always have the 'VerifiableCredential' type.
	Type []ssi.URI `json:"type"`
	// Issuer refers to the party that issued the credential
	Issuer ssi.URI `json:"issuer"`
	// IssuanceDate is a rfc3339 formatted datetime.
	IssuanceDate time.Time `json:"issuanceDate"`
	// ExpirationDate is a rfc3339 formatted datetime. It is optional
	ExpirationDate *time.Time `json:"expirationDate,omitempty"`
	// CredentialStatus holds information on how the credential can be revoked. It is optional
	CredentialStatus *CredentialStatus `json:"credentialStatus,omitempty"`
	// CredentialSchema holds information schema credential subject
	CredentialSchema *CredentialSchema `json:"credentialSchema,omitempty"`
	// CredentialSubject holds the actual data for the credential. It must be extracted using the UnmarshalCredentialSubject method and a custom type.
	CredentialSubject map[string]interface{} `json:"credentialSubject"`
	// Proof contains the cryptographic proof(s). It must be extracted using the Proofs method or UnmarshalProofValue method for non-generic proof fields.
	Proof []interface{} `json:"proof"`
}

// CredentialStatus defines the method on how to determine a credential is revoked.
type CredentialStatus struct {
	ID   url.URL `json:"id"`
	Type string  `json:"type"`
}

// CredentialSchema defines for schema subject.
type CredentialSchema struct {
	ID   ssi.URI        `json:"id"`
	Type ssi.SchemaType `json:"type"`
}

// Proofs returns the basic proofs for this credential. For specific proof contents, UnmarshalProofValue must be used.
func (vc VerifiableCredential) Proofs() ([]Proof, error) {
	var (
		target []Proof
		err    error
		asJSON []byte
	)
	asJSON, err = json.Marshal(vc.Proof)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(asJSON, &target)
	return target, err
}

func (vc VerifiableCredential) MarshalJSON() ([]byte, error) {
	type alias VerifiableCredential
	tmp := alias(vc)
	if data, err := json.Marshal(tmp); err != nil {
		return nil, err
	} else {
		return marshal.NormalizeDocument(data, pluralContext,
			marshal.Unplural(typeKey), marshal.Unplural(credentialSubjectKey), marshal.Unplural(proofKey))
	}
}

func (vc *VerifiableCredential) UnmarshalJSON(b []byte) error {
	type Alias VerifiableCredential
	normalizedVC, err := marshal.NormalizeDocument(b,
		pluralContext, marshal.Plural(typeKey), marshal.Unplural(credentialSubjectKey), marshal.Plural(proofKey))
	if err != nil {
		return err
	}
	tmp := Alias{}
	err = json.Unmarshal(normalizedVC, &tmp)
	if err != nil {
		return err
	}
	*vc = (VerifiableCredential)(tmp)
	return nil
}

// UnmarshalProofValue unmarshalls the proof to the given proof type. Always pass a slice as target since there could be multiple proofs.
// Each proof will result in a value, where null values may exist when the proof doesn't have the json member.
func (vc VerifiableCredential) UnmarshalProofValue(target interface{}) error {
	if asJSON, err := json.Marshal(vc.Proof); err != nil {
		return err
	} else {
		return json.Unmarshal(asJSON, target)
	}
}

// UnmarshalCredentialSubject unmarshalls the credentialSubject to the given credentialSubject type. Always pass a slice as target.
func (vc VerifiableCredential) UnmarshalCredentialSubject(target interface{}) error {
	if asJSON, err := json.Marshal(vc.CredentialSubject); err != nil {
		return err
	} else {
		return json.Unmarshal(asJSON, target)
	}
}

// IsType returns true when a credential contains the requested type
func (vc VerifiableCredential) IsType(vcType ssi.URI) bool {
	for _, t := range vc.Type {
		if t.String() == vcType.String() {
			return true
		}
	}

	return false
}

// ContainsContext returns true when a credential contains the requested context
func (vc VerifiableCredential) ContainsContext(context ssi.URI) bool {
	for _, c := range vc.Context {
		if c.String() == context.String() {
			return true
		}
	}

	return false
}

func GenerateCredentialID(issuer ssi.URI, id string) string {
	return fmt.Sprintf("%s#%s", issuer.String(), id)
}
