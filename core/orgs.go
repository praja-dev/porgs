package core

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/praja-dev/porgs"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

func handleOrgs(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		orgs, err := GetOrgs(ctx)
		if err != nil {
			porgs.ShowDefaultErrorPage(w, r)
			return
		}

		if len(orgs) == 0 {
			porgs.RenderView(w, r, porgs.View{Name: "core-orgs", Title: "Orgs", Data: nil})
			return
		}

		porgs.RenderView(w, r, porgs.View{Name: "core-orgs", Title: "Orgs", Data: orgs})
	})
}

func handleOrg(ctx context.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Error("handleOrg: can't parse ID", "id", idStr, "err", err)
			porgs.ShowDefaultErrorPage(w, r)
			return
		}

		org, err := GetOrg(ctx, int64(id))
		if err != nil {
			if errors.Is(err, porgs.ErrNotFound) {
				porgs.ShowErrorPage(w, r, porgs.ErrorPage{
					Msg:     fmt.Sprintf("There is no organization with an id of %d", id),
					BackURL: "/core/orgs",
					Title:   "Not Found",
				})
			} else {
				porgs.ShowDefaultErrorPage(w, r)
			}
			return
		}

		porgs.RenderView(w, r, porgs.View{Name: "core-org", Title: "Org", Data: org})
	})
}

func GetOrgs(ctx context.Context) ([]Org, error) {
	conn, err := porgs.DbConnPool.Take(ctx)
	if err != nil {
		return nil, err
	}
	defer porgs.DbConnPool.Put(conn)

	stmt, err := conn.Prepare(`SELECT id, created, updated,
       parent, sid, source, type, external_id, external_sid,
       name, trlx
		FROM org ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = stmt.Reset()
		if err != nil {
			slog.Error("GetOrgs: stmt reset", "err", err)
		}
	}()

	var orgs []Org
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return nil, err
		}
		if !hasRow {
			break
		}

		org := Org{
			ID:          stmt.GetInt64("id"),
			Created:     stmt.GetInt64("created"),
			Updated:     stmt.GetInt64("updated"),
			ParentID:    stmt.GetInt64("parent"),
			SequenceID:  stmt.GetInt64("sid"),
			Source:      stmt.GetInt64("source"),
			Type:        stmt.GetInt64("type"),
			ExternalID:  stmt.GetText("external_id"),
			ExternalSID: stmt.GetText("external_sid"),
			Name:        stmt.GetText("name"),
		}

		trlxJSON := stmt.GetText("trlx")
		err = json.Unmarshal([]byte(trlxJSON), &org.Trlx)
		if err != nil {
			slog.Error("GetOrgs: unmarshal Trlx", "err", err, "trlx", trlxJSON)
			return nil, err
		}

		orgs = append(orgs, org)
	}

	return orgs, nil
}

func GetOrg(ctx context.Context, id int64) (Org, error) {
	conn, err := porgs.DbConnPool.Take(ctx)
	if err != nil {
		return Org{}, err
	}
	defer porgs.DbConnPool.Put(conn)

	stmt, err := conn.Prepare(`SELECT id, created, updated,
	   parent, sid, source, type, external_id, external_sid,
	   name, trlx
		FROM org WHERE id=?`)
	if err != nil {
		return Org{}, err
	}
	defer func() {
		err = stmt.Reset()
		if err != nil {
			slog.Error("GetOrg: stmt reset", "err", err)
		}
	}()

	stmt.BindInt64(1, id)

	hasRow, err := stmt.Step()
	if err != nil {
		return Org{}, err
	}
	if !hasRow {
		return Org{}, porgs.ErrNotFound
	}

	org := Org{
		ID:          stmt.GetInt64("id"),
		Created:     stmt.GetInt64("created"),
		Updated:     stmt.GetInt64("updated"),
		ParentID:    stmt.GetInt64("parent"),
		SequenceID:  stmt.GetInt64("sid"),
		Source:      stmt.GetInt64("source"),
		Type:        stmt.GetInt64("type"),
		ExternalID:  stmt.GetText("external_id"),
		ExternalSID: stmt.GetText("external_sid"),
		Name:        stmt.GetText("name"),
	}

	trlxJSON := stmt.GetText("trlx")
	err = json.Unmarshal([]byte(trlxJSON), &org.Trlx)
	if err != nil {
		slog.Error("GetOrg: unmarshal Trlx", "err", err, "trlx", trlxJSON)
		return Org{}, err
	}

	return org, nil
}

// SaveOrg saves an org to the database.
func SaveOrg(org Org) error {
	if org.ID == 0 {
		return porgs.ErrNotImplemented
	}

	conn, err := porgs.DbConnPool.Take(context.Background())
	if err != nil {
		return err
	}
	defer porgs.DbConnPool.Put(conn)

	stmt, err := conn.Prepare(`INSERT INTO org
		(created, updated, parent, sid, source, type, external_id, external_sid,
		 name, trlx, id) VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	defer func() {
		err = stmt.Reset()
		if err != nil {
			slog.Error("SaveOrg: stmt reset", "err", err)
		}
		err = stmt.ClearBindings()
		if err != nil {
			slog.Error("SaveOrg: stmt clear bindings", "err", err)
		}
	}()

	now := time.Now().Unix()
	stmt.BindInt64(1, now) // created
	stmt.BindInt64(2, now) // updated
	stmt.BindInt64(3, org.ParentID)
	stmt.BindInt64(4, org.SequenceID)
	stmt.BindInt64(5, 0)
	stmt.BindInt64(6, 1)
	stmt.BindText(7, org.ExternalID)
	stmt.BindText(8, org.ExternalSID)
	stmt.BindText(9, org.Name)

	trlxJSON, err := json.Marshal(org.Trlx)
	if err != nil {
		slog.Error("SaveOrg: marshal Trlx", "err", err, "Trlx", org.Trlx)
		return err
	}
	stmt.BindText(10, string(trlxJSON))

	stmt.BindInt64(11, org.ID)

	_, err = stmt.Step()
	if err != nil {
		return err
	}

	return nil
}
