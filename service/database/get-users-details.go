package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetUserDetails(userhandle string) (UserDetails, error) {
	var details UserDetails
	err := db.c.QueryRow("SELECT users.handle, users.name, users.auth FROM users WHERE users.handle='" + userhandle + "'").Scan(&details)
	return details, err
}
