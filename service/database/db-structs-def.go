package database

type UserDetails struct {
	Handle string
	Name   string
	Auth   int64
}

type Follow struct {
	follower string
	followed string
	since    string
}

type Ban struct {
	banisher string
	banished string
	since    string
}

type PhotoDetails struct {
	id         int64
	author     string
	title      string
	uploadDate string
}

type Like struct {
	photoId int64
	liker   string
	since   string
}

type Comment struct {
	photoId int64
	id      int64
	author  string
	content string
	since   string
}
