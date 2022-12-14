package database

import (
	"database/sql"
	"errors"
)

// MaxSimilar is the maximum number of returned pairs on a request
const MaxSimilar int = 64

func (db *appdbimpl) FindSimilar(inputname string, basehandle string, basename string) ([]NameAndHandle, error) {
	var ids []NameAndHandle = make([]NameAndHandle, MaxSimilar)
	var res *sql.Rows
	var err error

	inputname = "%" + inputname + "%"

	if basehandle == "" && basename == "" {
		res, err = db.c.Query(`
			SELECT handle, name
			FROM users
			WHERE name LIKE ?
			ORDER BY handle ASC, name ASC
			LIMIT ?
			`, inputname, MaxSimilar)
	} else if basehandle != "" && basename != "" {
		res, err = db.c.Query(`
			SELECT handle, name
			FROM users
			WHERE name LIKE ? and name > ? and handle > ?
			ORDER BY handle ASC, name ASC
			LIMIT ?
			`, inputname, basename, basehandle, MaxSimilar)
	} else {
		return nil, errors.New("either basehandle and basename are both defined or neither is")
	}

	if err != nil {
		return nil, err
	}

	i := 0
	for res.Next() {
		err = res.Scan(&(ids[i].Handle), &(ids[i].Name)) // Since I can't do ids[i++]...
		i++                                              // This warning is outrageous, i++ is ugly by itself!
		if err != nil {
			return nil, err
		}
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	return ids[:i], err
}
