package main

import (
	"fmt"
	"github.com/eatonphil/gosqlite"
	"github.com/praja-dev/porgs"
)

func findUserBySessionToken(token string) (porgs.User, error) {
	// Create db connection
	conn, err := gosqlite.Open(porgs.BootConfig.DSN)
	if err != nil {
		return porgs.User{}, err
	}
	defer func() { _ = conn.Close() }()
	stmt, err := conn.Prepare("SELECT username FROM session where id=?", token)
	if err != nil {
		return porgs.User{}, err
	}
	defer func() { _ = stmt.Close() }()
	hasRow, err := stmt.Step()
	if err != nil {
		return porgs.User{}, err
	}
	if !hasRow {
		return porgs.User{}, fmt.Errorf("invalid session token")
	}
	var username string
	err = stmt.Scan(&username)
	if err != nil {
		return porgs.User{}, err
	}

	return porgs.User{
		Name: username,
	}, nil
}
