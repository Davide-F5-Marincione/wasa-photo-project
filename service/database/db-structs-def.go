package database

// UserDetails contains all the useful ready-at-hand
// information of a user
type UserDetails struct {
	Handle string
	Name   string
	Auth   int
}

// UserAndDatetime contains a handle and a datetime,
// this is useful for general lists of users... and related timestamps
type UserAndDatetime struct {
	Handle       string `json:"handle"`
	RelevantDate string `json:"relevantdate"`
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
	ID         int
	Author     string
	Title      string
	UploadDate string
}

// Like is needed to handle the db schema 'likes'
type Like struct {
	photoID int
	liker   string
	since   string
}

// Comment is needed to handle the db schema 'comments'
type Comment struct {
	PhotoID int
	ID      int
	Author  string
	Content string
	Since   string
}
