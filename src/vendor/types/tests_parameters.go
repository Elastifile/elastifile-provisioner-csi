package types

import (
	"strings"

	"helputils"
	"optional"
)

type TestParameters struct {
	Tools                 Tools            `yaml:"Tools,omitempty"`
	Files                 Files            `yaml:"Files,omitempty"`
	Timing                Timing           `yaml:"Timing,omitempty"`
	Performance           Performance      `yaml:"Performance,omitempty"`
	Snapshot              Snapshot         `yaml:"Snapshot,omitempty"`
	Iterations            optional.Int     `yaml:"Iterations,omitempty"`
	MaxSystemUsagePersent optional.Float64 `yaml:"MaxSystemUsagePersent,omitempty"`
	GoalCapacity          optional.Int     `yaml:"GoalCapacity,omitempty"`
	DcOptions             DcOptions        `yaml:"DcOptions,omitempty"`
	AsyncDr               AsyncDr          `yaml:"AsyncDr,omitempty"`
}

func (s TestParameters) ToString() string {
	return "{" + strings.Join(helputils.StructToStrings(s), " ") + "}"
}

type Timing struct {
	Timeout optional.String `yaml:"Timeout"`
	Ticker  optional.String `yaml:"Ticker"`
}

type Files struct {
	MaxFileSize           optional.String `yaml:"MaxFileSize"`
	MinFileSize           optional.String `yaml:"MinFileSize"`
	IncrementFileSize     optional.String `yaml:"IncrementFileSize"`
	NumberOfFilesToCreate optional.Int    `yaml:"NumberOfFilesToCreate"`
	TruncateSize          optional.String `yaml:"TruncateSize"`
}

type Snapshot struct {
	CreateCount                  optional.Int    `yaml:"CreateCount"`
	DeleteCount                  optional.Int    `yaml:"DeleteCount"`
	DelayBetweenOperations       optional.String `yaml:"DelayBetweenOperations"`
	DelayBetweenCreateOperations optional.String `yaml:"DelayBetweenCreateOperations"`
	DelayBetweenDeleteOperations optional.String `yaml:"DelayBetweenDeleteOperations"`
}

type AsyncDr struct {
	ChangeRate                optional.Int    `yaml:"ChangeRate"`
	WaitBeforeErunReplication optional.String `yaml:"WaitBeforeErunReplication"`
	DualReplication           optional.Bool   `yaml:"DualReplication"`
	ErunBackgroundIO          []optional.Bool `yaml:"ErunBackgroundIO"`
	DualReplicationCycle      []optional.Bool `yaml:"DualReplicationCycle"`
	ReplicationFailedCount    optional.Int    `yaml:"ReplicationFailedCount"`
	ErunOptions               ErunParams      `yaml:"ErunOptions"`
	DCOpts                    DcOptions       `yaml:"DCOpts"`
	NumberOfCycles            optional.Int    `yaml:"NumberOfCycles"`
	DCPairs                   DCPairs         `yaml:"DCPairs"`
}

type DCPairs struct {
	Count         optional.Int      `yaml:"Count"`
	Rpos          []optional.Int    `yaml:"Rpos"`
	Capacities    []optional.String `yaml:"Capacities"`
	NumberOfFiles []optional.Int    `yaml:"NumberOfFiles"`
}

type Performance struct {
	Iterations   optional.Int     `yaml:"Iterations"`
	InWorkingSet PerformanceCycle `yaml:"InWorkingSet"`
	AbovePstore  PerformanceCycle `yaml:"AbovePstore"`
	AboveDstore  PerformanceCycle `yaml:"AboveDstore"`
	AboveMD      PerformanceCycle `yaml:"AboveMD"`
}

type PerformanceCycle struct {
	Tools Tools `yaml:"Tools"`
}

type Tools struct {
	Erun    ErunParams    `yaml:"Erun"`
	Sfs2008 Sfs2008Params `yaml:"Sfs2008"`
}
type ErunParams struct {
	Clients           optional.Int      `yaml:"Clients"`
	QueueSize         optional.Int      `yaml:"QueueSize"`
	Readwrites        optional.String   `yaml:"Readwrites"`
	ReadwritesList    []optional.String `yaml:"ReadwritesList"`
	MaxFileSize       optional.String   `yaml:"MaxFileSize"`
	Duration          optional.String   `yaml:"Duration"`
	NumberOfFiles     optional.Int      `yaml:"NumberOfFiles"`
	NumberOfFilesList []optional.Int    `yaml:"NumberOfFilesList"`
	RecoveryTimeout   optional.String   `yaml:"RecoveryTimeout"`
	MaxIoSize         optional.String   `yaml:"MaxIoSize"`
	MinIoSize         optional.String   `yaml:"MinIoSize"`
	DataPayload       optional.Bool     `yaml:"DataPayload"`
	InitialWritePhase optional.Bool     `yaml:"InitialWritePhase"`
	InitialWriteStop  optional.Bool     `yaml:"InitialWriteStop"`
	MinUncomp         optional.Int      `yaml:"MinUncomp"`
	MaxUncomp         optional.Int      `yaml:"MaxUncomp"`
}

type Sfs2008Params struct {
	Runtime            optional.Int `yaml:"Runtime"`
	WarmupTime         optional.Int `yaml:"WarmupTime"`
	ProcessesPerClient optional.Int `yaml:"ProcessesPerClient"`
	IOPs               optional.Int `yaml:"IOPs"`
}

type DcOptions struct {
	Compression optional.Int `yaml:"Compression"`
	Dedup       optional.Int `yaml:"Dedup"`
}
