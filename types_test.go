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
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestParseURI(t *testing.T) {

	t.Run("for VC types", func(t *testing.T) {
		u, err := ParseURI("SomeType")

		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, "SomeType", u.String())
	})

	t.Run("for URI types", func(t *testing.T) {
		u, err := ParseURI("https://example.com/context/v1")

		if !assert.NoError(t, err) {
			return
		}

		assert.Equal(t, "https://example.com/context/v1", u.String())
	})

	t.Run("malformed input", func(t *testing.T) {
		_, err := ParseURI(string([]byte{0}))

		assert.Error(t, err)
	})
}

func TestURI_String(t *testing.T) {
	assert.Equal(t, "http://test", URI{url.URL{Scheme: "http", Host: "test"}}.String())
}
