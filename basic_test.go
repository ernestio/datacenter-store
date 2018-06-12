/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	aes "github.com/ernestio/crypto/aes"
	"github.com/nats-io/go-nats"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetHandler(t *testing.T) {
	setupNats()
	defer n.Close()
	_, _ = n.Subscribe("config.get.postgres", func(msg *nats.Msg) {
		_ = n.Publish(msg.Reply, []byte(`{"names":["users","datacenters","datacenters","services"],"password":"","url":"postgres://postgres@127.0.0.1","user":""}`))
	})
	createTestDB("test_projects")
	setupPg("test_projects")
	startHandler()

	Convey("Scenario: getting a project", t, func() {
		setupTestSuite()
		Convey("Given the project does not exist on the database", func() {
			msg, err := n.Request("datacenter.get", []byte(`{"id":32}`), time.Second)
			So(string(msg.Data), ShouldEqual, string(handler.NotFoundErrorMessage))
			So(err, ShouldBeNil)
		})

		Convey("Given the project exists on the database", func() {
			createEntities(1)
			e := Entity{}
			db.First(&e)
			id := fmt.Sprint(e.ID)

			msg, err := n.Request("datacenter.get", []byte(`{"id":`+id+`}`), time.Second)
			output := Entity{}
			err = json.Unmarshal(msg.Data, &output)
			So(err, ShouldBeNil)
			So(output.ID, ShouldEqual, e.ID)
			So(output.Name, ShouldEqual, e.Name)
			So(output.Type, ShouldEqual, e.Type)
			So(err, ShouldBeNil)
		})

		Convey("Given the project exists on the database and searching by name", func() {
			createEntities(1)
			e := Entity{}
			db.First(&e)

			msg, err := n.Request("datacenter.get", []byte(`{"name":"`+e.Name+`"}`), time.Second)
			output := Entity{}
			err = json.Unmarshal(msg.Data, &output)
			So(err, ShouldBeNil)
			So(output.ID, ShouldEqual, e.ID)
			So(output.Name, ShouldEqual, e.Name)
			So(output.Type, ShouldEqual, e.Type)
			So(err, ShouldBeNil)
		})
	})

	Convey("Scenario: deleting a project", t, func() {
		setupTestSuite()
		Convey("Given the project does not exist on the database", func() {
			msg, err := n.Request("datacenter.del", []byte(`{"id":32}`), time.Second)
			So(string(msg.Data), ShouldEqual, string(handler.NotFoundErrorMessage))
			So(err, ShouldBeNil)
		})

		Convey("Given the project exists on the database", func() {
			createEntities(1)
			last := Entity{}
			db.First(&last)
			id := fmt.Sprint(last.ID)

			msg, err := n.Request("datacenter.del", []byte(`{"id":`+id+`}`), time.Second)
			So(string(msg.Data), ShouldEqual, string(handler.DeletedMessage))
			So(err, ShouldBeNil)

			deleted := Entity{}
			db.First(&deleted, id)
			So(deleted.ID, ShouldEqual, 0)
		})
	})

	Convey("Scenario: project set", t, func() {
		setupTestSuite()
		Convey("Given we don't provide any id as part of the body", func() {
			Convey("Then it should return the created record and it should be stored on DB", func() {
				msg, err := n.Request("datacenter.set", []byte(`{"name":"test-101","aws_access_token_id":"foo","aws_secret_access_key":"bar", "type": "fake"}`), time.Second)
				output := Entity{}
				output.LoadFromInput(msg.Data)
				So(output.ID, ShouldNotEqual, nil)
				So(output.Name, ShouldEqual, "test-101")
				So(err, ShouldBeNil)

				stored := Entity{}
				db.First(&stored, output.ID)
				So(stored.Name, ShouldEqual, "test-101")
			})
		})

		Convey("Given we provide an unexisting id", func() {
			Convey("Then we should receive a not found message", func() {
				msg, err := n.Request("datacenter.set", []byte(`{"id": 1000, "name":"test-100", "type": "fake"}`), time.Second)
				So(string(msg.Data), ShouldEqual, string(handler.NotFoundErrorMessage))
				So(err, ShouldBeNil)
			})
		})

		Convey("Given we provide an existing id", func() {
			setupTestSuite()
			createEntities(1)
			e := Entity{}
			db.First(&e)
			id := fmt.Sprint(e.ID)
			Convey("Then we should receive an updated entity", func() {
				msg, err := n.Request("datacenter.set", []byte(`{"id": `+id+`, "name":"test-100", "type": "fake"}`), time.Second)
				output := Entity{}
				output.LoadFromInput(msg.Data)
				So(output.ID, ShouldNotEqual, nil)
				So(output.Name, ShouldEqual, "test-100")
				So(err, ShouldBeNil)

				stored := Entity{}
				db.First(&stored, output.ID)
				So(stored.Name, ShouldEqual, "test-100")
			})
		})
	})

	Convey("Scenario: find projects", t, func() {
		setupTestSuite()
		Convey("Given projects exist on the database", func() {
			createEntities(20)
			Convey("Then I should get a list of projects", func() {
				msg, _ := n.Request("datacenter.find", []byte(`{}`), time.Second)
				list := []Entity{}
				err = json.Unmarshal(msg.Data, &list)
				So(err, ShouldBeNil)
				So(len(list), ShouldEqual, 20)
			})
		})
	})

	Convey("Scenario: find projects by multiple ids", t, func() {
		setupTestSuite()
		Convey("Given projects exist on the database", func() {
			createEntities(20)
			Convey("Then I should get a list of projects", func() {
				msg, _ := n.Request("datacenter.find", []byte(`{}`), time.Second)
				list := []Entity{}
				err = json.Unmarshal(msg.Data, &list)
				So(err, ShouldBeNil)
				msg, _ = n.Request("datacenter.find", []byte(`{"ids":["`+fmt.Sprint(list[0].ID)+`","`+fmt.Sprint(list[1].ID)+`","`+fmt.Sprint(list[2].ID)+`"]}`), time.Second)
				err = json.Unmarshal(msg.Data, &list)
				So(err, ShouldBeNil)
				So(len(list), ShouldEqual, 3)
			})
		})
	})
}

