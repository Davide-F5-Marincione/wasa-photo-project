package database

// MaxIDs is the maximum number of returned ids on a request
const MaxPhotos int = 64

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetPhotosProfile(userhandle string) ([]int, error) {
	var ids []int = make([]int, MaxPhotos)

	res, err := db.c.Query(`
			SELECT id
			FROM photos
			WHERE
				author = ?
			ORDER BY id DESC
			LIMIT ?
			`, userhandle, MaxPhotos)

	if err != nil {
		return nil, err
	}

	i := 0
	for res.Next() {
		res.Scan(&ids[i]) // Since I can't do ids[i++]...
		i += 1            // This warning is outrageous, i++ is ugly by itself!
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	return ids, err
}

func (db *appdbimpl) GetPhotosProfileLimit(userhandle string, toplimit int) ([]int, error) {
	var ids []int = make([]int, MaxPhotos)

	res, err := db.c.Query(`
			SELECT id
			FROM photos
			WHERE
				author = ? AND id < ?
			ORDER BY id DESC
			LIMIT ?
			`, userhandle, toplimit, MaxPhotos)

	if err != nil {
		return nil, err
	}

	i := 0
	for res.Next() {
		res.Scan(&ids[i]) // Since I can't do ids[i++]...
		i += 1            // This warning is outrageous, i++ is ugly by itself!
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	return ids, err
}
