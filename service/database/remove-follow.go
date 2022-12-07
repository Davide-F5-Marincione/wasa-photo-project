package database

func (db *appdbimpl) RemoveFollow(follower string, followed string) error {
	_, err := db.c.Exec("DELETE FROM follows WHERE follower=? AND followed=?", follower, followed)
	return err
}
