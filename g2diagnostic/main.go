package g2diagnostic

import (
	"context"

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-observing/observer"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// The G2diagnostic interface is a Golang representation of Senzing's libg2diagnostic.h
type G2diagnostic interface {
	CheckDBPerf(ctx context.Context, secondsToRun int) (string, error)
	CloseEntityListBySize(ctx context.Context, entityListBySizeHandle uintptr) error
	Destroy(ctx context.Context) error
	FetchNextEntityBySize(ctx context.Context, entityListBySizeHandle uintptr) (string, error)
	FindEntitiesByFeatureIDs(ctx context.Context, features string) (string, error)
	GetAvailableMemory(ctx context.Context) (int64, error)
	GetDataSourceCounts(ctx context.Context) (string, error)
	GetDBInfo(ctx context.Context) (string, error)
	GetEntityDetails(ctx context.Context, entityID int64, includeInternalFeatures int) (string, error)
	GetEntityListBySize(ctx context.Context, entitySize int) (uintptr, error)
	GetEntityResume(ctx context.Context, entityID int64) (string, error)
	GetEntitySizeBreakdown(ctx context.Context, minimumEntitySize int, includeInternalFeatures int) (string, error)
	GetFeature(ctx context.Context, libFeatID int64) (string, error)
	GetGenericFeatures(ctx context.Context, featureType string, maximumEstimatedCount int) (string, error)
	GetLogicalCores(ctx context.Context) (int, error)
	GetMappingStatistics(ctx context.Context, includeInternalFeatures int) (string, error)
	GetPhysicalCores(ctx context.Context) (int, error)
	GetRelationshipDetails(ctx context.Context, relationshipID int64, includeInternalFeatures int) (string, error)
	GetResolutionStatistics(ctx context.Context) (string, error)
	GetTotalSystemMemory(ctx context.Context) (int64, error)
	Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error
	InitWithConfigID(ctx context.Context, moduleName string, iniParams string, initConfigID int64, verboseLogging int) error
	RegisterObserver(ctx context.Context, observer observer.Observer) error
	Reinit(ctx context.Context, initConfigID int64) error
	SetLogLevel(ctx context.Context, logLevel logger.Level) error
	UnregisterObserver(ctx context.Context, observer observer.Observer) error
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2diagnostic package found messages having the format "senzing-6003xxxx".
const ProductId = 6003

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2diagnostic package.
var IdMessages = map[int]string{
	1:    "Enter CheckDBPerf(%d).",
	2:    "Exit  CheckDBPerf(%d) returned (%s, %v).",
	3:    "Enter ClearLastException().",
	4:    "Exit  ClearLastException() returned (%v).",
	5:    "Enter CloseEntityListBySize().",
	6:    "Exit  CloseEntityListBySize() returned (%v).",
	7:    "Enter Destroy().",
	8:    "Exit  Destroy() returned (%v).",
	9:    "Enter FetchNextEntityBySize().",
	10:   "Exit  FetchNextEntityBySize() returned (%s, %v).",
	11:   "Enter FindEntitiesByFeatureIDs(%s).",
	12:   "Exit  FindEntitiesByFeatureIDs(%s) returned (%s, %v).",
	13:   "Enter GetAvailableMemory().",
	14:   "Exit  GetAvailableMemory() returned (%d, %v).",
	15:   "Enter GetDataSourceCounts().",
	16:   "Exit  GetDataSourceCounts() returned (%s, %v).",
	17:   "Enter GetDBInfo().",
	18:   "Exit  GetDBInfo()  returned (%s, %v).",
	19:   "Enter GetEntityDetails(%d, %d).",
	20:   "Exit  GetEntityDetails(%d, %d) returned (%s, %v).",
	21:   "Enter GetEntityListBySize(%d).",
	22:   "Exit  GetEntityListBySize(%d) returned (%v, %v).",
	23:   "Enter GetEntityResume(%d).",
	24:   "Exit  GetEntityResume(%d) returned (%s, %v).",
	25:   "Enter GetEntitySizeBreakdown(%d, %d).",
	26:   "Exit  GetEntitySizeBreakdown(%d, %d) returned (%s, %v).",
	27:   "Enter GetFeature(%d).",
	28:   "Exit  GetFeature(%d) returned (%s, %v).",
	29:   "Enter GetGenericFeatures(%s, %d).",
	30:   "Exit  GetGenericFeatures(%s, %d) returned (%s, %v).",
	31:   "Enter GetLastException().",
	32:   "Exit  GetLastException() returned (%s, %v).",
	33:   "Enter GetLastExceptionCode().",
	34:   "Exit  GetLastExceptionCode() returned (%d, %v).",
	35:   "Enter GetLogicalCores().",
	36:   "Exit  GetLogicalCores() returned (%d, %v).",
	37:   "Enter GetMappingStatistics(%d).",
	38:   "Exit  GetMappingStatistics(%d) returned (%s, %v).",
	39:   "Enter GetPhysicalCores().",
	40:   "Exit  GetPhysicalCores() returned (%d, %v).",
	41:   "Enter GetRelationshipDetails(%d, %d).",
	42:   "Exit  GetRelationshipDetails(%d, %d) returned (%s, %v).",
	43:   "Enter GetResolutionStatistics().",
	44:   "Exit  GetResolutionStatistics() returned (%s, %v).",
	45:   "Enter GetTotalSystemMemory().",
	46:   "Exit  GetTotalSystemMemory() returned (%d, %v).",
	47:   "Enter Init(%s, %s, %d).",
	48:   "Exit  Init(%s, %s, %d) returned (%v).",
	49:   "Enter InitWithConfigID(%s, %s, %d, %d).",
	50:   "Exit  InitWithConfigID(%s, %s, %d, %d) returned (%v).",
	51:   "Enter Reinit(%d).",
	52:   "Exit  Reinit(%d) returned (%v).",
	53:   "Enter SetLogLevel(%v).",
	54:   "Exit  SetLogLevel(%v) returned (%v).",
	4001: "Call to G2Diagnostic_checkDBPerf(%d) failed. Return code: %d",
	4002: "Call to G2Diagnostic_closeEntityListBySize() failed. Return code: %d",
	4003: "Call to G2Diagnostic_destroy() failed.  Return code: %d",
	4004: "Call to G2Diagnostic_fetchNextEntityBySize() failed.  Return code: %d",
	4005: "Call to G2Diagnostic_findEntitiesByFeatureIDs(%s) failed. Return code: %d",
	4006: "Call to G2Diagnostic_getDataSourceCounts() failed. Return code: %d",
	4007: "Call to G2Diagnostic_getDBInfo() failed. Return code: %d",
	4008: "Call to G2Diagnostic_getEntityDetails(%d, %d) failed. Return code: %d",
	4009: "Call to G2Diagnostic_getEntityListBySize(%d) failed. Return code: %d",
	4010: "Call to G2Diagnostic_getEntityResume(%d) failed. Return code: %d",
	4011: "Call to G2Diagnostic_getEntitySizeBreakdown(%d, %d) failed. Return code: %d",
	4012: "Call to G2Diagnostic_getFeature(%d) failed. Return code: %d",
	4013: "Call to G2Diagnostic_getGenericFeatures(%s, %d) failed. Return code: %d",
	4014: "Call to G2Diagnostic_getLastException() failed. Return code: %d",
	4015: "Call to G2Diagnostic_getMappingStatistics(%d) failed. Return code: %d",
	4016: "Call to G2Diagnostic_getRelationshipDetails(%d, %d) failed. Return code: %d",
	4017: "Call to G2Diagnostic_getResolutionStatistics() failed. Return code: %d",
	4018: "Call to G2Diagnostic_init(%s, %s, %d) failed. Return code: %d",
	4019: "Call to G2Diagnostic_initWithConfigID(%s, %s, %d, %d) failed. Return code: %d",
	4020: "Call to G2Diagnostic_reinit(%d) failed. Return Code: %d",
	5901: "During setup, call to messagelogger.NewSenzingApiLogger() failed.",
	5902: "During setup, call to g2eg2engineconfigurationjson.BuildSimpleSystemConfigurationJson() failed.",
	5903: "During setup, call to g2engine.Init() failed.",
	5904: "During setup, call to g2engine.PurgeRepository() failed.",
	5905: "During setup, call to g2engine.Destroy() failed.",
	5906: "During setup, call to g2config.Init() failed.",
	5907: "During setup, call to g2config.Create() failed.",
	5908: "During setup, call to g2config.AddDataSource() failed.",
	5909: "During setup, call to g2config.Save() failed.",
	5910: "During setup, call to g2config.Close() failed.",
	5911: "During setup, call to g2config.Destroy() failed.",
	5912: "During setup, call to g2configmgr.Init() failed.",
	5913: "During setup, call to g2configmgr.AddConfig() failed.",
	5914: "During setup, call to g2configmgr.SetDefaultConfigID() failed.",
	5915: "During setup, call to g2configmgr.Destroy() failed.",
	5916: "During setup, call to g2engine.Init() failed.",
	5917: "During setup, call to g2engine.AddRecord() failed.",
	5918: "During setup, call to g2engine.Destroy() failed.",
	5920: "During setup, call to setupSenzingConfig() failed.",
	5921: "During setup, call to setupPurgeRepository() failed.",
	5922: "During setup, call to setupAddRecords() failed.",
	5931: "During setup, call to g2engine.Init() failed.",
	5932: "During setup, call to g2engine.PurgeRepository() failed.",
	5933: "During setup, call to g2engine.Destroy() failed.",
	8001: "CheckDBPerf",
	8002: "CloseEntityListBySize",
	8003: "Destroy",
	8004: "FetchNextEntityBySize",
	8005: "FindEntitiesByFeatureIDs",
	8006: "GetAvailableMemory",
	8007: "GetDataSourceCounts",
	8008: "GetDBInfo",
	8009: "GetEntityDetails",
	8010: "GetEntityListBySize",
	8011: "GetEntityResume",
	8012: "GetEntitySizeBreakdown",
	8013: "GetFeature",
	8014: "GetGenericFeatures",
	8015: "GetLogicalCores",
	8016: "GetMappingStatistics",
	8017: "GetPhysicalCores",
	8018: "GetRelationshipDetails",
	8019: "GetResolutionStatistics",
	8020: "GetTotalSystemMemory",
	8021: "Init",
	8022: "InitWithConfigID",
	8023: "Reinit",
}

// Status strings for specific g2diagnostic messages.
var IdStatuses = map[int]string{}
