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

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// URI is a wrapper around url.URL to add json marshalling
type URI struct {
	url.URL
}

func (v URI) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

func (v *URI) UnmarshalJSON(bytes []byte) error {
	var value string
	if err := json.Unmarshal(bytes, &value); err != nil {
		return err
	}
	parsedUrl, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("could not parse URI: %w", err)
	}
	v.URL = *parsedUrl
	return nil
}

// ParseURI parses a raw URI. If it can't be parsed, an error is returned.
func ParseURI(input string) (*URI, error) {
	u, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	return &URI{URL: *u}, nil
}

func (v URI) String() string {
	return v.URL.String()
}