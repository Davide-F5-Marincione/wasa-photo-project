package database

func (db *appdbimpl) GetAuth(name string) (int, error) {
	var auth int
	err := db.c.QueryRow("SELECT auth FROM users WHERE name=?", name).Scan(&auth)
	return auth, err
}
