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
	"github.com/nats-io/nats"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetHandler(t *testing.T) {
	setupNats()
	_, _ = n.Subscribe("config.get.postgres", func(msg *nats.Msg) {
		_ = n.Publish(msg.Reply, []byte(`{"names":["users","datacenters","datacenters","services"],"password":"","url":"postgres://postgres@127.0.0.1","user":""}`))
	})
	setupPg()
	startHandler()

	Convey("Scenario: getting a datacenter", t, func() {
		setupTestSuite()
		Convey("Given the datacenter does not exist on the database", func() {
			msg, err := n.Request("datacenter.get", []byte(`{"id":"32"}`), time.Second)
			So(string(msg.Data), ShouldEqual, string(handler.NotFoundErrorMessage))
			So(err, ShouldEqual, nil)
		})

		Convey("Given the datacenter exists on the database", func() {
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
			So(output.Region, ShouldEqual, e.Region)
			So(output.Username, ShouldEqual, e.Username)
			So(output.Password, ShouldEqual, e.Password)
			So(output.VCloudURL, ShouldEqual, e.VCloudURL)
			So(output.VseURL, ShouldEqual, e.VseURL)
			So(output.ExternalNetwork, ShouldEqual, e.ExternalNetwork)
			So(output.AccessKeyID, ShouldEqual, e.AccessKeyID)
			So(output.SecretAccessKey, ShouldEqual, e.SecretAccessKey)
			So(err, ShouldEqual, nil)
		})

		Convey("Given the datacenter exists on the database and searching by name", func() {
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
			So(output.Region, ShouldEqual, e.Region)
			So(output.Username, ShouldEqual, e.Username)
			So(output.Password, ShouldEqual, e.Password)
			So(output.VCloudURL, ShouldEqual, e.VCloudURL)
			So(output.VseURL, ShouldEqual, e.VseURL)
			So(output.ExternalNetwork, ShouldEqual, e.ExternalNetwork)
			So(output.AccessKeyID, ShouldEqual, e.AccessKeyID)
			So(output.SecretAccessKey, ShouldEqual, e.SecretAccessKey)
			So(err, ShouldEqual, nil)
		})
	})

	Convey("Scenario: deleting a datacenter", t, func() {
		setupTestSuite()
		Convey("Given the datacenter does not exist on the database", func() {
			msg, err := n.Request("datacenter.del", []byte(`{"id":32}`), time.Second)
			So(string(msg.Data), ShouldEqual, string(handler.NotFoundErrorMessage))
			So(err, ShouldEqual, nil)
		})

		Convey("Given the datacenter exists on the database", func() {
			createEntities(1)
			last := Entity{}
			db.First(&last)
			id := fmt.Sprint(last.ID)

			msg, err := n.Request("datacenter.del", []byte(`{"id":`+id+`}`), time.Second)
			So(string(msg.Data), ShouldEqual, string(handler.DeletedMessage))
			So(err, ShouldEqual, nil)

			deleted := Entity{}
			db.First(&deleted, id)
			So(deleted.ID, ShouldEqual, 0)
		})
	})

	Convey("Scenario: datacenter set", t, func() {
		setupTestSuite()
		Convey("Given we don't provide any id as part of the body", func() {
			Convey("Then it should return the created record and it should be stored on DB", func() {
				msg, err := n.Request("datacenter.set", []byte(`{"name":"fred","aws_access_token_id":"foo","aws_secret_access_key":"bar"}`), time.Second)
				output := Entity{}
				output.LoadFromInput(msg.Data)
				So(output.ID, ShouldNotEqual, nil)
				So(output.Name, ShouldEqual, "fred")
				So(err, ShouldEqual, nil)

				stored := Entity{}
				db.First(&stored, output.ID)
				So(stored.Name, ShouldEqual, "fred")
			})
		})

		Convey("Given we provide an unexisting id", func() {
			Convey("Then we should receive a not found message", func() {
				msg, err := n.Request("datacenter.set", []byte(`{"id": 1000, "name":"fred"}`), time.Second)
				So(string(msg.Data), ShouldEqual, string(handler.NotFoundErrorMessage))
				So(err, ShouldEqual, nil)
			})
		})

		Convey("Given we provide an existing id", func() {
			setupTestSuite()
			createEntities(1)
			e := Entity{}
			db.First(&e)
			id := fmt.Sprint(e.ID)
			Convey("Then we should receive an updated entity", func() {
				msg, err := n.Request("datacenter.set", []byte(`{"id": `+id+`, "name":"fred"}`), time.Second)
				output := Entity{}
				output.LoadFromInput(msg.Data)
				So(output.ID, ShouldNotEqual, nil)
				So(output.Name, ShouldEqual, "fred")
				So(err, ShouldEqual, nil)

				stored := Entity{}
				db.First(&stored, output.ID)
				So(stored.Name, ShouldEqual, "fred")
			})
		})
	})

	Convey("Scenario: find datacenters", t, func() {
		setupTestSuite()
		Convey("Given datacenters exist on the database", func() {
			createEntities(20)
			Convey("Then I should get a list of datacenters", func() {
				msg, _ := n.Request("datacenter.find", []byte(`{"group_id":2}`), time.Second)
				list := []Entity{}
				err = json.Unmarshal(msg.Data, &list)
				So(err, ShouldBeNil)
				So(len(list), ShouldEqual, 1)
			})
		})
	})

	Convey("Scenario: find datacenters by multiple ids", t, func() {
		setupTestSuite()
		Convey("Given datacenters exist on the database", func() {
			createEntities(20)
			Convey("Then I should get a list of datacenters", func() {
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
	_, _ = n.Subscribe("config.get.postgres", func(msg *nats.Msg) {
		_ = n.Publish(msg.Reply, []byte(`{"names":["users","datacenters","datacenters","services"],"password":"","url":"postgres://postgres@127.0.0.1","user":""}`))
	})
	setupPg()
	startHandler()
	Convey("Scenario: update datacenters", t, func() {
		setupTestSuite()
		Convey("Given datacenters exist on the database", func() {
			createEntities(20)
			Convey("Then I should get a list of datacenters", func() {
				msg, _ := n.Request("datacenter.find", []byte(`{"group_id":2}`), time.Second)
				list := []Entity{}
				err = json.Unmarshal(msg.Data, &list)
				So(err, ShouldBeNil)
				So(len(list), ShouldEqual, 1)
				entity := list[0]
				entity.Name = "supu"
				entity.GroupID = 4
				entity.AccessKeyID = "blah"
				entity.SecretAccessKey = "blah"
				body, _ := json.Marshal(entity)
				msg, _ = n.Request("datacenter.set", body, time.Second)

				msg, _ = n.Request("datacenter.find", []byte(`{"name":"`+entity.Name+`"}`), time.Second)
				err = json.Unmarshal(msg.Data, &list)
				So(err, ShouldBeNil)
				So(len(list), ShouldEqual, 1)
				So(list[0].Name, ShouldEqual, entity.Name)
				So(list[0].GroupID, ShouldEqual, entity.GroupID)
				So(list[0].AccessKeyID, ShouldNotEqual, entity.AccessKeyID)
				So(list[0].SecretAccessKey, ShouldNotEqual, entity.SecretAccessKey)

				crypto := aes.New()
				key := os.Getenv("ERNEST_CRYPTO_KEY")
				token, err := crypto.Decrypt(list[0].AccessKeyID, key)
				So(err, ShouldBeNil)
				So(token, ShouldEqual, entity.AccessKeyID)
				secret, err := crypto.Decrypt(list[0].SecretAccessKey, key)
				So(err, ShouldBeNil)
				So(secret, ShouldEqual, entity.SecretAccessKey)

				encryptedAccessKeyID := token
				encryptedSecretAccessKey := secret
				entity = list[0]
				entity.Name = "supu"
				entity.GroupID = 3
				entity.AccessKeyID = ""
				entity.SecretAccessKey = ""
				body, _ = json.Marshal(entity)
				msg, _ = n.Request("datacenter.set", body, time.Second)
				msg, _ = n.Request("datacenter.find", []byte(`{"name":"`+entity.Name+`"}`), time.Second)
				err = json.Unmarshal(msg.Data, &list)
				So(err, ShouldBeNil)

				token, err = crypto.Decrypt(list[0].AccessKeyID, key)
				So(err, ShouldBeNil)
				So(encryptedAccessKeyID, ShouldEqual, encryptedAccessKeyID)
				secret, err = crypto.Decrypt(list[0].SecretAccessKey, key)
				So(err, ShouldBeNil)
				So(encryptedSecretAccessKey, ShouldEqual, encryptedSecretAccessKey)
			})
		})
	})

}
