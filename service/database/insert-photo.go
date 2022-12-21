package database

func (db *appdbimpl) InsertPhoto(author string, title string, file []byte) (int, error) {
	author_a, err := db.GetAuth(author)
	if err != nil {
		return -1, err
	}

	var id int
	err = db.c.QueryRow("INSERT INTO photos(author_auth, title, file) VALUES (?, ?, ?) RETURNING id", author_a, title, file).Scan(&id)
	return id, err
}
