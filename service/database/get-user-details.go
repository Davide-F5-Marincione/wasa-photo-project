package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetUserDetails(userhandle string) (UserDetails, error) {
	var details UserDetails
	err := db.c.QueryRow("SELECT handle, name, auth FROM users WHERE handle=?", userhandle).Scan(&details.Handle, &details.Name, &details.Auth)
	return details, err
}
