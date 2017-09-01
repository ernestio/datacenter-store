/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/lib/pq"
)

func setupTestSuite() {
	db.Unscoped().Delete(Entity{})
}

func createEntities(n int) {
	i := 0
	for i < n {
		x := strconv.Itoa(i)
		db.Create(&Entity{Name: "Test" + x, Credentials: Map{"access_key_id": "test-id", "secret_access_key": "test-key", "region": "eu-west-1"}})
		i++
	}
}

func createVcloudEntities(n int) {
	i := 0
	for i < n {
		x := strconv.Itoa(i)
		db.Create(&Entity{Name: "TestVcloud" + x, Type: "vcloud", Credentials: Map{"vcloud_url": "http://vcloud.com", "external_network": "ext-100", "username": "test", "password": "test"}})
		i++
	}
}

func createAWSEntities(n int) {
	i := 0
	for i < n {
		x := strconv.Itoa(i)
		db.Create(&Entity{Name: "TestAWS" + x, Type: "aws", Credentials: Map{"access_key_id": "test-id", "secret_access_key": "test-key", "region": "eu-west-1"}})
		i++
	}
}

func createTestDB(name string) error {
	db, derr := sql.Open("postgres", "user=postgres sslmode=disable")
	if derr != nil {
		return derr
	}

	_, derr = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", pq.QuoteIdentifier(name)))
	if derr != nil {
		return derr
	}

	_, derr = db.Exec(fmt.Sprintf("CREATE DATABASE %s", pq.QuoteIdentifier(name)))

	return derr
}
