package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"social/internal/models"
	"time"

	"github.com/gofrs/uuid"
)

func (grouprepo *Repository) Creategroup(grouppyload *models.GroupItem) error {
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
		grouppyload.CreatorID,
		grouppyload.Title,
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

	return nil
}
func (grouprepo *Repository) BringDiscovers(useriid string) ([]*models.GroupItem, error) {
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

	var Discoveerdgroups []*models.GroupItem
	for rows.Next() {
		g := &models.GroupItem{}
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
	fmt.Println("The rowwsss", rows)
	// requeststats tooo check if "pending"... and change the ui deepends on it
	if err = rows.Err(); err != nil {
		fmt.Println("the scaned rows", err)
		return nil, err
	}

	return Discoveerdgroups, nil
}
func (grourepo *Repository) JooinedGroups(useriid string) ([]*models.GroupItem, error) {
	query := `
		SELECT 
			g.id, g.user_id, g.title, g.description, g.created_at,
			(SELECT COUNT(*) FROM group_members WHERE group_id = g.id) as members_count
		FROM groups g
		INNER JOIN group_members gm ON g.id = gm.group_id
		WHERE gm.member_id = ?
		ORDER BY g.created_at DESC`

	rows, err := grourepo.DB.Database.Query(query, useriid)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to fetch joined groups: %w", err)
	}
	defer rows.Close()

	var groups []*models.GroupItem
	for rows.Next() {
		g := &models.GroupItem{}

		err := rows.Scan(
			&g.ID,
			&g.CreatorID,
			&g.Title,
			&g.Description,
			&g.CreatedAt,
			&g.MembersCount,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan joineed group: %w", err)
		}

		//while i bgring the jjoined groups so i will use IsMember to ...
		g.IsMember = true

		groups = append(groups, g)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: error during iteration: %w", err)
	}

	return groups, nil
}
func (grourepo *Repository) CheckifUsserMember(useriid string, groupid string) (bool, error) {
	var ismeember int
	// I have use EXists cause its Power
	query := ` SELECT EXISTS(
			SELECT 1 FROM group_members 
			WHERE group_id = ? AND member_id = ?
			)`
	err := grourepo.DB.Database.QueryRow(query,groupid, useriid).Scan(&ismeember)

	if err != nil {
		return false, fmt.Errorf("repository: failed to check group membership: %w", err)
	}

	return ismeember == 1, nil
}
func (grourepo *Repository) GouupInfoo(useriid string, groupid string) (*models.GroupDetail, error) {
	query := `
		SELECT 
			g.id, g.title, g.description, g.created_at,
			u.id, COALESCE(u.nickname, ''), u.firstname, u.lastname, COALESCE(u.avatarURL, ''),
			(SELECT COUNT(*) FROM group_members WHERE group_id = g.id) AS members_count
		FROM groups g
		JOIN users u ON g.user_id = u.id
		WHERE g.id = ?
		LIMIT 1`

	groupdetail := &models.GroupDetail{}
	err := grourepo.DB.Database.QueryRow(query, groupid).Scan(
		&groupdetail.ID,
		&groupdetail.Title,
		&groupdetail.Description,
		&groupdetail.CreatedAt,
		&groupdetail.CreatorID,
		&groupdetail.Creater.ID,
		&groupdetail.Creater.Firstname,
		&groupdetail.Creater.Lastname,
		&groupdetail.Creater.AvatarURL,
		&groupdetail.MembersCount,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("grouup not found")
		}
		return nil, fmt.Errorf("repository: failed to fetch group details:")
	}
	//  
	groupdetail.CreatorID = groupdetail.Creater.ID

	return groupdetail, nil
}
func (grourepo *Repository) GoupExist( groupid string) (bool, error) {
	var kayn bool
	query := `SELECT 1  FROM groups WHERE id = ? LIMIT 1`
	err := grourepo.DB.Database.QueryRow(query, groupid).Scan(&kayn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("grouup not found")
		}
		return false, fmt.Errorf("repository: failed to fetch group details:")
	}
	// group kayn
	return true, nil
}