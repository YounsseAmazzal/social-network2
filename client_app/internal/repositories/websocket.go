package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"social/internal/models"
	"strings"
	"time"
	"uuid"
)

func (groupmessage *Repository) CreateGroupMessage(groupID, senderID, content string) (string, error) {
	msgID := uuid.New().String()
	fmt.Println("msssssid", msgID)
	query := `
		INSERT INTO group_messages (id, sender_id, group_id, content, sent_at) 
		VALUES (?, ?, ?, ?, ?)`

	_, err := groupmessage.DB.Database.Exec(query,
		msgID,
		senderID,
		groupID,
		content,
		time.Now(),
	)
	// fmt.Println("there is an errror here ")
	if err != nil {
		return "", fmt.Errorf("repository: failed to save group message: %w", err)
	}
	return msgID, nil
}

func (r *Repository) GetGroupMetadata(groupID string) (*models.GroupItem, error) {
	// usinng COUNT to count number of group members
	query := `
		SELECT 
			g.id, g.title, 
			(SELECT COUNT(*) FROM group_members WHERE group_id = g.id) AS members_count,
			(SELECT GROUP_CONCAT(member_id) FROM group_members WHERE group_id = g.id) AS members_list
		FROM groups g
		WHERE g.id = ? 
		LIMIT 1`

	var group models.GroupItem
	var membersIDs sql.NullString

	err := r.DB.Database.QueryRow(query, groupID).Scan(
		&group.ID,
		&group.Title,
		&group.MembersCount,
		&membersIDs,
	)
	// procccess the errors
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("group not found")
		}
		return nil, fmt.Errorf("repository: failed to fetch group metadata: %w", err)
	}

	if membersIDs.Valid && membersIDs.String != "" {
		group.GroupMember = strings.Split(membersIDs.String, ",")
	} else {
		group.GroupMember = []string{}
	}

	return &group, nil
}
