package connman

import (
	"fmt"
)

// Spot hotspot general struct
type Spot struct {
	name    string
	pass    string
	up      bool
	deleted bool
}

var spotlightPath string = "/home/eco/go/src/connman/scripts/./spotlight.sh"

// NewSpot Spot struct const.
func NewSpot(name, pass string) *Spot {
	return &Spot{name, pass, false, false}
}

//Add add new hotspot network
func (s *Spot) Add() error {
	cmd := fmt.Sprintf(`sudo "%v" add "%v" "%v"`, spotlightPath, s.name, s.pass)
	out, err := exe(cmd, "add spot")
	if err != nil {
		return fmt.Errorf("error:%v, out:%v", err.Error(), out)
	}
	return nil
}

//Up hotspot up
func (s *Spot) Up() error {
	if s.deleted {
		return fmt.Errorf("spot:%v is deleted can not be up", s.name)
	}
	cmd := fmt.Sprintf(`sudo "%v" up "%v"`, spotlightPath, s.name)
	out, err := exe(cmd, "conn up")
	if err != nil {
		return fmt.Errorf("error:%v, out:%v", err.Error(), out)
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
	cmd := fmt.Sprintf(`sudo "%v" down "%v"`, spotlightPath, s.name)
	out, err := exe(cmd, "conn up")
	if err != nil {
		return fmt.Errorf("error:%v, out:%v", err.Error(), out)
	}
	s.up = false
	return nil
}

//Delete hotspot delete
func (s *Spot) Delete() error {
	if s.deleted {
		return nil
	}
	cmd := fmt.Sprintf(`sudo "%v" delete "%v"`, spotlightPath, s.name)
	out, err := exe(cmd, "conn delete")
	if err != nil {
		return fmt.Errorf("error:%v, out:%v", err.Error(), out)
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
