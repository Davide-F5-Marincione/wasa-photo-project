package database

func (db *appdbimpl) CheckFollow(follower string, followed string) bool {
	var f string
	err := db.c.QueryRow("SELECT follower FROM follows WHERE follower=? AND followed=?", follower, followed).Scan(&f)
	return err == nil
}
