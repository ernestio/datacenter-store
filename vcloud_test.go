/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/go-nats"
	. "github.com/smartystreets/goconvey/convey"
)

func TestVcloudDatacenter(t *testing.T) {
	setupNats()
	defer n.Close()
	_, _ = n.Subscribe("config.get.postgres", func(msg *nats.Msg) {
		_ = n.Publish(msg.Reply, []byte(`{"names":["users","datacenters","datacenters","services"],"password":"","url":"postgres://postgres@127.0.0.1","user":""}`))
	})
	createTestDB("test_vcloud")
	setupPg("test_vcloud")
	startHandler()

	Convey("Scenario: getting a vcloud project", t, func() {
		Convey("Given the project exists on the database", func() {
			createVcloudEntities(1)
			e := Entity{}
			db.Last(&e)
			id := fmt.Sprint(e.ID)

			msg, err := n.Request("datacenter.get", []byte(`{"id":`+id+`}`), time.Second)
			output := Entity{}
			_ = json.Unmarshal(msg.Data, &output)
			So(output.ID, ShouldEqual, e.ID)
			So(output.Name, ShouldEqual, e.Name)
			So(output.Type, ShouldEqual, e.Type)
			So(output.Credentials["vcloud_url"], ShouldEqual, "http://vcloud.com")
			So(output.Credentials["external_network"], ShouldEqual, "ext-100")
			So(output.Credentials["username"], ShouldEqual, "test")
			So(output.Credentials["password"], ShouldEqual, "test")
			So(err, ShouldBeNil)
		})

		Convey("Given the project exists on the database and searching by name", func() {
			e := Entity{}
			db.Last(&e)

			msg, err := n.Request("datacenter.get", []byte(`{"name":"`+e.Name+`"}`), time.Second)
			output := Entity{}
			_ = json.Unmarshal(msg.Data, &output)

			So(output.ID, ShouldEqual, e.ID)
			So(output.Name, ShouldEqual, e.Name)
			So(output.Type, ShouldEqual, e.Type)
			So(output.Credentials["vcloud_url"], ShouldEqual, "http://vcloud.com")
			So(output.Credentials["external_network"], ShouldEqual, "ext-100")
			So(output.Credentials["username"], ShouldEqual, "test")
			So(output.Credentials["password"], ShouldEqual, "test")
			So(err, ShouldBeNil)
		})

		Convey("Given the project exists on the database and searching with project.find by name", func() {
			e := Entity{}
			db.Last(&e)

			msg, err := n.Request("datacenter.find", []byte(`{"name":"`+e.Name+`"}`), time.Second)
			output := []Entity{}
			_ = json.Unmarshal(msg.Data, &output)

			So(len(output), ShouldEqual, 1)

			So(output[0].ID, ShouldEqual, e.ID)
			So(output[0].Name, ShouldEqual, e.Name)
			So(output[0].Type, ShouldEqual, e.Type)
			So(output[0].Credentials["vcloud_url"], ShouldEqual, "http://vcloud.com")
			So(output[0].Credentials["external_network"], ShouldEqual, "ext-100")
			So(output[0].Credentials["username"], ShouldEqual, "test")
			So(output[0].Credentials["password"], ShouldEqual, "test")
			So(err, ShouldBeNil)
		})
	})
}
