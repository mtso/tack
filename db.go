package tack

import "errors"
import "fmt"
import "time"

var ErrEnd = errors.New("END EXCHANGE")
var ErrNoTransaction = errors.New("NO TRANSACTION")
var ErrNotFound = errors.New("NULL")
var ErrInvalidArgc = errors.New("INVALID ARGUMENT COUNT")
var ErrNoMemory = errors.New("NO MEMORY")

type command func(args ...string) (string, error)

type memdata [2]int

type dataset map[string]entry

type db struct {
	maxMemory        uint64
	totalMemory      uint64
	dataMemory       uint64
	blockMemory      uint64
	blockDataMemory  uint64
	maxMemorySamples int

	store    dataset
	count    map[string]uint
	block    []dataset
	commands map[string]command
}

func (db *db) GetCommands() map[string]command {
	if db.commands != nil {
		return db.commands
	}
	db.commands = map[string]command{
		"SET": func(args ...string) (string, error) {
			if len(args) < 2 {
				return "", ErrInvalidArgc
			} else {
				db.stash(args[0])
				return db.set(args[0], args[1])
			}
		},
		"UNSET": func(args ...string) (string, error) {
			if len(args) < 1 {
				return "", ErrInvalidArgc
			} else {
				db.stash(args[0])
				return db.unset(args[0])
			}
		},
		"GET":        db.get,
		"NUMEQUALTO": db.numEqualTo,
		"BEGIN":      db.begin,
		"ROLLBACK":   db.rollback,
		"COMMIT":     db.commit,
		"END":        end,
		"MEMUSE":     db.memUse,
		"MEMUSEDATA": db.memUseData,
	}
	return db.commands
}

func CreateDb() *db {
	return &db{
		store:       make(dataset),
		count:       make(map[string]uint),
		totalMemory: 8 + 8 + 8 + 4,
	}
}

func (db *db) get(args ...string) (string, error) {
	entry, ok := db.store[args[0]]
	if !ok {
		return "", ErrNotFound
	}
	entry.setHit(time.Now().UnixNano())
	return entry.value, nil
}

func (db *db) memUse(_ ...string) (string, error) {
	return fmt.Sprintf("%v", db.totalMemory), nil
}

func (db *db) memUseData(_ ...string) (string, error) {
	return fmt.Sprintf("%v", db.dataMemory), nil
}

func (db *db) numEqualTo(args ...string) (string, error) {
	count, _ := db.count[args[0]]
	return fmt.Sprintf("%v", count), nil
}

func (db *db) set(name, value string) (string, error) {
	db.unset(name)

	e := entry{
		time.Now().UnixNano(),
		value,
	}
	totalMem, dataMem := calcMem(name, e)

	// Add count's key and value mem size.
	if _, exists := db.count[value]; !exists {
		valueMem := len(value)
		totalMem += uint64(valueMem + 4)
	}
	db.addMem(totalMem, dataMem)

	db.store[name] = e
	db.count[value] += 1
	return "", nil
}

func (db *db) unset(name string) (string, error) {
	found, exists := db.store[name]
	if !exists {
		return "", nil
	}

	totalMem, dataMem := calcMem(name, found)

	if ct, _ := db.count[found.value]; ct > 1 {
		db.count[found.value] -= 1
		delete(db.store, name)
	} else {
		// Add key and value size of entry in count table.
		totalMem += uint64(len(found.value) + 4)
		delete(db.count, found.value)
		delete(db.store, name)
	}

	db.addMem(-totalMem, -dataMem)
	return "", nil
}

func (db *db) addMem(total, data uint64) {
	db.totalMemory += total
	db.dataMemory += data
}

func calcMem(key string, ent entry) (total uint64, data uint64) {
	keyMem := len(key)
	// Key and value memory size.
	data = uint64(ent.getValueMem() + keyMem)
	// Total memory: key and entry (time and string) memory size.
	total = uint64(ent.getMem() + keyMem)
	return
}

func (db *db) stash(key string) {
	if len(db.block) < 1 {
		return
	}
	tx := db.block[0]
	if _, saved := tx[key]; !saved {
		// save
		prev, exists := db.store[key]
		if exists {
			totalMem, dataMem := calcMem(key, prev)
			db.addMem(totalMem, dataMem)
		} else {
			keyMem := uint64(len(key))
			db.addMem(keyMem, keyMem)
		}
		tx[key] = prev
	}
}

func (db *db) begin(_ ...string) (string, error) {
	if db.maxMemory > 0 && db.totalMemory < db.maxMemory {
		return "", ErrNoMemory
	}
	db.block = append([]dataset{make(dataset)}, db.block...)
	return "", nil
}

func (db *db) rollback(_ ...string) (string, error) {
	if len(db.block) < 1 {
		return "", ErrNoTransaction
	}
	tx := db.block[0]
	for k, e := range tx {
		if e.value == "" {
			db.unset(k)
			keyMem := uint64(len(k))
			db.addMem(-keyMem, -keyMem)
		} else {
			db.set(k, e.value)
			totalMem, dataMem := calcMem(k, e)
			db.addMem(-totalMem, -dataMem)
		}
	}
	db.block = db.block[1:]
	return "", nil
}

func (db *db) commit(_ ...string) (string, error) {
	blocklen := len(db.block)
	if blocklen < 1 {
		return "", ErrNoTransaction
	}

	done := make(chan byte)

	for _, tx := range db.block {
		go func() {
			for k, e := range tx {
				if e.value == "" {
					keyMem := uint64(len(k))
					db.addMem(-keyMem, -keyMem)
				} else {
					totalMem, dataMem := calcMem(k, e)
					db.addMem(-totalMem, -dataMem)
				}
			}
			done <- 1
		}()
	}

	for i := 0; i < blocklen; i++ {
		<-done
	}

	db.block = nil

	return "", nil
}

func end(_ ...string) (string, error) {
	return "", ErrEnd
}