func TestUpdateHandler(t *testing.T) {
	setupNats()
	defer n.Close()
	_, _ = n.Subscribe("config.get.postgres", func(msg *nats.Msg) {
		_ = n.Publish(msg.Reply, []byte(`{"names":["users","datacenters","datacenters","services"],"password":"","url":"postgres://postgres@127.0.0.1","user":""}`))
	})
	createTestDB("test_projects")
	setupPg("test_projects")
	startHandler()
	Convey("Scenario: update projects", t, func() {
		setupTestSuite()
		Convey("Given projects exist on the database", func() {
			createEntities(20)
			Convey("Then I should be able to create a project", func() {
				var list []Entity
				entity := Entity{
					Name: "supu",
					Credentials: Map{
						"access_key_id":     "blah",
						"secret_access_key": "blah",
					},
				}

				body, _ := json.Marshal(entity)
				_, _ = n.Request("datacenter.set", body, time.Second)

				msg, _ := n.Request("datacenter.find", []byte(`{"name":"`+entity.Name+`"}`), time.Second)
				err = json.Unmarshal(msg.Data, &list)
				So(err, ShouldBeNil)
				So(len(list), ShouldEqual, 1)
				So(list[0].Name, ShouldEqual, entity.Name)
				So(list[0].Credentials["access_key_id"], ShouldNotEqual, entity.Credentials["access_key_id"])
				So(list[0].Credentials["secret_access_key"], ShouldNotEqual, entity.Credentials["secret_access_key"])

				crypto := aes.New()
				key := os.Getenv("ERNEST_CRYPTO_KEY")
				token, err := crypto.Decrypt(list[0].Credentials["access_key_id"].(string), key)
				So(err, ShouldBeNil)
				So(token, ShouldEqual, entity.Credentials["access_key_id"])
				secret, err := crypto.Decrypt(list[0].Credentials["secret_access_key"].(string), key)
				So(err, ShouldBeNil)
				So(secret, ShouldEqual, entity.Credentials["secret_access_key"])
			})
		})
	})

}
