package database

func (db *appdbimpl) InsertBan(banisher string, banished string) error {
	_, err := db.c.Exec("INSERT INTO bans(banisher, banished) VALUES (?, ?)", banisher, banished)
	return err
}
