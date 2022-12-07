package database

func (db *appdbimpl) RemoveBan(banisher string, banished string) error {
	_, err := db.c.Exec("DELETE FROM bans WHERE banisher=? AND banished=?", banisher, banished)
	return err
}
