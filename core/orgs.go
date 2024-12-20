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

func handleOrgs() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		org, err := GetOrg(r.Context(), 1)
		if err != nil {
			porgs.ShowErrorPage(w, r, porgs.ErrorPage{
				Msg:     "There are no defined organizations. Please add one.",
				BackURL: "/home",
				Title:   "No Orgs",
			})
			return
		}

		porgs.RenderView(w, r, porgs.View{Name: "core-org", Title: "Orgs", Data: org})
	})
}

func handleOrg() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			slog.Error("core.handleOrg: parse id", "id", idStr, "err", err)
			porgs.ShowDefaultErrorPage(w, r)
			return
		}

		org, err := GetOrg(r.Context(), int64(id))
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

func GetSubOrgs(ctx context.Context, id int64) ([]Org, error) {
	conn, err := porgs.DbConnPool.Take(ctx)
	if err != nil {
		return nil, err
	}
	defer porgs.DbConnPool.Put(conn)

	stmt, err := conn.Prepare(`SELECT id, created, updated,
       parent, sid, source, type, external_id, external_sid,
       name, trlx
		FROM org WHERE parent = ? ORDER BY sid`)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = stmt.Reset()
		if err != nil {
			slog.Error("core.GetSubOrgs: stmt reset", "err", err)
		}
		err = stmt.ClearBindings()
		if err != nil {
			slog.Error("core.GetSubOrgs: stmt clear", "err", err)
		}
	}()

	stmt.BindInt64(1, id)

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
			TypeID:      stmt.GetInt64("type"),
			ExternalID:  stmt.GetText("external_id"),
			ExternalSID: stmt.GetText("external_sid"),
			Name:        stmt.GetText("name"),
		}

		trlxJSON := stmt.GetText("trlx")
		err = json.Unmarshal([]byte(trlxJSON), &org.Trlx)
		if err != nil {
			slog.Error("core.GetSubOrgs: unmarshal trlx", "trlx", trlxJSON, "err", err)
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
			slog.Error("core.GetOrg: stmt reset", "err", err)
		}
		err = stmt.ClearBindings()
		if err != nil {
			slog.Error("core.GetOrg: stmt clear", "err", err)
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
		TypeID:      stmt.GetInt64("type"),
		ExternalID:  stmt.GetText("external_id"),
		ExternalSID: stmt.GetText("external_sid"),
		Name:        stmt.GetText("name"),
	}

	trlxJSON := stmt.GetText("trlx")
	err = json.Unmarshal([]byte(trlxJSON), &org.Trlx)
	if err != nil {
		slog.Error("core.GetOrg: unmarshal trlx", "trlx", trlxJSON, "err", err)
		return Org{}, err
	}

	subOrgs, err := GetSubOrgs(ctx, id)
	if err != nil {
		return Org{}, err
	}
	org.SubOrgs = subOrgs

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
			slog.Error("core.SaveOrg: stmt reset", "err", err)
		}
		err = stmt.ClearBindings()
		if err != nil {
			slog.Error("core.SaveOrg: stmt clear", "err", err)
		}
	}()

	now := time.Now().Unix()
	stmt.BindInt64(1, now) // created
	stmt.BindInt64(2, now) // updated
	stmt.BindInt64(3, org.ParentID)
	stmt.BindInt64(4, org.SequenceID)
	stmt.BindInt64(5, 0)
	stmt.BindInt64(6, org.TypeID)
	stmt.BindText(7, org.ExternalID)
	stmt.BindText(8, org.ExternalSID)
	stmt.BindText(9, org.Name)

	trlxJSON, err := json.Marshal(org.Trlx)
	if err != nil {
		slog.Error("core.SaveOrg: marshal trlx", "trlx", org.Trlx, "err", err)
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

// GetOrgType retrieves the org_type with the given ID
func GetOrgType(ctx context.Context, id int64) (OrgType, error) {
	conn, err := porgs.DbConnPool.Take(ctx)
	if err != nil {
		return OrgType{}, err
	}
	defer porgs.DbConnPool.Put(conn)

	stmt, err := conn.Prepare(`SELECT id, created, updated, name, description, cxp FROM org_type WHERE id=?`)
	if err != nil {
		return OrgType{}, err
	}
	defer func() {
		err = stmt.Reset()
		if err != nil {
			slog.Error("core.GetOrgType: stmt reset", "err", err)
		}
		err = stmt.ClearBindings()
		if err != nil {
			slog.Error("core.GetOrgType: stmt clear", "err", err)
		}
	}()

	stmt.BindInt64(1, id)

	hasRow, err := stmt.Step()
	if err != nil {
		return OrgType{}, err
	}
	if !hasRow {
		return OrgType{}, porgs.ErrNotFound
	}

	orgType := OrgType{
		ID:          id,
		Created:     stmt.GetInt64("created"),
		Updated:     stmt.GetInt64("updated"),
		Name:        stmt.GetText("name"),
		Description: stmt.GetText("description"),
	}

	cxpJSON := stmt.GetText("cxp")
	err = json.Unmarshal([]byte(cxpJSON), &orgType.XProps)
	if err != nil {
		slog.Error("core.GetOrgType: unmarshal cxp", "cxp", cxpJSON, "err", err)
		return OrgType{}, err
	}

	return orgType, nil
}
