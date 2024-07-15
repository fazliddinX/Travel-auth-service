package postgres

import (
	pb "Auth-service/genproto/auth_service"
	"Auth-service/models"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func (u *UserRepo) Create(in *pb.RegisterUserReq) (*pb.RegisterUserRes, error) {
	user := pb.RegisterUserRes{}
	err := u.DB.QueryRow("INSERT INTO users (user_name, email, password_hash, full_name) VALUES ($1, $2, $3, $4) returning id, created_at",
		in.UserName, in.Email, in.Password, in.FullName).Scan(user.Id, user.CreatedAt)
	user.Email = in.Email
	user.FullName = in.FullName
	user.UserName = in.UserName
	return &user, err
}

func (u *UserRepo) Login(in *pb.LoginRequest) (models.LoginUser, error) {
	var user models.LoginUser
	err := u.DB.QueryRow("select id, full_name, age, email from users where email = $1 and password_hash = $2 and deleted_at = 0", in.Email, in.Password).
		Scan(&user.Id, &user.Name, &user.Age, &user.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return models.LoginUser{}, errors.New("user not found")
	}
	return user, nil
}

func (u *UserRepo) GetProfile(id *pb.Id) (*pb.Profile, error) {
	user := pb.Profile{}
	err := u.DB.QueryRow("select id, user_name, email, full_name ,bio, countriesVisited, created_at, updated_at from users"+
		"where id = $1 and deleted_at = 0", id.Id).Scan(&user.Id, &user.UserName, &user.Email, &user.FullName, &user.CountriesVisited, &user.CreatedAt, &user.UpdatedAt)
	return &user, err
}

func (u *UserRepo) UpdateProfile(in *pb.UpdateUser) (*pb.Profile, error) {
	_, err := u.DB.Exec("update users set  full_name = $1, bio = $2, CountriesVisited = $3 where id = $4 and deleted_at = 0",
		in.FullName, in.Bio, in.CountriesVisited, in.Id)
	if err != nil {
		return nil, err
	}

	user := pb.Profile{}
	err = u.DB.QueryRow("select id, user_name, email, full_name ,bio, countriesVisited, created_at, updated_at from users"+
		"where id = $1 and deleted_at = 0", in.Id).Scan(&user.Id, &user.UserName, &user.Email, &user.FullName, &user.CountriesVisited, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) GetUsers(filter *pb.FilterGet) (*pb.Users, error) {
	query := "SELECT id, user_name, full_name, countries_visited FROM users WHERE deleted_at = 0"
	args := []interface{}{}

	if filter.UserName != "" {
		query += " AND user_name ILIKE $1"
		args = append(args, "%"+filter.UserName+"%")
	}
	if filter.FullName != "" {
		query += " AND full_name ILIKE $2"
		args = append(args, "%"+filter.FullName+"%")
	}
	if filter.CountriesVisited != 0 {
		query += " AND countries_visited = $3"
		args = append(args, filter.CountriesVisited)
	}
	query += " LIMIT $4 OFFSET $5"
	args = append(args, filter.Limit, filter.Offset)

	rows, err := u.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		user := &pb.User{}
		err := rows.Scan(&user.Id, &user.UserName, &user.FullName, &user.CountriesVisited)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &pb.Users{Users: users}, nil
}

func (u *UserRepo) Delete(id *pb.Id) (*pb.Success, error) {
	_, err := u.DB.Exec(`UPDATE users SET 
		deleted_at = date_part('epoch', current_timestamp)::INT 
		WHERE comments_id = $1 AND deleted_at = 0`, id.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Success{Successful: "deleted successfully"}, nil
}

func (u *UserRepo) PasswordRecovery(email string) (bool, error) {
	var id string
	err := u.DB.QueryRow("select password_hash from users where email = $1", email).Scan(&id)
	if err != nil {
		return false, err
	}
	fmt.Printf("Password recovery email sent to: %s\n", email)
	return true, nil
}

func (u *UserRepo) ActivityProfile(id *pb.Id) (*pb.UserActivities, error) {
	var a int64
	err := u.DB.QueryRow("select countries_visited from users where id = $1", id.Id).Scan(&a)
	if err != nil {
		return nil, err
	}
	activities := &pb.UserActivities{
		CountriesVisited: a,
	}
	return activities, nil
}

func (u *UserRepo) Follow(followReq *pb.FollowRequest) (*pb.FollowResponse, error) {
	_, err := u.DB.Exec("INSERT INTO follows (follower_id, following_id) VALUES ($1, $2)",
		followReq.FollowerId, followReq.FollowingId)
	if err != nil {
		return nil, err
	}
	followResp := &pb.FollowResponse{
		FollowerId:  followReq.FollowerId,
		FollowingId: followReq.FollowingId,
		FollowedAt:  time.Now().String(),
	}
	return followResp, nil
}

func (u *UserRepo) Unfollow(followReq *pb.FollowRequest) error {
	_, err := u.DB.Exec("DELETE FROM follows WHERE follower_id = $1", followReq.FollowerId)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepo) GetFollowers(in *pb.FilterFollowers) (*pb.Followers, error) {
	query := `SELECT u.id, u.email, u.full_name FROM users u JOIN followers f ON u.id = f.follower_id
                                  where f.following_id = $1 AND deleted_at = 0`
	args := []interface{}{in.Id}
	if in.UserName != "" {
		query += " AND u.user_name ILIKE $2"
		args = append(args, "%"+in.UserName+"%")
	}
	if in.Email != "" {
		query += " AND u.email ILIKE $3"
		args = append(args, "%"+in.Email+"%")
	}
	query += " LIMIT $4 OFFSET $5"
	args = append(args, in.Limit, in.Offset)

	rows, err := u.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []*pb.Follower
	for rows.Next() {
		follower := &pb.Follower{}
		err = rows.Scan(&follower.Id, &follower.UserName, &follower.Email)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}
	return &pb.Followers{Followers: followers}, nil
}

func (u *UserRepo) AddTravelStories(story *pb.ResTravelStories) (*pb.TravelStories, error) {
	//var id string
	//err := u.DB.QueryRow("INSERT INTO travel_stories (title, content, location, tags, author_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id",
	//	story.Title, story.Content, story.Location, pq.Array(story.Tags), story.AuthorId).Scan(&id)
	//if err != nil {
	//	return nil, err
	//}
	return nil, nil
}
