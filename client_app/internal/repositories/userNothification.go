package repositories

import (
	"errors"
	"fmt"
	"social/internal/models"

	"github.com/gofrs/uuid"
)

func (nothificationrepo *Repository) UserNothification(liimit int, userid string) ([]*models.Notification, error) {
	if liimit <= 0 {
		liimit = 20
	}
	query := `
		SELECT 
			id, 
			recipient_id, 
			COALESCE(sender_id, '') AS sender_id, 
			notification_category, 
			COALESCE(target_reference_id, '') AS target_id, 
			COALESCE(target_reference_type, '') AS target_type, 
			display_message, 
			read_status, 
			created_at
		FROM notifications
		WHERE recipient_id = ?
		ORDER BY created_at DESC
		LIMIT ?`

	rows, err := nothificationrepo.DB.Database.Query(query, userid, liimit)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to fetch notifications: %w", err)
	}
	defer rows.Close()

	var list []*models.Notification

	for rows.Next() {
		n := &models.Notification{}
		var statusInt int

		err := rows.Scan(
			&n.ID,
			&n.RecipientID,
			&n.SenderID,
			&n.Category,
			&n.TargetID,
			&n.TargetType,
			&n.Message,
			&statusInt,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan notification row: %w", err)
		}

		n.IsRead = (statusInt == 1)

		list = append(list, n)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: error during notification iteration: %w", err)
	}
	fmt.Println("the useeer nothiiiif ",list)
	return list, nil
}

func (nothificationrepo *Repository) CreateNotification(n *models.Notification) error {
	newID, err := uuid.NewV4()
	if err != nil {
		return errors.New("repository: failed to generate uuid")
	}
	if n.ID == "" {
		n.ID = newID.String()
	}

	query := `
		INSERT INTO notifications (
			id, 
			recipient_id, 
			sender_id, 
			notification_category, 
			target_reference_id, 
			target_reference_type, 
			display_message, 
			read_status
		) VALUES (?, ?, ?, ?, ?, ?, ?, 0)`

	_, err = nothificationrepo.DB.Database.Exec(query,
		n.ID,
		n.RecipientID,
		n.SenderID,
		n.Category,
		n.TargetID,
		n.TargetType,
		n.Message,
	)

	if err != nil {
		return fmt.Errorf("repository: failed to insert notification: %w", err)
	}

	return nil
}
