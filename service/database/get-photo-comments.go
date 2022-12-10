package database

import "database/sql"

// MaxComments represents the maximum number of comments one may get
const MaxComments int = 64

// CommentShow is just Comment without the PhotoID (since it is redundant in this case)
type CommentShow struct {
	ID      int    `json:"comment-id"`
	Author  string `json:"comment-author"`
	Content string `json:"comment-text"`
	Since   string `json:"comment-date"`
}

func (db *appdbimpl) GetPhotoComments(photoid int, commentlimit int) ([]CommentShow, error) {
	var ids []CommentShow = make([]CommentShow, MaxComments)
	var res *sql.Rows
	var err error

	if commentlimit < 1 {
		res, err = db.c.Query(`
			SELECT id, author, content, since
			FROM comments
			WHERE
				photoId = ?
			ORDER BY id DESC
			LIMIT ?
			`, photoid, MaxComments)
	} else {
		res, err = db.c.Query(`
			SELECT id, author, content, since
			FROM comments
			WHERE
				photoId = ? AND id < ?
			ORDER BY id DESC
			LIMIT ?
			`, photoid, commentlimit, MaxComments)
	}

	if err != nil {
		return nil, err
	}

	i := 0
	for res.Next() {
		err = res.Scan(&(ids[i].ID), &(ids[i].Author), &(ids[i].Content), &(ids[i].Since)) // Since I can't do ids[i++]...
		i++                                                                                // This warning is outrageous, i++ is ugly by itself!
		if err != nil {
			return nil, err
		}
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	return ids[:i], err
}
