package database

func (db *appdbimpl) GetUserDetailsAuth(auth int) (UserDetails, error) {
	var details UserDetails
	err := db.c.QueryRow("SELECT handle, name, auth FROM users WHERE auth=?", auth).Scan(&details.Handle, &details.Name, &details.Auth)
	return details, err
}
