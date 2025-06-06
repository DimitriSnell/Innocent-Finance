package main

import "github.com/google/uuid"

type Donator struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewDonator(n string) *Donator {
	result := &Donator{}
	result.Name = n
	result.ID = GenerateID()
	return result
}

func (d Donator) GetName() string {
	return d.Name
}

func GenerateID() string {
	return uuid.New().String()
}
