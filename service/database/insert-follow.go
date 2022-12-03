package database

// GetName is an example that shows you how to query data
func (db *appdbimpl) InsertFollow(follower string, followed string) error {
	_, err := db.c.Exec("INSERT INTO follows(follower, followed) VALUES (?, ?)", follower, followed)
	return err
}
