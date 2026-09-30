package repositories

import (
	"errors"
	"fmt"
	"social/internal/models"
	"time"

	"github.com/gofrs/uuid"
)

func (grouprepo *Repository) Creategroup(grouppyload *models.Group) error {
	newID, err := uuid.NewV4()
	if err != nil {
		return errors.New("repository: failed to generate uuid")
	}
	//  cause i would interact with two tables i have  use transactions 
	  
	grouppyload.ID = newID.String()
	tx, err := grouprepo.DB.Database.Begin()
	if err != nil {
		return fmt.Errorf("repository: failed to start transaction: %w", err)
	}
	defer tx.Rollback()
	queryGroup := `
		INSERT INTO groups (id, user_id, title,  description, created_at)
		VALUES (?, ?, ?, ?, ?)`

	_, err = tx.Exec(queryGroup,
		newID,
		grouppyload.Title,
		grouppyload.CreatorID,
		grouppyload.Description,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("repository: failed to insert group: %w", err)
	}

	queryMember := `INSERT INTO group_members (group_id, member_id) VALUES (?, ?)`
	_, err = tx.Exec(queryMember, newID, grouppyload.CreatorID)
	if err != nil {
		return fmt.Errorf("repository: failed to add creator as member: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("repository: failed to commit transaction: %w", err)
	}

	return  nil
}
func (grouprepo *Repository) BringDiscovers(useriid string)([]*models.Group, error) {
	fmt.Println("HE enter to bring data")
	query := `
		SELECT 
			g.id, g.user_id, g.title, g.description, g.created_at,
			(SELECT COUNT(*) FROM group_members WHERE group_id = g.id) as members_count,
			COALESCE(gr.status, 'none') as request_status
		FROM groups g
		LEFT JOIN group_requests gr ON gr.group_id = g.id AND gr.sender_id = ?
		WHERE NOT EXISTS (
			SELECT 1 FROM group_members gm 
			WHERE gm.group_id = g.id AND gm.member_id = ?
		)
		ORDER BY g.created_at DESC`

	rows, err := grouprepo.DB.Database.Query(query, useriid, useriid)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to fetch discovered groups: %w", err)
	}
	defer rows.Close()

	var Discoveerdgroups []*models.Group
	for rows.Next() {
		g := &models.Group{}
		var requestStatus string

		// (Scanning)
		err := rows.Scan(
			&g.ID, &g.CreatorID, &g.Title, &g.Description, 
			&g.CreatedAt, &g.MembersCount, &requestStatus,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan group: %w", err)
		}
		Discoveerdgroups = append(Discoveerdgroups, g)
	}
	fmt.Println("The rowwsss",rows)
	// requeststats tooo check if "pending"... and change the ui deepends on it
	if err = rows.Err(); err != nil {
		fmt.Println("the scaned rows",err) 
		return nil, err
	}

	return Discoveerdgroups, nil
}
