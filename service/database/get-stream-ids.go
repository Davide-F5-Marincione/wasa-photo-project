package database

import "database/sql"

// MaxIDs is the maximum number of returned ids on a request
const MaxIDs int = 64

// GetName is an example that shows you how to query data

func (db *appdbimpl) GetStream(userhandle string, toplimit int) ([]int, error) {
	var ids []int = make([]int, MaxIDs)
	var res *sql.Rows
	var err error

	if toplimit < 1 {
		res, err = db.c.Query(`
			SELECT id
			FROM photos
			WHERE
				author IN (
						SELECT followed
						FROM follows
						WHERE follower = ?
					EXCEPT
						SELECT banisher
						FROM bans
						WHERE banished = ?
					)
				OR author = ?
			ORDER BY id DESC
			LIMIT ?
			`, userhandle, userhandle, userhandle, MaxIDs)
	} else {
		res, err = db.c.Query(`
			SELECT id
			FROM photos
			WHERE
				(author IN (
						SELECT followed
						FROM follows
						WHERE follower = ?
					EXCEPT
						SELECT banisher
						FROM bans
						WHERE banished = ?
					)
				OR author = ?) AND id < ?
			ORDER BY id DESC
			LIMIT ?
			`, userhandle, userhandle, userhandle, toplimit, MaxIDs)
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
