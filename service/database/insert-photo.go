package database

func (db *appdbimpl) InsertPhoto(author string, title string, file []byte) (int, error) {
	var id int
	err := db.c.QueryRow("INSERT INTO photos(author, title, file) VALUES (?, ?, ?) RETURNING id", author, title, file).Scan(&id)
	return id, err
}
