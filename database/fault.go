package database

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Jarvis-Sui/chaos-os/binding"
	"github.com/Jarvis-Sui/chaos-os/util"
	"github.com/araddon/dateparse"
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
	table = FaultTable{dbFile, tableName, nil}

	if _, err := os.Stat(dbFile); err != nil {
		if _, err := os.Create(dbFile); err != nil {
			logrus.WithField("err", err).Errorf("failed to create db file: %s. ", dbFile)
			os.Exit(1)
		}
	} else {
		return // database file already exists
	}
	table.Open()
	defer table.Close()

	_, err := table.conn.Exec(fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s(
		id  			TEXT,
		action 			TEXT,
		fault_type 		TEXT,
		command 		TEXT,
		status  		TEXT,
		reason 			TEXT,
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
		CREATE INDEX IF NOT EXISTS %[1]s_action_idx ON %[1]s (action);
	`, tableName)); err != nil {
		logrus.WithField("err", err).Errorf("failed to create index %[1]s_action_idx for %[1]s", tableName)
		os.Exit(1)
	} else {
		logrus.Infof("created index %[1]s_action_idx for %[1]s", tableName)
	}

	if _, err := table.conn.Exec(fmt.Sprintf(`
		CREATE INDEX IF NOT EXISTS %[1]s_status_idx ON %[1]s (status);
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

	sql := fmt.Sprintf(`
		INSERT INTO %s
		(id, action, fault_type, command, status, create_time, update_time)
		VALUES
		(?, ?, ?, ?, ?, ?, ?)
	`, ft.TableName)

	stmt, err := ft.conn.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(fault.Uid, "Create", fault.Type, fault.Command, fault.Status, fault.CreateTime, time.Now())
	return err
}

func (ft *FaultTable) UpdateFaultStatus(uid string, status string, reason string) error {
	ft.Open()
	defer ft.Close()

	sql := fmt.Sprintf(`UPDATE %s SET status=?, reason=?, update_time=? WHERE id=?`, ft.TableName)

	stmt, err := ft.conn.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(status, reason, time.Now(), uid)
	return err
}

func (ft *FaultTable) GetFaultById(uid string) (*binding.Fault, error) {
	ft.Open()
	defer ft.Close()

	sql := fmt.Sprintf(`SELECT id, fault_type, command, status, reason, create_time, update_time FROM %s WHERE id=?`, ft.TableName)
	stmt, err := ft.conn.Prepare(sql)
	if err != nil {
		return nil, err
	}

	if rows, err := stmt.Query(uid); err != nil {
		return nil, err
	} else {
		defer rows.Close()

		if rows.Next() {
			fault, err := rowToFault(rows)
			return fault, err
		} else {
			return nil, nil
		}
	}
}

func (ft *FaultTable) GetFaults(id string, status binding.FaultStatus, limit int) []*binding.Fault {
	ft.Open()
	defer ft.Close()

	sql := fmt.Sprintf(`SELECT id, fault_type, command, status, reason, create_time, update_time FROM %s`, ft.TableName)

	if id != "" {
		sql += fmt.Sprintf(" WHERE id='%s'", id)
	} else {
		if status != binding.FS_UNSET {
			sql += fmt.Sprintf(" WHERE status='%s'", status)
		}
		sql += " ORDER BY create_time DESC"
		if limit != 0 {
			sql += fmt.Sprintf(" LIMIT %d", limit)
		}
	}

	faults := make([]*binding.Fault, 0)
	if rows, err := ft.conn.Query(sql); err != nil {
		return faults
	} else {
		for rows.Next() {
			fault, _ := rowToFault(rows)
			faults = append(faults, fault)
		}
		return faults
	}
}

func rowToFault(rows *sql.Rows) (*binding.Fault, error) {
	var uid string
	var ftype binding.FaultType
	var status binding.FaultStatus
	var command string
	var reason string
	var createTime string
	var updateTime string

	if err := rows.Scan(&uid, &ftype, &command, &status, &reason, &createTime, &updateTime); err != nil {
		return nil, err
	}

	fault := binding.Fault{
		Uid:        uid,
		Type:       ftype,
		Status:     status,
		Command:    command,
		Reason:     reason,
		CreateTime: dateparse.MustParse(createTime),
		UpdateTime: dateparse.MustParse(updateTime),
	}

	return &fault, nil
}
