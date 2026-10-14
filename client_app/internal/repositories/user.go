package repositories

import (
	"database/sql"
	"fmt"
	"html"
	"log"
	"social/internal/models"
	"social/internal/utils"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

func (ur *Repository) CreateUser(payload *models.UserInfo) (*models.User, error) {
	newID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("repository: failed to generate uuid: %w", err)
	}
	user := &models.User{
		ID:        newID.String(),
		Firstname: strings.TrimSpace(payload.Firstname),
		Lastname:  strings.TrimSpace(payload.Lastname),
		Email:     strings.ToLower(strings.TrimSpace(payload.Email)),
		Password:  payload.Password,
		Age:       payload.Age,
		Gender:    payload.Gender,
		AboutMe:   html.EscapeString(payload.AboutMe),
		IsPrivate: payload.IsPrivate == 1,
	}
	nickname := strings.TrimSpace(payload.Nickname)
	if nickname == "" {
		if len(user.Firstname) > 0 {
			user.Nickname = strings.ToLower(string(user.Firstname[0]) + user.Lastname)
		} else {
			user.Nickname = "user_" + user.ID[:8]
		}
	} else {
		user.Nickname = strings.ToLower(nickname)
	}

	user.AvatarURL = payload.Avatar
	// log.Fatal("the usseer avaaataaar", payload)
	// log.Fatal("the usseer avaaataaar", payload.Avatar)
	if user.AvatarURL == "" {
		user.AvatarURL = utils.Defauteavaaatar
	}

	query := `
		INSERT INTO users (
			id, nickname, firstname, lastname, age, 
			gender, email, password, avatarURL, about_me, is_private
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = ur.DB.Database.Exec(query,
		user.ID,
		user.Nickname,
		user.Firstname,
		user.Lastname,
		user.Age,
		user.Gender,
		user.Email,
		user.Password,
		user.AvatarURL,
		user.AboutMe,
		payload.IsPrivate,
		//autmaaateee privete (boool )
	)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to insert user: %w", err)
	}
	return user, nil
}
func (ur *Repository) CheckuserIfExxist(usernameoremail string) (*models.User, bool) {
	var UserInfo models.User
	//fmt.Println("the indefair", usernameoremail)
	usernameoremail = strings.ToLower(usernameoremail)
	query := `
		SELECT 
			id, nickname, firstname, lastname, age, gender, 
			email, avatarURL, COALESCE(about_me, ''), COALESCE(is_private, 0), password 
		FROM users 
		WHERE nickname = ? OR email = ? 
		LIMIT 1`

	err := ur.DB.Database.QueryRow(query, usernameoremail, usernameoremail).Scan(
		&UserInfo.ID,
		&UserInfo.Nickname,
		&UserInfo.Firstname,
		&UserInfo.Lastname,
		&UserInfo.Age,
		&UserInfo.Gender,
		&UserInfo.Email,
		&UserInfo.AvatarURL,
		&UserInfo.AboutMe,
		&UserInfo.IsPrivate,
		&UserInfo.Password,
	)
	//fmt.Println("that is the id of user", UserInfo.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		log.Printf("Repository Error Fach bghit njib user:%v", err)
		return nil, false
	}

	return &UserInfo, true
}
func (ur *Repository) IsNicknameTaken(nickname string) (bool, error) {
	var exists int
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE nickname = ?)`
	err := ur.DB.Database.QueryRow(query, strings.ToLower(nickname)).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}
