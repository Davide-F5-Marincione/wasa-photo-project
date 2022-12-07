package database

func (db *appdbimpl) RemovePhoto(id int) error {
	_, err := db.c.Exec("DELETE FROM photos WHERE id=?", id)
	return err
}
