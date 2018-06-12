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

func TestAWSDatacenter(t *testing.T) {
	setupNats()
	defer n.Close()
	_, _ = n.Subscribe("config.get.postgres", func(msg *nats.Msg) {
		_ = n.Publish(msg.Reply, []byte(`{"names":["users","datacenters","datacenters","services"],"password":"","url":"postgres://postgres@127.0.0.1","user":""}`))
	})

	createTestDB("test_aws")
	setupPg("test_aws")
	startHandler()

	Convey("Scenario: getting a aws project", t, func() {
		Convey("Given the project exists on the database", func() {
			createAWSEntities(1)
			e := Entity{}
			db.Last(&e)
			id := fmt.Sprint(e.ID)

			msg, err := n.Request("datacenter.get", []byte(`{"id":`+id+`}`), time.Second)
			output := Entity{}
			err = json.Unmarshal(msg.Data, &output)
			So(err, ShouldBeNil)
			So(output.ID, ShouldEqual, e.ID)
			So(output.Name, ShouldEqual, e.Name)
			So(output.Type, ShouldEqual, e.Type)
			So(output.Credentials["region"], ShouldEqual, "eu-west-1")
			So(output.Credentials["access_key_id"], ShouldEqual, "test-id")
			So(output.Credentials["secret_access_key"], ShouldEqual, "test-key")
			So(err, ShouldBeNil)
		})

		Convey("Given the project exists on the database and searching by name", func() {
			e := Entity{}
			db.Last(&e)

			msg, err := n.Request("datacenter.get", []byte(`{"name":"`+e.Name+`"}`), time.Second)
			output := Entity{}
			err = json.Unmarshal(msg.Data, &output)
			So(err, ShouldBeNil)

			So(output.ID, ShouldEqual, e.ID)
			So(output.Name, ShouldEqual, e.Name)
			So(output.Type, ShouldEqual, e.Type)
			So(output.Credentials["region"], ShouldEqual, "eu-west-1")
			So(output.Credentials["access_key_id"], ShouldEqual, "test-id")
			So(output.Credentials["secret_access_key"], ShouldEqual, "test-key")
			So(err, ShouldBeNil)
		})

	})

}
