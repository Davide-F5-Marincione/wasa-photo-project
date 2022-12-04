package database

func (db *appdbimpl) InsertFollow(follower string, followed string) error {
	_, err := db.c.Exec("INSERT INTO follows(follower, followed) VALUES (?, ?)", follower, followed)
	return err
}
