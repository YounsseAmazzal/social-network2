package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"social/internal/models"
	"time"

	"github.com/gofrs/uuid"
)

func (eventrepo *Repository) CreateEvent(event *models.Eventplyload, userid string) (*models.GroupEvent, error) {
	// var GroupEvent models.GroupEvent
	groupEvent := &models.GroupEvent{
		Title:       event.Title,
		Description: event.Description,
		GroupID:     event.GrouupID,
		CreatorID:   userid,
		EventDate:   event.EventDate,
		CreatedAt:   time.Now(),
	}
	fmt.Println("THe event before insertation", event)
	newID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("repository: failed to generate uuid: %w", err)
	}
	groupEvent.ID = newID.String()
	tx, err := eventrepo.DB.Database.Begin()
	if err != nil {
		return nil, fmt.Errorf("repository: failed to start transaction: %w", err)
	}
	defer tx.Rollback()
	query := `
		INSERT INTO group_events (
			id, group_id, member_id, title, description, 
			event_date, created_at, total_going, total_not_going
		) VALUES (?, ?, ?, ?, ?, ?, ?, 0, 0)`

	_, err = tx.Exec(query,
		groupEvent.ID,
		groupEvent.GroupID,
		groupEvent.CreatorID,
		groupEvent.Title,
		groupEvent.Description,
		groupEvent.EventDate,
		groupEvent.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("repository: failed to insert event: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("repository: failed to commit event: %w", err)
	}
	return groupEvent, nil
}

func (eventrepo *Repository) GetGroupEvents(groupid, userid string, limit, offset int) ([]*models.GroupEvent, error) {
	query := `
		SELECT 
			e.id, e.title, e.description, e.event_date, e.created_at,
			e.total_going, e.total_not_going, 
			COALESCE(ev.status, '') AS viewer_vote,
			u.id, u.firstname, u.lastname, COALESCE(u.nickname, ''), COALESCE(u.avatarURL, '')
		FROM group_events e
		LEFT JOIN group_participation ev ON ev.event_id = e.id AND ev.member_id = ?
		JOIN users u ON u.id = e.member_id
		WHERE e.group_id = ?
		ORDER BY e.event_date ASC
		LIMIT ? OFFSET ?`
	rows, err := eventrepo.DB.Database.Query(query, userid, groupid, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to query group events: %w", err)
	}
	defer rows.Close()

	var Allevents []*models.GroupEvent

	for rows.Next() {
		event := &models.GroupEvent{}

		err := rows.Scan(
			&event.ID, &event.Title, &event.Description, &event.EventDate, &event.CreatedAt,
			&event.TotalGoing, &event.TotalNotGoing,
			&event.Userrate,
			&event.Creator.ID, &event.Creator.Firstname, &event.Creator.Lastname,
			&event.Creator.Nickname, &event.Creator.AvatarURL,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan event: %w", err)
		}

		if event.Creator.AvatarURL != "" {
			event.Creator.AvatarURL ="/socialup/auth/"+ event.Creator.AvatarURL
		}

		Allevents = append(Allevents, event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return Allevents, nil
}

func (grourepo *Repository) EventEXIST(evenid  string) (bool, error) {
	var kayn bool
	query := `SELECT 1  FROM group_events WHERE id = ? LIMIT 1`
	err := grourepo.DB.Database.QueryRow(query, evenid).Scan(&kayn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("event not found")
		}
		return false, fmt.Errorf("repository: failed to fetch event details:")
	}
	// event kayn
	return true, nil
}

func (r *Repository) Eventparticipate(eventid, userid, status string) error {
	fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$the status of participation :",status)
	// remove the participation
	if status == "remove" || status == "" {
		query := `DELETE FROM group_participation WHERE event_id = ? AND member_id = ?`
	_, err := r.DB.Database.Exec( query, eventid, userid)
		if err != nil {
			return fmt.Errorf("repository: failed to remove event rate: %w", err)
		}
		return nil
	}
	// log.Fatal("he passed here ")
	// to let the user change his opinion of going or not a have useed ON CONFLICT
	//Why use  a logic  called UPSERT (UPdate and INSERT)
	query := `
		INSERT INTO group_participation (event_id, member_id, status)
		VALUES (?, ?, ?)
		ON CONFLICT(event_id, member_id) DO UPDATE SET status = excluded.status`

		fmt.Println("the status",status)
		fmt.Println("event id ",eventid)

		fmt.Println("userid",userid)
	_, err := r.DB.Database.Exec( query, eventid, userid, status)
	if err != nil {
	log.Fatal("he cames here " ,err)

		return fmt.Errorf("repository: failed to process event vote: %w", err)
	}
	return nil
}