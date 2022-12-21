package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) CheckAuth(auth int) string {
	var a string
	err := db.c.QueryRow("SELECT name FROM users WHERE auth=?", auth).Scan(&a)
	if err != nil {
		return ""
	} else {
		return a
	}
}
