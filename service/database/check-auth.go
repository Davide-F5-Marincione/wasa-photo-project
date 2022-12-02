package database

import "strconv"

// GetName is an example that shows you how to query data
func (db *appdbimpl) CheckAuthFree(auth int64) bool {
	err := db.c.QueryRow("SELECT * FROM users WHERE users.auth=" + strconv.Itoa(int(auth)))
	return err != nil
}
