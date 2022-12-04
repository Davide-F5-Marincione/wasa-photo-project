package database

func (db *appdbimpl) InsertComment(author string, content string, photoid int) (int, error) {
	var id int
	err := db.c.QueryRow("INSERT INTO comments(author, content, photoId) VALUES (?, ?, ?) RETURNING id", author, content, photoid).Scan(&id)
	return id, err
}
