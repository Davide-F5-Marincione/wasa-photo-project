package database

func (db *appdbimpl) InsertLike(liker string, photoid int) error {
	_, err := db.c.Exec("INSERT INTO likes(liker, photoId) VALUES (?, ?)", liker, photoid)
	return err
}
