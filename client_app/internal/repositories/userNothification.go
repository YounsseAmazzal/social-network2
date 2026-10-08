package repositories

import (
	"fmt"
	"social/internal/models"
)

func (r *Repository) UserNothification( liimit int,userid string) ([]*models.Notification, error) {
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

	rows, err := r.DB.Database.Query(query, userid, liimit)
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

	return list, nil
}
