package tack

// import "fmt"

type db struct {
	store  map[interface{}]interface{}
	count  map[interface{}]int
	isUndo bool
}

type Command func(args ...interface{}) interface{}

func CreateDb() *db {
	return &db{
		store: make(map[interface{}]interface{}),
		count: make(map[interface{}]int),
	}
}

func MakeHandler() map[string]Command {
	db := CreateDb()
	return map[string]Command{
		"GET":        db.Get,
		"SET":        db.Set,
		"NUMEQUALTO": db.NumEqualTo,
	}
}

func (db *db) Get(args ...interface{}) (value interface{}) {
	value, ok := db.store[args[0]]
	if !ok {
		return "NULL"
	}
	return
}

func (db *db) NumEqualTo(args ...interface{}) (count interface{}) {
	count, _ = db.count[args[0]]
	return
}

func (db *db) Set(args ...interface{}) interface{} {
	db.store[args[0]] = args[1]
	db.count[args[1]] += 1
	return nil
}
