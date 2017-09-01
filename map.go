/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Map : holds a map[string]interface{} value that can be loaded/serialized to a JSONB field
type Map map[string]interface{}

// Value : returns a valid []byte json object
func (m Map) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan : serializes the jsonb object to a map[string]interface{}
func (m *Map) Scan(src interface{}) error {
	var ok bool
	var i interface{}
	var source []byte

	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return errors.New("type assertion .([]byte) & .(string) failed")
	}

	if string(source) == "null" {
		source = []byte("{}")
	}

	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*m, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("type assertion .(map[string]interface{}) failed")
	}

	return nil
}
