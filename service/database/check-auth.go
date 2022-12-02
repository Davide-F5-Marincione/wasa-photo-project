package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) CheckAuthFree(auth int) bool {
	err := db.c.QueryRow("SELECT * FROM users WHERE auth=?", auth)
	return err != nil
}
