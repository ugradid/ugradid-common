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
	"github.com/ugradid/ugradid-common/did"
	"regexp"
	"strconv"
	"strings"
)

const (
	VersionRxStr           = `^[0-9]+\.[0-9]+$`
	idRxStr = `^(did:(?:ugra):\S+)\;id=(\S+);version=(\d+\.\d+)$`
	VersionPathResource    = "version"
	ResourceIDPathResource = "id"
	FragSep                = ";"
	FragAssignment         = "="
)

var IDRx = regexp.MustCompile(idRxStr)

// VersionFromStr parses a version string into a Version object. Returns an error if the version
// string does not match the VersionRxStr regular expression, i.e. "<major>.<minor>", where major
// and minor are integers.
func VersionFromStr(versionStr string) (Version, error) {
	schemaVersionRx := regexp.MustCompile(VersionRxStr)
	if !schemaVersionRx.MatchString(versionStr) {
		return Version{}, UnRecognisedVersionError{versionStr}
	}

	vnums := strings.Split(versionStr, ".")
	major, err := strconv.Atoi(vnums[0])
	if err != nil {
		return Version{}, err
	}
	minor, err := strconv.Atoi(vnums[1])
	if err != nil {
		return Version{}, err
	}
	return Version{
		Major: major,
		Minor: minor,
	}, nil
}

// ExtractSchemaVersionFromID parses the schema URI (did:ugra:<authorDID>;id=<uuid>;version=<version>)
// and returns the schema version.
func ExtractSchemaVersionFromID(schemaID string) (Version, error) {
	if !IDRx.MatchString(schemaID) {
		idFormatErr := IDFormatErr{schemaID}
		return Version{}, idFormatErr
	}
	vpr := VersionPathResource + FragAssignment
	vstr := schemaID[strings.Index(schemaID, vpr)+len(vpr):]
	return VersionFromStr(vstr)
}

// ExtractSchemaResourceID parses the schema URI (did:ugra:<authorDID>;id=<uuid>;version=<version>)
// and returns the resource ID.
func ExtractSchemaResourceID(schemaID string) (string, error) {
	if !IDRx.MatchString(schemaID) {
		idFormatErr := IDFormatErr{schemaID}
		return "", idFormatErr
	}

	idIdentifier := ResourceIDPathResource + FragAssignment
	rid := schemaID[strings.Index(schemaID, idIdentifier)+
		len(idIdentifier) : strings.Index(schemaID, FragSep+VersionPathResource)]

	return rid, nil
}

// ExtractSchemaAuthorDID parses the schema URI (did:ugra:<authorDID>;id=<uuid>;version=<version>)
// and returns the author's DID.
func ExtractSchemaAuthorDID(schemaID string) (*did.DID, error) {
	if !IDRx.MatchString(schemaID) {
		idFormatErr := IDFormatErr{schemaID}
		return nil, idFormatErr
	}
	didStr := schemaID[:strings.Index(schemaID, FragSep+ResourceIDPathResource)]
	return did.ParseDID(didStr)
}

func incrementIntAsString(input string) (string, error) {
	i, err := strconv.Atoi(input)
	if err != nil {
		return "", fmt.Errorf("Could not parse input %s to int", input)
	}
	return strconv.Itoa(i + 1), nil
}

func incrementMinorVersion(previousVersion string) (string, error) {
	majorMinor := strings.Split(previousVersion, ".")
	if len(majorMinor) != 2 {
		return "", fmt.Errorf("input not as expected for previous version: %s", previousVersion)
	}

	major := majorMinor[0]
	minor := majorMinor[1]

	i, err := incrementIntAsString(minor)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", major, i), nil
}

func incrementMajorVersion(previousVersion string) (string, error) {
	majorMinor := strings.Split(previousVersion, ".")
	if len(majorMinor) != 2 {
		return "", fmt.Errorf("Input not as expected for previous version: %s", previousVersion)
	}

	major := majorMinor[0]

	i, err := incrementIntAsString(major)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.%s", i, "0"), nil
}