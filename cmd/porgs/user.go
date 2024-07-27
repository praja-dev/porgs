package main

import (
	"context"
	"fmt"
	"github.com/praja-dev/porgs"
)

func findUserBySessionToken(token string) (porgs.User, error) {
	conn, err := porgs.DbConnPool.Take(context.Background())
	if err != nil {
		return porgs.User{}, err
	}
	defer porgs.DbConnPool.Put(conn)

	stmt, err := conn.Prepare("SELECT username FROM session where id=?")
	if err != nil {
		return porgs.User{}, err
	}
	defer func() { _ = stmt.ClearBindings(); _ = stmt.Reset() }()

	stmt.BindText(1, token)
	hasRow, err := stmt.Step()
	if err != nil {
		return porgs.User{}, err
	}
	if !hasRow {
		return porgs.User{}, fmt.Errorf("invalid session token")
	}
	username := stmt.GetText("username")

	return porgs.User{
		Name: username,
	}, nil
}
