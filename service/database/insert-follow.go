package database

func (db *appdbimpl) InsertFollow(follower string, followed string) error {
	follower_a, err := db.GetAuth(follower)
	if err != nil {
		return err
	}

	followed_a, err := db.GetAuth(followed)
	if err != nil {
		return err
	}

	_, err = db.c.Exec("INSERT INTO follows(follower_auth, followed_auth) VALUES (?, ?)", follower_a, followed_a)
	return err
}
