package db

var defaultMetaTable = table{
	Name:   "@meta",
	Prefix: 0,
	Cols:   []string{"name", "val"},
	Types:  []dbType{dbString, dbBytes},
	PKeys:  1,
}

var defaultTBDefTable = table{
	Name:   "@tbdef",
	Prefix: 1,
	Cols:   []string{"name", "def"},
	Types:  []dbType{dbString, dbBytes},
	PKeys:  1,
}

var defaultUsersTable = table{
	Name:   "@users",
	Prefix: 2,
	Cols:   []string{"name", "pass"},
	Types:  []dbType{dbString, dbBytes},
	PKeys:  1,
}
