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
	ssi "github.com/ugradid/ugradid-common"
	"time"
)

// Proof represents a credential/presentation proof as defined by the Linked Data Proofs 1.0 specification (https://w3c-ccg.github.io/ld-proofs/).
// The proof value must be implemented in a custom type since the specification doesn't define the json object for this.
// For example: a jws for detached JSON Web Signatures uses the 'jws' json field
type Proof struct {
	// Type defines the specific proof type used.
	// For example, an Ed25519Signature2018 type indicates that the proof includes a digital signature produced by an ed25519 cryptographic key.
	Type ssi.ProofType `json:"type"`
	// ProofPurpose defines the intent for the proof, the reason why an entity created it.
	// Acts as a safeguard to prevent the proof from being misused for a purpose other than the one it was intended for.
	// For example, a proof can be used for purposes of authentication, for asserting control of a Verifiable Credential (assertionMethod), and several others.
	ProofPurpose string `json:"proofPurpose"`
	// VerificationMethod points to the ID that can be used to verify the proof, eg: a public key.
	VerificationMethod ssi.URI `json:"verificationMethod"`
	// Created notes when the proof was created using a iso8601 string
	Created time.Time `json:"created"`
	// Domain specifies the restricted domain of the proof
	Domain *string `json:"domain,omitempty"`
}

// JSONWebSignature2020Proof is a VC proof with a signature according to JsonWebSignature2020
type JSONWebSignature2020Proof struct {
	Proof
	Jws string `json:"jws"`
}
