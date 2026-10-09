package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"social/internal/models"
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
func (followrepo *Repository) BringPendingFollowers(userId string) ([]*models.User, error) {
	// treat with empty value with using COALESCE
	query := `
		SELECT 
			u.id, u.firstname, u.lastname, 
			COALESCE(u.nickname, ''), 
			COALESCE(u.avatarURL, '')
		FROM followers f
		JOIN users u ON u.id = f.follower_id
		WHERE f.following_id = ? AND f.status = 'followpending'
		ORDER BY f.created_at DESC`

	rows, err := followrepo.DB.Database.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to query pending requests: %w", err)
	}
	defer rows.Close()
	var requests []*models.User

	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Firstname,
			&user.Lastname,
			&user.Nickname,
			&user.AvatarURL,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan pending request: %w", err)
		}
		requests = append(requests, user)
	}
	//check errors here
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: error during rows iteration: %w", err)
	}

	return requests, nil
}

func (followrepo *Repository) CheckPending(followerID, targetID string) (bool, error) {
	var ThereIsPanding int
	//query for pending
	//fmt.Println("THE user ID ",followerID)
	//fmt.Println("THE Target ID ",targetID)
	query := `
		SELECT 1 FROM followers 
		WHERE follower_id = ? AND following_id = ? AND status = 'followpending' 
		LIMIT 1`

	err := followrepo.DB.Database.QueryRow(query, followerID, targetID).Scan(&ThereIsPanding)
	///the query not found
	if err != nil {
		if err == sql.ErrNoRows {
			return false, err
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
	//fmt.Println("RESULT  OF ACCEPPTING FOLLOWING",result)
	//for check if the request handled or not
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repository: error checking affected rows: %w", err)
	}
	// fmt.Printf("DEBUG: Looking for follower: %s, following: %s, with status: followpending\n", followerID, userID)

	if rowsAffected == 0 {
		// //fmt.Println("ERRRROR OF ROWAFFECTED")
		return errors.New("no pending follow request found to accept")
	}

	return nil
}

// get the users with the relation like all ,followers amd following
func (followrepo *Repository) GetUserList(action string, userID string) ([]*models.User, error) {
	var query string
	// //fmt.Println("$$$$$$$$$$$bring followers from db")
	var args []interface{}
	//define the query depend on the actions
	switch action {
	case "Userfollowers":
		// my followers
		query = followrepo.getFollowersQuery()
		args = []interface{}{userID, userID}

	case "Userfollowing":
		query = followrepo.getFollowingQuery()
		args = []interface{}{userID, userID}

	case "All":
		// all users
		query = followrepo.getDiscoverQuery()
		args = []interface{}{userID, userID, userID}

	default:
		return nil, fmt.Errorf("invalid action: %s", action)
	}

	rows, err := followrepo.DB.Database.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to fetch user list: %w", err)
	}
	defer rows.Close()
	// discover every result
	// //fmt.Println("%%%%%%%%%%%%%%%%")
	//fmt.Println("The followers rows ",rows)
	return followrepo.scanUsers(rows)
}

// isolated queries
func (r *Repository) getFollowersQuery() string {
	return `
		SELECT u.id, u.firstname, u.lastname, COALESCE(u.nickname, ''), 
		       COALESCE(u.avatarURL, ''), u.is_private,
		       CASE 
		           WHEN me.status = 'folowAccepted' THEN 'mutual'
		           WHEN me.status = 'followpending'  THEN 'requested'
		           ELSE 'follower'
		       END as Rolation
		FROM followers f
		JOIN users u ON u.id = f.follower_id
		LEFT JOIN followers me ON me.following_id = u.id AND me.follower_id = ?
		WHERE f.following_id = ? AND f.status = 'folowAccepted'`
}

func (r *Repository) getFollowingQuery() string {
	return `
		SELECT u.id, u.firstname, u.lastname, COALESCE(u.nickname, ''), 
		       COALESCE(u.avatarURL, ''), u.is_private,
		       CASE 
		           WHEN them.status = 'folowAccepted' THEN 'mutual'
		           ELSE 'following'
		       END as Rolation
		FROM followers f
		JOIN users u ON u.id = f.following_id
		LEFT JOIN followers them ON them.follower_id = u.id AND them.following_id = ?
		WHERE f.follower_id = ? AND f.status = 'folowAccepted'`
}

func (r *Repository) getDiscoverQuery() string {
	return `
		SELECT u.id, u.firstname, u.lastname, COALESCE(u.nickname, ''), 
		       COALESCE(u.avatarURL, ''), u.is_private,
		       CASE 
		           WHEN f1.status = 'folowAccepted' AND f2.status = 'folowAccepted' THEN 'mutual'
		           WHEN f1.status = 'folowAccepted' THEN 'following'
		           WHEN f1.status = 'followpending' THEN 'requested'
		           WHEN f2.status = 'folowAccepted' THEN 'follower'
		           ELSE 'none'
		       END as Rolation
		FROM users u
		LEFT JOIN followers f1 ON f1.following_id = u.id AND f1.follower_id = ?
		LEFT JOIN followers f2 ON f2.follower_id = u.id AND f2.following_id = ?
		WHERE u.id != ?`
}

// scanUsers to reaad the rowws and change them to a structs
func (followrepo *Repository) scanUsers(rows *sql.Rows) ([]*models.User, error) {
	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		err := rows.Scan(
			&u.ID, &u.Firstname, &u.Lastname, &u.Nickname,
			&u.AvatarURL, &u.IsPrivate, &u.Rolation,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	//fmt.Println("users ",users)
	return users, nil
}
