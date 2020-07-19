package testutil

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-testfixtures/testfixtures/v3"
)

// DataTester provide helper method to load test data
type DataTester struct {
	TestData map[string][]byte
	DB       *sql.DB
	DBType   string
}

// SetupDB setup test db
func (tester *DataTester) SetupDB(db *sql.DB, dbType string) error {
	if err := db.Ping(); err != nil {
		return err
	}
	tester.DB = db
	tester.DBType = dbType
	return nil
}

// LoadTestData from join files from given directory path
func (tester *DataTester) LoadTestData(dataDir string) error {
	return filepath.Walk(dataDir, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filename := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))

			data, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			tester.TestData[filename] = data
		}
		return nil
	})
}

// UnmarshalTestData unmarshal testdata to slice of instance
func (tester DataTester) UnmarshalTestData(dataName string, targer interface{}) error {
	return json.Unmarshal(tester.TestData[dataName], targer)
}

// LoadFixtures setup test fixtures
func (tester *DataTester) LoadFixtures(fixtureDir string) error {
	if tester.DB == nil || tester.DBType == "" {
		return errors.New("should setup data and its type first")
	}

	fixture, err := testfixtures.New(
		testfixtures.Database(tester.DB),
		testfixtures.Dialect(tester.DBType),
		testfixtures.Directory(fixtureDir))
	if err != nil {
		return err
	}
	return fixture.Load()
}

// ClearDB clear database
func (tester DataTester) ClearDB(db *sql.DB, tables ...tabler) error {
	for _, table := range tables {
		tableName := strings.Split(table.TableName(), ";")[0]

		_, err := db.Exec("DELETE FROM " + tableName)
		if err != nil {
			return err
		}
	}
	return nil
}

type tabler interface {
	TableName() string
}
