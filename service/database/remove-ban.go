package database

func (db *appdbimpl) RemoveBan(banisher string, banished string) error {
	banisher_a, err := db.GetAuth(banisher)
	if err != nil {
		return err
	}

	banished_a, err := db.GetAuth(banished)
	if err != nil {
		return err
	}

	_, err = db.c.Exec("DELETE FROM bans WHERE banisher_auth=? AND banished_auth=?", banisher_a, banished_a)
	return err
}
