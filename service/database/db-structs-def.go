package database

type UserDetails struct {
	Handle string
	Name   string
	Auth   int
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
	id         int
	author     string
	title      string
	uploadDate string
}

type Like struct {
	photoId int
	liker   string
	since   string
}

type Comment struct {
	photoId int
	id      int
	author  string
	content string
	since   string
}
