package core

import (
	"context"
	"encoding/json"
	"github.com/praja-dev/porgs"
	"log/slog"
	"net/http"
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

func GetOrgs(_ context.Context) ([]Org, error) {
	var orgs []Org
	return orgs, nil
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
