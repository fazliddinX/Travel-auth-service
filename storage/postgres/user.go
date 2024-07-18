package postgres

import (
	pb "auth-service/generated/user"
	"auth-service/models"
	"database/sql"
	"time"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (repo *UserRepo) CreateUser(user models.RegisterRequest) (*models.RegisterResponse, error) {
	var (
		response  models.RegisterResponse
		createdAt time.Time
	)
	err := repo.DB.QueryRow(`
		INSERT INTO users (
			username,
			email,
			password_hash,
			full_name
		)
		VALUES (
			$1, $2, $3, $4
		)
		RETURNING
			id,
			username,
			email,
			full_name,
			created_at
	`, user.Username, user.Email, user.Password, user.FullName).
		Scan(&response.ID, &response.Username, &response.Email, &response.FullName, &createdAt)

	if err != nil {
		return nil, err
	}

	response.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	return &response, nil
}

func (repo *UserRepo) GetUserByEmail(email string) (*models.LoginResponse, error) {
	var response models.LoginResponse

	err := repo.DB.QueryRow(`
		SELECT
			id,
			username,
			email,
			password_hash
		FROM
			users
		WHERE
			deleted_at = 0 AND email = $1
	`, email).Scan(&response.ID, &response.Username, &response.Email, &response.Password)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (repo *UserRepo) UpdatePassword(resetPassword models.UpdatePassword) (*models.Success, error) {
	_, err := repo.DB.Exec(`
		UPDATE users
		SET 
			password_hash = $1
		WHERE 
			id = $2 AND deleted_at = 0
	`, resetPassword.NewPassword, resetPassword.ID)

	if err != nil {
		return &models.Success{
			Message: "Error in updating password",
		}, err
	}

	return &models.Success{
		Message: "Password reset successfully",
	}, nil
}

func (repo *UserRepo) EmailExists(email string) (bool, error) {
	var exists bool
	err := repo.DB.QueryRow(`
		SELECT
			EXISTS (
				SELECT 1
				FROM users
				WHERE email = $1
			)
	`, email).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (repo *UserRepo) GetUserInfo(id string) (*pb.UserInfoResponse, error) {
	var response pb.UserInfoResponse

	err := repo.DB.QueryRow(`
		SELECT
			id,
			username,
			full_name
		FROM
			users
		WHERE
			deleted_at = 0 AND id = $1
	`, id).Scan(&response.Id, &response.Username, &response.FullName)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (repo *UserRepo) GetUserProfile(id string) (*pb.GetProfileResponse, error) {
	var (
		response  pb.GetProfileResponse
		bio       sql.NullString
		createdAt time.Time
		updatedAt time.Time
	)

	err := repo.DB.QueryRow(`
		SELECT
			id,
			username,
			email,
			full_name,
			bio,
			countries_visited,
			created_at,
			updated_at
		FROM
			users
		WHERE
			id = $1 AND deleted_at = 0
	`, id).Scan(&response.Id, &response.Username, &response.Email, &response.FullName, &bio, &response.CountriesVisited, &createdAt, &updatedAt)

	if err != nil {
		return nil, err
	}

	response.Bio = bio.String
	if !bio.Valid {
		response.Bio = "No Bio"
	}
	response.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	response.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")

	return &response, nil
}

func (repo *UserRepo) UpdateUserProfile(req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	var (
		response  pb.UpdateProfileResponse
		updatedAt time.Time
	)

	err := repo.DB.QueryRow(`
		UPDATE users
		SET 
			full_name = $1,
			bio = $2,
			countries_visited = $3,
			updated_at = CURRENT_TIMESTAMP
		WHERE
			id = $4 AND deleted_at = 0
		RETURNING
			id,
			username,
			email,
			full_name,
			bio,
			countries_visited,
			updated_at
	`, req.FullName, req.Bio, req.CountriesVisited, req.Id).Scan(&response.Id, &response.Username, &response.Email, &response.FullName, &response.Bio, &response.CountriesVisited, &updatedAt)

	if err != nil {
		return nil, err
	}

	response.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")
	return &response, nil
}

func (repo *UserRepo) GetUsers(req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	offset := (req.Page - 1) * req.Limit

	rows, err := repo.DB.Query(`
		SELECT 
			id, 
			username, 
			full_name, 
			countries_visited
		FROM 
			users
		WHERE 
			deleted_at = 0
		ORDER BY 
			username
		LIMIT $1
		OFFSET $2
	`, req.Limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*pb.User
	for rows.Next() {
		var user pb.User
		if err := rows.Scan(&user.Id, &user.Username, &user.FullName, &user.CountriesVisited); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	var total int32
	err = repo.DB.QueryRow(`
		SELECT 
			COUNT(*) 
		FROM 
			users
		WHERE
			deleted_at = 0
	`).Scan(&total)

	if err != nil {
		return nil, err
	}

	response := &pb.ListUsersResponse{
		Users: users,
		Total: total,
		Page:  req.Page,
		Limit: req.Limit,
	}

	return response, nil
}

func (repo *UserRepo) DeleteUser(id string) (*pb.DeleteUserResponse, error) {
	res, err := repo.DB.Exec(`
		UPDATE users
		SET
			deleted_at = $1
		WHERE
			deleted_at = 0 AND id = $2
	`, time.Now().Unix(), id)

	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, sql.ErrNoRows
	}

	return &pb.DeleteUserResponse{
		Message: "User successfully deleted",
	}, nil
}

func (repo *UserRepo) FollowingUser(req *pb.FollowUserRequest) (*pb.FollowUserResponse, error) {
	var response pb.FollowUserResponse

	err := repo.DB.QueryRow(`
		INSERT INTO followers (
			follower_id,
			following_id
		)
		VALUES (
			$1, $2
		)
		RETURNING
			follower_id,
			following_id,
			followed_at
	`, req.FollowerId, req.FollowingId).Scan(&response.FollowerId, &response.FollowingId, &response.FollowingAt)

	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (repo *UserRepo) GetFollowers(req *pb.ListFollowersRequest) (*pb.ListFollowersResponse, error) {
	var followers []*pb.Follower
	offset := (req.Page - 1) * req.Limit

	rows, err := repo.DB.Query(`
		SELECT
			u.id,
			u.username,
			u.full_name
		FROM
			users u
		INNER JOIN
			followers f ON u.id = f.follower_id
		WHERE
			f.following_id = $1 AND u.deleted_at = 0
		OFFSET $2
		LIMIT $3
	`, req.UserId, offset, req.Limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var follower pb.Follower
		if err := rows.Scan(&follower.Id, &follower.Username, &follower.FullName); err != nil {
			return nil, err
		}
		followers = append(followers, &follower)
	}

	var total int32
	err = repo.DB.QueryRow(`
		SELECT
			COUNT(*)
		FROM
			users u
		INNER JOIN
			followers f ON u.id = f.follower_id
		WHERE
			f.following_id = $1 AND u.deleted_at = 0`,
		req.UserId).Scan(&total)

	if err != nil {
		return nil, err
	}
	response := &pb.ListFollowersResponse{
		Followers: followers,
		Total:     total,
		Page:      req.Page,
		Limit:     req.Limit,
	}

	return response, nil
}
