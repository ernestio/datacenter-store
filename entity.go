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
	"github.com/nats-io/go-nats"
	"github.com/r3labs/natsdb"
)

// Entity : the database mapped entity
type Entity struct {
	ID          uint     `json:"id" gorm:"primary_key"`
	IDs         []string `json:"ids,omitempty" sql:"-"`
	Name        string   `json:"name" gorm:"unique_index"`
	Names       []string `json:"names,omitempty" sql:"-"`
	Type        string   `json:"type"`
	Credentials Map      `json:"credentials" gorm:"type: jsonb not null default '{}'::jsonb"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `json:"-" sql:"index"`
}

// TableName : set Entity's table name to be datacenters
func (Entity) TableName() string {
	return "projects"
}

// Find : based on the defined fields for the current entity
// will perform a search on the database
func (e *Entity) Find() []interface{} {
	entities := []Entity{}
	if len(e.IDs) > 0 {
		db.Where("id in (?)", e.IDs).Find(&entities)
	} else if len(e.Names) > 0 {
		db.Where("name in (?)", e.Names).Find(&entities)
	} else if e.Name != "" {
		db.Where("name = ?", e.Name).Find(&entities)
	} else {
		db.Find(&entities)
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
	e.Name = stored.Name
	e.Type = stored.Type
	e.Credentials = stored.Credentials
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
	e.Credentials = make(Map)

	e.MapInput(body)
	stored := Entity{}
	db.First(&stored, e.ID)
	stored.Name = e.Name

	ec, err := encryptCredentials(e.Credentials)
	if err != nil {
		return err
	}

	for k, v := range ec {
		stored.Credentials[k] = v
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
	ec, err := encryptCredentials(e.Credentials)
	if err != nil {
		return err
	}

	e.Credentials = ec
	db.Save(&e)

	return nil
}

func encryptCredentials(c Map) (Map, error) {
	for k, v := range c {
		if k == "region" || k == "vdc" || k == "username" || k == "vcloud_url" {
			continue
		}

		xc, ok := v.(string)
		if !ok {
			continue
		}

		x, err := crypt(xc)
		if err != nil {
			return c, err
		}

		c[k] = x
	}

	return c, nil
}
