package main

import (
	"crypto/subtle"
	"encoding/base64"
	"github.com/praja-dev/porgs"
	"golang.org/x/crypto/argon2"
	"log/slog"
	"net/http"
	"time"
)

// vmLogin is the view model for the login screen
type vmLogin struct {
	Usr string
	Msg string
}

// Configuration for Argon2 hashing
const (
	a2Time    = 1
	a2Memory  = 64 * 1024
	a2Threads = 4
	a2KeyLen  = 32
)

const MsgInvalidCredentials = "Invalid username or password"

func handleLoginGet() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		porgs.RenderView(w, r, porgs.View{Name: "main-login", Title: "Login | PORGS", Data: vmLogin{}})
	})
}

func handleLoginPost() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// # Extract username and password from the form
		if err := r.ParseForm(); err != nil {
			slog.Error("main.handleLoginPost: parse form", "err", err)
			porgs.ShowDefaultErrorPage(w, r)
			return
		}
		usr := r.PostFormValue("username")
		pwd := r.PostFormValue("password")

		// # Find the user record with this username
		conn, err := porgs.DbConnPool.Take(r.Context())
		if err != nil {
			slog.Error("main.handleLoginPost: get db conn", "usr", usr, "err", err)
			porgs.ShowDefaultErrorPage(w, r)
			return
		}
		defer porgs.DbConnPool.Put(conn)

		stSelect, err := conn.Prepare("SELECT password, salt FROM user WHERE username = ?")
		if err != nil {
			slog.Error("main.handleLoginPost: stmt prepare", "usr", usr, "err", err)
			porgs.ShowDefaultErrorPage(w, r)
			return
		}
		defer func() { _ = stSelect.Reset(); _ = stSelect.ClearBindings() }()

		stSelect.BindText(1, usr)

		hasRow, err := stSelect.Step()
		if err != nil {
			slog.Error("main.handleLoginPost: stmt step", "usr", usr, "err", err)
			porgs.ShowDefaultErrorPage(w, r)
			return
		}
		if !hasRow {
			slog.Info("main.handleLoginPost", "usr", usr, "msg", "not found")
			porgs.RenderView(w, r, porgs.View{Name: "main-login", Title: "Login | PORGS", Data: vmLogin{
				Usr: usr,
				Msg: MsgInvalidCredentials,
			}})
			return
		}

		// # Check if the password matches
		password := stSelect.GetText("password")
		salt := stSelect.GetText("salt")
		match, err := pwdMatch(pwd, password, salt)
		if err != nil {
			slog.Error("main.handleLoginPost: match pwd", "usr", usr, "err", err)
			porgs.ShowDefaultErrorPage(w, r)
			return
		}
		if !match {
			slog.Info("main.handleLoginPost: match pwd", "usr", usr, "msg", "incorrect password")
			porgs.RenderView(w, r, porgs.View{Name: "main-login", Title: "Login | PORGS", Data: vmLogin{
				Usr: usr,
				Msg: MsgInvalidCredentials,
			}})
			return
		}

		// # Generate session token
		token, err := porgs.RandomBase64String(16)
		if err != nil {
			slog.Error("main.handleLoginPost: generate token", "usr", usr, "err", err)
			porgs.ShowDefaultErrorPage(w, r)
			return
		}

		// # Save the session token
		now := time.Now().UTC().Unix()
		stInsert, err := conn.Prepare("INSERT INTO session (id, created, updated, username) VALUES (?, ?, ?, ?)")
		if err != nil {
			slog.Error("main.handleLoginPost: save session: stmt prepare", "usr", usr, "err", err)
			porgs.ShowDefaultErrorPage(w, r)
			return
		}
		defer func() { _ = stInsert.Reset(); _ = stInsert.ClearBindings() }()

		stInsert.BindText(1, token)
		stInsert.BindInt64(2, now)
		stInsert.BindInt64(3, now)
		stInsert.BindText(4, usr)

		_, err = stInsert.Step()
		if err != nil {
			slog.Error("main.handleLoginPost: save session: stmt exec", "usr", usr, "err", err)
			porgs.ShowDefaultErrorPage(w, r)
			return
		}

		// # Set the session token in an HttpOnly cookie
		cookie := http.Cookie{
			Name:     porgs.CookieNameSession,
			Path:     "/",
			Value:    token,
			MaxAge:   int(24 * time.Hour),
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		// # Redirect to the home page
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	})
}

// pwdMatch checks if the given password match the stored one
func pwdMatch(plainPwd string, pwdField string, saltField string) (bool, error) {
	hashedSavedPWD, err := base64.RawStdEncoding.DecodeString(pwdField)
	if err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(saltField)
	if err != nil {
		return false, err
	}
	hashedInputPwd := argon2.IDKey([]byte(plainPwd), salt, a2Time, a2Memory, a2Threads, a2KeyLen)

	if subtle.ConstantTimeCompare(hashedInputPwd, hashedSavedPWD) == 0 {
		return false, nil
	}

	return true, nil
}

func handleLogout() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value("user").(porgs.User)

		// # Delete all user sessions for this user,
		// # so she will be logged out from all devices.
		conn, err := porgs.DbConnPool.Take(r.Context())
		if err != nil {
			slog.Error("main.handleLogout: get db conn", "usr", user, "err", err)
			porgs.ShowDefaultErrorPage(w, r)
			return
		}
		defer porgs.DbConnPool.Put(conn)
		stmt, err := conn.Prepare("DELETE FROM session WHERE username = ?")
		if err != nil {
			slog.Error("main.handleLogout: stmt prepare", "usr", user, "err", err)
			porgs.ShowDefaultErrorPage(w, r)
			return
		}
		defer func() { _ = stmt.Reset(); _ = stmt.ClearBindings() }()
		stmt.BindText(1, user.Name)
		_, err = stmt.Step()
		if err != nil {
			slog.Error("main.handleLogout: stmt exec", "usr", user, "err", err)
			porgs.ShowDefaultErrorPage(w, r)
			return
		}

		// # Delete cookie
		cookie := http.Cookie{
			Name:     porgs.CookieNameSession,
			Path:     "/",
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		// # Redirect to root page
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
