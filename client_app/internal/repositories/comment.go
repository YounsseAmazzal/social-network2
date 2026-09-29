package repositories

import (
	"errors"
	"fmt"
	"social/internal/models"

	"github.com/gofrs/uuid"
)

func (postr *Repository) CreaatComment(commentpyyload *models.Comment) error {
	fmt.Println("commment content ",commentpyyload.Content)
	fmt.Println("commment content ",commentpyyload.CreatedAt)
	fmt.Println("commment content ",commentpyyload.Image)
	fmt.Println("commment content ",commentpyyload.CreaterID)
	fmt.Println("commment content ",commentpyyload.PostID)
	fmt.Println("commment content ",commentpyyload.ID)
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
		commentpyyload.ID, commentpyyload.Content,commentpyyload.CreaterID,
		commentpyyload.PostID, commentpyyload.Image,
	)
	if err != nil {
		fmt.Println("THe commment error",err)
		return errors.New("repository: failed to insert Comments")
	}
	// acccept all the channgments
	if err := tx.Commit(); err != nil {
		return errors.New("Comments repository: failed to commit transtaion ")
	}

	return nil
}
