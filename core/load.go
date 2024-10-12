package core

import (
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func loadData() {
	loadDir := os.Getenv("PORGS_LOAD_DIR")
	if loadDir == "" {
		return
	}
	slog.Info("core: load data", "PORGS_LOAD_DIR", loadDir)

	// Check if loadDir is valid
	info, err := os.Stat(loadDir)
	if err != nil {
		slog.Error("core: loadData: check directory",
			"PORGS_LOAD_DIR", loadDir, "error", err)
		os.Exit(3)
	}
	if !info.IsDir() {
		slog.Error("core: loadData: check directory",
			"PORGS_LOAD_DIR", loadDir, "error", "not a directory")
		os.Exit(3)
	}

	loadOrgs(loadDir)
}

func loadOrgs(directory string) {
	level := 0
	for {
		orgs, err := readOrgCSVsForLevel(directory, level)
		if err != nil {
			slog.Error("loadOrgs", "err", err)
			os.Exit(3)
		}

		if len(orgs) == 0 {
			return
		}

		if level == 0 && len(orgs) != 1 {
			slog.Error("loadOrgs", "err", "there can only be one level 0 (root) organization")
			os.Exit(3)
		}

		for _, org := range orgs {
			err = SaveOrg(org)
			if err != nil {
				slog.Error("loadOrgs", "err", err)
				os.Exit(3)
			}
		}

		level++
	}
}

func readOrgCSVsForLevel(directory string, level int) ([]Org, error) {
	matches, err := filepath.Glob(filepath.Join(directory, fmt.Sprintf("L%d-*.csv", level)))
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, nil
	}

	var orgs []Org
	for _, file := range matches {
		orgsInFile, err := readOrgCSV(file)
		if err != nil {
			return nil, err
		}
		orgs = append(orgs, orgsInFile...)
	}

	return orgs, nil
}

func readOrgCSV(filePath string) ([]Org, error) {
	// # Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	// # Read the header
	reader := csv.NewReader(file)
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("read header: %w", err)
	}
	numFields := len(header)
	if numFields < 6 {
		return nil, fmt.Errorf("at least 6 columns expected, got: %d", numFields)
	}

	// # Validate header has the required 6 fields
	header[0] = strings.Trim(header[0], "\uFEFF") // Remove BOM
	if header[0] != "PID" {
		return nil, fmt.Errorf("PID column expected as field #1, got: %s", header[0])
	}
	if header[1] != "SID" {
		return nil, fmt.Errorf("SID column expected as field #2, got: %s", header[1])
	}
	if header[2] != "ID" {
		return nil, fmt.Errorf("ID column expected as field #3, got: %s", header[2])
	}
	if header[3] != "EID" {
		return nil, fmt.Errorf("EID (ExternalID) column expected as field #4, got: %s", header[3])
	}
	if header[4] != "ESID" {
		return nil, fmt.Errorf("ESID (ExternalSID) column expected as field #5, got: %s", header[4])
	}
	if header[5] != "NAME" {
		return nil, fmt.Errorf("NAME column expected as field #6, got: %s", header[5])
	}

	// # Figure out the column indexes for the various translations of the NAME property
	// # e.g. NAME_SI, NAME_TA, NAME_FR, ...
	indexOfNameByLang := make(map[string]int)
	rxName := regexp.MustCompile(`NAME_(.+)`)
	for i := 6; i < numFields; i++ {
		fld := header[i]
		if fld == "" {
			return nil, fmt.Errorf("header: empty field Name at column %d", i+1)
		}

		// TODO: Handle fields beyond NAME
		if !strings.HasPrefix(fld, "NAME_") {
			slog.Warn("readOrgCSV: invalid field name", "col", i+1, "field", fld)
			continue
		}

		matches := rxName.FindStringSubmatch(fld)

		if len(matches) != 2 {
			return nil, fmt.Errorf("header: invalid field Name at column %d: %s", i+1, fld)
		}
		idx := matches[1]
		indexOfNameByLang[idx] = i
	}
	slog.Debug("indexOfNameByLang", "indexOfNameByLang", indexOfNameByLang)

	// # Read the data
	var orgs []Org
	line := 1
	for {
		rec, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("read line %d: %w", line, err)
		}
		if len(rec) < 8 {
			return nil, fmt.Errorf("8 columns expected, got: %d", len(rec))
		}

		org := Org{}
		pidVal := rec[0]
		if pidVal != "" {
			pid, err := strconv.Atoi(pidVal)
			if err != nil {
				return nil, fmt.Errorf("line %d: invalid PID: %w", line, err)
			}
			org.ParentID = int64(pid)
		}

		sidVal := rec[1]
		sid, err := strconv.Atoi(sidVal)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid SID: %w", line, err)
		}
		org.SequenceID = int64(sid)

		idVal := rec[2]
		id, err := strconv.Atoi(idVal)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid OID: %w", line, err)
		}
		org.ID = int64(id)

		org.ExternalID = rec[3]
		org.ExternalSID = rec[4]

		org.Name = rec[5]
		if rec[5] == "" {
			return nil, fmt.Errorf("name is required. line: %d, file: %s", line, filePath)
		}

		trlx := make(map[string]OrgProps)
		for k, v := range indexOfNameByLang {
			trlx[k] = OrgProps{Name: rec[v]}
		}
		org.Trlx = trlx

		orgs = append(orgs, org)
		line++
	}

	return orgs, nil
}
