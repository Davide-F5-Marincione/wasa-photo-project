package database

func (db *appdbimpl) GetBlobPhoto(id int) ([]byte, error) {
	var blob []byte
	err := db.c.QueryRow("SELECT file FROM photos WHERE id=?", id).Scan(&blob)
	return blob, err
}
