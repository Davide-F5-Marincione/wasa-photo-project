package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) UpdateUsername(handle string, newname string) error {
	_, err := db.c.Exec("UPDATE users SET name = ? WHERE handle = ?", newname, handle)
	return err
}
