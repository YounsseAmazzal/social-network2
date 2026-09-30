package repositories

import (
	"database/sql"
	"errors"
	"fmt"
)

func (raterepo *Repository) GetPostrating(poostid, useerid string) (bool, int, error) {
	var rate int
	err := raterepo.DB.Database.QueryRow("SELECT vote FROM post_rating WHERE post_id = ? AND user_id = ?", poostid, useerid).Scan(&rate)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil
		}
		return false, 0, fmt.Errorf("faild to geeet rating %w", err)
	}
	return true, rate, nil
}
func (raterepo *Repository) Like(poostid, useerid string) (bool, int, error) {
	var rate int
	err := raterepo.DB.Database.QueryRow("SELECT vote FROM post_rating WHERE post_id = ? AND user_id = ?", poostid, useerid).Scan(&rate)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0, nil
		}
		return false, 0, fmt.Errorf("faild to geeet rating %w", err)
	}
	return true, rate, nil
}
func (likerepo *Repository) PostLike(rate int ,useerid string, postid string ) error {
	var rateexxist int
	/// check if there is a rprevious rate
	querySelect := `SELECT vote FROM post_rating WHERE post_id = ? AND user_id = ?`
	err := likerepo.DB.Database.QueryRow(querySelect, postid, useerid).Scan(&rateexxist)

	if err != nil {
		/// if there is not rate add onnnne 
		if errors.Is(err, sql.ErrNoRows) {
			queryInsert := `INSERT INTO post_rating (post_id, user_id, vote) VALUES (?, ?, ?)`
			_, err = likerepo.DB.Database.Exec(queryInsert, postid, useerid, rate)
			if err != nil {
				return fmt.Errorf("repository: failed to insert vote: %w", err)
			}
			return nil
		}
		return fmt.Errorf("repository: failed to check existing rate", err)
	}

	// if i foound a rate
	if rateexxist == rate {
		//(toggle off ) if the user click on the same button twince
		queryDelete := `DELETE FROM post_rating WHERE post_id = ? AND user_id = ?`
		_, err = likerepo.DB.Database.Exec(queryDelete, postid, useerid)
		if err != nil {
			return fmt.Errorf("repository: failed to remove vote: %w", err)
		}
	} else {
		// (Switching) the user choice anthor state like/ddislike
		queryUpdate := `UPDATE post_rating SET vote = ? WHERE post_id = ? AND user_id = ?`
		_, err = likerepo.DB.Database.Exec(queryUpdate, rate, postid, useerid)
		if err != nil {
			return fmt.Errorf("repository: failed to update vote: %w", err)
		}
	}

	return nil
}
