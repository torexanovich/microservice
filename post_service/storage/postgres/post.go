package postgres

import (
	"fmt"
	"log"
	"time"

	p "gitlab.com/micro/post_service/genproto/post"
)

func (r *PostRepo) CreatePost(post *p.PostRequest) (*p.GetPostResponse, error) {
	var res p.GetPostResponse
	err := r.db.QueryRow(`
		insert into 
			posts(title, description, user_id) 
		values
			($1, $2, $3) 
		returning 
			id, title, description, likes, user_id, created_at, updated_at`, post.Title, post.Description, post.UserId).Scan(&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to create post")
		return &p.GetPostResponse{}, err
	}

	return &res, nil
}

func (r *PostRepo) GetPostById(post *p.IdRequest) (*p.GetPostResponse, error) {
	var res p.GetPostResponse
	err := r.db.QueryRow(`
		select 
			id, title, description, likes, user_id, created_at, updated_at 
		from 
			posts 
		where 
			id = $1 and deleted_at is null`, post.Id).Scan(&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to get post")
		return &p.GetPostResponse{}, err
	}

	return &res, nil
}

func (r *PostRepo) GetPostByUserId(id *p.IdUser) (*p.Posts, error) {
	res := p.Posts{}
	rows, err := r.db.Query(`
		select 
			id, title, description, likes, user_id, created_at, updated_at 
		from 
			posts 
		where 
			user_id = $1 and deleted_at is null`, id.Id)

	if err != nil {
		log.Println("failed to get post by user_id")
		return &p.Posts{}, err
	}

	for rows.Next() {
		post := p.GetPostResponse{}

		err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Description,
			&post.Likes,
			&post.UserId,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

		if err != nil {
			log.Println("failed to scanning post")
			return &p.Posts{}, err
		}

		res.Posts = append(res.Posts, &post)
	}

	return &res, nil
}

func (r *PostRepo) GetPostForUser(id *p.IdUser) (*p.Posts, error) {
	res := p.Posts{}
	rows, err := r.db.Query(`
		select 
			id, title, description, likes, user_id, created_at, updated_at 
		from 
			posts 
		where
			user_id = $1`, id.Id)

	if err != nil {
		log.Println("failed to get post")
		return &p.Posts{}, err
	}

	for rows.Next() {
		post := p.GetPostResponse{}

		err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Description,
			&post.Likes,
			&post.UserId,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

		if err != nil {
			log.Println("failed to scanning post")
			return &p.Posts{}, nil
		}

		res.Posts = append(res.Posts, &post)
	}

	return &res, nil
}

func (r *PostRepo) GetPostForComment(post *p.IdRequest) (*p.GetPostResponse, error) {
	res := p.GetPostResponse{}
	err := r.db.QueryRow(`
		select 
			id, title, description, likes, user_id, created_at, updated_at 
		from 
			posts 
		where 
			id = $1 and deleted_at is null`, post.Id).Scan(&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to get post for comment")
		return &p.GetPostResponse{}, err
	}

	return &res, nil
}

func (r *PostRepo) SearchByTitle(title *p.Title) (*p.Posts, error) {
	res := p.Posts{}
	query := fmt.Sprint("select id, title, description, likes, user_id, created_at, updated_at from posts where title ilike '%" + title.Title + "%' and deleted_at is null")

	rows, err := r.db.Query(query)
	if err != nil {
		log.Println("failed to search post")
		return &p.Posts{}, nil
	}

	for rows.Next() {
		post := p.GetPostResponse{}

		err = rows.Scan(
			&post.Id,
			&post.Title,
			&post.Description,
			&post.Likes,
			&post.UserId,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			log.Println("failed to scanning post")
			return &p.Posts{}, nil
		}

		res.Posts = append(res.Posts, &post)
	}

	return &res, nil
}

func (r *PostRepo) LikePost(l *p.LikeRequest) (*p.GetPostResponse, error) {
	res := p.GetPostResponse{}
	if l.IsLiked {
		err := r.db.QueryRow(`
			update 
				posts 
			set 
				likes = likes + 1 
			where 
				id = $1 
			returning 
				id, title, description, likes, user_id, created_at, updated_at`, l.PostId).Scan(
			&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)
		if err != nil {
			log.Println("failed to like post")
			return &p.GetPostResponse{}, err
		}
	} else {
		err := r.db.QueryRow(`
			select 
				id, title, description, likes + 1, user_id, created_at, updated_at 
			from 
				posts 
			where 
				id = $1`, l.PostId).Scan(&res.Id, &res.Title, &res.Description, &res.Likes, &res.UserId, &res.CreatedAt, &res.UpdatedAt)

		if err != nil {
			log.Println("failed to like post")
			return &p.GetPostResponse{}, err
		}
	}

	return &res, nil
}

func (r *PostRepo) UpdatePost(post *p.UpdatePostRequest) error {
	res, err := r.db.Exec(`
		update
			posts 
		set 
			title = $1, description = $2, updated_at = $3 
		where 
			id = $4`, post.Title, post.Description, time.Now(), post.Id)
	if err != nil {
		log.Println("failed to update post")
		return err
	}

	fmt.Println(res.RowsAffected())

	return nil
}

func (r *PostRepo) DeletePost(id *p.IdRequest) (*p.GetPostResponse, error) {
	post := p.GetPostResponse{}
	err := r.db.QueryRow(`
		update 
			posts 
		set 
			deleted_at = $1 
		where 
			id = $2 
		returning 
			id, title, description, likes, user_id, created_at, updated_at`, time.Now(), id.Id).Scan(&post.Id, &post.Title, &post.Description, &post.Likes, &post.UserId, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		log.Println("failed to delete post")
		return &p.GetPostResponse{}, err
	}

	return &post, nil
}
