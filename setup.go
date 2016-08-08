/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"log"
	"os"
	"time"

	ecc "github.com/ernestio/ernest-config-client"
)

var c *ecc.Config

func setupNats() {
	c = ecc.NewConfig(os.Getenv("NATS_URI"))
	n = c.Nats()
}

func setupPg() {
	db = c.Postgres("datacenters")
	for true {
		if err = db.AutoMigrate(&Entity{}).Error; err != nil {
			log.Println("could not connect run migrations. retrying")
			time.Sleep(time.Second * 10)
			continue
		}
		return
	}
}
