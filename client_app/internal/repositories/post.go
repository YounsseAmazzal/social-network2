package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"social/internal/models"
	"social/internal/utils"

	"github.com/gofrs/uuid"
)

func (postr *Repository) CreaatePost(postPyyload *models.PostCreation) error {
	newID, err := uuid.NewV4()
	if err != nil {
		return errors.New("repository: failed to generate uuid")
	}
	postPyyload.ID = newID.String()
	//because i treat and involes with two tables post and post_shared_users that why i use tronsrtion

	// also i don't use the original db to let the changmets does save imedytly  and pernametly
	// that why i use transtations
	tx, err := postr.DB.Database.Begin()
	if err != nil {
		return fmt.Errorf("Post repository: failed to start transaction")
	}
	defer tx.Rollback()
	postQuery := `
		INSERT INTO posts (id, title, description, authorID, Image, Privecytype) 
		VALUES (?, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(postQuery,
		postPyyload.ID, postPyyload.Title, postPyyload.Description,
		postPyyload.CreaterID, postPyyload.Image, postPyyload.Postprivecy,
	)
	if err != nil {
		return errors.New("repository: failed to insert post")
	}
	if postPyyload.Postprivecy == "followers_chosen" && len(postPyyload.SharedWith) > 0 {
		audienceQuery := `INSERT INTO post_shared_users (post_id, user_id) VALUES (?, ?)`
		/// insteaad of send 10 sql commmends to db i have use tx.prepare that send send data once and make a loop to them
		// create an statment that close with rollback and commit
		stmt, err := tx.Prepare(audienceQuery)
		if err != nil {
			return errors.New("repository: failed to prepare audience statement")
		}
		defer stmt.Close()

		for _, targetUserID := range postPyyload.SharedWith {
			if _, err := stmt.Exec(postPyyload.ID, targetUserID); err != nil {
				return fmt.Errorf("repository: failed to insert audience user %s: %w", targetUserID, err)
			}
		}
	}
	// acccept all the channgments
	if err := tx.Commit(); err != nil {
		return errors.New("Post repository: failed to commit transtaion ")
	}

	return nil
}
func (postr *Repository) GetAllVisiblePosts(viewerID string) ([]*models.Post, error) {
	var posts []*models.Post
	// i make this query with 4 special rules
	//1 public posts
	/// My own posts
	// Followers only
	/// specifec users

	query := `
		SELECT 
			p.id, p.title, u.nickname AS authorName, p.createDate,
			(SELECT COUNT(*) FROM comments WHERE postID = p.id) AS Commentsnum,
			COALESCE(p.Image, '') AS image_url,
			COUNT(DISTINCT CASE WHEN v.rate = 1 THEN v.user_id END) AS likesNum,
			COUNT(DISTINCT CASE WHEN v.rate = 0 THEN v.user_id END) AS dislikesNum,
			COALESCE(my_vote.rate, -1) AS viewer_vote_status
		FROM posts p
		JOIN users u ON p.authorID = u.id
		LEFT JOIN comments c ON p.id = c.postID
		LEFT JOIN post_rating v ON p.id = v.post_id
		LEFT JOIN post_rating my_vote ON p.id = my_vote.post_id AND my_vote.user_id = ?
		WHERE 
			p.Privecytype = 'publicpost'                             
			OR p.authorID = ?                                   
			OR (p.Privecytype = 'only_followers' AND EXISTS (   
				SELECT 1 FROM followers f 
				WHERE f.following_id = p.authorID AND f.follower_id = ? AND f.status = 'folowAccepted'
			))
			OR (p.Privecytype = 'followers_chosen' AND EXISTS (          
				SELECT 1 FROM post_shared_users psu 
				WHERE psu.post_id = p.id AND psu.user_id = ?
			))
		GROUP BY p.id
		ORDER BY p.createDate DESC`

	rows, err := postr.DB.Database.Query(query, viewerID, viewerID, viewerID, viewerID)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to fetch posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		p := &models.Post{}
		var voteStatus int

		err := rows.Scan(
			&p.ID, &p.Title, &p.AuthorName, &p.CreateDate,
			&p.TotalComments, &p.Image_URL, &p.Likes, &p.Dislikes, &voteStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan post: %w", err)
		}

		// Logic Processing
		p.CreateDate = utils.RegleThedate(p.CreateDate)
		if p.Image_URL != "" {
			p.Image_URL = "/socialup/" + p.Image_URL
		}

		// Map -1 (no vote) to nil, and 0/1 to the pointer
		if voteStatus != -1 {
			p.VoteStatus = &voteStatus
		}

		posts = append(posts, p)
	}

	return posts, nil
}
func (postr *Repository) CheckPostOnGroupposts(postid string) (bool, error) {
	query := ` SELECT EXISTS( SELECT 1 FROM posts_of_group WHERE id = ? );`
	var kaayn bool
	err := postr.DB.Database.QueryRow(query, postid).Scan(&kaayn)
	if err != nil {
		return false, err
	}
	return kaayn, nil
}
func (postrepo *Repository) CheeeckPostidvalid(postid string) (bool, error) {
	//fmt.Println("POSSSSSSSSSSSSST  ID ",postid)
	var kayn int
	query := `SELECT COUNT(id) FROM posts WHERE id = ?`

	err := postrepo.DB.Database.QueryRow(query, postid).Scan(&kayn)
	if err != nil {
		log.Fatal("errrrr", err)
		return false, fmt.Errorf("database error: %w", err)
	}
	if kayn == 0 {
		return false, nil
	}
	return true, nil

}
func (r *Repository) GetCompletePostByID(postID, viewerID string) (*models.PostDetails, error) {
	var post models.PostDetails
	var userVote sql.NullString
	query := `
		SELECT
			p.id, p.title, p.description, p.authorID, u.nickname,
			COALESCE(u.avatarURL, ''), p.createDate, COALESCE(p.Image, ''),
			(SELECT COUNT(*) FROM post_rating WHERE post_id = p.id AND rate = 1) AS likesnum,
			(SELECT COUNT(*) FROM post_rating WHERE post_id = p.id AND rate = 0) AS dislikesnum,
			(SELECT CASE WHEN rate = 1 THEN 'like' WHEN rate = 0 THEN 'dislike' END
			 FROM post_rating WHERE post_id = p.id AND user_id = ?) AS user_vote
		FROM posts p
		JOIN users u ON u.id = p.authorID
		WHERE p.id = ?
		AND (
			p.Privecytype = 'publicpost'
			OR p.authorID = ?
			OR (p.Privecytype = 'only_followers' AND EXISTS (
				SELECT 1 FROM followers f WHERE f.following_id = p.authorID AND f.follower_id = ? AND f.status = 'folowAccepted'
			))
			OR (p.Privecytype = 'followers_chosen' AND EXISTS (
				SELECT 1 FROM post_shared_users pau WHERE pau.post_id = p.id AND pau.user_id = ?
			))
		)`

	err := r.DB.Database.QueryRow(query, viewerID, postID, viewerID, viewerID, viewerID).Scan(
		&post.ID, &post.Title, &post.Description, &post.ID, &post.AuthorName,
		&post.AuthorAvatar, &post.CreateDate, &post.Image_URL, &post.Likes, &post.Dislikes,
		&userVote,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("post not found or access denied")
		}
		return nil, fmt.Errorf("repository: failed to fetch complete post: %w", err)
	}

	//  Image formattiiinng
	if post.Image_URL != "" {
		post.Image_URL = "/socialup" + post.Image_URL
	}
	if post.AuthorAvatar != "" {
		post.AuthorAvatar = "/socialup/" + post.AuthorAvatar
	}

	// Post-processing: User vote pointer logic
	if userVote.Valid {
		voteValue := userVote.String
		post.Raaating = &voteValue
	}

	// Post-processing: Date formatting (Optional, depends on your frontend needs)
	post.CreateDate = utils.RegleThedate(post.CreateDate)

	return &post, nil
}
