package database

func (db *appdbimpl) CheckFollow(follower string, followed string) bool {
	err := db.c.QueryRow("SELECT follower, followed FROM follows WHERE follower=? AND followed=?", follower, followed).Err()
	return err == nil
}
