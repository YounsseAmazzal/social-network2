package repositories

import (
	"fmt"
	"social/internal/models"
	"time"
)

func (messagesreepo *Repository) GetGroupMessagesCursor(groupID string, limit int, akhirmessagetime string) ([]models.Groupmessages, error) {
	var query string
	var args []any

	if akhirmessagetime == ""  || akhirmessagetime=="0"{
		// bring the new messsages0
		query = `
			SELECT m.id, m.sender_id, m.group_id, m.content, m.sent_at,
			       u.nickname, COALESCE(u.avatarURL, ''), u.firstname || ' ' || u.lastname
			FROM group_messages m
			JOIN users u ON u.id = m.sender_id
			WHERE m.group_id = ?
			ORDER BY m.sent_at DESC
			LIMIT ?`
		args = []any{groupID, limit}
	} else {
		// oollllllllld messages
		query = `
			SELECT m.id, m.sender_id, m.group_id, m.content, m.sent_at,
			       u.nickname, COALESCE(u.avatarURL, ''), u.firstname || ' ' || u.lastname
			FROM group_messages m
			JOIN users u ON u.id = m.sender_id
			WHERE m.group_id = ? AND m.sent_at < ?
			ORDER BY m.sent_at DESC
			LIMIT ?`
		args = []any{groupID, akhirmessagetime, limit}
	}

	rows, err := messagesreepo.DB.Database.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// scan
	var messages []models.Groupmessages
	for rows.Next() {
		var m models.Groupmessages
		var sentAt time.Time
		var nickname string
		err := rows.Scan(
			&m.ID, &m.SenderID, &m.GroupID, &m.GRoup_Message, &sentAt,
			&nickname, &m.UserAvatar, &m.FullName,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan message: %w", err)
		}
		m.Send_date = sentAt.Format("2006-01-02 15:04:05")
		if m.UserAvatar != "" {
			m.UserAvatar = m.UserAvatar
		}

		messages = append(messages, m)
	}

	return messages, nil
}
