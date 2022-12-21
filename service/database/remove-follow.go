package database

func (db *appdbimpl) RemoveFollow(follower string, followed string) error {
	follower_a, err := db.GetAuth(follower)
	if err != nil {
		return err
	}

	followed_a, err := db.GetAuth(followed)
	if err != nil {
		return err
	}

	_, err = db.c.Exec("DELETE FROM follows WHERE follower_auth=? AND followed_auth=?", follower_a, followed_a)
	return err
}
