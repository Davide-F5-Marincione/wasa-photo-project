package database

func (db *appdbimpl) GetUserDetailsAuth(auth int) (UserDetails, error) {
	var details UserDetails
	err := db.c.QueryRow("SELECT name, auth FROM users WHERE auth=?", auth).Scan(&details.Name, &details.Auth)
	return details, err
}
