package tack

import "fmt"

type Db struct {
	store map[string]interface{}
	count map[interface{}]int
}

func (db *Db) Set(name string, value interface{}) {
	db.store[name] = value
	db.count[value] += 1
}
