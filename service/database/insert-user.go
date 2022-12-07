package database

func (db *appdbimpl) InsertUser(details UserDetails) error {
	_, err := db.c.Exec("INSERT INTO users(handle, name, auth) VALUES (?, ?, ?)", details.Handle, details.Name, details.Auth)
	return err
}
