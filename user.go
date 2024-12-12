package rtc

import "context"

type User struct {
	ID       int
	Name     string
	Email    *string
	GithubID int
}

// Validate requires that name must not empty and github ID > 0.
func (u *User) Validate() bool {
	return u.Name != "" && u.GithubID > 0
}

type FilterUser struct {
	ID       *int
	Email    *string
	GithubID *int
}

type UserService interface {
	FindUserByID(ctx context.Context, id int) (User, error)
	FindUserByGithubID(ctx context.Context, githubID int) (User, error)
	FindUsers(ctx context.Context, filter FilterUser) ([]User, error)
	CreateUser(ctx context.Context, user User) (int, error)
}
