package connman

import (
	"fmt"
)

// Spot hotspot general struct
type Spot struct {
	name    string
	pass    string
	exists  bool
	up      bool
	deleted bool
}

// NewSpot Spot struct const.
func NewSpot(name, pass string) *Spot {
	exists, err := spotExists(name)
	if err != nil {
		return nil
	}
	return &Spot{name: name, pass: pass, exists: exists, up: false, deleted: false}
}

//Add add new hotspot network
func (s *Spot) Add() error {
	if s.exists {
		return nil
	}
	err := addSpot(s.name, s.pass)
	if err != nil {
		return err
	}
	return nil
}

//Up hotspot up
func (s *Spot) Up() error {
	if s.deleted {
		return fmt.Errorf("spot:%v is deleted can not be up", s.name)
	}
	err := setSpot("up", s.name)
	if err != nil {
		return err
	}
	s.up = true
	return nil
}

//State spot state up, down, deleted
func (s *Spot) State() string {
	if s.deleted {
		return "deleted"
	}
	if s.up {
		return "up"
	}
	return "down"
}

//Down hotspot down
func (s *Spot) Down() error {
	if s.deleted {
		return fmt.Errorf("spot:%v is deleted can not be down", s.name)
	}
	err := setSpot("down", s.name)
	if err != nil {
		return err
	}
	s.up = false
	return nil
}

//Delete hotspot delete
func (s *Spot) Delete() error {
	if s.deleted {
		return nil
	}
	err := setSpot("delete", s.name)
	if err != nil {
		return err
	}
	s.deleted = true
	return nil
}

//Revive revive a deleted hotspot
func (s *Spot) Revive() error {
	if !s.deleted {
		return nil
	}
	err := s.Add()
	if err != nil {
		return err
	}
	s.up = false
	s.deleted = false
	return nil
}

//HotspotOn hotspot on wifi network off
func hotspotOn(spot *Spot) error {
	err := spot.Add()
	if err != nil {
		return err
	}
	err = spot.Up()
	if err != nil {
		return err
	}
	return nil
}

//HotspotOff hotspot off wifi network on
func hotspotOff(spot *Spot) error {
	err := spot.Add()
	if err != nil {
		return err
	}
	err = spot.Down()
	if err != nil {
		return err
	}
	err = connectAvailable()
	if err != nil {
		return err
	}
	return nil
}
