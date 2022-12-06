package database

func (db *appdbimpl) InsertComment(author string, content string, photoid int) (int, error) {
	var rowid int
	err := db.c.QueryRow("INSERT INTO comments(author, content, photoId) VALUES (?, ?, ?) RETURNING ROWID", author, content, photoid).Scan(&rowid)
	if err != nil {
		return -1, err
	}
	var id int
	err = db.c.QueryRow("SELECT id FROM comments WHERE ROWID = ?", rowid).Scan(&id)
	return id, err
}
