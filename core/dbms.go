package core

type DBMS struct {
	Dir       string
	Databases map[string]database
}
