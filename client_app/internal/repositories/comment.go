package repositories

import (
	"errors"
	"fmt"
	"social/internal/models"
	"social/internal/utils"

	"github.com/gofrs/uuid"
)

func (postr *Repository) CreaatComment(commentpyyload *models.Comment) error {
	// //fmt.Println("commment content ",commentpyyload.Content)
	// //fmt.Println("commment content ",commentpyyload.CreatedAt)
	// //fmt.Println("commment content ",commentpyyload.Image)
	// //fmt.Println("commment content ",commentpyyload.CreaterID)
	// //fmt.Println("commment content ",commentpyyload.PostID)
	// //fmt.Println("commment content ",commentpyyload.ID)
	newID, err := uuid.NewV4()
	if err != nil {
		return errors.New("repository: failed to generate uuid")
	}
	commentpyyload.ID = newID.String()
	//because i treat and involes with two tables post and post_shared_users that why i use tronsrtion

	// also i don't use the original db to let the changmets does save imedytly  and pernametly
	// that why i use transtations
	tx, err := postr.DB.Database.Begin()
	if err != nil {
		return fmt.Errorf("Comments repository: failed to start transaction")
	}
	defer tx.Rollback()
	postQuery := `INSERT INTO comments (id, content, createrID, postID, comment_image) VALUES (?, ?, ?, ?, ?)`

	_, err = tx.Exec(postQuery,
		commentpyyload.ID, commentpyyload.Content, commentpyyload.CreaterID,
		commentpyyload.PostID, commentpyyload.Image,
	)
	if err != nil {
		//fmt.Println("THe commment error", err)
		return errors.New("repository: failed to insert Comments")
	}
	// acccept all the channgments
	if err := tx.Commit(); err != nil {
		return errors.New("Comments repository: failed to commit transtaion ")
	}

	return nil
}

func (comment *Repository) GetPostCommensts(userId string, posstId string, limit,offfeset int) ([]*models.CommentResponse, error) {
	var comments []*models.CommentResponse
	query := `
		SELECT 
			c.id, c.content, c.createrID, c.createDate, COALESCE(c.comment_image, ''), 
			u.nickname, COALESCE(u.avatarURL, '')
		FROM comments c
		JOIN users u ON c.createrID = u.id
		WHERE c.postID = ?
		ORDER BY c.createDate DESC
		LIMIT ? OFFSET ?`

	rows, err := comment.DB.Database.Query(query, posstId, limit,offfeset)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to fetch comments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		c := &models.CommentResponse{}

		err := rows.Scan(
			&c.ID,
			&c.Content,
			&c.CreaterID,
			&c.CreatedAt,
			&c.Image,
			&c.CreaterName,
			&c.Createravatar,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan comment: %w", err)
		}

		c.CreatedAt = utils.RegleThedate(c.CreatedAt)
		if c.Image != "" {
			c.Image = "/socialup/" + c.Image
		}
		if c.Createravatar != "" {
			c.Createravatar = "/socialup/" + c.Createravatar
		}

		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	// //fmt.Println("WE BRING THE POSTS COMMENTS HERE ",comments)
	return comments, nil
}
