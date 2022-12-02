package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) InsertUser(details UserDetails) error {
	_, err := db.c.Exec("INSERT users(handle, name, auth) VALUES ('?', '?', ?)", details.Handle, details.Name, details.Auth)
	return err
}
