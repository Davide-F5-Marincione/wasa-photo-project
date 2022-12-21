package database

func (db *appdbimpl) InsertUser(details UserDetails) error {
	_, err := db.c.Exec("INSERT INTO users(name, auth) VALUES (?, ?)", details.Name, details.Auth)
	return err
}
