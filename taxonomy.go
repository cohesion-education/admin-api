package main

import mgo "gopkg.in/mgo.v2"

type taxonomy struct {
	name     string
	id       int32
	parentID int32
	childIDs []int32
}

type taxonomyRepository struct {
	session *mgo.Session
}

func (repo *taxonomyRepository) list() ([]*taxonomy, error) {
	var list []*taxonomy

	return list, nil
}

func (t *taxonomy) Parent() *taxonomy {
	if t.parentID != -1 {

	}

	return nil
}

func (t *taxonomy) Children() []*taxonomy {
	return nil
}