func (ur *Repository) Checksession(tocccken string) (string, error) {
	var userID string
	var expireAt time.Time
	//queery to data base
	query := "SELECT user_id, expire_at FROM sessions WHERE token = ?"
	err := ur.DB.Database.QueryRow(query, tocccken).Scan(&userID, &expireAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("session not found")
		}
		return "", fmt.Errorf("database error: %w", err)
	}

	if time.Now().After(expireAt) {
		//cleaaant the dataa base
		_ = ur.DeleteSession(tocccken)
		return "", fmt.Errorf("session expired")
	}

	return userID, nil
}
func (ur *Repository) DeleteSession(token string) error {
	query := "DELETE FROM sessions WHERE token = ?"
	_, err := ur.DB.Database.Exec(query, token)
	if err != nil {
		return fmt.Errorf("repository: failed to delete session: %w", err)
	}
	return nil
}
func (ur *Repository) GetUserByToken(token string) (*models.User, error) {
	var UseerId string
	var expireAt time.Time
	//taaaake the sesssion and ensure that it is good
	query := "SELECT user_id, expire_at FROM sessions WHERE token = ?"
	err := ur.DB.Database.QueryRow(query, token).Scan(&UseerId, &expireAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("database error during session lookup: %w", err)
	}
	// //fmt.Println("$$$$$$$$user from session")
	// the expireetion
	if time.Now().After(expireAt) {
		//deeelte the session automatlcty
		_, _ = ur.DB.Database.Exec("DELETE FROM sessions WHERE token = ?", token)
		return nil, fmt.Errorf("session expired")
	}

	exaactuser, err := ur.BringUserByid(UseerId)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user for active session: %w", err)
	}
	return exaactuser, nil
}
func (ur *Repository) BringUserByid(usseridd string) (*models.User, error) {
	var useeer models.User
	// i prefer that way to query to data ((Backticks
	query := `
		SELECT 
			id, nickname, firstname, lastname, age, gender, 
			email, avatarURL, COALESCE(about_me, ''), COALESCE(is_private, 0) 
		FROM users 
		WHERE id = ? 
		LIMIT 1`

	err := ur.DB.Database.QueryRow(query, usseridd).Scan(
		&useeer.ID,
		&useeer.Nickname,
		&useeer.Firstname,
		&useeer.Lastname,
		&useeer.Age,
		&useeer.Gender,
		&useeer.Email,
		&useeer.AvatarURL,
		&useeer.AboutMe,
		&useeer.IsPrivate,
	)
	if useeer.AvatarURL!=""{
		useeer.AvatarURL="/socialup/auth/"+useeer.AvatarURL
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("repository: failed to get user by id: %w", err)
	}

	return &useeer, nil
}
func (ur *Repository) TakeUserProfile(useeerid string, Profileid string) (*models.ProfileResponse, error) {
	//fmt.Println("he is teh profile owner Profile",Profileid)
	userinfoo, err := ur.fetchBasicUserInfo(Profileid)
	if err != nil {
		// log.Fatal("heeeeessddfsdsdfdsf")
		return nil, fmt.Errorf("Proofie not found %w", &err)
	}
	isOwner := useeerid == Profileid
	rolaation := ur.checkFollowStatus(useeerid, Profileid)
	profile := &models.ProfileResponse{
		User:        *userinfoo,
		IsPrivate:   userinfoo.IsPrivate,
		IsOwner:     isOwner,
		IsFollowing: rolaation=="folowAccepted",
		FollowStatus: rolaation,
		FollowersCount: userinfoo.FollowersCount,
		FollowingCount: userinfoo.FollowingCount,
		PostsCount: userinfoo.PostsCount,
	}
	profile.User.Rolation=rolaation
	//return profile info for now without checking following and bringing the posts
	return profile, nil
	// return
}
func (ur *Repository) UpdateProfileSettings(userID, nickname, aboutMe string, isPrivate bool) error {
	//with updation of the profille should i concider the following requests are accepted

	tx, err := ur.DB.Database.Begin()
	if err != nil {
		return err
	}

	// Rollback will trigger if an error occurs before Commit()
	defer tx.Rollback()

	query := `
		UPDATE users 
		SET nickname = ?, about_me = ?, is_private = ? 
		WHERE id = ?`

	privateInt := 0
	if isPrivate {
		privateInt = 1
	}

	_, err = tx.Exec(query, nickname, aboutMe, privateInt, userID)
	if err != nil {
		return err
	}

	if !isPrivate {
		acceptQuery := `
			UPDATE followers 
			SET status = 'folowAccepted' 
			WHERE following_id = ? AND status = 'followpending'`
		
		_, err = tx.Exec(acceptQuery, userID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
// GetFollowStatus   (accepted, pending, none)
func (ur *Repository) checkFollowStatus(followerID, followingID string) string {
	var status string
	query := `SELECT status FROM followers WHERE follower_id = ? AND following_id = ?`
	
	err := ur.DB.Database.QueryRow(query, followerID, followingID).Scan(&status)
	
	if err != nil {
		///no rolaation
		return "none"
	}
	
	//  retuuurn 'folowAccepted' or 'followpending'
	return status
}
func (ur *Repository) fetchBasicUserInfo(ProoofileId string) (*models.User, error) {
	//fmt.Println("The profile id",ProoofileId)
	var usserr models.User
	quuery := `
		SELECT 
			u.id, u.nickname, u.firstname, u.lastname, u.age, u.gender, u.email, u.avatarURL,
			COALESCE(u.about_me, ''), 
			COALESCE(u.is_private, 0),
			(SELECT COUNT(*) FROM followers WHERE following_id = u.id AND status = 'folowAccepted') AS followers_count,
			(SELECT COUNT(*) FROM followers WHERE follower_id = u.id AND status = 'folowAccepted') AS following_count,
			(SELECT COUNT(*) FROM posts WHERE authorID = u.id AND Privecytype = 'publicpost') AS public_posts_count
		FROM users u 
		WHERE u.id = ? 
		LIMIT 1`
	err := ur.DB.Database.QueryRow(quuery, ProoofileId).Scan(
		&usserr.ID, &usserr.Nickname, &usserr.Firstname, &usserr.Lastname,
		&usserr.Age, &usserr.Gender, &usserr.Email, &usserr.AvatarURL,
		&usserr.AboutMe, &usserr.IsPrivate,
		&usserr.FollowersCount,
				&usserr.FollowingCount,
						&usserr.PostsCount,
	)
	if usserr.AvatarURL!=""{
		usserr.AvatarURL="/socialup/auth/"+usserr.AvatarURL
	}
	// //fmt.Println("!!!!!!!!!!!!!!I find it!!!!!!!!!!!",err)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return &usserr, err
}
func (ur *Repository) UserExistsByID(userID string) (bool, error) {
	var exists int
		query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)`

	err := ur.DB.Database.QueryRow(query, userID).Scan(&exists)
	
	if err != nil {
		return false, fmt.Errorf("repository: error checking user existence: %w", err)
	}
	return exists == 1, nil
}