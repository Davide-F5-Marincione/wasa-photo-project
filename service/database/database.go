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
	CheckAuth(auth int) string
	GetUserDetails(username string) (UserDetails, error)
	GetUserDetailsAuth(auth int) (UserDetails, error)
	UpdateUsername(currname string, newname string) error
	FindSimilar(inputname string, basename string) ([]string, error)

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
	GetStream(username string, toplimit int) ([]int, error)
	GetFollowers(username string, basename string) ([]UserAndDatetime, error)
	GetFollowing(username string, basename string) ([]UserAndDatetime, error)
	GetPhotosProfile(username string, toplimit int) ([]int, error)

	GetPhotoComments(photoid int, commentlimit int) ([]CommentShow, error)
	GetPhotoLikes(photoid int, basename string) ([]UserAndDatetime, error)

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
			name TEXT PRIMARY KEY,
			auth INTEGER NOT NULL UNIQUE,
			registerDate TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
			lastLogin TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL
		);`)
	if err != nil {
		return nil, err
	}

	err = addTable(db, "follows",
		`CREATE TABLE follows (
			follower TEXT,
			followed TEXT,
			follower_auth INTEGER NOT NULL,
			followed_auth INTEGER NOT NULL,
			since TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
			FOREIGN KEY(follower_auth) REFERENCES users(auth),
			FOREIGN KEY(followed_auth) REFERENCES users(auth),
			PRIMARY KEY (follower_auth, followed_auth)
		);`)
	if err != nil {
		return nil, err
	}

	err = addTable(db, "bans",
		`CREATE TABLE bans (
			banisher TEXT,
			banished TEXT,
			banisher_auth INTEGER NOT NULL,
			banished_auth INTEGER NOT NULL,
			since TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
			FOREIGN KEY(banisher_auth) REFERENCES users(auth),
			FOREIGN KEY(banished_auth) REFERENCES users(auth),
			PRIMARY KEY (banisher_auth, banished_auth)
		);`)
	if err != nil {
		return nil, err
	}

	err = addTable(db, "photos",
		`CREATE TABLE photos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			author TEXT,
			author_auth INTEGER NOT NULL,
			title TEXT NOT NULL,
			uploadDate TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
			file BLOB NOT NULL,
			commentsCounter INTEGER DEFAULT 1 NOT NULL,
			FOREIGN KEY(author_auth) REFERENCES users(auth)
		);`)
	if err != nil {
		return nil, err
	}

	err = addTable(db, "likes",
		`CREATE TABLE likes (
			photoId INTEGER NOT NULL,
			liker TEXT,
			liker_auth INTEGER NOT NULL,
			since TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
			FOREIGN KEY(liker_auth) REFERENCES users(auth),
			FOREIGN KEY(photoId) REFERENCES photos(id),
			PRIMARY KEY (photoId, liker_auth)
		);`)
	if err != nil {
		return nil, err
	}

	err = addTable(db, "comments",
		`CREATE TABLE comments (
			photoId INTEGER NOT NULL,
			id INTEGER DEFAULT 0 NOT NULL,
			author TEXT,
			author_auth INTEGER NOT NULL,
			content TEXT NOT NULL,
			since TEXT DEFAULT CURRENT_TIMESTAMP NOT NULL,
			FOREIGN KEY(author_auth) REFERENCES users(auth),
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

	err = addTrigger(db, "followNameAdd", // To delete comments and likes
		`CREATE TRIGGER followNameAdd
			AFTER INSERT ON follows
		BEGIN
			UPDATE follows
				SET follower = (SELECT name FROM users WHERE auth = NEW.follower_auth),
					followed = (SELECT name FROM users WHERE auth = NEW.followed_auth)
			WHERE ROWID = new.ROWID;
		END;`)

	if err != nil {
		return nil, err
	}

	err = addTrigger(db, "banNameAdd", // To delete comments and likes
		`CREATE TRIGGER banNameAdd
			AFTER INSERT ON bans
		BEGIN
			UPDATE bans
				SET banisher = (SELECT name FROM users WHERE auth = NEW.banisher_auth),
				banished = (SELECT name FROM users WHERE auth = NEW.banished_auth)
			WHERE ROWID = new.ROWID;
		END;`)

	if err != nil {
		return nil, err
	}

	err = addTrigger(db, "photoAuthorAdd", // To delete comments and likes
		`CREATE TRIGGER photoAuthorAdd
			AFTER INSERT ON photos
		BEGIN
			UPDATE photos
				SET author = (SELECT name FROM users WHERE auth = NEW.author_auth)
			WHERE ROWID = new.ROWID;
		END;`)

	if err != nil {
		return nil, err
	}

	err = addTrigger(db, "likeNameAdd", // To delete comments and likes
		`CREATE TRIGGER likeNameAdd
			AFTER INSERT ON likes
		BEGIN
			UPDATE likes
				SET liker = (SELECT name FROM users WHERE auth = NEW.liker_auth)
			WHERE ROWID = new.ROWID;
		END;`)

	if err != nil {
		return nil, err
	}

	err = addTrigger(db, "commentAuthorAdd", // To delete comments and likes
		`CREATE TRIGGER commentAuthorAdd
			AFTER INSERT ON comments
		BEGIN
			UPDATE comments
				SET author = (SELECT name FROM users WHERE auth = NEW.author_auth)
			WHERE ROWID = new.ROWID;
		END;`)

	if err != nil {
		return nil, err
	}

	err = addTrigger(db, "nameChangeCascade", // To delete comments and likes
		`CREATE TRIGGER nameChangeCascade
			AFTER UPDATE OF name ON users
		BEGIN
			UPDATE follows
				SET follower = NEW.name
				WHERE follower_auth = NEW.auth;
			UPDATE follows
				SET followed = NEW.name
				WHERE followed_auth = NEW.auth;

			UPDATE bans
				SET banisher = NEW.name
				WHERE banisher_auth = NEW.auth;
			UPDATE bans
				SET banished = NEW.name
				WHERE banished_auth = NEW.auth;

			UPDATE photos
				SET author = NEW.name
				WHERE author_auth = NEW.auth;

			UPDATE likes
				SET liker = NEW.name
				WHERE liker_auth = NEW.auth;
			
			UPDATE comments
				SET author = NEW.name
				WHERE author_auth = NEW.auth;
		END;`)

	if err != nil {
		return nil, err
	}

	// In case we may want to also delete users, still haven't designed option to do so
	err = addTrigger(db, "userDelCascade", // First delete comments and likes, then photos, then follows and bans
		`CREATE TRIGGER userDelCascade
			BEFORE DELETE ON users
		BEGIN
			DELETE FROM comments WHERE author = OLD.name;
			DELETE FROM likes WHERE liker = OLD.name;
			DELETE FROM photos WHERE author = OLD.name;
			DELETE FROM follows WHERE follower = OLD.name OR followed = OLD.name;
			DELETE FROM bans WHERE banisher = OLD.name OR banished = OLD.name;
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
