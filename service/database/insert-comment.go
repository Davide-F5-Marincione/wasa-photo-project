package database

func (db *appdbimpl) InsertComment(author string, content string, photoid int) (int, error) {
	// We have to do this strange stuff since the true comment's id is put after the insertion
	var rowid int
	err := db.c.QueryRow("INSERT INTO comments(author, content, photoId) VALUES (?, ?, ?) RETURNING ROWID", author, content, photoid).Scan(&rowid)
	if err != nil {
		return -1, err
	}

	var id int
	err = db.c.QueryRow("SELECT id FROM comments WHERE ROWID = ?", rowid).Scan(&id)
	return id, err
}
