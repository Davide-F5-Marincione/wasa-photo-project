package database

func (db *appdbimpl) CheckBan(banisher string, banished string) bool {
	err := db.c.QueryRow("SELECT banisher, banished FROM bans WHERE banisher=? AND banished=?", banisher, banished).Err()
	return err == nil
}
