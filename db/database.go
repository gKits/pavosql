package db

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/gKits/PavoSQL/store"
)

type Database struct {
	Name string
	kv   store.Store
}

func (db *Database) Get(tbName string, r *Row) error {
	tb, err := db.GetTable(tbName)
	if err != nil {
		return fmt.Errorf("table '%s' not found: %v", tbName, err)
	}

	tb.checkRow(r)

	return nil
}

func (db *Database) Insert(tbName string) {}

func (db *Database) Update() {}

func (db *Database) Delete(tbName string) {}

func (db *Database) GetTable(tbName string) (table, error) {
	k := defaultTBDefTable.encodeKey([]byte(tbName))

	v, err := db.kv.Get(k)
	if err != nil {
		return table{}, err
	}

	tb := table{}
	if err := json.Unmarshal(v, &tb); err != nil {
		return table{}, err
	}

	return tb, nil
}

func (db *Database) CreateTable(tb table) error {
	v, err := json.Marshal(tb)
	if err != nil {
		return err
	}

	pref, err := db.nextPrefix()
	if err != nil {
		return err
	}

	tb.Prefix = pref

	k := defaultTBDefTable.encodeKey([]byte(tb.Name))

	if err := db.kv.Set(k, v); err != nil {
		return err
	}

	return nil
}

func (db *Database) DeleteTable(tbName string) (bool, error) {
	k := defaultTBDefTable.encodeKey([]byte(tbName))

	del, err := db.kv.Del(k)
	if err != nil {
		return false, err
	}

	return del, nil
}

func (db *Database) prefix() (uint32, error) {
	k := defaultMetaTable.encodeKey([]byte("prefix"))

	_, err := db.kv.Get(k)
	if err != nil {
		return 0, err
	}

	var pref uint32

	return pref, nil
}

func (db *Database) nextPrefix() (uint32, error) {
	k := defaultMetaTable.encodeKey([]byte("prefix"))

	v, err := db.kv.Get(k)
	if err != nil {
		return 0, err
	}

	var pref uint32

	v = binary.BigEndian.AppendUint32(v, pref+1)

	if err := db.kv.Set(k, v); err != nil {
		return 0, err
	}

	return pref + 1, nil
}
