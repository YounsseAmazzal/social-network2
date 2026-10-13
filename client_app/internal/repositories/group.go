package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"social/internal/models"
	"social/internal/utils"
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
	// fmt.Println("HE enter to bring data")
	query := `
		SELECT 
			g.id, g.user_id, g.title, g.description, g.created_at,
			(SELECT COUNT(*) FROM group_members WHERE group_id = g.id) as members_count,
			COALESCE(gr.status, 'none') as request_status
		FROM groups g
		LEFT JOIN group_Join_requests gr ON gr.group_id = g.id AND gr.sender_id = ?
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
	// fmt.Println("The rowwsss", rows)
	// requeststats tooo check if "pending"... and change the ui deepends on it
	if err = rows.Err(); err != nil {
		// fmt.Println("the scaned rows", err)
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
	err := grourepo.DB.Database.QueryRow(query, groupid, useriid).Scan(&ismeember)

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
		&groupdetail.Creater.ID,
		&groupdetail.Creater.Nickname,
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
func (grourepo *Repository) GoupExist(groupid string) (bool, error) {
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
func (grourepo *Repository) GrouupMembers(groupid string) ([]models.User, error) {
	// the user should be memeber
	query := `
		SELECT 
			u.id, u.firstname, u.lastname, 
			COALESCE(u.nickname, ''), 
			COALESCE(u.avatarURL, '')
		FROM group_members groupmembers
		JOIN users u ON u.id = groupmembers.member_id
		WHERE groupmembers.group_id = ?`

	rows, err := grourepo.DB.Database.Query(query, groupid)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to query group members: %w", err)
	}
	defer rows.Close()
	var members []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID,
			&u.Firstname,
			&u.Lastname,
			&u.Nickname,
			&u.AvatarURL,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan member: %w", err)
		}
		if u.AvatarURL != "" {
			u.AvatarURL = "/socialup/auth/" + u.AvatarURL
		}
		members = append(members, u)
	}
	//check the errors
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}
func (grouprepo *Repository) CreatePostGroup(postGroup *models.GroupPostCreation) error {
	newID, err := uuid.NewV4()
	if err != nil {
		return errors.New("repository: failed to generate uuid")
	}
	postGroup.GroupPost.ID = newID.String()
	now := time.Now()
	tx, err := grouprepo.DB.Database.Begin()
	/// dayer 9aahwa ocha3eel garouu
	if err != nil {
		return fmt.Errorf("repository: failed to start transaction: %w", err)
	}
	defer tx.Rollback()
	query := `
		INSERT INTO posts_of_group (id, group_id, member_id, title, content, image, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(query,
		postGroup.GroupPost.ID,
		postGroup.GrouupID,
		postGroup.GroupPost.CreaterID,
		postGroup.GroupPost.Title,
		postGroup.GroupPost.Description,
		postGroup.GroupPost.Image,
		now,
	)
	if err != nil {
		return fmt.Errorf("repository: failed to insert group post: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("repository: failed to commit transaction: %w", err)
	}
	return nil
}
func (grouprepo *Repository) CreateGroupComment(Commentgp *models.Comment) error {
	now := time.Now()
	newID, err := uuid.NewV4()
	if err != nil {
		return errors.New("repository: failed to generate uuid")
	}
	tx, err := grouprepo.DB.Database.Begin()
	/// dayer 9aahwa ocha3eel garouu
	if err != nil {
		return fmt.Errorf("repository: failed to start transaction: %w", err)
	}
	Commentgp.ID = newID.String()
	defer tx.Rollback()
	query := `
		INSERT INTO group_comments (id, group_post_id, member_id, content, image, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err = tx.Exec(query,
		Commentgp.ID,
		Commentgp.PostID,
		Commentgp.CreaterID,
		Commentgp.Content,
		Commentgp.Image,
		now,
	)
	if err != nil {
		return fmt.Errorf("repository: failed to insert group post commment : %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("repository: failed to commit transaction: %w", err)
	}
	return nil
}
func (grouprepo *Repository) GetAllGroupPosts(GrouuupID string, liimit, offset int) ([]models.Post, error) {

	// query := ` SELECT * FROM  posts_of_group WHERE group_id =? `
	// rows , err:= grouprepo.DB.Database.Query(query,GrouuupID)
	// for rows.Next(){
	// 	if err!=sql.ErrNoRows{
	// 		return [],fmt.Errorf("reposiotory: failed to ")
	// 	}

	// }
	// return
	query := `
		SELECT
			p.id, p.title, p.content, p.member_id, COALESCE(p.image, '') AS image_path, p.created_at,
			u.nickname, COALESCE(u.avatarURL, ''),
			(SELECT COUNT(*) FROM group_comments WHERE group_post_id = p.id) AS comments_count,
			(SELECT COUNT(*) FROM post_group_rating WHERE post_id = p.id AND rate = 1) AS likes_count,
			(SELECT COUNT(*) FROM post_group_rating WHERE post_id = p.id AND rate = 0) AS dislikes_count
		FROM posts_of_group p
		JOIN users u ON u.id = p.member_id
		WHERE p.group_id = ?
		ORDER BY p.created_at DESC
		LIMIT ? OFFSET ?`

	rows, err := grouprepo.DB.Database.Query(query, GrouuupID, liimit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to fetch group posts: %w", err)
	}
	defer rows.Close()

	var Grouupposts []models.Post
	for rows.Next() {
		grouppost := models.Post{}
		err := rows.Scan(
			&grouppost.ID, &grouppost.Title, &grouppost.Description, &grouppost.Author.ID, &grouppost.Image_URL, &grouppost.CreateDate,
			&grouppost.Author.Firstname, &grouppost.Author.AvatarURL,
			&grouppost.TotalComments, &grouppost.Likes, &grouppost.Dislikes,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan grouup post")
		}
		//for date
		grouppost.CreateDate = utils.RegleThedate(grouppost.CreateDate)
		// log.Fatal("the image of group post",grouppost.Image_URL)
		if grouppost.Image_URL != "" {
			grouppost.Image_URL = "/socialup/posts/" + grouppost.Image_URL
		}
		// log.Fatal("usssseeer avaaatar",grouppost.Author.AvatarURL)
		if grouppost.Author.AvatarURL != "" {
			grouppost.Author.AvatarURL = "/socialup/auth/" + grouppost.Author.AvatarURL
		}
		// log.Fatal("the image of group post", grouppost.Image_URL)
		grouppost.CreaterID = grouppost.Author.ID
		grouppost.CreaterID = grouppost.Author.Nickname
		grouppost.CreaterID = grouppost.Author.AvatarURL
		Grouupposts = append(Grouupposts, grouppost)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return Grouupposts, nil
}

func (grouprepo *Repository) PostGrouupDetailed(usserid, poostid string) (*models.GrouupPostDetails, error) {
	var GroupPost models.GrouupPostDetails
	var ratiing sql.NullString
	//
	query := `
	SELECT
	gp.id, gp.member_id, u.nickname, COALESCE(u.avatarURL, ''),
	gp.title, gp.content, COALESCE(gp.image, ''), gp.created_at,
	(SELECT COUNT(*) FROM post_group_rating WHERE post_id = gp.id AND rate = 1) AS likes_count,
	(SELECT COUNT(*) FROM post_group_rating WHERE post_id = gp.id AND rate = 0) AS dislikes_count,
	(SELECT CASE WHEN rate = 1 THEN 'like' WHEN rate = 0 THEN 'dislike' END 
	FROM post_group_rating WHERE post_id = gp.id AND user_id = ?) AS viewer_vote
	FROM posts_of_group gp
	JOIN users u ON u.id = gp.member_id
	WHERE gp.id = ? 
	AND EXISTS (
		SELECT 1 FROM group_members gm 
		WHERE gm.group_id = gp.group_id AND gm.member_id = ?
		)`
	// query selection with check the member ship so i don't need to check it again
	err := grouprepo.DB.Database.QueryRow(query, usserid, poostid, usserid).Scan(
		&GroupPost.ID,
		&GroupPost.Author.ID,
		&GroupPost.AuthorName,
		&GroupPost.AuthorAvatar,
		&GroupPost.Title,
		&GroupPost.Description,
		&GroupPost.Image_URL,
		&GroupPost.CreateDate,
		&GroupPost.Likes,
		&GroupPost.Dislikes,
		&ratiing,
	)
	if err != nil {
		// log.Fatal(err)
		// post not found or user not member
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("post not fouund or you don't have peermission to view it")
		}
		return nil, fmt.Errorf("repository: failed to fetch group post: %w", err)
	}

	GroupPost.CreateDate = utils.RegleThedate(GroupPost.CreateDate)

	if GroupPost.Image_URL != "" {
		GroupPost.Image_URL = "/socialup/posts/" + GroupPost.Image_URL
	}
	if GroupPost.AuthorAvatar != "" {
		GroupPost.AuthorAvatar = "/socialup/auth/" + GroupPost.AuthorAvatar
	}
	// (((pointer logic )))
	if ratiing.Valid {
		voteVal := ratiing.String
		GroupPost.Raaating = &voteVal
	}
	// fmt.Println("THe post detaileed", GroupPost)
	return &GroupPost, nil
}
func (grouprepo *Repository) PostGroupComments(usserid, poostid string, limit, offset int) ([]*models.GroupCommentResponse, error) {
	var groupcomments []*models.GroupCommentResponse
	query := `
		SELECT
			c.id, c.content, c.member_id, c.created_at, COALESCE(c.image, ''),
			u.nickname, COALESCE(u.avatarURL, '')
		FROM group_comments c
		JOIN users u ON u.id = c.member_id
		JOIN posts_of_group gp ON gp.id = c.group_post_id
		WHERE c.group_post_id = ?
		AND EXISTS (
			SELECT 1 FROM group_members gm 
			WHERE gm.group_id = gp.group_id AND gm.member_id = ?
		)
		ORDER BY c.created_at DESC
		LIMIT ? OFFSET ?`
	// the vieweer should be member of the grouuuup
	rows, err := grouprepo.DB.Database.Query(query, poostid, usserid, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to fetch group comments: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		commentgp := &models.GroupCommentResponse{}

		err := rows.Scan(
			&commentgp.ID,
			&commentgp.Content,
			&commentgp.CreaterID,
			&commentgp.CreatedAt,
			&commentgp.Image,
			&commentgp.CreaterName,
			&commentgp.Createravatar,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan group comment: %w", err)
		}
		commentgp.CreatedAt = utils.RegleThedate(commentgp.CreatedAt)
		if commentgp.Image != "" {
			commentgp.Image = commentgp.Image
		}
		if commentgp.Createravatar != "" {
			commentgp.Createravatar = commentgp.Createravatar
		}
		groupcomments = append(groupcomments, commentgp)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return groupcomments, nil
}

func (grouprepo *Repository) GroupPostLike(rate int, useerid string, postid string) error {
	var rateexxist int
	/// check if there is a rprevious rate
	querySelect := `SELECT rate FROM post_group_rating WHERE post_id = ? AND user_id = ?`
	err := grouprepo.DB.Database.QueryRow(querySelect, postid, useerid).Scan(&rateexxist)

	if err != nil {
		/// if there is not rate add onnnne
		if errors.Is(err, sql.ErrNoRows) {
			queryInsert := `INSERT INTO post_group_rating (post_id, user_id, rate) VALUES (?, ?, ?)`
			_, err = grouprepo.DB.Database.Exec(queryInsert, postid, useerid, rate)
			if err != nil {
				return fmt.Errorf("repository: failed to insert rate: %w", err)
			}
			return nil
		}
		return fmt.Errorf("repository: failed to check existing rate", err)
	}

	// if i foound a rate
	if rateexxist == rate {
		//(toggle off ) if the user click on the same button twince
		queryDelete := `DELETE FROM post_group_rating WHERE post_id = ? AND user_id = ?`
		_, err = grouprepo.DB.Database.Exec(queryDelete, postid, useerid)
		if err != nil {
			return fmt.Errorf("repository: failed to remove vote: %w", err)
		}
	} else {
		// (Switching) the user choice anthor state like/ddislike
		queryUpdate := `UPDATE post_group_rating SET rate = ? WHERE post_id = ? AND user_id = ?`
		_, err = grouprepo.DB.Database.Exec(queryUpdate, rate, postid, useerid)
		if err != nil {
			return fmt.Errorf("repository: failed to update rate: %w", err)
		}
	}

	return nil
}

func (grouprepo *Repository) GroupJoinRequests(groupid, senderid string) error {
	querry := `
	INSERT OR IGNORE INTO group_Join_requests (group_id, sender_id) VALUES (?, ?)
	`
	_, err := grouprepo.DB.Database.Exec(querry, groupid, senderid)
	return err
}
func (grouprepo *Repository) GetGroupAdmin(grouupid string) (string, error) {
	var admiinId string
	// logicly the admin of the group is the creater =>the first user of the group table
	query := `SELECT user_id FROM groups WHERE id = ? LIMIT 1`

	err := grouprepo.DB.Database.QueryRow(query, grouupid).Scan(&admiinId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("group not found")
		}
		return "", fmt.Errorf("repository: failed to fetch group admin: %w", err)
	}
	return admiinId, nil
}
func (grouprepo *Repository) GetJoinRequesters(groupID string) ([]*models.User, error) {
	query := `
		SELECT 
			u.id, u.firstname, u.lastname, 
			COALESCE(u.nickname, ''), 
			COALESCE(u.avatarURL, '')
		FROM group_Join_requests gr
		JOIN users u ON u.id = gr.sender_id
		WHERE gr.group_id = ? AND gr.status = 'pending'`

	rows, err := grouprepo.DB.Database.Query(query, groupID)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to fetch pending members: %w", err)
	}
	defer rows.Close()

	var pendingUsers []*models.User

	for rows.Next() {
		u := &models.User{}
		err := rows.Scan(
			&u.ID,
			&u.Firstname,
			&u.Lastname,
			&u.Nickname,
			&u.AvatarURL,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan pending user: %w", err)
		}
		if u.AvatarURL != "" {
			u.AvatarURL = u.AvatarURL
		}

		pendingUsers = append(pendingUsers, u)
	}

	return pendingUsers, nil
}
func (grouprepo *Repository) IsJoinRequestPending(gouupid, userid string) (bool, error) {
	var exists int
	query := `
		SELECT EXISTS(
			SELECT 1 FROM group_Join_requests 
			WHERE group_id = ? AND sender_id = ? AND status = 'pending'
		)`
	err := grouprepo.DB.Database.QueryRow(query, gouupid, userid).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("repository: failed to check pending request: %w", err)
	}
	return exists == 1, nil
}

func (grouprepo *Repository) JoinRequestProcces(grouupid, userrid, decision string) error {
	tx, err := grouprepo.DB.Database.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	if decision == "joinaccept" {
		insertMember := `INSERT INTO group_members (group_id, member_id) VALUES (?, ?)`
		if _, err := tx.Exec(insertMember, grouupid, userrid); err != nil {
			return fmt.Errorf("failed to add member: %w", err)
		}
	}

	deleteRequest := `DELETE FROM group_Join_requests WHERE group_id = ? AND sender_id = ?`
	if _, err := tx.Exec(deleteRequest, grouupid, userrid); err != nil {
		return fmt.Errorf("failed to clear request: %w", err)
	}
	//onees the admin take the decesions remove the nothification
	queryCleanNotif := `
		DELETE FROM notifications 
		WHERE sender_id = ? 
		  AND target_reference_id = ? 
		  AND notification_category = 'group_join_request'`
	
	if _, err := tx.Exec(queryCleanNotif, userrid, grouupid); err != nil {
		fmt.Println("Warning: failed to cleanup notification:", err)
	}

	return tx.Commit()
}
func (grouprepo *Repository) GetFollowersToinvite(groupID, userID string) ([]*models.User, error) {
	query := `
		SELECT u.id, u.firstname, u.lastname, COALESCE(u.nickname, ''), COALESCE(u.avatarURL, '')
		FROM users u
		WHERE u.id IN (
			SELECT following_id FROM followers WHERE follower_id = ? AND status = 'folowAccepted'
			UNION
			SELECT follower_id FROM followers WHERE following_id = ? AND status = 'folowAccepted'
		)
		AND NOT EXISTS (
			SELECT 1 FROM group_members gm 
			WHERE gm.member_id = u.id AND gm.group_id = ?
		)
		AND NOT EXISTS (
			SELECT 1 FROM group_Send_invetations gi 
			WHERE gi.invitee_id = u.id AND gi.group_id = ? AND gi.status = 'pending'
		)`

	rows, err := grouprepo.DB.Database.Query(query, userID, userID, groupID, groupID)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to fetch invitable followers: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		err := rows.Scan(
			&u.ID,
			&u.Firstname,
			&u.Lastname,
			&u.Nickname,
			&u.AvatarURL,
		)
		if err != nil {
			return nil, fmt.Errorf("repository: failed to scan invitable user: %w", err)
		}
		if u.AvatarURL != "" {
			u.AvatarURL = u.AvatarURL
		}

		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
