/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"
	"time"

	"github.com/nats-io/nats"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAWSDatacenter(t *testing.T) {
	setupNats()
	n.Subscribe("config.get.postgres", func(msg *nats.Msg) {
		n.Publish(msg.Reply, []byte(`{"names":["users","datacenters","datacenters","services"],"password":"","url":"postgres://postgres@127.0.0.1","user":""}`))
	})
	setupPg()
	startHandler()

	Convey("Scenario: creating a vcloud datacenter", t, func() {
		setupTestSuite()
		Convey("Given i create a datacenter with type vcloud", func() {
			Convey("Then it should return the created record and it should be stored on DB", func() {
				msg, err := n.Request("datacenter.set", []byte(`{"name":"awsDC", "type":"aws"}`), time.Second)
				output := Entity{}
				output.LoadFromInput(msg.Data)
				So(output.ID, ShouldNotEqual, nil)
				So(output.Name, ShouldEqual, "awsDC")
				So(output.Type, ShouldEqual, "aws")
				So(err, ShouldEqual, nil)
			})
		})
	})

}
