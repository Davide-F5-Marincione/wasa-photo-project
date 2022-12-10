/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	InsertUser(details UserDetails) error
	CheckAuthFree(auth int) bool
	GetUserDetails(userhandle string) (UserDetails, error)
	GetUserDetailsAuth(auth int) (UserDetails, error)
	UpdateUsername(handle string, newname string) error

	InsertBan(banisher string, banished string) error
	CheckBan(banisher string, banished string) bool
	RemoveBan(banisher string, banished string) error

	InsertFollow(follower string, followed string) error
	CheckFollow(follower string, followed string) bool
	RemoveFollow(follower string, followed string) error

	InsertPhoto(author string, title string, file []byte) (int, error)
	GetPhotoDetails(id int) (PhotoDetails, error)
	GetBlobPhoto(id int) ([]byte, error)
	RemovePhoto(id int) error

	CheckLike(liker string, photoid int) bool
	InsertLike(liker string, photoid int) error
	RemoveLike(liker string, photoid int) error

	InsertComment(author string, content string, photoid int) (int, error)
	GetComment(photoid int, id int) (Comment, error)
	RemoveComment(photoid int, id int) error

	// Batches gets
	GetStream(userhandle string, toplimit int) ([]int, error)
	GetFollowers(userhandle string, basehandle string) ([]UserAndDatetime, error)
	GetFollowing(userhandle string, basehandle string) ([]UserAndDatetime, error)
	GetPhotosProfile(userhandle string, toplimit int) ([]int, error)

	GetPhotoComments(photoid int, commentlimit int) ([]CommentShow, error)
	GetPhotoLikes(photoid int, basehandle string) ([]UserAndDatetime, error)

	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

func addTable(db *sql.DB, tablename string, tabledef string) error {
	// Check if table exists. If not, the database is empty, and we need to create the structure
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='" + tablename + "';").Scan(&tablename)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(tabledef)
		if err != nil {
			return fmt.Errorf("error creating database table "+tablename+": %w", err)
		}
	}
	return nil
}

func addTrigger(db *sql.DB, triggername string, triggerdef string) error {
	// Check if trigger exists. If not, create it!
	err := db.QueryRow("SELECT name FROM sqlite_master WHERE type='trigger' AND name='" + triggername + "';").Scan(&triggername)
	if errors.Is(err, sql.ErrNoRows) {
		_, err = db.Exec(triggerdef)
		if err != nil {
			return fmt.Errorf("error creating database trigger "+triggername+": %w", err)
		}
	}
	return nil
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	/// TABLES
	err := addTable(db, "users",
		`CREATE TABLE users (
			handle TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			auth INTEGER NOT NULL UNIQUE,
			registerDate TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
			lastLogin TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL
		);`)
	if err != nil {
		return nil, err
	}

	err = addTable(db, "follows",
		`CREATE TABLE follows (
			follower TEXT NOT NULL,
			followed TEXT NOT NULL,
			since TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
			FOREIGN KEY(follower) REFERENCES users(handle),
			FOREIGN KEY(followed) REFERENCES users(handle),
			PRIMARY KEY (follower, followed)
		);`)
	if err != nil {
		return nil, err
	}

	err = addTable(db, "bans",
		`CREATE TABLE bans (
			banisher TEXT NOT NULL,
			banished TEXT NOT NULL,
			since TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
			FOREIGN KEY(banisher) REFERENCES users(handle),
			FOREIGN KEY(banished) REFERENCES users(handle),
			PRIMARY KEY (banisher, banished)
		);`)
	if err != nil {
		return nil, err
	}

	err = addTable(db, "photos",
		`CREATE TABLE photos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author TEXT NOT NULL,
			title TEXT NOT NULL,
			uploadDate TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
			file BLOB NOT NULL,
			commentsCounter INTEGER DEFAULT 1 NOT NULL,
			FOREIGN KEY(author) REFERENCES users(handle)
		);`)
	if err != nil {
		return nil, err
	}

	err = addTable(db, "likes",
		`CREATE TABLE likes (
			photoId INTEGER NOT NULL,
			liker TEXT NOT NULL,
			since TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
			FOREIGN KEY(liker) REFERENCES users(handle),
			FOREIGN KEY(photoId) REFERENCES photos(id),
			PRIMARY KEY (photoId, liker)
		);`)
	if err != nil {
		return nil, err
	}

	err = addTable(db, "comments",
		`CREATE TABLE comments (
			photoId INTEGER NOT NULL,
			id INTEGER DEFAULT 0 NOT NULL,
			author TEXT NOT NULL,
			content TEXT NOT NULL,
			since TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
			FOREIGN KEY(author) REFERENCES users(handle),
			FOREIGN KEY(photoId) REFERENCES photos(id),
			PRIMARY KEY (photoId, id)
		);`)
	if err != nil {
		return nil, err
	}

	/// TRIGGERS
	err = addTrigger(db, "commentsIncr", // To increment the comment id and add a date
		`CREATE TRIGGER commentsIncr
			AFTER INSERT ON comments
		BEGIN
			UPDATE comments
				SET id = (SELECT photos.commentsCounter
					FROM photos
					WHERE photos.id = NEW.photoId)
				WHERE ROWID = new.ROWID;
			UPDATE photos
				SET commentsCounter = photos.commentsCounter + 1
				WHERE id = NEW.photoId;
		END;`)

	if err != nil {
		return nil, err
	}

	err = addTrigger(db, "photoDelCascade", // To delete comments and likes
		`CREATE TRIGGER photoDelCascade
			BEFORE DELETE ON photos
		BEGIN
			DELETE FROM comments WHERE photoId = OLD.id;
			DELETE FROM likes WHERE photoId = OLD.id;
		END;`)

	if err != nil {
		return nil, err
	}

	// In case we may want to also delete users, still haven't designed option to do so
	err = addTrigger(db, "userDelCascade", // First delete comments and likes, then photos, then follows and bans
		`CREATE TRIGGER userDelCascade
			BEFORE DELETE ON users
		BEGIN
			DELETE FROM comments WHERE author = OLD.handle;
			DELETE FROM likes WHERE liker = OLD.handle;
			DELETE FROM photos WHERE author = OLD.handle;
			DELETE FROM follows WHERE follower = OLD.handle OR followed = OLD.handle;
			DELETE FROM bans WHERE banisher = OLD.handle OR banished = OLD.handle;
		END;`)

	if err != nil {
		return nil, err
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
