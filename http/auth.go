package http

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/go-github/v66/github"
	"github.com/huyfm/rtc"
	"golang.org/x/oauth2"
)

func (s *Server) handleOauthGithub(w http.ResponseWriter, r *http.Request) {
	// Create random OAuth2 state to prevent csrf attack.
	buf := make([]byte, 64)
	rand.Read(buf)
	state := hex.EncodeToString(buf)
	// Store OAuth2 state in request's cookie.
	sess := Session{State: state}
	if err := s.setSession(w, sess); err != nil {
		rtc.Logger.Warn().Err(err).Msg("")
	}
	// Redirect agent to Github OAuth2 consent page asking for permission.
	http.Redirect(w, r, s.OAuth2Svc.AuthCodeURL(state), http.StatusFound)
}

func (s *Server) handleOauthGithubCallback(w http.ResponseWriter, r *http.Request) {
	sess, err := s.session(r)
	if err != nil {
		Error(w, r, err)
		return
	}

	state, code := r.FormValue("state"), r.FormValue("code")
	// Check OAuth2 state to prevent csrf attack.
	if state != sess.State {
		Error(w, r, errors.New("state not matched"))
		return
	}

	// Exchange access code for access token.
	tok, err := s.OAuth2Svc.Exchange(r.Context(), code)
	if err != nil {
		Error(w, r, fmt.Errorf("can't exchange code for token err=%s", err))
		return
	}

	// Prepare Github OAuth2 client with access token.
	c := oauth2.NewClient(r.Context(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: tok.AccessToken}))
	gitClient := github.NewClient(c)

	// Query user info from Github.
	info, _, err := gitClient.Users.Get(r.Context(), "")
	if err != nil {
		Error(w, r, errors.New("can't get github user"))
		return
	}
	if info.Name == nil || info.ID == nil {
		Error(w, r, errors.New("empty github name/ID"))
		return
	}

	// Search githubID in database. If not found, create new user.
	var userID int
	if user, err := s.UserSrv.FindUserByGithubID(r.Context(), int(*info.ID)); err != nil {
		userID = user.ID
	} else {
		user := rtc.User{
			Name:     *info.Name,
			Email:    info.Email,
			GithubID: int(*info.ID),
		}
		if id, err := s.UserSrv.CreateUser(r.Context(), user); err != nil {
			Error(w, r, err)
			return
		} else {
			userID = id
		}
	}

	// Set UserID, clear state in session cookie.
	sess = Session{UserID: userID}
	s.setSession(w, sess)

	// Redirect to homepage
	http.Redirect(w, r, "/", http.StatusFound)
}
