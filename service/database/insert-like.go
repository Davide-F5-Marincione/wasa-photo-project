package database

func (db *appdbimpl) InsertLike(liker string, photoid int) error {
	liker_a, err := db.GetAuth(liker)
	if err != nil {
		return err
	}

	_, err = db.c.Exec("INSERT INTO likes(liker_auth, photoId) VALUES (?, ?)", liker_a, photoid)
	return err
}
