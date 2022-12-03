package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) CheckAuthFree(auth int) bool {
	var a int
	err := db.c.QueryRow("SELECT auth FROM users WHERE auth=?", auth).Scan(&a)
	return err != nil
}
