package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	u "gitlab.com/micro/user_service/genproto/user"
)

func (r *UserRepo) CreateUser(user *u.UserRequest) (*u.UserResponse, error) {
	var res u.UserResponse
	err := r.db.QueryRow(`
		insert into 
			users(id, first_name, last_name, email, password, refresh_token) 
		values
			($1, $2, $3, $4, $5, $6) 
		returning 
			id, first_name, last_name, email, created_at, updated_at, refresh_token`, user.Id, user.FirstName, user.LastName, user.Email, user.Password, user.RefreshToken).Scan(
		&res.Id, &res.FirstName, &res.LastName, &res.Email, &res.CreatedAt, &res.UpdatedAt, &res.RefreshToken)
	if err != nil {
		log.Println("failed to create user")
		return &u.UserResponse{}, err
	}

	return &res, nil
}

func (r *UserRepo) GetUserById(user *u.IdRequest) (*u.UserResponse, error) {
	var res u.UserResponse
	err := r.db.QueryRow(`
		select 
			id, first_name, last_name, email, password, created_at, updated_at
		from 
			users 
		where id = $1 and deleted_at is null`, user.Id).Scan(
		&res.Id, &res.FirstName, &res.LastName, &res.Email, &res.Password, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to get user")
		return &u.UserResponse{}, err
	}

	return &res, nil
}

func (r *UserRepo) GetUserForClient(user_id *u.IdRequest) (*u.UserResponse, error) {
	var res u.UserResponse
	err := r.db.QueryRow(`
		select 
			id, first_name, last_name, email, created_at, updated_at 
		from 
			users 
		where id = $1`, user_id.Id).Scan(
		&res.Id, &res.FirstName, &res.LastName, &res.Email, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to get user for client")
		return &u.UserResponse{}, err
	}

	return &res, nil
}

func (r *UserRepo) GetAllUsers(req *u.AllUsersRequest) (*u.Users, error) {
	var res u.Users
	offset := (req.Page - 1) * req.Limit
	rows, err := r.db.Query(`
		select 
			id, first_name, last_name, email, created_at, updated_at 
		from 
			users 
		where 
			deleted_at is null 
		limit $1 offset $2`, req.Limit, offset)

	if err != nil {
		log.Println("failed to get all users")
		return &u.Users{}, err
	}

	for rows.Next() {
		temp := u.UserResponse{}

		err = rows.Scan(
			&temp.Id,
			&temp.FirstName,
			&temp.LastName,
			&temp.Email,
			&temp.CreatedAt,
			&temp.UpdatedAt,
		)
		if err != nil {
			log.Println("failed to scanning user")
			return &u.Users{}, err
		}

		res.Users = append(res.Users, &temp)
	}

	return &res, nil
}

func (r *UserRepo) SearchUsersByName(req *u.SearchUsers) (*u.Users, error) {
	var res u.Users
	query := fmt.Sprint("select id, first_name, last_name, email, created_at, updated_at from users where first_name ilike '%" + req.FirstName + "%' and deleted_at is null")

	rows, err := r.db.Query(query)
	if err != nil {
		log.Println("failed to searching user")
		return &u.Users{}, err
	}

	for rows.Next() {
		temp := u.UserResponse{}

		err = rows.Scan(
			&temp.Id,
			&temp.FirstName,
			&temp.LastName,
			&temp.Email,
			&temp.CreatedAt,
			&temp.UpdatedAt,
		)
		if err != nil {
			log.Println("failed to searching user")
			return &u.Users{}, err
		}

		res.Users = append(res.Users, &temp)
	}

	return &res, nil
}

func (r *UserRepo) UpdateUser(user *u.UpdateUserRequest) error {
	res, err := r.db.Exec(`
		update
			users
		set
			first_name = $1, last_name = $2, email = $3, password = $4, updated_at = $5
		where 
			id = $5`, user.FirstName, user.LastName, user.Email, user.Password, time.Now(), user.Id)

	if err != nil {
		log.Println("failed to update user")
		return err
	}

	fmt.Println(res.RowsAffected())
	return nil
}

