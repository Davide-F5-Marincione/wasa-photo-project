package database

func (db *appdbimpl) RemoveComment(photoid int, id int) error {
	_, err := db.c.Exec("DELETE FROM comments WHERE photoId=? AND id=?", photoid, id)
	return err
}
