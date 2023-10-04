package dbms

type DBMS struct {
	Dir string
}

func (dbms *DBMS) CreateDatabase() error {
	return nil
}

func (dbms *DBMS) DeleteDatabase() error {
	return nil
}
