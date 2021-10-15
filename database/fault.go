package database

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/Jarvis-Sui/chaos-os/binding"
	"github.com/Jarvis-Sui/chaos-os/util"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

type FaultTable struct {
	dbFile    string
	TableName string
	conn      *sql.DB
}

var once sync.Once
var table FaultTable

func GetFaultTable() *FaultTable {
	once.Do(initFaultTable)
	return &table
}

func initFaultTable() {
	dbFile := util.GetDBFilePath()
	tableName := "fault"

	if _, err := os.Stat(dbFile); err != nil {
		if _, err := os.Create(dbFile); err != nil {
			logrus.WithField("err", err).Errorf("failed to create db file: %s. ", dbFile)
			os.Exit(1)
		}
	}

	table = FaultTable{dbFile, tableName, nil}
	table.Open()
	defer table.Close()

	_, err := table.conn.Exec(fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s(
		id  			TEXT,
		command 		TEXT,
		fault_type 		TEXT,
		params 			TEXT,
		status  		TEXT,
		create_time  	TEXT,
		update_time  	TEXT
	)`, tableName), tableName)

	if err != nil {
		logrus.WithField("err", err).Errorf("failed to create table %s", tableName)
		os.Exit(1)
	}

	if _, err := table.conn.Exec(fmt.Sprintf(`
		CREATE UNIQUE INDEX IF NOT EXISTS %[1]s_id_idx ON %[1]s (id);
	`, tableName)); err != nil {
		logrus.WithField("err", err).Errorf("failed to create index %[1]s_id_idx for %[1]s", tableName)
		os.Exit(1)
	} else {
		logrus.Infof("created index %[1]s_id_idx for %[1]s", tableName)
	}

	if _, err := table.conn.Exec(fmt.Sprintf(`
		CREATE UNIQUE INDEX IF NOT EXISTS %[1]s_command_idx ON %[1]s (command);
	`, tableName)); err != nil {
		logrus.WithField("err", err).Errorf("failed to create index %[1]s_command_idx for %[1]s", tableName)
		os.Exit(1)
	} else {
		logrus.Infof("created index %[1]s_command_idx for %[1]s", tableName)
	}

	if _, err := table.conn.Exec(fmt.Sprintf(`
		CREATE UNIQUE INDEX IF NOT EXISTS %[1]s_status_idx ON %[1]s (status);
	`, tableName)); err != nil {
		logrus.WithField("err", err).Errorf("failed to create index %[1]s_status_idx for %[1]s", tableName)
		os.Exit(1)
	} else {
		logrus.Infof("created index %[1]s_status_idx for %[1]s", tableName)
	}
}

func (ft *FaultTable) Open() {
	conn, err := sql.Open("sqlite3", ft.dbFile)
	if err != nil {
		logrus.WithField("err", err).Errorf("failed to open db file %s", ft.dbFile)
		os.Exit(1)
	}
	ft.conn = conn
}

func (ft *FaultTable) Close() {
	if ft.conn != nil {
		if err := ft.conn.Close(); err == nil {
			ft.conn = nil
		} else {
			logrus.WithField("err", err).Errorf("failed to close db connection %s", ft.dbFile)
			os.Exit(1)
		}
	}
}

func (ft *FaultTable) AddFault(fault *binding.Fault) error {
	ft.Open()
	defer ft.Close()
	return nil
}

func (ft *FaultTable) UpdateFaultStatus(uid string, status string, reason string) error {
	return nil
}

func (ft *FaultTable) GetFaultById(uid string) {

}
