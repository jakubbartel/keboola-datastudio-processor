package keboola

import "os"

var (
	DataDir     = os.Getenv("KBC_DATADIR")
	RunID       = os.Getenv("KBC_RUNID")
	ProjectID   = os.Getenv("KBC_PROJECTID")
	StackID     = os.Getenv("KBC_STACKID")
	ConfigID    = os.Getenv("KBC_CONFIGID")
	ComponentID = os.Getenv("KBC_COMPONENTID")

	InFilesDir   = DataDir + "in/files/"
	OutFilesDir  = DataDir + "out/files/"
	InTablesDir  = DataDir + "in/tables/"
	OutTablesDir = DataDir + "out/tables/"
)
