package repositories

import (
	"fmt"
	"log"
	"social/internal/models"
	"social/internal/utils"
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
func (r *Repository) GetChatUsers(userID string) ([]*models.Userchat, error) {
	var users []*models.Userchat

	// the users and the last chat between them 
	// subquery to bring the last message with the users 
	// rollle :there a converstion with the user or following
	// actice chat appearss first
	query := `
	SELECT DISTINCT
		u.id, u.nickname, COALESCE(u.avatarURL, ''),
		COALESCE(m.content, '') AS last_message,
		COALESCE(m.createDate, '') AS last_message_time
	FROM users u
	LEFT JOIN (
		SELECT 
			CASE WHEN senderID = ? THEN receiverID ELSE senderID END AS otherUserID,
			MAX(createDate) AS maxDate
		FROM messages
		WHERE senderID = ? OR receiverID = ?
		GROUP BY otherUserID
	) latest ON u.id = latest.otherUserID
	LEFT JOIN messages m ON (
		(m.senderID = latest.otherUserID OR m.receiverID = latest.otherUserID) 
		AND m.createDate = latest.maxDate
	)
	WHERE u.id != ?
	AND (
		latest.otherUserID IS NOT NULL
		OR EXISTS (
			SELECT 1 FROM followers f
			WHERE f.status = 'folowAccepted'
			AND ((f.following_id = u.id AND f.follower_id = ?) OR (f.follower_id = u.id AND f.following_id = ?))
		)
	)
	ORDER BY 
		(last_message_time = '') ASC,
		last_message_time DESC, 
		u.nickname ASC`

	rows, err := r.DB.Database.Query(query, userID, userID, userID, userID, userID, userID)
	if err != nil {
		log.Fatal("there is an errrrrrrr",err)
		return nil, fmt.Errorf("repository: failed to fetch chat users: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		u := &models.Userchat{}
		err := rows.Scan(
			&u.ID, 
			&u.Nickname, 
			&u.Avatar_URL, 
			&u.LastChatMessage, 
			&u.LastChatDate,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan chat user: %w", err)
		}
		if u.LastChatDate != "" {
			u.LastChatDate = utils.RegleThedate(u.LastChatDate)
		}
		// log.Fatal("the user avaatr",u.Avatar_URL)
		if u.Avatar_URL != "" {
			u.Avatar_URL = "/socialup/auth/"+ u.Avatar_URL
		}
		
		users = append(users, u)
	}
	
	// log.Fatal("there is an errrrrrrr",err, users)
	return users, nil
}

func (r *Repository) CoulldChat(usserrrid, taalkerid string) (bool, error) {
	var exists int
	/// ensure that there a rolation between the users following or followers
	query := `
		SELECT EXISTS (
			SELECT 1 FROM followers 
			WHERE status = 'folowAccepted' 
			AND (
				(follower_id = ? AND following_id = ?) 
				OR 
				(follower_id = ? AND following_id = ?)
			)
		)`
	err := r.DB.Database.QueryRow(query, usserrrid, taalkerid, taalkerid, usserrrid).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("repository: failed to check chat permission: %w", err)
	}
	return exists == 1, nil
}


func (r *Repository) GetPrivateMessages(userID, takerwithid string, limit int, lastTime string) ([]*models.Message, error) {
	var query string
	var args []any
	/// query for the messages from the siddes 
	baseQuery := `
		SELECT id, senderID, receiverID, content, createDate 
		FROM messages 
		WHERE ((senderID = ? AND receiverID = ?) OR (senderID = ? AND receiverID = ?))`

	if lastTime == "" || lastTime == "0" {
		query = baseQuery + ` ORDER BY createDate DESC LIMIT ?`
		args = []any{userID, takerwithid, takerwithid, userID, limit}
	} else {
			// bring the new ones
		query = baseQuery + ` AND createDate < ? ORDER BY createDate DESC LIMIT ?`
		args = []any{userID, takerwithid, takerwithid, userID, lastTime, limit}
	}

	rows, err := r.DB.Database.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to fetch private messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message

	for rows.Next() {
		var sentAt time.Time
		m := &models.Message{}
		err := rows.Scan(&m.ID, &m.Seender_id, &m.Receiver_ID, &m.Message, &m.Send_at)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan message: %w", err)
		}
		// i forgot to format date
		m.Send_at = sentAt.Format("2006-01-02 15:04:05")
		messages = append(messages, m)
	}
	/// reverce from the old the new cause i bring by DESC 
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}