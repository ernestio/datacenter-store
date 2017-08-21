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
	_, _ = n.Subscribe("config.get.postgres", func(msg *nats.Msg) {
		_ = n.Publish(msg.Reply, []byte(`{"names":["users","datacenters","datacenters","services"],"password":"","url":"postgres://postgres@127.0.0.1","user":""}`))
	})
	createTestDB("test_vcloud")
	setupPg("test_vcloud")
	startHandler()

	Convey("Scenario: getting a vcloud datacenter", t, func() {
		Convey("Given the datacenter exists on the database", func() {
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
			So(output.Region, ShouldEqual, e.Region)
			So(output.Username, ShouldEqual, e.Username)
			So(output.Password, ShouldEqual, e.Password)
			So(output.VCloudURL, ShouldEqual, e.VCloudURL)
			So(output.VseURL, ShouldEqual, e.VseURL)
			So(output.ExternalNetwork, ShouldEqual, e.ExternalNetwork)
			So(output.AccessKeyID, ShouldEqual, e.AccessKeyID)
			So(output.SecretAccessKey, ShouldEqual, e.SecretAccessKey)
			So(err, ShouldBeNil)
		})

		Convey("Given the datacenter exists on the database and searching by name", func() {
			e := Entity{}
			db.Last(&e)

			msg, err := n.Request("datacenter.get", []byte(`{"name":"`+e.Name+`"}`), time.Second)
			output := Entity{}
			_ = json.Unmarshal(msg.Data, &output)

			So(output.ID, ShouldEqual, e.ID)
			So(output.Name, ShouldEqual, e.Name)
			So(output.Type, ShouldEqual, e.Type)
			So(output.Region, ShouldEqual, e.Region)
			So(output.Username, ShouldEqual, e.Username)
			So(output.Password, ShouldEqual, e.Password)
			So(output.VCloudURL, ShouldEqual, e.VCloudURL)
			So(output.VseURL, ShouldEqual, e.VseURL)
			So(output.ExternalNetwork, ShouldEqual, e.ExternalNetwork)
			So(output.AccessKeyID, ShouldEqual, e.AccessKeyID)
			So(output.SecretAccessKey, ShouldEqual, e.SecretAccessKey)
			So(err, ShouldBeNil)
		})

		Convey("Given the datacenter exists on the database and searching with datacenter.find by name", func() {
			e := Entity{}
			db.Last(&e)

			msg, err := n.Request("datacenter.find", []byte(`{"name":"`+e.Name+`"}`), time.Second)
			output := []Entity{}
			_ = json.Unmarshal(msg.Data, &output)

			So(len(output), ShouldEqual, 1)

			So(output[0].ID, ShouldEqual, e.ID)
			So(output[0].Name, ShouldEqual, e.Name)
			So(output[0].Type, ShouldEqual, e.Type)
			So(output[0].Region, ShouldEqual, e.Region)
			So(output[0].Username, ShouldEqual, e.Username)
			So(output[0].Password, ShouldEqual, e.Password)
			So(output[0].VCloudURL, ShouldEqual, e.VCloudURL)
			So(output[0].VseURL, ShouldEqual, e.VseURL)
			So(output[0].ExternalNetwork, ShouldEqual, e.ExternalNetwork)
			So(output[0].AccessKeyID, ShouldEqual, e.AccessKeyID)
			So(output[0].SecretAccessKey, ShouldEqual, e.SecretAccessKey)
			So(err, ShouldBeNil)
		})
	})
}
