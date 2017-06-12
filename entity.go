/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	aes "github.com/ernestio/crypto/aes"
	"github.com/nats-io/nats"
	"github.com/r3labs/natsdb"
)

// Entity : the database mapped entity
type Entity struct {
	ID              uint   `json:"id" gorm:"primary_key"`
	GroupID         uint   `json:"group_id" gorm:"unique_index:idx_per_group"`
	Name            string `json:"name" gorm:"unique_index:idx_per_group"`
	Type            string `json:"type" gorm:"unique_index:idx_per_group"`
	Region          string `json:"region"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	VCloudURL       string `json:"vcloud_url"`
	VseURL          string `json:"vse_url"`
	ExternalNetwork string `json:"external_network"`
	AccessKeyID     string `json:"aws_access_key_id"`
	SecretAccessKey string `json:"aws_secret_access_key"`
	SubscriptionID  string `json:"azure_subscription_id"`
	ClientID        string `json:"azure_client_id"`
	ClientSecret    string `json:"azure_client_secret"`
	TenantID        string `json:"azure_tenant_id"`
	Environment     string `json:"azure_environment"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time `json:"-" sql:"index"`
}

// TableName : set Entity's table name to be datacenters
func (Entity) TableName() string {
	return "datacenters"
}

// Find : based on the defined fields for the current entity
// will perform a search on the database
func (e *Entity) Find() []interface{} {
	entities := []Entity{}
	if e.Name != "" && e.GroupID != 0 {
		db.Where("name = ?", e.Name).Where("group_id = ?", e.GroupID).Find(&entities)
	} else {
		if e.Name != "" {
			db.Where("name = ?", e.Name).Find(&entities)
		} else if e.GroupID != 0 {
			db.Where("group_id = ?", e.GroupID).Find(&entities)
		} else {
			db.Find(&entities)
		}
	}

	list := make([]interface{}, len(entities))
	for i, s := range entities {
		list[i] = s
	}

	return list
}

// MapInput : maps the input []byte on the current entity
func (e *Entity) MapInput(body []byte) {
	if err := json.Unmarshal(body, &e); err != nil {
		log.Println("Invalid input " + err.Error())
	}
}

// HasID : determines if the current entity has an id or not
func (e *Entity) HasID() bool {
	if e.ID == 0 {
		return false
	}
	return true
}

// LoadFromInput : Will load from a []byte input the database stored entity
func (e *Entity) LoadFromInput(msg []byte) bool {
	e.MapInput(msg)
	var stored Entity
	if e.ID != 0 {
		db.First(&stored, e.ID)
	} else if e.Name != "" {
		db.Where("name = ?", e.Name).First(&stored)
	}
	if &stored == nil {
		return false
	}
	if ok := stored.HasID(); !ok {
		return false
	}

	e.ID = stored.ID
	e.GroupID = stored.GroupID
	e.Name = stored.Name
	e.Type = stored.Type
	e.Username = stored.Username
	e.Password = stored.Password
	e.ExternalNetwork = stored.ExternalNetwork
	e.AccessKeyID = stored.AccessKeyID
	e.SecretAccessKey = stored.SecretAccessKey
	e.Region = stored.Region
	e.VCloudURL = stored.VCloudURL
	e.VseURL = stored.VseURL
	e.SubscriptionID = stored.SubscriptionID
	e.ClientID = stored.ClientID
	e.ClientSecret = stored.ClientSecret
	e.TenantID = stored.TenantID
	e.Environment = stored.Environment
	e.CreatedAt = stored.CreatedAt
	e.UpdatedAt = stored.UpdatedAt

	return true
}

// LoadFromInputOrFail : Will try to load from the input an existing entity,
// or will call the handler to Fail the nats message
func (e *Entity) LoadFromInputOrFail(msg *nats.Msg, h *natsdb.Handler) bool {
	stored := &Entity{}
	ok := stored.LoadFromInput(msg.Data)
	if !ok {
		h.Fail(msg)
	}
	*e = *stored

	return ok
}

// Update : It will update the current entity with the input []byte
func (e *Entity) Update(body []byte) error {
	e.MapInput(body)
	stored := Entity{}
	db.First(&stored, e.ID)
	stored.Name = e.Name
	stored.GroupID = e.GroupID

	if e.Username != "" {
		stored.Username, _ = crypt(e.Username)
	}
	if e.Password != "" {
		stored.Password, _ = crypt(e.Password)
	}
	if e.AccessKeyID != "" {
		stored.AccessKeyID, _ = crypt(e.AccessKeyID)
	}
	if e.SecretAccessKey != "" {
		stored.SecretAccessKey, _ = crypt(e.SecretAccessKey)
	}
	if e.SubscriptionID != "" {
		stored.SubscriptionID, _ = crypt(e.SubscriptionID)
	}
	if e.ClientID != "" {
		stored.ClientID, _ = crypt(e.ClientID)
	}
	if e.ClientSecret != "" {
		stored.ClientSecret, _ = crypt(e.ClientSecret)
	}
	if e.TenantID != "" {
		stored.TenantID, _ = crypt(e.TenantID)
	}
	if e.Environment != "" {
		stored.Environment, _ = crypt(e.Environment)
	}

	db.Save(&stored)
	e = &stored

	return nil
}

// Delete : Will delete from database the current Entity
func (e *Entity) Delete() error {
	db.Unscoped().Delete(&e)

	return nil
}

func crypt(s string) (string, error) {
	crypto := aes.New()
	key := os.Getenv("ERNEST_CRYPTO_KEY")
	if s != "" {
		encrypted, err := crypto.Encrypt(s, key)
		if err != nil {
			return "", err
		}
		s = encrypted
	}

	return s, nil
}

// Save : Persists current entity on database
func (e *Entity) Save() error {
	var err error

	e.Username, err = crypt(e.Username)
	if err != nil {
		return err
	}

	e.Password, err = crypt(e.Password)
	if err != nil {
		return err
	}

	e.AccessKeyID, err = crypt(e.AccessKeyID)
	if err != nil {
		return err
	}

	e.SecretAccessKey, err = crypt(e.SecretAccessKey)
	if err != nil {
		return err
	}

	e.SubscriptionID, err = crypt(e.SubscriptionID)
	if err != nil {
		return err
	}

	e.ClientID, err = crypt(e.ClientID)
	if err != nil {
		return err
	}

	e.ClientSecret, err = crypt(e.ClientSecret)
	if err != nil {
		return err
	}

	e.TenantID, err = crypt(e.TenantID)
	if err != nil {
		return err
	}

	e.Environment, err = crypt(e.Environment)
	if err != nil {
		return err
	}

	db.Save(&e)

	return nil
}
