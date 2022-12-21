package database

import (
	"database/sql"
)

// MaxSimilar is the maximum number of returned pairs on a request
const MaxSimilar int = 64

func (db *appdbimpl) FindSimilar(inputname string, basename string) ([]string, error) {
	var ids []string = make([]string, MaxSimilar)
	var res *sql.Rows
	var err error

	inputname = "%" + inputname + "%"

	if basename == "" {
		res, err = db.c.Query(`
			SELECT name
			FROM users
			WHERE name LIKE ?
			ORDER BY name ASC
			LIMIT ?
			`, inputname, MaxSimilar)
	} else {
		res, err = db.c.Query(`
			SELECT name
			FROM users
			WHERE name LIKE ? and name > ?
			ORDER BY name ASC
			LIMIT ?
			`, inputname, basename, MaxSimilar)
	}

	if err != nil {
		return nil, err
	}

	i := 0
	for res.Next() {
		err = res.Scan(&ids[i]) // Since I can't do ids[i++]...
		i++                     // This warning is outrageous, i++ is ugly by itself!
		if err != nil {
			return nil, err
		}
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	return ids[:i], err
}
