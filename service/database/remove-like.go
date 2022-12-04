package database

func (db *appdbimpl) RemoveLike(liker string, photoid int) error {
	_, err := db.c.Exec("DELETE FROM likes WHERE liker=? AND photoId=?", liker, photoid)
	return err
}
