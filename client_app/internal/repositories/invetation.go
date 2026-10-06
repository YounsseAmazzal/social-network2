package repositories

import (
	"errors"
	"fmt"
)

func (groupinvetasion *Repository) SendInvetation(grouupid, inviterID, inviteeID string) error {
	if inviterID == inviteeID {
		return errors.New("you cannot invite yourself to the group")
	}
	// i haaaave use ON CONFLICT to treat with the repeted requests
	query := `
		INSERT INTO group_Send_invetations (group_id, inviter_id, invitee_id, status)
		VALUES (?, ?, ?, 'pending')
		ON CONFLICT(group_id, invitee_id) 
		DO UPDATE SET status = 'pending', inviter_id = excluded.inviter_id 
		WHERE status = 'declined'`

	result, err := groupinvetasion.DB.Database.Exec(query, grouupid, inviterID, inviteeID)
	if err != nil {
		return fmt.Errorf("repository: failed to send group invitation: %w", err)
	}

	// ensure the invetaation sent to the user
	//  any line shouldn't affecteeed by user is member or have a invetation
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("invitation already sent or user is already a member")
	}

	return nil
}
