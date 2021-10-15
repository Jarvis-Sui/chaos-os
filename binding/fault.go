package binding

import "time"

type FaultType string
type FaultStatus string

const (
	FT_UNSET    FaultType = ""
	FT_NETDELAY FaultType = "net_delay"
	FT_NETLOSS  FaultType = "net_loss"
)

const (
	FS_READY     FaultStatus = "Ready"
	FS_RUNNING   FaultStatus = "Running"
	FS_ERROR     FaultStatus = "Error"
	FS_DESTROYED FaultStatus = "Destoyed"
)

type Fault struct {
	Uid        string      `json:"id"`
	Type       FaultType   `json:"type"`
	Status     FaultStatus `json:"status"`
	Command    string      `json:"command"`
	Reason     string      `json:"reason"`
	CreateTime time.Time   `json:"create_time"`
	UpdateTime time.Time   `json:"update_time"`
}