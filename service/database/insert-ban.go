package database

func (db *appdbimpl) InsertBan(banisher string, banished string) error {
	banisher_a, err := db.GetAuth(banisher)
	if err != nil {
		return err
	}

	banished_a, err := db.GetAuth(banished)
	if err != nil {
		return err
	}

	_, err = db.c.Exec("INSERT INTO bans(banisher_auth, banished_auth) VALUES (?, ?)", banisher_a, banished_a)
	return err
}
