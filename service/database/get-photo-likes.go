package database

import "database/sql"

// MaxLikers represents the maximum number of likers we return each time
const MaxLikers int = 64

func (db *appdbimpl) GetPhotoLikes(photoid int, basename string) ([]UserAndDatetime, error) {
	var ids []UserAndDatetime = make([]UserAndDatetime, MaxLikers)
	var res *sql.Rows
	var err error

	if basename == "" {
		res, err = db.c.Query(`
			SELECT liker, since
			FROM likes
			WHERE photoId = ?
			ORDER BY liker ASC
			LIMIT ?
			`, photoid, MaxLikers)
	} else {
		res, err = db.c.Query(`
			SELECT liker, since
			FROM likes
			WHERE photoId = ? AND liker > ?
			ORDER BY liker ASC
			LIMIT ?
			`, photoid, basename, MaxLikers)
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
