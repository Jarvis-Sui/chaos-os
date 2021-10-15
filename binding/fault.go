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
	FS_RUNNING   FaultStatus = "Running"
	FS_READY     FaultStatus = "Ready"
	FS_ERROR     FaultStatus = "Error"
	FS_DESTROYED FaultStatus = "Destoyed"
)

type Fault struct {
	Uid        string
	Type       FaultType
	Status     FaultStatus
	Command    string
	CreateTime time.Time
}