func (r *UserRepo) DeleteUser(user *u.IdRequest) (*u.UserResponse, error) {
	temp := u.UserResponse{}
	err := r.db.QueryRow(`
		update 
			users
		set 
			deleted_at = $1 
		where 
			id = $2 
		returning 
			id, first_name, last_name, email, created_at, updated_at`, time.Now(), user.Id).Scan(
		&temp.Id, &temp.FirstName, &temp.LastName, &temp.Email, &temp.CreatedAt, &temp.UpdatedAt)

	if err != nil {
		log.Println("failed to delete user")
		return &u.UserResponse{}, err
	}

	return &temp, nil
}

func (r *UserRepo) CheckField(req *u.CheckFieldReq) (*u.CheckFieldResp, error) {
	query := fmt.Sprintf("SELECT 1 FROM users WHERE %s=$1", req.Field)
	res := &u.CheckFieldResp{}
	temp := -1
	err := r.db.QueryRow(query, req.Value).Scan(&temp)
	if err != nil {
		res.Exists = false
		return res, nil
	}
	if temp == 0 {
		res.Exists = true
	} else {
		res.Exists = false
	}
	return res, nil
}

func (r *UserRepo) Login(req *u.LoginReq) (*u.LoginResp, error) {
	var res u.LoginResp
	err := r.db.QueryRow(`
		select 
			id, first_name, last_name, email, password, created_at, updated_at
		from 
			users 
		where email = $1 and deleted_at is null`, req.Email).Scan(
		&res.Id, &res.FirstName, &res.LastName, &res.Email, &res.Password, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to get user by email")
		return &u.LoginResp{}, err
	}

	return &res, nil
}

func (r *UserRepo) GetByEmail(req *u.EmailReq) (*u.UserResponse, error) {
	var res u.UserResponse
	err := r.db.QueryRow(`
		select 
			id, first_name, last_name, email, password, created_at, updated_at
		from 
			users 
		where email = $1 and deleted_at is null`, req.Email).Scan(
		&res.Id, &res.FirstName, &res.LastName, &res.Email, &res.Password, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to get user by email")
		return &u.UserResponse{}, err
	}

	return &res, nil
}


func (r *UserRepo) GetAdmin(req *u.GetAdminReq) (*u.GetAdminRes, error) {
	var res u.GetAdminRes
		err := r.db.QueryRow(`
		SELECT 
			id,
			admin_name,
			admin_password, 
			created_at, 
			updated_at
		FROM 
			admin 
		WHERE 
			deleted_at 
			IS NULL AND 
			admin_name=$1`, req.Name).Scan(
		&res.Id,
		&res.Name,
		&res.Password,
		&res.CreatedAt,
		&res.UpdatedAt,	
	)

	if err == sql.ErrNoRows {
		fmt.Println("Error while getting admin no rows")
		return &res, nil
	}
	if err != nil {
		return &u.GetAdminRes{}, err
	}
	return &res, nil
}

func (r *UserRepo) GetModerator(req *u.GetModeratorReq) (*u.GetModeratorRes, error) {
	res := u.GetModeratorRes{}
	err := r.db.QueryRow(`SELECT 
			id, 
			name, 
			password,
			created_at,
			updated_at
		FROM 
			moderator
		WHERE deleted_at IS NULL AND name=$1`, req.Name).Scan(
		&res.Id,
		&res.Name,
		&res.Password,
		&res.CreatedAt,
		&res.UpdatedAt)

	if err == sql.ErrNoRows {
		fmt.Println("Error while getting admin no rows")
		return &res, nil
	}
	if err != nil {
		fmt.Println("error while getting moderator")
		return &u.GetModeratorRes{}, err
	}
	return &res, nil
}