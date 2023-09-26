package core

import (
	store "github.com/gKits/PavoSQL/backend"
)

type database struct {
	path      string
	kv        store.KVStore
	tables    map[string]*table
	metaTable table
	defTable  table
}

func (db *database) get() {}

func (db *database) insert() {}

func (db *database) update() {}

func (db *database) delete() {}

func (db *database) createTable(t table) error {
	return nil
}

func (db *database) deleteTable() {}

func (db *database) alterTable() {}
