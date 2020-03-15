package routes

import (
	"encoding/json"
	"github.com/segmentio/ksuid"
	"log"
	"net/http"
	"time"
)

const SessionTimeout = 86400 // 1 day in seconds
const CookieSessionToken = "session_token"

type loginBody struct {
	Token  string `json:"auth_token"`
	UserId string `json:"user_id"`
}

func (app *App) RegisterLoginRoutes() {
	app.router.HandleFunc("/login", app.Login).Methods("POST")
	app.router.HandleFunc("/logout", app.Logout).Methods("POST")
	//app.router.Use(app.authCheckMiddleware)
}

func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	var loginBody loginBody

	err := json.NewDecoder(r.Body).Decode(&loginBody)
	if err != nil {
		log.Printf("App.Login: couldnt parse body, %+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := app.store.UserProvider.GetOneUser(loginBody.UserId)
	if err != nil {
		log.Printf("App.Login: couldnt parse body, %+v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionToken, err := ksuid.NewRandom()
	if err != nil {
		log.Printf("App.Login: couldnt gen ksuid, %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = app.cache.Do("SETEX", sessionToken.String(), SessionTimeout, user.ID)
	if err != nil {
		log.Printf("App.Login: couldnt store session-key, %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    CookieSessionToken,
		Value:   sessionToken.String(),
		Expires: time.Now().Add(SessionTimeout * time.Second),
	})
	w.WriteHeader(http.StatusAccepted)
}

func (app *App) Logout(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie("session_token")
	_, err := app.cache.Do("DEL", c.Value)
	if err != nil {
		log.Printf("App.Logout: couldnt delete session-key, %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (app *App) authCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.RequestURI == "/login" { // short-circuit for /login route
			next.ServeHTTP(w, r)
			return
		}
		c, err := r.Cookie(CookieSessionToken)
		if err != nil {
			if err == http.ErrNoCookie {
				log.Printf("App.authCheckMiddleware: no cookie, %+v\n", err)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			log.Printf("App.authCheckMiddleware: error - , %+v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sessionToken := c.Value

		response, err := app.cache.Do("GET", sessionToken)
		if err != nil {
			log.Printf("App.authCheckMiddleware: couldnt get session-token, %+v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if response == nil {
			log.Printf("App.authCheckMiddleware: unknown session-token, %+v\n", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// TODO logic for auto-refresh (security concern???)
		// If TTL is < 3600s then create new sessions token and attach cookie
		next.ServeHTTP(w, r)
	})
}
