package database

func (db *appdbimpl) GetPhotoDetails(id int) (PhotoDetails, error) {
	var details PhotoDetails
	err := db.c.QueryRow("SELECT id, author, title, uploadDate FROM photos WHERE id=?", id).Scan(&details.ID, &details.Author, &details.Title, &details.UploadDate)
	return details, err
}
