package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/securecookie"
	"github.com/huyfm/rtc"
	"github.com/huyfm/rtc/http/html"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Server struct {
	Addr   string
	Router *chi.Mux
	// TLS cert.
	TLSCert    string
	TLSPrivKey string
	// Services.
	OAuth2Svc *oauth2.Config
	CookieSvc *securecookie.SecureCookie
	UserSrv   rtc.UserService
}

// NewServer creates a new server from Config.
// Its services must be assigned manually.
func NewServer(conf *rtc.Config) *Server {
	s := &Server{}
	s.Addr = ":" + conf.SrvPort
	s.Router = chi.NewRouter()
	s.SetRoutes()
	s.TLSCert = conf.TLSCert
	s.TLSPrivKey = conf.TLSPrivKey
	s.OAuth2Svc = &oauth2.Config{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		Endpoint:     github.Endpoint,
	}
	s.CookieSvc = securecookie.New([]byte(conf.HashKey), nil)
	return s
}

func (s *Server) SetRoutes() {
	// Apply middlewares.
	s.Router.Use(s.authenticate)

	// Authentication endpoints.
	s.Router.Get("/", s.handleIndex)
	s.Router.Get("/oauth/github", s.handleOauthGithub)
	s.Router.Get("/oauth/github/callback", s.handleOauthGithubCallback)
}

func (s *Server) Open() error {
	rtc.Logger.Info().Msgf("serving at %s", s.Addr)
	if s.TLSCert == "" || s.TLSPrivKey == "" {
		rtc.Logger.Info().Msg("not found TLS cert, use http")
		return http.ListenAndServe(s.Addr, s.Router)
	}
	rtc.Logger.Info().Msg("found TLS cert, use https")
	return http.ListenAndServeTLS(s.Addr, s.TLSCert, s.TLSPrivKey, s.Router)
}

// Handle "/"
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	user := rtc.UserInContext(r.Context())
	tmpl := html.IndexPage(user)
	if err := tmpl.Render(r.Context(), w); err != nil {
		Error(w, r, err)
	}
}

// authenticate gets userID from session cookie and puts the current user
// to the request's context.
func (s *Server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, err := s.session(r)
		if err != nil || sess.UserID == 0 {
			next.ServeHTTP(w, r)
			return
		}
		// Check session's userID in database.
		user, err := s.UserSrv.FindUserByID(r.Context(), sess.UserID)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		// Add user to request's context.
		ctx := rtc.ContextWithUser(r.Context(), &user)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// setSession stores session in response's cookie.
func (s *Server) setSession(w http.ResponseWriter, sess Session) error {
	// Serialize and sign session cookie.
	val, err := s.MarshalSession(sess)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    val,
		Path:     "/",
		MaxAge:   24 * 3600,
		HttpOnly: true,
		Secure:   true,
	})
	return nil
}

// session retrieves session cookie from request.
func (s *Server) session(r *http.Request) (Session, error) {
	sessCookie, err := r.Cookie(sessionName)
	if err != nil {
		return Session{}, err
	}
	var sess Session
	if err := s.UnmarshalSession(sessCookie.Value, &sess); err != nil {
		return Session{}, err
	}
	return sess, nil
}

func (s *Server) MarshalSession(sess Session) (string, error) {
	return s.CookieSvc.Encode(sessionName, sess)
}

func (s *Server) UnmarshalSession(sessVal string, sess *Session) error {
	return s.CookieSvc.Decode(sessionName, sessVal, sess)
}
