package pg

import (
	"context"
	"strconv"
	"strings"

	"github.com/huyfm/rtc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{db: db}
}

func (s *UserService) FindUserByID(ctx context.Context, id int) (rtc.User, error) {
	return findUserByID(ctx, s.db, id)
}

func (s *UserService) FindUserByGithubID(ctx context.Context, githubID int) (rtc.User, error) {
	users, err := findUsers(ctx, s.db, rtc.FilterUser{GithubID: &githubID})
	if err != nil {
		return rtc.User{}, nil
	}
	return users[0], nil
}

func (s *UserService) FindUsers(ctx context.Context, filter rtc.FilterUser) ([]rtc.User, error) {
	return findUsers(ctx, s.db, filter)
}

func (s *UserService) CreateUser(ctx context.Context, user rtc.User) (int, error) {
	return createUser(ctx, s.db, user)
}

func findUserByID(ctx context.Context, db *pgxpool.Pool, id int) (rtc.User, error) {
	var user rtc.User
	err := db.QueryRow(ctx, `
		SELECT * FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Name, &user.Email, &user.GithubID)
	return user, err
}

func findUsers(ctx context.Context, db *pgxpool.Pool, filter rtc.FilterUser) ([]rtc.User, error) {
	// Build up WHERE clause.
	count := 0
	where, args := []string{"1 = 1"}, []interface{}{}
	if filter.ID != nil {
		count++
		where, args = append(where, "id = $"+strconv.Itoa(count)), append(args, filter.ID)
	}
	if filter.Email != nil {
		count++
		where, args = append(where, "email = $"+strconv.Itoa(count)), append(args, filter.Email)
	}
	if filter.GithubID != nil {
		count++
		where, args = append(where, "github_id = $"+strconv.Itoa(count)), append(args, filter.GithubID)
	}

	rows, err := db.Query(ctx, `
		SELECT id, name, email, github_id
		FROM users	
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id`,
		args,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]rtc.User, 0)
	for rows.Next() {
		var user rtc.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.GithubID); err != nil {
			return nil, err
		}
		res = append(res, user)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return res, nil
}

func createUser(ctx context.Context, db *pgxpool.Pool, user rtc.User) (int, error) {
	var id int
	err := db.QueryRow(ctx, `
		INSERT INTO users ( name, email, github_id) 
		VALUES ($1, $2, $3)
		RETURNING id`,
		user.Name, user.Email, user.GithubID,
	).Scan(&id)
	return id, err
}
