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

func (nothificationrepo *Repository) CountermesaageNotRead(recipientID string) ([]models.CountermessageNotread, error) {
	query := `
		SELECT 
			sender_id, 
			COUNT(*) AS unread_total
		FROM notifications
		WHERE recipient_id = ?
		  AND notification_category = 'direct_message'
		  AND read_status = 0
		  AND sender_id IS NOT NULL AND sender_id != ''
		GROUP BY sender_id
		ORDER BY MAX(created_at) DESC`

	rows, err := nothificationrepo.DB.Database.Query(query, recipientID)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to fetch unread chat counts: %w", err)
	}
	defer rows.Close()

	var counters []models.CountermessageNotread

	for rows.Next() {
		var item models.CountermessageNotread
		if err := rows.Scan(&item.SenderID, &item.Count); err != nil {
			return nil, fmt.Errorf("repository: failed to scan unread counter: %w", err)
		}
		counters = append(counters, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return counters, nil
}

func (r *Repository) CounterGroupmesaageNotRead(recipientID string) ([]models.CountermessageGroupNotread, error) {

	query := `
		SELECT 
			target_reference_id, 
			COUNT(*) AS unread_total
		FROM notifications
		WHERE recipient_id = ?
		  AND notification_category = 'group_chat_message'
		  AND read_status = 0
		  AND target_reference_id IS NOT NULL AND target_reference_id != ''
		GROUP BY target_reference_id
		ORDER BY MAX(created_at) DESC`

	rows, err := r.DB.Database.Query(query, recipientID)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to fetch unread group counts: %w", err)
	}
	defer rows.Close()

	var counters []models.CountermessageGroupNotread

	for rows.Next() {
		var item models.CountermessageGroupNotread
		if err := rows.Scan(&item.GroupID, &item.Count); err != nil {
			return nil, fmt.Errorf("repository: failed to scan group counter: %w", err)
		}
		counters = append(counters, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return counters, nil
}