package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) InsertBan(banisher string, banished string) error {
	_, err := db.c.Exec("INSERT INTO bans(banisher, banished) VALUES (?, ?)", banisher, banished)
	return err
}
