package main

import (
	"errors"
)

type Relationship string

type Family struct {
	Members map[Relationship]Person
}

type Person struct {
	FirstName string
	LastName  string
	Age       int
}

const (
	Father      = Relationship("father")
	Mother      = Relationship("mother")
	Child       = Relationship("child")
	GrandMother = Relationship("grandMother")
	GrandFather = Relationship("grandFather")
)

var (
	ErrRelationshipExists = errors.New("Rel already exists")
)

func (f *Family) addMember(r Relationship, p Person) error {
	if f.Members == nil {
		f.Members = make(map[Relationship]Person)
	}

	if _, ok := f.Members[r]; ok {
		return ErrRelationshipExists
	}
	f.Members[r] = p
	return nil

}

func Sum(vals ...int) int {
	var total int
	for _, v := range vals {
		total += v
	}
	return total
}


