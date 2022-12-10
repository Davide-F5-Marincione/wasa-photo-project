package database

import "database/sql"

// MaxFollowers is the maximum number of returned ids on a request
const MaxFollowers int = 64

func (db *appdbimpl) GetFollowers(userhandle string, basehandle string) ([]UserAndDatetime, error) {
	var ids []UserAndDatetime = make([]UserAndDatetime, MaxFollowers)
	var res *sql.Rows
	var err error

	if basehandle == "" {
		res, err = db.c.Query(`
			SELECT follower, since
			FROM follows
			WHERE followed=?
			ORDER BY follower ASC
			LIMIT ?
			`, userhandle, MaxFollowers)
	} else {
		res, err = db.c.Query(`
			SELECT follower, since
			FROM follows
			WHERE followed=? and followed > ?
			ORDER BY follower ASC
			LIMIT ?
			`, userhandle, basehandle, MaxFollowers)
	}

	if err != nil {
		return nil, err
	}

	i := 0
	for res.Next() {
		err = res.Scan(&(ids[i].Handle), &(ids[i].RelevantDate)) // Since I can't do ids[i++]...
		i++                                                      // This warning is outrageous, i++ is ugly by itself!
		if err != nil {
			return nil, err
		}
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	return ids[:i], err
}
