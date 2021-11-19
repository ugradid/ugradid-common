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

package ssi

type KeyType string

// JsonWebKey2020 is a VerificationMethod type.
// https://w3c-ccg.github.io/lds-jws2020/
const JsonWebKey2020 = KeyType("JsonWebKey2020")

// ED25519VerificationKey2018 is the Ed25519VerificationKey2018 verification key type as specified here:
// https://w3c-ccg.github.io/lds-ed25519-2018/
const ED25519VerificationKey2018 = KeyType("Ed25519VerificationKey2018")

// ECDSASECP256K1VerificationKey2019 is the EcdsaSecp256k1VerificationKey2019 verification key type as specified here:
// https://w3c-ccg.github.io/lds-ecdsa-secp256k1-2019/
const ECDSASECP256K1VerificationKey2019 = KeyType("EcdsaSecp256k1VerificationKey2019")

// RSAVerificationKey2018 is the RsaVerificationKey2018 verification key type as specified here:
// https://w3c-ccg.github.io/lds-rsa2018/
const RSAVerificationKey2018 = KeyType("RsaVerificationKey2018")

type ProofType string

// JsonWebSignature2020 is a Proof type.
// https://w3c-ccg.github.io/lds-jws2020
const JsonWebSignature2020 = ProofType("JsonWebSignature2020")
