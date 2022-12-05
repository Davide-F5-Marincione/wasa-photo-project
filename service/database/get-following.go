package database

// MaxIDs is the maximum number of returned ids on a request
const MaxFollowing int = 64

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetFollowing(userhandle string) ([]UserAndDatetime, error) {
	var ids []UserAndDatetime = make([]UserAndDatetime, MaxFollowing)

	res, err := db.c.Query(`
			SELECT followed, since
			FROM follows
			WHERE follower=?
			ORDER BY followed ASC
			LIMIT ?
			`, userhandle, MaxFollowing)

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

func (db *appdbimpl) GetFollowingLimit(userhandle string, basehandle string) ([]UserAndDatetime, error) {
	var ids []UserAndDatetime = make([]UserAndDatetime, MaxFollowing)

	res, err := db.c.Query(`
			SELECT followed, since
			FROM follows
			WHERE follower=? and follower > ?
			ORDER BY followed ASC
			LIMIT ?
			`, userhandle, basehandle, MaxFollowing)

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
