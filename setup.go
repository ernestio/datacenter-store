/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/nats-io/nats"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Sets up the nats client based on NATS_URI environment variable
func setupNats() {
	var natsURL = os.Getenv("NATS_URI")
	n, err = nats.Connect(natsURL)
	if err != nil {
		panic(err)
	}
}

func setupPg() {
	var cfg map[string]interface{}
	var err error
	var resp *nats.Msg

	resp, err = n.Request("config.get.postgres", nil, time.Second)
	if err != nil {
		log.Println("could not load config.")
		panic(err)
	}

	err = json.Unmarshal(resp.Data, &cfg)
	if err != nil {
		log.Println("could not read config.")
		panic(err)
	}

	for i := 0; i < 3; i++ {
		pgURL := fmt.Sprintf("%s/%s?sslmode=disable", cfg["url"], "datacenters")
		db, err = gorm.Open("postgres", pgURL)
		if err != nil {
			log.Println("could not connect to postgres. retrying")
			time.Sleep(time.Second * 10)
			continue
		}

		err = db.AutoMigrate(&Entity{}).Error
		if err != nil {
			log.Println("could not connect run migrations. retrying")
			time.Sleep(time.Second * 10)
			continue
		}
		return
	}

	panic(err)
}
