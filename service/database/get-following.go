package database

import "database/sql"

// MaxFollowing is the maximum number of returned ids on a request
const MaxFollowing int = 64

// GetName is an example that shows you how to query data

func (db *appdbimpl) GetFollowing(username string, basename string) ([]UserAndDatetime, error) {
	var ids []UserAndDatetime = make([]UserAndDatetime, MaxFollowing)
	var res *sql.Rows
	var err error

	if basename == "" {
		res, err = db.c.Query(`
			SELECT followed, since
			FROM follows
			WHERE follower=?
			ORDER BY followed ASC
			LIMIT ?
			`, username, MaxFollowing)
	} else {
		res, err = db.c.Query(`
			SELECT followed, since
			FROM follows
			WHERE follower=? and follower > ?
			ORDER BY followed ASC
			LIMIT ?
			`, username, basename, MaxFollowing)
	}

	if err != nil {
		return nil, err
	}

	i := 0
	for res.Next() {
		err = res.Scan(&(ids[i].Name), &(ids[i].RelevantDate)) // Since I can't do ids[i++]...
		i++                                                    // This warning is outrageous, i++ is ugly by itself!
		if err != nil {
			return nil, err
		}
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	return ids[:i], err
}
