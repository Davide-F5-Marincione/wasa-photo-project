package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) GetUserDetails(username string) (UserDetails, error) {
	var details UserDetails
	err := db.c.QueryRow("SELECT name, auth FROM users WHERE name=?", username).Scan(&details.Name, &details.Auth)
	return details, err
}
