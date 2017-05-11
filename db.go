package tack

import "errors"

var ErrEnd = errors.New("END EXCHANGE")
var ErrNoTransaction = errors.New("NO TRANSACTION")
var ErrNotFound = errors.New("NULL")

type command func(args ...string) (string, error)

type dataset map[string]entry

type db struct {
	maxMemory uint64
	currentMemory uint64
	maxMemorySamples int

	store dataset
	count map[entry]int
	block []dataset
	commands map[string]command
}

func (db *db) GetCommands() map[string]command {
	if db.commands == nil {
		db.commands = map[string]command{
			"GET":        db.get,
			// "SET":        db.stashSet,
			// "NUMEQUALTO": db.numEqualTo,
			// "UNSET":      db.stashUnset,
			// "BEGIN":      db.begin,
			// "ROLLBACK":   db.rollback,
			// "COMMIT":     db.commit,
			"END":        end,
		}
	}
	return db.commands
}

func CreateDb() *db {
	return &db{
		store: make(dataset),
		count: make(map[entry]int),
	}
}

// func MakeHandler() map[string]command {
// 	db := CreateDb()
// 	return map[string]command{
// 		"GET":        db.get,
// 		"SET":        db.stashSet,
// 		"NUMEQUALTO": db.numEqualTo,
// 		"UNSET":      db.stashUnset,
// 		"BEGIN":      db.begin,
// 		"ROLLBACK":   db.rollback,
// 		"COMMIT":     db.commit,
// 		"END":        end,
// 	}
// }

// func (db *db) stash(key interface{}) (_ interface{}) {
// 	if len(db.block) < 1 {
// 		return
// 	}
// 	if _, exists := db.block[0][key]; !exists {
// 		db.block[0][key] = db.store[key]
// 	}
// 	return
// }

func (db *db) get(args ...string) (string, error) {
	entry, ok := db.store[args[0]]
	if !ok {
		return "", ErrNotFound
	}
	return entry.value, nil
}

// func (db *db) numEqualTo(args ...interface{}) (count interface{}) {
// 	count, _ = db.count[args[0]]
// 	return
// }

// func (db *db) stashSet(args ...interface{}) (_ interface{}) {
// 	db.stash(args[0])
// 	db.set(args...)
// 	return
// }

// // func (db *db) set(name )

// func (db *db) set(args ...interface{}) (_ interface{}) {
// 	// db.unset(args[0])
// 	db.store[args[0]] = args[1]
// 	db.count[args[1]] += 1

// 	// if _, exists := db.count[args[1]]; !exists {
		
// 	// }
// 	// // count exists
// 	return
// }

// func (db *db) stashUnset(args ...interface{}) (_ interface{}) {
// 	db.stash(args[0])
// 	db.unset(args...)
// 	return
// }

// func (db *db) unset(args ...interface{}) (_ interface{}) {
// 	v := db.get(args[0])
// 	if db.count[v] > 1 {
// 		db.count[v] -= 1
// 		delete(db.store, args[0])
// 	} else {
// 		delete(db.count, v)
// 		delete(db.store, args[0])
// 	}
// 	return
// }

// func (db *db) begin(_ ...interface{}) (_ interface{}) {
// 	db.block = append([]dataset{make(dataset)}, db.block...)
// 	return
// }

// func (db *db) rollback(_ ...interface{}) (_ interface{}) {
// 	if db.block == nil {
// 		return ErrNoTransaction
// 	}
// 	tx := db.block[0]
// 	for k, v := range tx {
// 		if v == nil {
// 			db.unset(k)
// 		} else {
// 			db.set(k, v)
// 		}
// 	}
// 	db.block = db.block[1:]
// 	return
// }

// func (db *db) commit(_ ...interface{}) (_ interface{}) {
// 	if len(db.block) == 0 {
// 		return ErrNoTransaction
// 	}
// 	db.block = nil
// 	return
// }

func end(_ ...string) (string, error) {
	return "", ErrEnd
}
