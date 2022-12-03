package database

// UserDetails contains all the useful ready-at-hand
// information of a user
type UserDetails struct {
	Handle string
	Name   string
	Auth   int
}

// Follow is needed to handle the db schema 'follows'
type Follow struct {
	follower string
	followed string
	since    string
}

// Ban is needed to handle the db schema 'bans'
type Ban struct {
	banisher string
	banished string
	since    string
}

// PhotoDetails contains all the useful ready-at-hand
// information of a photo
type PhotoDetails struct {
	id         int
	author     string
	title      string
	uploadDate string
}

// Like is needed to handle the db schema 'likes'
type Like struct {
	photoID int
	liker   string
	since   string
}

// Comment is needed to handle the db schema 'comments'
type Comment struct {
	photoID int
	id      int
	author  string
	content string
	since   string
}
