/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"runtime"

	"github.com/jinzhu/gorm"
	"github.com/nats-io/nats"
	"github.com/r3labs/natsdb"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var n *nats.Conn
var db *gorm.DB
var err error
var handler natsdb.Handler

func startHandler() {
	handler = natsdb.Handler{
		NotFoundErrorMessage:   []byte(`{"error":"not found"}`),
		UnexpectedErrorMessage: []byte(`{"error":"unexpected"}`),
		DeletedMessage:         []byte(`"deleted"`),
		Nats:                   n,
		NewModel: func() natsdb.Model {
			return &Entity{}
		},
	}

	n.Subscribe("datacenter.get", handler.Get)
	n.Subscribe("datacenter.del", handler.Del)
	n.Subscribe("datacenter.set", handler.Set)
	n.Subscribe("datacenter.find", handler.Find)
}

func main() {
	setupNats()
	setupPg()
	startHandler()

	runtime.Goexit()
}
