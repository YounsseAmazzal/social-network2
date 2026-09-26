package repositories

import (
	"database/sql"
	"errors"
	"fmt"
)

func (followrepo *Repository) Follwing(followerID, targetID string) (string, error) {

	var isPrivate bool
	err := followrepo.DB.Database.QueryRow("SELECT is_private FROM users WHERE id = ?", targetID).Scan(&isPrivate)
	if err != nil {
		return "", fmt.Errorf("target user not found: %w", err)
	}

	inistatus := "folowAccepted"
	if isPrivate {
		inistatus = "followpending"
	}
	//  Uasaaage INSERT OR IGNORE ooor  ON CONFLICT
	/// good practice SELECT
	query := `
	INSERT INTO followers (follower_id, following_id, status)
	VALUES (?, ?, ?)
	ON CONFLICT(follower_id, following_id) DO UPDATE SET
	status = excluded.status WHERE status = 'followdeclined' OR status = 'followpending'`

	_, err = followrepo.DB.Database.Exec(query, followerID, targetID, inistatus)
	if err != nil {
		return "", fmt.Errorf("failed to process follow request: %w", err)
	}

	return inistatus, nil
}
func (followrepo *Repository) Unfollowing(followerID, targetID string) error {
	if followerID == targetID {
		return errors.New("cannot unfollow yourself")
	}
	query := `DELETE FROM followers WHERE follower_id = ? AND following_id = ?`
	result, err := followrepo.DB.Database.Exec(query, followerID, targetID)
	if err != nil {
		return fmt.Errorf("repository: failed to delete follow record: %w", err)
	}
	//when i apply the DELETE the result means that the commends successd so
	// that why i have used RowsAffected to check if really delete it
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: error checking affected rows: %w", err)
	}

	// thaaaat means he doesn't follow him from first
	if rowsAffected == 0 {
		return errors.New("you are not following this user")
	}

	return nil
}

func (followrepo *Repository) CheckPending(followerID, targetID string) (bool, error) {
	var ThereIsPanding int
	//query for pending
	query := `
		SELECT 1 FROM followers 
		WHERE follower_id = ? AND following_id = ? AND status = 'followpending' 
		LIMIT 1`

	err := followrepo.DB.Database.QueryRow(query, followerID, targetID).Scan(&ThereIsPanding)
	///the query not found
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("repository: failed to check pending status: %w", err)
	}

	return true, nil
}
func (followrepo *Repository) AcceptFollowRequest(followerID, userID string) error {
	// should i accept only the pendding requests that exist
	query := `
		UPDATE followers 
		SET status = 'folowAccepted' 
		WHERE follower_id = ? AND following_id = ? AND status = 'followpending'`

	result, err := followrepo.DB.Database.Exec(query, followerID, userID)
	if err != nil {
		return fmt.Errorf("repository: failed to update follow status: %w", err)
	}
	//for check if the request handled or not 
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: error checking affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no pending follow request found to accept")
	}

	return nil
}