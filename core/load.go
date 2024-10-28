package core

import (
	"encoding/csv"
	"fmt"
	"github.com/praja-dev/porgs"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var rxOrgFileName = regexp.MustCompile(`L\d-(\d+)-`)
var rxOrgName = regexp.MustCompile(`NAME_(.+)`)

func loadData() {
	loadDir, ok := porgs.Args["--load"]
	if !ok {
		return
	}
	if loadDir == "" {
		slog.Error("core.loadData", "err",
			"--load arg value should be the source directory path - e.g. --load=~/src/lk-data/admin")
		os.Exit(3)
	}
	slog.Info("core.loadData", "path", loadDir)

	// Check if loadDir is valid
	info, err := os.Stat(loadDir)
	if err != nil {
		slog.Error("core.loadData: check directory", "err", err)
		os.Exit(3)
	}
	if !info.IsDir() {
		slog.Error("core.loadData: check directory", "err", "not a directory")
		os.Exit(3)
	}
	slog.Info("core.loadData: check directory: ok")

	loadOrgs(loadDir)
}

func loadOrgs(directory string) {
	level := 0
	for {
		orgs, err := readOrgCSVsForLevel(directory, level)
		if err != nil {
			slog.Error("core.loadOrgs: read cvs files for 1 level", "level", level, "err", err)
			os.Exit(3)
		}
		slog.Info("core.loadOrgs: read cvs files for 1 level: ok", "level", level, "orgsCount", len(orgs))

		if level == 0 {
			if len(orgs) != 1 {
				slog.Error("core.loadOrgs: check root org", "err", "there can only be one level 0 (root) organization")
				os.Exit(3)
			} else {
				slog.Info("core.loadOrgs: check root org: ok")
			}
		}

		if len(orgs) == 0 {
			slog.Info("core.loadOrgs: ok", "highestLevel", level-1)
			return
		}

		for _, org := range orgs {
			err = SaveOrg(org)
			if err != nil {
				slog.Error("core.loadOrgs: save org", "level", level, "org", org, "err", err)
				os.Exit(3)
			}
		}
		slog.Info("core.loadOrgs: save: ok", "level", level, "orgsCount", len(orgs))

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
	fileName := filepath.Base(filePath)

	// # Get the organization type from the filename
	orgTypeID, err := getOrgTypeFromFileName(fileName)
	if err != nil {
		return nil, err
	}
	orgType, err := GetOrgType(porgs.Context, orgTypeID)
	if err != nil {
		return nil, err
	}

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
		return nil, fmt.Errorf("core.readOrgCSV %s: header: %w", fileName, err)
	}
	numFields := len(header)
	if numFields < 6 {
		return nil, fmt.Errorf("core.readOrgCSV %s: header: column count: 6 or more expected, only got: %d",
			fileName, numFields)
	}

	// # Validate header has the required 6 fields
	header[0] = strings.Trim(header[0], "\uFEFF") // Remove BOM
	if header[0] != "PID" {
		return nil, fmt.Errorf("core.readOrgCSV %s: header: coulmn 1: should be named PID, got \"%s\"", fileName, header[0])
	}
	if header[1] != "SID" {
		return nil, fmt.Errorf("core.readOrgCSV %s: header: coulmn 2: should be named SID, got \"%s\"", fileName, header[1])
	}
	if header[2] != "ID" {
		return nil, fmt.Errorf("core.readOrgCSV %s: header: coulmn 3: should be named ID, got \"%s\"", fileName, header[2])
	}
	if header[3] != "EID" {
		return nil, fmt.Errorf("core.readOrgCSV %s: header: coulmn 4: should be named EID, got \"%s\"", fileName, header[3])
	}
	if header[4] != "ESID" {
		return nil, fmt.Errorf("core.readOrgCSV %s: header: coulmn 5: should be named ESID, got \"%s\"", fileName, header[4])
	}
	if header[5] != "NAME" {
		return nil, fmt.Errorf("core.readOrgCSV %s: header: coulmn 6: should be named NAME, got \"%s\"", fileName, header[5])
	}

	// # Figure out the column indexes for the various translations of the NAME property
	// # e.g. NAME_SI, NAME_TA, NAME_FR, ...
	indexOfNameByLang := make(map[string]int)
	for i := 6; i < numFields; i++ {
		fld := header[i]
		if fld == "" {
			return nil, fmt.Errorf("core.readOrgCSV %s: header: column %d: name empty", fileName, i+1)
		}

		// TODO: Handle fields beyond NAME
		if !strings.HasPrefix(fld, "NAME_") {
			slog.Warn("core.readOrgCSV: header", "fileName",
				fileName, "column", i+1, "err", fmt.Sprintf("name unrecognized: %s", fld))
			//core.readOrgCSV
			continue
		}

		matches := rxOrgName.FindStringSubmatch(fld)
		if len(matches) != 2 {
			return nil, fmt.Errorf("core.readOrgCSV %s: header: column %d: name has invalid language suffix: %s",
				fileName, i+1, fld)
		}
		idx := matches[1]
		indexOfNameByLang[idx] = i
	}
	slog.Info("core.readOrgCSV: header: ok", "fileName", fileName, "indexOfNameByLang", indexOfNameByLang)

	// # Read the data
	var orgs []Org
	line := 1
	for {
		rec, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("core.readOrgCSV %s: data: row %d: %w", fileName, line, err)
		}
		// TODO: Compare with the actual number of columns found when parsing the header
		if len(rec) < 8 {
			return nil, fmt.Errorf("core.readOrgCSV %s: data: row %d: column count should be 8 or more, found %d",
				fileName, line, len(rec))
		}

		org := Org{TypeID: orgType.ID}
		pidVal := rec[0]
		if pidVal != "" {
			pid, err := strconv.Atoi(pidVal)
			if err != nil {
				return nil, fmt.Errorf("core.readOrgCSV %s: data: row %d: column %d: field PID: %w",
					fileName, line, 1, err)
			}
			org.ParentID = int64(pid)
		}

		sidVal := rec[1]
		sid, err := strconv.Atoi(sidVal)
		if err != nil {
			return nil, fmt.Errorf("core.readOrgCSV %s: data: row %d: column %d: field SID: %w",
				fileName, line, 2, err)
		}
		org.SequenceID = int64(sid)

		idVal := rec[2]
		id, err := strconv.Atoi(idVal)
		if err != nil {
			return nil, fmt.Errorf("core.readOrgCSV %s: data: row %d: column %d: field ID: %w",
				fileName, line, 3, err)
		}
		org.ID = int64(id)

		org.ExternalID = rec[3]
		org.ExternalSID = rec[4]

		org.Name = rec[5]
		if rec[5] == "" {
			return nil, fmt.Errorf("core.readOrgCSV %s: data: row %d: column %d: field NAME: is empty",
				fileName, line, 6)
		}

		trlx := make(map[string]OrgProps)
		for k, v := range indexOfNameByLang {
			trlx[strings.ToLower(k)] = OrgProps{Name: rec[v]}
		}
		trlx["en"] = OrgProps{Name: org.Name}
		org.Trlx = trlx

		orgs = append(orgs, org)
		line++
	}

	slog.Info("core.readOrgCSV: data: ok", "fileName", fileName, "count", len(orgs))
	return orgs, nil
}

func getOrgTypeFromFileName(fileName string) (int64, error) {
	matches := rxOrgFileName.FindStringSubmatch(fileName)
	if len(matches) < 2 {
		return 0, fmt.Errorf("invalid filename pattern: %s", fileName)
	}

	orgType, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("can't convert to integer: %s", matches[1])
	}

	return int64(orgType), nil
}
