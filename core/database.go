package core

import (
	"github.com/gKits/PavoSQL/store"
)

type database struct {
	path      string
	kv        store.Store
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
