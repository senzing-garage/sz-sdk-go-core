package g2engine

import (
	"context"

	"github.com/senzing/go-logging/logger"
	"github.com/senzing/go-observing/observer"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

// Senzing flags.
type FlagMask int64

// The G2engine interface is a Golang representation of Senzing's libg2.h
type G2engine interface {
	AddRecord(ctx context.Context, dataSourceCode string, recordID string, jsonData string, loadID string) error
	AddRecordWithInfo(ctx context.Context, dataSourceCode string, recordID string, jsonData string, loadID string, flags int64) (string, error)
	AddRecordWithInfoWithReturnedRecordID(ctx context.Context, dataSourceCode string, jsonData string, loadID string, flags int64) (string, string, error)
	AddRecordWithReturnedRecordID(ctx context.Context, dataSourceCode string, jsonData string, loadID string) (string, error)
	CheckRecord(ctx context.Context, record string, recordQueryList string) (string, error)
	CloseExport(ctx context.Context, responseHandle uintptr) error
	CountRedoRecords(ctx context.Context) (int64, error)
	DeleteRecord(ctx context.Context, dataSourceCode string, recordID string, loadID string) error
	DeleteRecordWithInfo(ctx context.Context, dataSourceCode string, recordID string, loadID string, flags int64) (string, error)
	Destroy(ctx context.Context) error
	ExportConfig(ctx context.Context) (string, error)
	ExportConfigAndConfigID(ctx context.Context) (string, int64, error)
	ExportCSVEntityReport(ctx context.Context, csvColumnList string, flags int64) (uintptr, error)
	ExportJSONEntityReport(ctx context.Context, flags int64) (uintptr, error)
	FetchNext(ctx context.Context, responseHandle uintptr) (string, error)
	FindInterestingEntitiesByEntityID(ctx context.Context, entityID int64, flags int64) (string, error)
	FindInterestingEntitiesByRecordID(ctx context.Context, dataSourceCode string, recordID string, flags int64) (string, error)
	FindNetworkByEntityID_V2(ctx context.Context, entityList string, maxDegree int, buildOutDegree int, maxEntities int, flags int64) (string, error)
	FindNetworkByEntityID(ctx context.Context, entityList string, maxDegree int, buildOutDegree int, maxEntities int) (string, error)
	FindNetworkByRecordID_V2(ctx context.Context, recordList string, maxDegree int, buildOutDegree int, maxEntities int, flags int64) (string, error)
	FindNetworkByRecordID(ctx context.Context, recordList string, maxDegree int, buildOutDegree int, maxEntities int) (string, error)
	FindPathByEntityID_V2(ctx context.Context, entityID1 int64, entityID2 int64, maxDegree int, flags int64) (string, error)
	FindPathByEntityID(ctx context.Context, entityID1 int64, entityID2 int64, maxDegree int) (string, error)
	FindPathByRecordID_V2(ctx context.Context, dataSourceCode1 string, recordID1 string, dataSourceCode2 string, recordID2 string, maxDegree int, flags int64) (string, error)
	FindPathByRecordID(ctx context.Context, dataSourceCode1 string, recordID1 string, dataSourceCode2 string, recordID2 string, maxDegree int) (string, error)
	FindPathExcludingByEntityID_V2(ctx context.Context, entityID1 int64, entityID2 int64, maxDegree int, excludedEntities string, flags int64) (string, error)
	FindPathExcludingByEntityID(ctx context.Context, entityID1 int64, entityID2 int64, maxDegree int, excludedEntities string) (string, error)
	FindPathExcludingByRecordID_V2(ctx context.Context, dataSourceCode1 string, recordID1 string, dataSourceCode2 string, recordID2 string, maxDegree int, excludedRecords string, flags int64) (string, error)
	FindPathExcludingByRecordID(ctx context.Context, dataSourceCode1 string, recordID1 string, dataSourceCode2 string, recordID2 string, maxDegree int, excludedRecords string) (string, error)
	FindPathIncludingSourceByEntityID_V2(ctx context.Context, entityID1 int64, entityID2 int64, maxDegree int, excludedEntities string, requiredDsrcs string, flags int64) (string, error)
	FindPathIncludingSourceByEntityID(ctx context.Context, entityID1 int64, entityID2 int64, maxDegree int, excludedEntities string, requiredDsrcs string) (string, error)
	FindPathIncludingSourceByRecordID_V2(ctx context.Context, dataSourceCode1 string, recordID1 string, dataSourceCode2 string, recordID2 string, maxDegree int, excludedRecords string, requiredDsrcs string, flags int64) (string, error)
	FindPathIncludingSourceByRecordID(ctx context.Context, dataSourceCode1 string, recordID1 string, dataSourceCode2 string, recordID2 string, maxDegree int, excludedRecords string, requiredDsrcs string) (string, error)
	GetActiveConfigID(ctx context.Context) (int64, error)
	GetEntityByEntityID_V2(ctx context.Context, entityID int64, flags int64) (string, error)
	GetEntityByEntityID(ctx context.Context, entityID int64) (string, error)
	GetEntityByRecordID_V2(ctx context.Context, dataSourceCode string, recordID string, flags int64) (string, error)
	GetEntityByRecordID(ctx context.Context, dataSourceCode string, recordID string) (string, error)
	GetRecord_V2(ctx context.Context, dataSourceCode string, recordID string, flags int64) (string, error)
	GetRecord(ctx context.Context, dataSourceCode string, recordID string) (string, error)
	GetRedoRecord(ctx context.Context) (string, error)
	GetRepositoryLastModifiedTime(ctx context.Context) (int64, error)
	GetVirtualEntityByRecordID_V2(ctx context.Context, recordList string, flags int64) (string, error)
	GetVirtualEntityByRecordID(ctx context.Context, recordList string) (string, error)
	HowEntityByEntityID_V2(ctx context.Context, entityID int64, flags int64) (string, error)
	HowEntityByEntityID(ctx context.Context, entityID int64) (string, error)
	Init(ctx context.Context, moduleName string, iniParams string, verboseLogging int) error
	InitWithConfigID(ctx context.Context, moduleName string, iniParams string, initConfigID int64, verboseLogging int) error
	PrimeEngine(ctx context.Context) error
	Process(ctx context.Context, record string) error
	ProcessRedoRecord(ctx context.Context) (string, error)
	ProcessRedoRecordWithInfo(ctx context.Context, flags int64) (string, string, error)
	ProcessWithInfo(ctx context.Context, record string, flags int64) (string, error)
	ProcessWithResponse(ctx context.Context, record string) (string, error)
	ProcessWithResponseResize(ctx context.Context, record string) (string, error)
	PurgeRepository(ctx context.Context) error
	ReevaluateEntity(ctx context.Context, entityID int64, flags int64) error
	ReevaluateEntityWithInfo(ctx context.Context, entityID int64, flags int64) (string, error)
	ReevaluateRecord(ctx context.Context, dataSourceCode string, recordID string, flags int64) error
	ReevaluateRecordWithInfo(ctx context.Context, dataSourceCode string, recordID string, flags int64) (string, error)
	RegisterObserver(ctx context.Context, observer observer.Observer) error
	Reinit(ctx context.Context, initConfigID int64) error
	ReplaceRecord(ctx context.Context, dataSourceCode string, recordID string, jsonData string, loadID string) error
	ReplaceRecordWithInfo(ctx context.Context, dataSourceCode string, recordID string, jsonData string, loadID string, flags int64) (string, error)
	SearchByAttributes_V2(ctx context.Context, jsonData string, flags int64) (string, error)
	SearchByAttributes(ctx context.Context, jsonData string) (string, error)
	SetLogLevel(ctx context.Context, logLevel logger.Level) error
	Stats(ctx context.Context) (string, error)
	UnregisterObserver(ctx context.Context, observer observer.Observer) error
	WhyEntities_V2(ctx context.Context, entityID1 int64, entityID2 int64, flags int64) (string, error)
	WhyEntities(ctx context.Context, entityID1 int64, entityID2 int64) (string, error)
	WhyEntityByEntityID_V2(ctx context.Context, entityID int64, flags int64) (string, error)
	WhyEntityByEntityID(ctx context.Context, entityID int64) (string, error)
	WhyEntityByRecordID_V2(ctx context.Context, dataSourceCode string, recordID string, flags int64) (string, error)
	WhyEntityByRecordID(ctx context.Context, dataSourceCode string, recordID string) (string, error)
	WhyRecords_V2(ctx context.Context, dataSourceCode1 string, recordID1 string, dataSourceCode2 string, recordID2 string, flags int64) (string, error)
	WhyRecords(ctx context.Context, dataSourceCode1 string, recordID1 string, dataSourceCode2 string, recordID2 string) (string, error)
}

// ----------------------------------------------------------------------------
// Constants
// ----------------------------------------------------------------------------

// Identfier of the g2engine package found messages having the format "senzing-6004xxxx".
const ProductId = 6004

// Flags used by the Senzing G2Engine.
// These flags are single-bit flags.
const (
	G2_EXPORT_INCLUDE_RESOLVED         FlagMask = 0         // 0 we should include entities with "resolved" relationships
	G2_EXPORT_INCLUDE_POSSIBLY_SAME    FlagMask = 1 << iota // 1 we should include entities with "possibly same" relationships
	G2_EXPORT_INCLUDE_POSSIBLY_RELATED                      // 2 we should include entities with "possibly related" relationships
	G2_EXPORT_INCLUDE_NAME_ONLY                             // 3 we should include entities with "name only" relationships
	G2_EXPORT_INCLUDE_DISCLOSED                             // 4 we should include entities with "disclosed" relationships
	G2_EXPORT_INCLUDE_SINGLETONS                            // 5 we should include singleton entities

	/* flags for outputting entity relation data  */
	G2_ENTITY_INCLUDE_POSSIBLY_SAME_RELATIONS    // 6 include "possibly same" relationships on entities
	G2_ENTITY_INCLUDE_POSSIBLY_RELATED_RELATIONS // 7 include "possibly related" relationships on entities
	G2_ENTITY_INCLUDE_NAME_ONLY_RELATIONS        // 8 include "name only" relationships on entities
	G2_ENTITY_INCLUDE_DISCLOSED_RELATIONS        // 9 include "disclosed" relationships on entities

	/* flags for outputting entity feature data  */
	G2_ENTITY_INCLUDE_ALL_FEATURES            // 10 include all features for entities
	G2_ENTITY_INCLUDE_REPRESENTATIVE_FEATURES // 11 include only representative features on entities

	/* flags for getting extra information about an entity  */
	G2_ENTITY_INCLUDE_ENTITY_NAME            // 12 include the name of the entity
	G2_ENTITY_INCLUDE_RECORD_SUMMARY         // 13 include the record summary of the entity
	G2_ENTITY_INCLUDE_RECORD_DATA            // 14 include the basic record data for the entity
	G2_ENTITY_INCLUDE_RECORD_MATCHING_INFO   // 15 include the record matching info for the entity
	G2_ENTITY_INCLUDE_RECORD_JSON_DATA       // 16 include the record json data for the entity
	G2_ENTITY_INCLUDE_RECORD_FORMATTED_DATA  // 17 include the record formattted data for the entity
	G2_ENTITY_INCLUDE_RECORD_FEATURE_IDS     // 18 include the features identifiers for the records
	G2_ENTITY_INCLUDE_RELATED_ENTITY_NAME    // 19 include the name of the related entities
	G2_ENTITY_INCLUDE_RELATED_MATCHING_INFO  // 20 include the record matching info of the related entities
	G2_ENTITY_INCLUDE_RELATED_RECORD_SUMMARY // 21 include the record summary of the related entities
	G2_ENTITY_INCLUDE_RELATED_RECORD_DATA    // 22 include the basic record of the related entities

	/* flags for extra feature data  */
	G2_ENTITY_OPTION_INCLUDE_INTERNAL_FEATURES // 23 include internal features
	G2_ENTITY_OPTION_INCLUDE_FEATURE_STATS     // 24 include statistics on features

	/* flags for finding entity path data  */
	G2_FIND_PATH_PREFER_EXCLUDE // 25 excluded entities are still allowed, but not preferred

	/* flags for including search result feature scores  */
	G2_INCLUDE_FEATURE_SCORES // 26 include feature scores
	G2_SEARCH_INCLUDE_STATS   //  27 include statistics from search results
)

// Flags used by the Senzing G2Engine.
// These flags combine single-bit flags.
const (
	/* flags for exporting entity data  */
	G2_EXPORT_INCLUDE_ALL_ENTITIES      = G2_EXPORT_INCLUDE_RESOLVED | G2_EXPORT_INCLUDE_SINGLETONS
	G2_EXPORT_INCLUDE_ALL_RELATIONSHIPS = G2_EXPORT_INCLUDE_POSSIBLY_SAME | G2_EXPORT_INCLUDE_POSSIBLY_RELATED | G2_EXPORT_INCLUDE_NAME_ONLY | G2_EXPORT_INCLUDE_DISCLOSED

	/* flags for outputting entity relation data  */
	G2_ENTITY_INCLUDE_ALL_RELATIONS  = G2_ENTITY_INCLUDE_POSSIBLY_SAME_RELATIONS | G2_ENTITY_INCLUDE_POSSIBLY_RELATED_RELATIONS | G2_ENTITY_INCLUDE_NAME_ONLY_RELATIONS | G2_ENTITY_INCLUDE_DISCLOSED_RELATIONS
	G2_SEARCH_INCLUDE_FEATURE_SCORES = G2_INCLUDE_FEATURE_SCORES // include feature scores from search results

	/* flags for searching for entities  */
	G2_SEARCH_INCLUDE_RESOLVED         = G2_EXPORT_INCLUDE_RESOLVED
	G2_SEARCH_INCLUDE_POSSIBLY_SAME    = G2_EXPORT_INCLUDE_POSSIBLY_SAME
	G2_SEARCH_INCLUDE_POSSIBLY_RELATED = G2_EXPORT_INCLUDE_POSSIBLY_RELATED
	G2_SEARCH_INCLUDE_NAME_ONLY        = G2_EXPORT_INCLUDE_NAME_ONLY
	G2_SEARCH_INCLUDE_ALL_ENTITIES     = G2_SEARCH_INCLUDE_RESOLVED | G2_SEARCH_INCLUDE_POSSIBLY_SAME | G2_SEARCH_INCLUDE_POSSIBLY_RELATED | G2_SEARCH_INCLUDE_NAME_ONLY

	/* recommended settings */
	G2_RECORD_DEFAULT_FLAGS       = G2_ENTITY_INCLUDE_RECORD_JSON_DATA                                                                                                                                                                                                                                                                                                                   // the recommended default flag values for getting records
	G2_ENTITY_DEFAULT_FLAGS       = G2_ENTITY_INCLUDE_ALL_RELATIONS | G2_ENTITY_INCLUDE_REPRESENTATIVE_FEATURES | G2_ENTITY_INCLUDE_ENTITY_NAME | G2_ENTITY_INCLUDE_RECORD_SUMMARY | G2_ENTITY_INCLUDE_RECORD_DATA | G2_ENTITY_INCLUDE_RECORD_MATCHING_INFO | G2_ENTITY_INCLUDE_RELATED_ENTITY_NAME | G2_ENTITY_INCLUDE_RELATED_RECORD_SUMMARY | G2_ENTITY_INCLUDE_RELATED_MATCHING_INFO // the recommended default flag values for getting entities
	G2_ENTITY_BRIEF_DEFAULT_FLAGS = G2_ENTITY_INCLUDE_RECORD_MATCHING_INFO | G2_ENTITY_INCLUDE_ALL_RELATIONS | G2_ENTITY_INCLUDE_RELATED_MATCHING_INFO                                                                                                                                                                                                                                   // the recommended default flag values for a brief entity result
	G2_EXPORT_DEFAULT_FLAGS       = G2_EXPORT_INCLUDE_ALL_ENTITIES | G2_EXPORT_INCLUDE_ALL_RELATIONSHIPS | G2_ENTITY_INCLUDE_ALL_RELATIONS | G2_ENTITY_INCLUDE_REPRESENTATIVE_FEATURES | G2_ENTITY_INCLUDE_ENTITY_NAME | G2_ENTITY_INCLUDE_RECORD_DATA | G2_ENTITY_INCLUDE_RECORD_MATCHING_INFO | G2_ENTITY_INCLUDE_RELATED_MATCHING_INFO                                                // the recommended default flag values for exporting entities
	G2_FIND_PATH_DEFAULT_FLAGS    = G2_ENTITY_INCLUDE_ALL_RELATIONS | G2_ENTITY_INCLUDE_ENTITY_NAME | G2_ENTITY_INCLUDE_RECORD_SUMMARY | G2_ENTITY_INCLUDE_RELATED_MATCHING_INFO                                                                                                                                                                                                         // the recommended default flag values for finding entity paths
	G2_WHY_ENTITY_DEFAULT_FLAGS   = G2_ENTITY_DEFAULT_FLAGS | G2_ENTITY_INCLUDE_RECORD_FEATURE_IDS | G2_ENTITY_OPTION_INCLUDE_INTERNAL_FEATURES | G2_ENTITY_OPTION_INCLUDE_FEATURE_STATS | G2_INCLUDE_FEATURE_SCORES                                                                                                                                                                     // the recommended default flag values for why-analysis on entities
	G2_HOW_ENTITY_DEFAULT_FLAGS   = G2_ENTITY_DEFAULT_FLAGS | G2_ENTITY_INCLUDE_RECORD_FEATURE_IDS | G2_ENTITY_OPTION_INCLUDE_INTERNAL_FEATURES | G2_ENTITY_OPTION_INCLUDE_FEATURE_STATS | G2_INCLUDE_FEATURE_SCORES                                                                                                                                                                     // the recommended default flag values for how-analysis on entities

	G2_SEARCH_BY_ATTRIBUTES_ALL            = G2_SEARCH_INCLUDE_ALL_ENTITIES | G2_ENTITY_INCLUDE_REPRESENTATIVE_FEATURES | G2_ENTITY_INCLUDE_ENTITY_NAME | G2_ENTITY_INCLUDE_RECORD_SUMMARY | G2_SEARCH_INCLUDE_FEATURE_SCORES                               // the recommended flag values for searching by attributes, returning all matching entities
	G2_SEARCH_BY_ATTRIBUTES_STRONG         = G2_SEARCH_INCLUDE_RESOLVED | G2_SEARCH_INCLUDE_POSSIBLY_SAME | G2_ENTITY_INCLUDE_REPRESENTATIVE_FEATURES | G2_ENTITY_INCLUDE_ENTITY_NAME | G2_ENTITY_INCLUDE_RECORD_SUMMARY | G2_SEARCH_INCLUDE_FEATURE_SCORES // the recommended flag values for searching by attributes, returning only strongly matching entities
	G2_SEARCH_BY_ATTRIBUTES_MINIMAL_ALL    = G2_SEARCH_INCLUDE_ALL_ENTITIES                                                                                                                                                                                 // return minimal data with all matches
	G2_SEARCH_BY_ATTRIBUTES_MINIMAL_STRONG = G2_SEARCH_INCLUDE_RESOLVED | G2_SEARCH_INCLUDE_POSSIBLY_SAME                                                                                                                                                   // return minimal data with only the strongest matches
	G2_SEARCH_BY_ATTRIBUTES_DEFAULT_FLAGS  = G2_SEARCH_BY_ATTRIBUTES_ALL                                                                                                                                                                                    // the recommended default flag values for search-by-attributes
)

// ----------------------------------------------------------------------------
// Variables
// ----------------------------------------------------------------------------

// Message templates for the g2engine package.
var IdMessages = map[int]string{
	1:    "Enter AddRecord(%s, %s, %s, %s).",
	2:    "Exit  AddRecord(%s, %s, %s, %s) returned (%v).",
	3:    "Enter AddRecordWithInfo(%s, %s, %s, %s, %d).",
	4:    "Exit  AddRecordWithInfo(%s, %s, %s, %s, %d) returned (%s, %v).",
	5:    "Enter AddRecordWithInfoWithReturnedRecordID(%s, %s, %s, %d).",
	6:    "Exit  AddRecordWithInfoWithReturnedRecordID(%s, %s, %s, %d) returned (%s, %s, %v).",
	7:    "Enter AddRecordWithReturnedRecordID(%s, %s, %s).",
	8:    "Exit  AddRecordWithReturnedRecordID(%s, %s, %s) returned (%s, %v).",
	9:    "Enter CheckRecord(%s, %s).",
	10:   "Exit  CheckRecord(%s, %s) returned (%s, %v).",
	11:   "Enter ClearLastException().",
	12:   "Exit  ClearLastException() returned (%v).",
	13:   "Enter CloseExport(%v).",
	14:   "Exit  CloseExport(%v) returned (%v).",
	15:   "Enter CountRedoRecords().",
	16:   "Exit  CountRedoRecords() returned (%d, %v).",
	17:   "Enter DeleteRecord(%s, %s, %s).",
	18:   "Exit  DeleteRecord(%s, %s, %s) returned (%v).",
	19:   "Enter DeleteRecordWithInfo(%s, %s, %s, %d).",
	20:   "Exit  DeleteRecordWithInfo(%s, %s, %s, %d) returned (%s, %v).",
	21:   "Enter Destroy().",
	22:   "Exit  Destroy() returned (%v).",
	23:   "Enter ExportConfigAndConfigID().",
	24:   "Exit  ExportConfigAndConfigID() returned (%s, %d, %v).",
	25:   "Enter ExportConfig().",
	26:   "Exit  ExportConfig() returned (%s, %v).",
	27:   "Enter ExportCSVEntityReport(%s, %d).",
	28:   "Exit  ExportCSVEntityReport(%s, %d) returned (%v, %v).",
	29:   "Enter ExportJSONEntityReport(%d).",
	30:   "Exit  ExportJSONEntityReport(%d) returned (%v, %v).",
	31:   "Enter FetchNext(%v).",
	32:   "Exit  FetchNext(%v) returned (%s, %v).",
	33:   "Enter FindInterestingEntitiesByEntityID(%d, %d).",
	34:   "Exit  FindInterestingEntitiesByEntityID(%d, %d) returned (%s, %v).",
	35:   "Enter FindInterestingEntitiesByRecordID(%s, %s, %d).",
	36:   "Exit  FindInterestingEntitiesByRecordID(%s, %s, %d) returned (%s, %v).",
	37:   "Enter FindNetworkByEntityID(%s, %d, %d, %d).",
	38:   "Exit  FindNetworkByEntityID(%s, %d, %d, %d) returned (%s, %v).",
	39:   "Enter FindNetworkByEntityID_V2(%s, %d, %d, %d, %d).",
	40:   "Exit  FindNetworkByEntityID_V2(%s, %d, %d, %d, %d) returned (%s, %v).",
	41:   "Enter FindNetworkByRecordID(%s, %d, %d, %d).",
	42:   "Exit  FindNetworkByRecordID(%s, %d, %d, %d) returned (%s, %v).",
	43:   "Enter FindNetworkByRecordID_V2(%s, %d, %d, %d, %d).",
	44:   "Exit  FindNetworkByRecordID_V2(%s, %d, %d, %d, %d) returned (%s, %v).",
	45:   "Enter FindPathByEntityID(%d, %d, %d).",
	46:   "Exit  FindPathByEntityID(%d, %d, %d) returned (%s, %v).",
	47:   "Enter FindPathByEntityID_V2(%d, %d, %d, %d).",
	48:   "Exit  FindPathByEntityID_V2(%d, %d, %d, %d) returned (%s, %v).",
	49:   "Enter FindPathByRecordID(%s, %s, %s, %s, %d).",
	50:   "Exit  FindPathByRecordID(%s, %s, %s, %s, %d) returned (%s, %v).",
	51:   "Enter FindPathByRecordID_V2(%s, %s, %s, %s, %d, %d).",
	52:   "Exit  FindPathByRecordID_V2(%s, %s, %s, %s, %d, %d) returned (%s, %v).",
	53:   "Enter FindPathExcludingByEntityID(%d, %d, %d, %s).",
	54:   "Exit  FindPathExcludingByEntityID(%d, %d, %d, %s) returned (%s, %v).",
	55:   "Enter FindPathExcludingByEntityID_V2(%d, %d, %d, %s, %d).",
	56:   "Exit  FindPathExcludingByEntityID_V2(%d, %d, %d, %s, %d) returned (%s, %v).",
	57:   "Enter FindPathExcludingByRecordID(%s, %s, %s, %s, %d, %s).",
	58:   "Exit  FindPathExcludingByRecordID(%s, %s, %s, %s, %d, %s) returned (%s, %v).",
	59:   "Enter FindPathExcludingByRecordID_V2(%s, %s, %s, %s, %d, %s, %d).",
	60:   "Exit  FindPathExcludingByRecordID_V2(%s, %s, %s, %s, %d, %s, %d) returned (%v).",
	61:   "Enter FindPathIncludingSourceByEntityID(%d, %d, %d, %s, %s).",
	62:   "Exit  FindPathIncludingSourceByEntityID(%d, %d, %d, %s, %s) returned (%s, %v).",
	63:   "Enter FindPathIncludingSourceByEntityID_V2(%d, %d, %d, %s, %s, %d).",
	64:   "Exit  FindPathIncludingSourceByEntityID_V2(%d, %d, %d, %s, %s, %d) returned (%s, %v).",
	65:   "Enter FindPathIncludingSourceByRecordID(%s, %s, %s, %s, %d, %s, %s).",
	66:   "Exit  FindPathIncludingSourceByRecordID(%s, %s, %s, %s, %d, %s, %s) returned (%s, %v).",
	67:   "Enter FindPathIncludingSourceByRecordID_V2(%s, %s, %s, %s, %d, %s, %s, %d).",
	68:   "Exit  FindPathIncludingSourceByRecordID_V2(%s, %s, %s, %s, %d, %s, %s, %d) returned (%s, %v).",
	69:   "Enter GetActiveConfigID().",
	70:   "Exit  GetActiveConfigID() returned (%d, %v).",
	71:   "Enter GetEntityByEntityID(%d).",
	72:   "Exit  GetEntityByEntityID(%d) returned (%s, %v).",
	73:   "Enter GetEntityByEntityID_V2(%d, %d).",
	74:   "Exit  GetEntityByEntityID_V2(%d, %d) returned (%s, %v).",
	75:   "Enter GetEntityByRecordID(%s, %s).",
	76:   "Exit  GetEntityByRecordID(%s, %s) returned (%s, %v).",
	77:   "Enter GetEntityByRecordID_V2(%s, %s, %d).",
	78:   "Exit  GetEntityByRecordID_V2(%s, %s, %d) returned (%s, %v).",
	79:   "Enter GetLastException().",
	80:   "Exit  GetLastException() returned (%s, %v).",
	81:   "Enter GetLastExceptionCode().",
	82:   "Exit  GetLastExceptionCode() returned (%d, %v).",
	83:   "Enter GetRecord(%s, %s).",
	84:   "Exit  GetRecord(%s, %s) returned (%s, %v).",
	85:   "Enter GetRecord_V2(%s, %s, %d).",
	86:   "Exit  GetRecord_V2(%s, %s, %d) returned (%s, %v).",
	87:   "Enter GetRedoRecord().",
	88:   "Exit  GetRedoRecord() returned (%s, %v).",
	89:   "Enter GetRepositoryLastModifiedTime().",
	90:   "Exit  GetRepositoryLastModifiedTime() returned (%d, %v).",
	91:   "Enter GetVirtualEntityByRecordID(%s).",
	92:   "Exit  GetVirtualEntityByRecordID(%s) returned (%s, %v).",
	93:   "Enter GetVirtualEntityByRecordID_V2(%s, %d).",
	94:   "Exit  GetVirtualEntityByRecordID_V2(%s, %d) returned (%s, %v).",
	95:   "Enter HowEntityByEntityID(%d).",
	96:   "Exit  HowEntityByEntityID(%d) returned (%s, %v).",
	97:   "Enter HowEntityByEntityID_V2(%d, %d).",
	98:   "Exit  HowEntityByEntityID_V2(%d, %d) returned (%s, %v).",
	99:   "Enter Init(%s, %s, %d).",
	100:  "Exit  Init(%s, %s, %d) returned (%v).",
	101:  "Enter InitWithConfigID(%s, %s, %d, %d).",
	102:  "Exit  InitWithConfigID(%s, %s, %d, %d) returned (%v).",
	103:  "Enter PrimeEngine().",
	104:  "Exit  PrimeEngine() returned (%v).",
	105:  "Enter Process(%s).",
	106:  "Exit  Process(%s) returned (%v).",
	107:  "Enter ProcessRedoRecord().",
	108:  "Exit  ProcessRedoRecord() returned (%s, %v).",
	109:  "Enter ProcessRedoRecordWithInfo(%d).",
	110:  "Exit  ProcessRedoRecordWithInfo(%d) returned (%s, %s, %v).",
	111:  "Enter ProcessWithInfo(%s, %d).",
	112:  "Exit  ProcessWithInfo(%s, %d) returned (%s, %v).",
	113:  "Enter ProcessWithResponse(%s).",
	114:  "Exit  ProcessWithResponse(%s) returned (%s, %v).",
	115:  "Enter ProcessWithResponseResize(%s).",
	116:  "Exit  ProcessWithResponseResize(%s) returned (%s, %v).",
	117:  "Enter PurgeRepository().",
	118:  "Exit  PurgeRepository() returned (%v).",
	119:  "Enter ReevaluateEntity(%d, %d).",
	120:  "Exit  ReevaluateEntity(%d, %d) returned (%v).",
	121:  "Enter ReevaluateEntityWithInfo(%d, %d).",
	122:  "Exit  ReevaluateEntityWithInfo(%d, %d) returned (%s, %v).",
	123:  "Enter ReevaluateRecord(%s, %s, %d).",
	124:  "Exit  ReevaluateRecord(%s, %s, %d) returned (%v).",
	125:  "Enter ReevaluateRecordWithInfo(%s, %s, %d).",
	126:  "Exit  ReevaluateRecordWithInfo(%s, %s, %d) returned (%s, %v).",
	127:  "Enter Reinit(%d).",
	128:  "Exit  Reinit(%d) returned (%v).",
	129:  "Enter ReplaceRecord(%s, %s, %s, %s).",
	130:  "Exit  ReplaceRecord(%s, %s, %s, %s) returned (%v).",
	131:  "Enter ReplaceRecordWithInfo(%s, %s, %s, %s, %d).",
	132:  "Exit  ReplaceRecordWithInfo(%s, %s, %s, %s, %d) returned (%s, %v).",
	133:  "Enter SearchByAttributes(%s).",
	134:  "Exit  SearchByAttributes(%s) returned (%s, %v).",
	135:  "Enter SearchByAttributes_V2(%s, %d).",
	136:  "Exit  SearchByAttributes_V2(%s, %d) returned (%s, %v).",
	137:  "Enter SetLogLevel(%v).",
	138:  "Exit  SetLogLevel(%v) returned (%v).",
	139:  "Enter Stats().",
	140:  "Exit  Stats() returned (%s, %v).",
	141:  "Enter WhyEntities(%d, %d).",
	142:  "Exit  WhyEntities(%d, %d) returned (%s, %v).",
	143:  "Enter WhyEntities_V2(%d, %d, %d).",
	144:  "Exit  WhyEntities_V2(%d, %d, %d) returned (%s, %v).",
	145:  "Enter WhyEntityByEntityID(%d).",
	146:  "Exit  WhyEntityByEntityID(%d) returned (%s, %v).",
	147:  "Enter WhyEntityByEntityID_V2(%d, %d).",
	148:  "Exit  WhyEntityByEntityID_V2(%d, %d) returned (%s, %v).",
	149:  "Enter WhyEntityByRecordID(%s, %s).",
	150:  "Exit  WhyEntityByRecordID(%s, %s) returned (%s, %v).",
	151:  "Enter WhyEntityByRecordID_V2(%s, %s, %d).",
	152:  "Exit  WhyEntityByRecordID_V2(%s, %s, %d) returned (%s, %v).",
	153:  "Enter WhyRecords(%s, %s, %s, %s).",
	154:  "Exit  WhyRecords(%s, %s, %s, %s) returned (%s, %v).",
	155:  "Enter WhyRecords_V2(%s, %s, %s, %s, %d).",
	156:  "Exit  WhyRecords_V2(%s, %s, %s, %s, %d) returned (%s, %v).",
	4001: "Call to G2_addRecord(%s, %s, %s, %s) failed. Return code: %d",
	4002: "Call to G2_addRecordWithInfo(%s, %s, %s, %s, %d) failed. Return code: %d",
	4003: "Call to G2_addRecordWithInfoWithReturnedRecordID(%s, %s, %s, %d) failed. Return code: %d",
	4004: "Call to G2_addRecordWithReturnedRecordID(%s, %s, %s) failed. Return code: %d",
	4005: "Call to G2_checkRecord(%s, %s) failed. Return code: %d",
	4006: "Call to G2_closeExport(%v) failed. Return code: %d",
	4007: "Call to G2_deleteRecord(%s, %s, %s) failed. Return code: %d",
	4008: "Call to G2_deleteRecordWithInfo(%s, %s, %s, %d) failed. Return code: %d",
	4009: "Call to G2_destroy() failed. Return code: %d",
	4010: "Call to G2_exportConfigAndConfigID() failed. Return code: %d",
	4011: "Call to G2_exportConfig() failed. Return code: %d",
	4012: "Call to G2_exportCSVEntityReport(%s, %d) failed. Return code: %d",
	4013: "Call to G2_exportJSONEntityReport(%d) failed. Return code: %d",
	4014: "Call to G2_fetchNext(%v) failed. Return code: %d",
	4015: "Call to G2_findInterestingEntitiesByEntityID(%d, %d) failed. Return code: %d",
	4016: "Call to G2_findInterestingEntitiesByRecordID(%s, %s, %d) failed. Return code: %d",
	4017: "Call to G2_findNetworkByEntityID(%s, %d, %d, %d) failed. Return code: %d",
	4018: "Call to G2_findNetworkByEntityID_V2(%s, %d, %d, %d, %d) failed. Return code: %d",
	4019: "Call to G2_findNetworkByRecordID(%s, %d, %d, %d) failed. Return code: %d",
	4020: "Call to G2_findNetworkByRecordID_V2(%s, %d, %d, %d, %d) failed. Return code: %d",
	4021: "Call to G2_findPathByEntityID(%d, %d, %d) failed. Return code: %d",
	4022: "Call to G2_findPathByEntityID_V2(%d, %d, %d, %d) failed. Return code: %d",
	4023: "Call to G2_findPathByRecordID(%s, %s, %s, %s, %d) failed. Return code: %d",
	4024: "Call to G2_findPathByRecordID_V2(%s, %s, %s, %s, %d, %d) failed. Return code: %d",
	4025: "Call to G2_findPathExcludingByEntityID(%d, %d, %d, %s) failed. Return code: %d",
	4026: "Call to G2_findPathExcludingByEntityID_V2(%d, %d, %d, %s, %d) failed. Return code: %d",
	4027: "Call to G2_findPathExcludingByRecordID(%s, %s, %s, %s %d, %s) failed. Return code: %d",
	4028: "Call to G2_findPathExcludingByRecordID_V2(%s, %s, %s, %s %d, %s, %d) failed. Return code: %d",
	4029: "Call to G2_findPathIncludingSourceByEntityID(%d, %d, %d, %s, %s) failed. Return code: %d",
	4030: "Call to G2_findPathIncludingSourceByEntityID_V2(%d, %d, %d, %s, %s, %d) failed. Return code: %d",
	4031: "Call to G2_findPathIncludingSourceByRecordID(%s, %s, %s, %s %d, %s, %s) failed. Return code: %d",
	4032: "Call to G2_findPathIncludingSourceByRecordID_V2(%s, %s, %s, %s %d, %s, %s, %d) failed. Return code: %d",
	4033: "Call to G2_getActiveConfigID() failed. Return code: %d",
	4034: "Call to G2_getEntityByEntityID(%d) failed. Return code: %d",
	4035: "Call to G2_getEntityByEntityID_V2(%d, %d) failed. Return code: %d",
	4036: "Call to G2_getEntityByRecordID(%s, %s) failed. Return code: %d",
	4037: "Call to G2_getEntityByRecordID_V2(%s, %s, %d) failed. Return code: %d",
	4038: "Call to G2_getLastException() failed. Return code: %d",
	4039: "Call to G2_getRecord(%s, %s) failed. Return code: %d",
	4040: "Call to G2_getRecord_V2(%s, %s, %d) failed. Return code: %d",
	4041: "Call to G2_getRedoRecord() failed. Return code: %d",
	4042: "Call to G2_getRepositoryLastModifiedTime() failed. Return code: %d",
	4043: "Call to G2_getVirtualEntityByRecordID(%s) failed. Return code: %d",
	4044: "Call to G2_getVirtualEntityByRecordID_V2(%s, %d) failed. Return code: %d",
	4045: "Call to G2_howEntityByEntityID(%d) failed. Return code: %d",
	4046: "Call to G2_howEntityByEntityID_V2(%d, %d) failed. Return code: %d",
	4047: "Call to G2_init(%s, %s, %d) failed. Return code: %d",
	4048: "Call to G2_initWithConfigID(%s, %s, %d, %d) failed. Return code: %d",
	4049: "Call to G2_primeEngine() failed. Return code: %d",
	4050: "Call to G2_process(%s) failed. Return code: %d",
	4051: "Call to G2_processRedoRecord() failed. Return code: %d",
	4052: "Call to G2_processRedoRecordWithInfo(%d) failed. Return code: %d",
	4053: "Call to G2_processWithInfo(%s, %d) failed. Return code: %d",
	4054: "Call to G2_processWithResponse(%s) failed. Return code: %d",
	4055: "Call to G2_processWithResponseResize(%s) failed. Return code: %d",
	4056: "Call to G2_purgeRepository() failed. Return code: %d",
	4057: "Call to G2_reevaluateEntity(%d, %d) failed. Return code: %d",
	4058: "Call to G2_reevaluateEntityWithInfo(%d, %d) failed. Return code: %d",
	4059: "Call to G2_reevaluateRecord(%s, %s, %d) failed. Return code: %d",
	4060: "Call to G2_reevaluateRecordWithInfo(%s, %s, %d) failed. Return code: %d",
	4061: "Call to G2_reinit(%d) failed. Return code: %d",
	4062: "Call to G2_replaceRecord(%s, %s, %s, %s) failed. Return code: %d",
	4063: "Call to G2_replaceRecordWithInfo(%s, %s, %s, %s, %d) failed. Return code: %d",
	4064: "Call to G2_searchByAttributes(%s) failed. Return code: %d",
	4065: "Call to G2_searchByAttributes_V2(%s, %d) failed. Return code: %d",
	4066: "Call to G2_stats() failed. Return code: %d",
	4067: "Call to G2_whyEntities(%d, %d) failed. Return code: %d",
	4068: "Call to G2_whyEntities_V2(%d, %d, %d) failed. Return code: %d",
	4069: "Call to G2_whyEntityByEntityID(%d) failed. Return code: %d",
	4070: "Call to G2_whyEntityByEntityID_V2(%d, %d) failed. Return code: %d",
	4071: "Call to G2_whyEntityByRecordID(%s, %s) failed. Return code: %d",
	4072: "Call to G2_whyEntityByRecordID_V2(%s, %s, %d) failed. Return code: %d",
	4073: "Call to G2_whyRecords(%s, %s, %s, %s) failed. Return code: %d",
	4074: "Call to G2_whyRecords_V2(%s, %s, %s, %s, %d) failed. Return code: %d",
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
	8001: "AddRecord",
	8002: "AddRecordWithInfo",
	8003: "AddRecordWithInfoWithReturnedRecordID",
	8004: "AddRecordWithReturnedRecordID",
	8005: "CheckRecord",
	8006: "CloseExport",
	8007: "CountRedoRecords",
	8008: "DeleteRecord",
	8009: "DeleteRecordWithInfo",
	8010: "Destroy",
	8011: "ExportConfig",
	8012: "ExportConfigAndConfigID",
	8013: "ExportCSVEntityReport",
	8014: "ExportJSONEntityReport",
	8015: "FetchNext",
	8016: "FindInterestingEntitiesByEntityID",
	8017: "FindInterestingEntitiesByRecordID",
	8018: "FindNetworkByEntityID",
	8019: "FindNetworkByEntityID_V2",
	8020: "FindNetworkByRecordID",
	8021: "FindNetworkByRecordID_V2",
	8022: "FindPathByEntityID",
	8023: "FindPathByEntityID_V2",
	8024: "FindPathByRecordID",
	8025: "FindPathByRecordID_V2",
	8026: "FindPathExcludingByEntityID",
	8027: "FindPathExcludingByEntityID_V2",
	8028: "FindPathExcludingByRecordID",
	8029: "FindPathExcludingByRecordID_V2",
	8030: "FindPathIncludingSourceByEntityID",
	8031: "FindPathIncludingSourceByEntityID_V2",
	8032: "FindPathIncludingSourceByRecordID",
	8033: "FindPathIncludingSourceByRecordID_V2",
	8034: "GetActiveConfigID",
	8035: "GetEntityByEntityID",
	8036: "GetEntityByEntityID_V2",
	8037: "GetEntityByRecordID",
	8038: "GetEntityByRecordID_V2",
	8039: "GetRecord",
	8040: "GetRecord_V2",
	8041: "GetRedoRecord",
	8042: "GetRepositoryLastModifiedTime",
	8043: "GetVirtualEntityByRecordID",
	8044: "GetVirtualEntityByRecordID_V2",
	8045: "HowEntityByEntityID",
	8046: "HowEntityByEntityID_V2",
	8047: "Init",
	8048: "InitWithConfigID",
	8049: "PrimeEngine",
	8050: "Process",
	8051: "ProcessRedoRecord",
	8052: "ProcessRedoRecordWithInfo",
	8053: "ProcessWithInfo",
	8054: "ProcessWithResponse",
	8055: "ProcessWithResponseResize",
	8056: "PurgeRepository",
	8057: "ReevaluateEntity",
	8058: "ReevaluateEntityWithInfo",
	8059: "ReevaluateRecord",
	8060: "ReevaluateRecordWithInfo",
	8061: "Reinit",
	8062: "ReplaceRecord",
	8063: "ReplaceRecordWithInfo",
	8064: "SearchByAttributes",
	8065: "SearchByAttributes_V2",
	8066: "Stats",
	8067: "WhyEntities",
	8068: "WhyEntities_V2",
	8069: "WhyEntityByEntityID",
	8070: "WhyEntityByEntityID_V2",
	8071: "WhyEntityByRecordID",
	8072: "WhyEntityByRecordID_V2",
	8073: "WhyRecords",
	8074: "WhyRecords_V2",
}

// Status strings for specific g2engine messages.
var IdStatuses = map[int]string{}
