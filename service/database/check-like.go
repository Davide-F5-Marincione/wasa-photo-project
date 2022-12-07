package database

func (db *appdbimpl) CheckLike(liker string, photoid int) bool {
	var l string
	err := db.c.QueryRow("SELECT liker FROM likes WHERE liker=? AND photoId=?", liker, photoid).Scan(&l)
	return err == nil
}
