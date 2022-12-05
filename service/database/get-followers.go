package database

// MaxIDs is the maximum number of returned ids on a request
const MaxFollowers int = 64

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetFollowers(userhandle string) ([]UserAndDatetime, error) {
	var ids []UserAndDatetime = make([]UserAndDatetime, MaxFollowers)

	res, err := db.c.Query(`
			SELECT follower, since
			FROM follows
			WHERE followed=?
			ORDER BY follower ASC
			LIMIT ?
			`, userhandle, MaxFollowers)

	if err != nil {
		return nil, err
	}

	i := 0
	for res.Next() {
		res.Scan(&(ids[i].Handle), &(ids[i].RelevantDate)) // Since I can't do ids[i++]...
		i += 1                                             // This warning is outrageous, i++ is ugly by itself!
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	return ids, err
}

func (db *appdbimpl) GetFollowersLimit(userhandle string, basehandle string) ([]UserAndDatetime, error) {
	var ids []UserAndDatetime = make([]UserAndDatetime, MaxFollowers)

	res, err := db.c.Query(`
			SELECT follower, since
			FROM follows
			WHERE followed=? and followed > ?
			ORDER BY follower ASC
			LIMIT ?
			`, userhandle, basehandle, MaxFollowers)

	if err != nil {
		return nil, err
	}

	i := 0
	for res.Next() {
		res.Scan(&(ids[i].Handle), &(ids[i].RelevantDate)) // Since I can't do ids[i++]...
		i += 1                                             // This warning is outrageous, i++ is ugly by itself!
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	return ids, err
}
