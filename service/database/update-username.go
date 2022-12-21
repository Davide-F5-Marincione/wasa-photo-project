package database

func (db *appdbimpl) UpdateUsername(currname string, newname string) error {
	_, err := db.c.Exec("UPDATE users SET name = ? WHERE name = ?", newname, currname)
	return err
}
