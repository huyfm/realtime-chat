package pg_test

import (
	"context"
	"testing"
	"time"

	"github.com/huyfm/rtc"
	"github.com/huyfm/rtc/pg"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matryer/is"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestUser_CreateUser_OK(t *testing.T) {
	db := MustOpenDB(t)
	s := pg.NewUserService(db)
	is := is.New(t)

	user := rtc.User{Name: "huy"}
	id, err := s.CreateUser(context.Background(), user)
	is.NoErr(err)
	is.Equal(id, 1)

	found, err := s.FindUserByID(context.Background(), 1)
	is.NoErr(err)
	user.ID = 1
	is.Equal(found, user) // must found user with id 1
}

func MustOpenDB(t testing.TB) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:17.2-alpine",
		postgres.WithInitScripts("../schema.sql"),
		postgres.WithDatabase("db"),
		postgres.WithUsername("admin"),
		postgres.WithPassword("admin"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatal(err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	pgpool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatal(err)
	}
	return pgpool
}
