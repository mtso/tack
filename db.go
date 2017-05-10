package tack

// import "fmt"

type db struct {
	store map[string]interface{}
	count map[interface{}]int
	isUndo bool
}

func CreateDb() *db {
	return &db{
		store: make(map[string]interface{}),
		count: make(map[interface{}]int),
	}
}

func (db *db) Get(name string) (value interface{}) {
	value, ok := db.store[name]
	if !ok {
		return "NULL"
	}
	return
}

func (db *db) NumEqualTo(value interface{}) (count int) {
	count, _ = db.count[value]
	return
}

func (db *db) Set(name string, value interface{}) {
	db.store[name] = value
	db.count[value] += 1
}
