/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats"
	. "github.com/smartystreets/goconvey/convey"
)

func TestVcloudDatacenter(t *testing.T) {
	setupNats()
	n.Subscribe("config.get.postgres", func(msg *nats.Msg) {
		n.Publish(msg.Reply, []byte(`{"names":["users","datacenters","datacenters","services"],"password":"","url":"postgres://postgres@127.0.0.1","user":""}`))
	})
	setupPg()
	startHandler()

	Convey("Scenario: getting a vcloud datacenter", t, func() {
		Convey("Given the datacenter exists on the database", func() {
			createVcloudEntities(1)
			e := Entity{}
			db.Last(&e)
			id := fmt.Sprint(e.ID)

			msg, err := n.Request("datacenter.get", []byte(`{"id":`+id+`}`), time.Second)
			output := Entity{}
			json.Unmarshal(msg.Data, &output)
			So(output.ID, ShouldEqual, e.ID)
			So(output.Name, ShouldEqual, e.Name)
			So(output.Type, ShouldEqual, e.Type)
			So(output.Region, ShouldEqual, e.Region)
			So(output.Username, ShouldEqual, e.Username)
			So(output.Password, ShouldEqual, e.Password)
			So(output.VCloudURL, ShouldEqual, e.VCloudURL)
			So(output.VseURL, ShouldEqual, e.VseURL)
			So(output.ExternalNetwork, ShouldEqual, e.ExternalNetwork)
			So(output.Token, ShouldEqual, e.Token)
			So(output.Secret, ShouldEqual, e.Secret)
			So(err, ShouldEqual, nil)
		})

		Convey("Given the datacenter exists on the database and searching by name", func() {
			e := Entity{}
			db.Last(&e)

			msg, err := n.Request("datacenter.get", []byte(`{"name":"`+e.Name+`"}`), time.Second)
			output := Entity{}
			json.Unmarshal(msg.Data, &output)

			So(output.ID, ShouldEqual, e.ID)
			So(output.Name, ShouldEqual, e.Name)
			So(output.Type, ShouldEqual, e.Type)
			So(output.Region, ShouldEqual, e.Region)
			So(output.Username, ShouldEqual, e.Username)
			So(output.Password, ShouldEqual, e.Password)
			So(output.VCloudURL, ShouldEqual, e.VCloudURL)
			So(output.VseURL, ShouldEqual, e.VseURL)
			So(output.ExternalNetwork, ShouldEqual, e.ExternalNetwork)
			So(output.Token, ShouldEqual, e.Token)
			So(output.Secret, ShouldEqual, e.Secret)
			So(err, ShouldEqual, nil)
		})

	})
}
