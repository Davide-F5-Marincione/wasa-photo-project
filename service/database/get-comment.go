package database

func (db *appdbimpl) GetComment(photoid int, id int) (Comment, error) {
	var comment Comment
	err := db.c.QueryRow("SELECT photoId, id, author, content, since FROM comments WHERE photoId=? AND id=?", photoid, id).Scan(&comment.PhotoID, &comment.ID, &comment.Author, &comment.Content, &comment.Since)
	return comment, err
}
