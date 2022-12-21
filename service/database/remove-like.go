package database

func (db *appdbimpl) RemoveLike(liker string, photoid int) error {
	liker_a, err := db.GetAuth(liker)
	if err != nil {
		return err
	}

	_, err = db.c.Exec("DELETE FROM likes WHERE liker_auth=? AND photoId=?", liker_a, photoid)
	return err
}
