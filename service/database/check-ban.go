package database

func (db *appdbimpl) CheckBan(banisher string, banished string) bool {
	var b string
	err := db.c.QueryRow("SELECT banisher FROM bans WHERE banisher=? AND banished=?", banisher, banished).Scan(&b)
	return err == nil
}
