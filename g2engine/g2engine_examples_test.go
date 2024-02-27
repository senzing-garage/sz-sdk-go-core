//go:build linux

package g2engine

import (
	"context"
	"fmt"

	jutil "github.com/senzing-garage/go-common/jsonutil"
	"github.com/senzing-garage/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2engine_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2engine.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleG2engine_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2config/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2engine.SetObserverOrigin(ctx, origin)
	result := g2engine.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleG2engine_AddRecord() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	loadID := "G2Engine_test"
	err := g2engine.AddRecord(ctx, dataSourceCode, recordID, jsonData, loadID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_AddRecord_secondRecord() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1002"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`
	loadID := "G2Engine_test"
	err := g2engine.AddRecord(ctx, dataSourceCode, recordID, jsonData, loadID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_AddRecordWithInfo() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "TEST"
	recordID := "ABC123"
	jsonData := `{"DATA_SOURCE": "TEST", "RECORD_ID": "ABC123", "NAME_FULL": "JOE SCHMOE", "DATE_OF_BIRTH": "12/11/1978", "EMAIL_ADDRESS": "joeschmoe@nowhere.com"}`
	loadID := "G2Engine_test"
	flags := int64(0)
	result, err := g2engine.AddRecordWithInfo(ctx, dataSourceCode, recordID, jsonData, loadID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.Redact(result, "ENTITY_ID")))
	// Output: {"AFFECTED_ENTITIES":[{"ENTITY_ID":null}],"DATA_SOURCE":"TEST","INTERESTING_ENTITIES":{"ENTITIES":[]},"RECORD_ID":"ABC123"}
}

func ExampleG2engine_CloseExport() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	flags := int64(0)
	responseHandle, err := g2engine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	g2engine.CloseExport(ctx, responseHandle)
	// Output:
}

func ExampleG2engine_CountRedoRecords() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.CountRedoRecords(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: 2
}

func ExampleG2engine_ExportConfig() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.ExportConfig(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 42))
	// Output: {"G2_CONFIG":{"CFG_ETYPE":[{"ETYPE_ID":...
}

func ExampleG2engine_ExportConfigAndConfigID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	_, configId, err := g2engine.ExportConfigAndConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configId > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_ExportCSVEntityReport() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	csvColumnList := ""
	flags := int64(0)
	responseHandle, err := g2engine.ExportCSVEntityReport(ctx, csvColumnList, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseHandle > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_ExportCSVEntityReportIterator() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	csvColumnList := ""
	flags := int64(0)
	for result := range g2engine.ExportCSVEntityReportIterator(ctx, csvColumnList, flags) {
		if result.Error != nil {
			fmt.Println(result.Error)
			break
		}
		fmt.Println(result.Value)
	}
	// Output: RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL,MATCH_KEY,DATA_SOURCE,RECORD_ID
}

func ExampleG2engine_ExportJSONEntityReport() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	flags := int64(0)
	responseHandle, err := g2engine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseHandle > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_ExportJSONEntityReportIterator() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	flags := int64(0)
	for result := range g2engine.ExportJSONEntityReportIterator(ctx, flags) {
		if result.Error != nil {
			fmt.Println(result.Error)
			break
		}
		fmt.Println(result.Value)
	}
	// Output:
}

func ExampleG2engine_FetchNext() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	flags := int64(0)
	responseHandle, err := g2engine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err = g2engine.CloseExport(ctx, responseHandle)
	}()

	jsonEntityReport := ""
	for {
		jsonEntityReportFragment, err := g2engine.FetchNext(ctx, responseHandle)
		if err != nil {
			fmt.Println(err)
		}
		if len(jsonEntityReportFragment) == 0 {
			break
		}
		jsonEntityReport += jsonEntityReportFragment
	}

	fmt.Println(len(jsonEntityReport) >= 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_FindInterestingEntitiesByEntityID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.FindInterestingEntitiesByEntityID(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_FindInterestingEntitiesByRecordID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := int64(0)
	result, err := g2engine.FindInterestingEntitiesByRecordID(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_FindNetworkByEntityID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1001") + `}, {"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1002") + `}]}`
	maxDegree := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	result, err := g2engine.FindNetworkByEntityID(ctx, entityList, maxDegree, buildOutDegree, maxEntities)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.NormalizeAndSort(jutil.Flatten(jutil.Redact(result, "FIRST_SEEN_DT", "LAST_SEEN_DT")))))
	// Output: {"ENTITIES":[{"RELATED_ENTITIES":[],"RESOLVED_ENTITY":{"ENTITY_ID":100001,"ENTITY_NAME":"Robbie Smith","LAST_SEEN_DT":null,"RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","FIRST_SEEN_DT":null,"LAST_SEEN_DT":null,"RECORD_COUNT":5}]}}],"ENTITY_PATHS":[]}
}

func ExampleG2engine_FindNetworkByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1001") + `}, {"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1002") + `}]}`
	maxDegree := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	flags := int64(0)
	result, err := g2engine.FindNetworkByEntityID_V2(ctx, entityList, maxDegree, buildOutDegree, maxEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleG2engine_FindNetworkByRecordID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`
	maxDegree := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	result, err := g2engine.FindNetworkByRecordID(ctx, recordList, maxDegree, buildOutDegree, maxEntities)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.NormalizeAndSort(jutil.Flatten(jutil.Redact(result, "FIRST_SEEN_DT", "LAST_SEEN_DT")))))
	// Output: {"ENTITIES":[{"RELATED_ENTITIES":[],"RESOLVED_ENTITY":{"ENTITY_ID":100001,"ENTITY_NAME":"Robbie Smith","LAST_SEEN_DT":null,"RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","FIRST_SEEN_DT":null,"LAST_SEEN_DT":null,"RECORD_COUNT":5}]}}],"ENTITY_PATHS":[]}
}

func ExampleG2engine_FindNetworkByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`
	maxDegree := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	flags := int64(0)
	result, err := g2engine.FindNetworkByRecordID_V2(ctx, recordList, maxDegree, buildOutDegree, maxEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleG2engine_FindPathByEntityID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := int64(1)
	result, err := g2engine.FindPathByEntityID(ctx, entityID1, entityID2, maxDegree)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[100001]}],"ENTITIES":[{"RE...
}

func ExampleG2engine_FindPathByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := int64(1)
	flags := int64(0)
	result, err := g2engine.FindPathByEntityID_V2(ctx, entityID1, entityID2, maxDegree, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[100001]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleG2engine_FindPathByRecordID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := int64(1)
	result, err := g2engine.FindPathByRecordID(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 87))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[100001...
}

func ExampleG2engine_FindPathByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := int64(1)
	flags := int64(0)
	result, err := g2engine.FindPathByRecordID_V2(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[100001]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleG2engine_FindPathExcludingByEntityID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	result, err := g2engine.FindPathExcludingByEntityID(ctx, entityID1, entityID2, maxDegree, excludedEntities)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[100001]}],"ENTITIES":[{"RE...
}

func ExampleG2engine_FindPathExcludingByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	flags := int64(0)
	result, err := g2engine.FindPathExcludingByEntityID_V2(ctx, entityID1, entityID2, maxDegree, excludedEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[100001]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleG2engine_FindPathExcludingByRecordID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := int64(1)
	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003"}]}`
	result, err := g2engine.FindPathExcludingByRecordID(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedRecords)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[100001]}],"ENTITIES":[{"RE...
}

func ExampleG2engine_FindPathExcludingByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := int64(1)
	excludedRecords := `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003"}]}`
	flags := int64(0)
	result, err := g2engine.FindPathExcludingByRecordID_V2(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedRecords, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[100001]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleG2engine_FindPathIncludingSourceByEntityID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	result, err := g2engine.FindPathIncludingSourceByEntityID(ctx, entityID1, entityID2, maxDegree, excludedEntities, requiredDsrcs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 106))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[]}],"ENTITIES":[{"RESOLVE...
}

func ExampleG2engine_FindPathIncludingSourceByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityID2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	flags := int64(0)
	result, err := g2engine.FindPathIncludingSourceByEntityID_V2(ctx, entityID1, entityID2, maxDegree, excludedEntities, requiredDsrcs, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleG2engine_FindPathIncludingSourceByRecordID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	result, err := g2engine.FindPathIncludingSourceByRecordID(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedEntities, requiredDsrcs)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 119))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"E...
}

func ExampleG2engine_FindPathIncludingSourceByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	maxDegree := int64(1)
	excludedEntities := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	requiredDsrcs := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	flags := int64(0)
	result, err := g2engine.FindPathIncludingSourceByRecordID_V2(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedEntities, requiredDsrcs, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleG2engine_GetActiveConfigID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.GetActiveConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_GetEntityByEntityID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	result, err := g2engine.GetEntityByEntityID(ctx, entityID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":100001,"ENTITY_N...
}

func ExampleG2engine_GetEntityByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.GetEntityByEntityID_V2(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":100001}}
}

func ExampleG2engine_GetEntityByRecordID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	result, err := g2engine.GetEntityByRecordID(ctx, dataSourceCode, recordID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 35))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleG2engine_GetEntityByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := int64(0)
	result, err := g2engine.GetEntityByRecordID_V2(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":100001}}
}

func ExampleG2engine_GetRecord() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	result, err := g2engine.GetRecord(ctx, dataSourceCode, recordID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.Normalize(result)))
	// Output: {"DATA_SOURCE":"CUSTOMERS","JSON_DATA":{"ADDR_LINE1":"123 Main Street, Las Vegas NV 89132","ADDR_TYPE":"MAILING","AMOUNT":"100","DATA_SOURCE":"CUSTOMERS","DATE":"1/2/18","DATE_OF_BIRTH":"12/11/1978","EMAIL_ADDRESS":"bsmith@work.com","PHONE_NUMBER":"702-919-1300","PHONE_TYPE":"HOME","PRIMARY_NAME_FIRST":"Robert","PRIMARY_NAME_LAST":"Smith","RECORD_ID":"1001","RECORD_TYPE":"PERSON","STATUS":"Active"},"RECORD_ID":"1001"}
}

func ExampleG2engine_GetRecord_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := int64(0)
	result, err := g2engine.GetRecord_V2(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.Normalize(result)))
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}
}

func ExampleG2engine_GetRedoRecord() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.GetRedoRecord(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"REASON":"deferred delete","DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002","ENTITY_TYPE":"GENERIC","DSRC_ACTION":"X"}
}

func ExampleG2engine_GetRepositoryLastModifiedTime() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.GetRepositoryLastModifiedTime(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_GetVirtualEntityByRecordID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`
	result, err := g2engine.GetVirtualEntityByRecordID(ctx, recordList)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":100002,"ENTITY_N...
}

func ExampleG2engine_GetVirtualEntityByRecordID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`
	flags := int64(0)
	result, err := g2engine.GetVirtualEntityByRecordID_V2(ctx, recordList, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":100002}}
}

func ExampleG2engine_HowEntityByEntityID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	result, err := g2engine.HowEntityByEntityID(ctx, entityID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.NormalizeAndSort(jutil.Flatten(jutil.Redact(result, "RECORD_ID", "INBOUND_FEAT_USAGE_TYPE")))))
	// Output: {"HOW_RESULTS":{"FINAL_STATE":{"NEED_REEVALUATION":1,"VIRTUAL_ENTITIES":[{"MEMBER_RECORDS":[{"INTERNAL_ID":100002,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100003,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100004,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100005,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":200009,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100002-S4"}]},"RESOLUTION_STEPS":[{"INBOUND_VIRTUAL_ENTITY_ID":"V100002-S1","MATCH_INFO":{"ERRULE_CODE":"SF1_SNAME_CFF_CSTAB","FEATURE_SCORES":{"ADDRESS":[{"CANDIDATE_FEAT":"123 Main Street, Las Vegas NV 89132","CANDIDATE_FEAT_ID":100003,"CANDIDATE_FEAT_USAGE_TYPE":"MAILING","FULL_SCORE":33,"INBOUND_FEAT":"1515 Adela Ln Las Vegas NV 89132","INBOUND_FEAT_ID":100040,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FF","SCORE_BUCKET":"NO_CHANCE"},{"CANDIDATE_FEAT":"123 Main Street, Las Vegas NV 89132","CANDIDATE_FEAT_ID":100003,"CANDIDATE_FEAT_USAGE_TYPE":"MAILING","FULL_SCORE":42,"INBOUND_FEAT":"1515 Adela Lane Las Vegas NV 89111","INBOUND_FEAT_ID":100020,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FF","SCORE_BUCKET":"NO_CHANCE"}],"DOB":[{"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_ID":100002,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":88,"INBOUND_FEAT":"11/12/1979","INBOUND_FEAT_ID":100039,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FMES","SCORE_BUCKET":"PLAUSIBLE"},{"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_ID":100002,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":95,"INBOUND_FEAT":"11/12/1978","INBOUND_FEAT_ID":100019,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FMES","SCORE_BUCKET":"CLOSE"}],"EMAIL":[{"CANDIDATE_FEAT":"bsmith@work.com","CANDIDATE_FEAT_ID":100005,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"INBOUND_FEAT":"bsmith@work.com","INBOUND_FEAT_ID":100005,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"F1","SCORE_BUCKET":"SAME"}],"NAME":[{"CANDIDATE_FEAT":"Robert Smith","CANDIDATE_FEAT_ID":100001,"CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GENERATION_MATCH":-1,"GNR_FN":97,"GNR_GN":95,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_ID":100018,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"NAME","SCORE_BUCKET":"CLOSE"}],"PHONE":[{"CANDIDATE_FEAT":"702-919-1300","CANDIDATE_FEAT_ID":100004,"CANDIDATE_FEAT_USAGE_TYPE":"HOME","FULL_SCORE":100,"INBOUND_FEAT":"702-919-1300","INBOUND_FEAT_ID":100004,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FF","SCORE_BUCKET":"SAME"}],"RECORD_TYPE":[{"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_ID":100016,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_ID":100016,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FVME","SCORE_BUCKET":"SAME"}]},"MATCH_KEY":"+NAME+DOB+PHONE+EMAIL"},"RESULT_VIRTUAL_ENTITY_ID":"V100002-S2","STEP":2,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":100002,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100004,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100002-S1"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":200009,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V200009"}},{"INBOUND_VIRTUAL_ENTITY_ID":"V100002-S3","MATCH_INFO":{"ERRULE_CODE":"SF1_PNAME_CSTAB","FEATURE_SCORES":{"DOB":[{"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_ID":100002,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"INBOUND_FEAT":"12/11/1978","INBOUND_FEAT_ID":100002,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FMES","SCORE_BUCKET":"SAME"}],"EMAIL":[{"CANDIDATE_FEAT":"bsmith@work.com","CANDIDATE_FEAT_ID":100005,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"INBOUND_FEAT":"bsmith@work.com","INBOUND_FEAT_ID":100005,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"F1","SCORE_BUCKET":"SAME"}],"NAME":[{"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_ID":100032,"CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GENERATION_MATCH":-1,"GNR_FN":90,"GNR_GN":88,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Robbie Smith","INBOUND_FEAT_ID":100048,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"NAME","SCORE_BUCKET":"CLOSE"},{"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_ID":100032,"CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GENERATION_MATCH":-1,"GNR_FN":90,"GNR_GN":88,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Robert Smith","INBOUND_FEAT_ID":100001,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"NAME","SCORE_BUCKET":"CLOSE"},{"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_ID":100032,"CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GENERATION_MATCH":-1,"GNR_FN":93,"GNR_GN":93,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_ID":100018,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"NAME","SCORE_BUCKET":"CLOSE"}],"RECORD_TYPE":[{"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_ID":100016,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_ID":100016,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FVME","SCORE_BUCKET":"SAME"}]},"MATCH_KEY":"+NAME+DOB+EMAIL"},"RESULT_VIRTUAL_ENTITY_ID":"V100002-S4","STEP":4,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":100002,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100004,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100005,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":200009,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100002-S3"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":100003,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100003"}},{"INBOUND_VIRTUAL_ENTITY_ID":"V100004","MATCH_INFO":{"ERRULE_CODE":"CNAME_CFF_CEXCL","FEATURE_SCORES":{"ADDRESS":[{"CANDIDATE_FEAT":"1515 Adela Lane Las Vegas NV 89111","CANDIDATE_FEAT_ID":100020,"CANDIDATE_FEAT_USAGE_TYPE":"HOME","FULL_SCORE":96,"INBOUND_FEAT":"1515 Adela Ln Las Vegas NV 89132","INBOUND_FEAT_ID":100040,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FF","SCORE_BUCKET":"CLOSE"}],"DOB":[{"CANDIDATE_FEAT":"11/12/1978","CANDIDATE_FEAT_ID":100019,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":92,"INBOUND_FEAT":"11/12/1979","INBOUND_FEAT_ID":100039,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FMES","SCORE_BUCKET":"CLOSE"}],"NAME":[{"CANDIDATE_FEAT":"Bob Smith","CANDIDATE_FEAT_ID":100018,"CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GENERATION_MATCH":-1,"GNR_FN":92,"GNR_GN":85,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"B Smith","INBOUND_FEAT_ID":100038,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"NAME","SCORE_BUCKET":"CLOSE"}],"RECORD_TYPE":[{"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_ID":100016,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_ID":100016,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FVME","SCORE_BUCKET":"SAME"}]},"MATCH_KEY":"+NAME+DOB+ADDRESS"},"RESULT_VIRTUAL_ENTITY_ID":"V100002-S1","STEP":1,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":100002,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100002"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":100004,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100004"}},{"INBOUND_VIRTUAL_ENTITY_ID":"V100005","MATCH_INFO":{"ERRULE_CODE":"CNAME_CFF","FEATURE_SCORES":{"ADDRESS":[{"CANDIDATE_FEAT":"123 Main Street, Las Vegas NV 89132","CANDIDATE_FEAT_ID":100003,"CANDIDATE_FEAT_USAGE_TYPE":"MAILING","FULL_SCORE":91,"INBOUND_FEAT":"123 E Main St Henderson NV 89132","INBOUND_FEAT_ID":100049,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FF","SCORE_BUCKET":"CLOSE"},{"CANDIDATE_FEAT":"1515 Adela Lane Las Vegas NV 89111","CANDIDATE_FEAT_ID":100020,"CANDIDATE_FEAT_USAGE_TYPE":"HOME","FULL_SCORE":21,"INBOUND_FEAT":"123 E Main St Henderson NV 89132","INBOUND_FEAT_ID":100049,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FF","SCORE_BUCKET":"NO_CHANCE"},{"CANDIDATE_FEAT":"1515 Adela Ln Las Vegas NV 89132","CANDIDATE_FEAT_ID":100040,"CANDIDATE_FEAT_USAGE_TYPE":"HOME","FULL_SCORE":33,"INBOUND_FEAT":"123 E Main St Henderson NV 89132","INBOUND_FEAT_ID":100049,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FF","SCORE_BUCKET":"NO_CHANCE"}],"NAME":[{"CANDIDATE_FEAT":"Bob Smith","CANDIDATE_FEAT_ID":100018,"CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GENERATION_MATCH":-1,"GNR_FN":97,"GNR_GN":95,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Robbie Smith","INBOUND_FEAT_ID":100048,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"NAME","SCORE_BUCKET":"CLOSE"},{"CANDIDATE_FEAT":"Robert Smith","CANDIDATE_FEAT_ID":100001,"CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GENERATION_MATCH":-1,"GNR_FN":97,"GNR_GN":95,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Robbie Smith","INBOUND_FEAT_ID":100048,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"NAME","SCORE_BUCKET":"CLOSE"}],"RECORD_TYPE":[{"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_ID":100016,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_ID":100016,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FVME","SCORE_BUCKET":"SAME"}]},"MATCH_KEY":"+NAME+ADDRESS"},"RESULT_VIRTUAL_ENTITY_ID":"V100002-S3","STEP":3,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":100002,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100004,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":200009,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100002-S2"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":100005,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100005"}}]}}
}

func ExampleG2engine_HowEntityByEntityID_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.HowEntityByEntityID_V2(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.NormalizeAndSort(jutil.Flatten(jutil.Redact(result, "RECORD_ID")))))
	// Output: {"HOW_RESULTS":{"FINAL_STATE":{"NEED_REEVALUATION":1,"VIRTUAL_ENTITIES":[{"MEMBER_RECORDS":[{"INTERNAL_ID":100002,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100003,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100004,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100005,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":200009,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100002-S4"}]},"RESOLUTION_STEPS":[{"INBOUND_VIRTUAL_ENTITY_ID":"V100002-S1","MATCH_INFO":{"ERRULE_CODE":"SF1_SNAME_CFF_CSTAB","MATCH_KEY":"+NAME+DOB+PHONE+EMAIL"},"RESULT_VIRTUAL_ENTITY_ID":"V100002-S2","STEP":2,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":100002,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100004,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100002-S1"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":200009,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V200009"}},{"INBOUND_VIRTUAL_ENTITY_ID":"V100002-S3","MATCH_INFO":{"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL"},"RESULT_VIRTUAL_ENTITY_ID":"V100002-S4","STEP":4,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":100002,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100004,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100005,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":200009,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100002-S3"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":100003,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100003"}},{"INBOUND_VIRTUAL_ENTITY_ID":"V100004","MATCH_INFO":{"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+ADDRESS"},"RESULT_VIRTUAL_ENTITY_ID":"V100002-S1","STEP":1,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":100002,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100002"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":100004,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100004"}},{"INBOUND_VIRTUAL_ENTITY_ID":"V100005","MATCH_INFO":{"ERRULE_CODE":"CNAME_CFF","MATCH_KEY":"+NAME+ADDRESS"},"RESULT_VIRTUAL_ENTITY_ID":"V100002-S3","STEP":3,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":100002,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100004,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":200009,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100002-S2"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":100005,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100005"}}]}}
}

func ExampleG2engine_PrimeEngine() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.PrimeEngine(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_SearchByAttributes() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	jsonData := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	result, err := g2engine.SearchByAttributes(ctx, jsonData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.Redact(jutil.Flatten(jutil.NormalizeAndSort(result)), "FIRST_SEEN_DT", "LAST_SEEN_DT")))
	// Output: {"RESOLVED_ENTITIES":[{"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":100001,"ENTITY_NAME":"Robbie Smith","FEATURES":{"ADDRESS":[{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","FEAT_DESC_VALUES":[{"FEAT_DESC":"123 E Main St Henderson NV 89132","LIB_FEAT_ID":100049},{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","LIB_FEAT_ID":100003}],"LIB_FEAT_ID":100003,"USAGE_TYPE":"MAILING"},{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","FEAT_DESC_VALUES":[{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","LIB_FEAT_ID":100020},{"FEAT_DESC":"1515 Adela Ln Las Vegas NV 89132","LIB_FEAT_ID":100040}],"LIB_FEAT_ID":100020,"USAGE_TYPE":"HOME"}],"DOB":[{"FEAT_DESC":"11/12/1979","FEAT_DESC_VALUES":[{"FEAT_DESC":"11/12/1979","LIB_FEAT_ID":100039}],"LIB_FEAT_ID":100039},{"FEAT_DESC":"12/11/1978","FEAT_DESC_VALUES":[{"FEAT_DESC":"11/12/1978","LIB_FEAT_ID":100019},{"FEAT_DESC":"12/11/1978","LIB_FEAT_ID":100002}],"LIB_FEAT_ID":100002}],"DRLIC":[{"FEAT_DESC":"112233 NV","FEAT_DESC_VALUES":[{"FEAT_DESC":"112233 NV","LIB_FEAT_ID":100050}],"LIB_FEAT_ID":100050}],"EMAIL":[{"FEAT_DESC":"bsmith@work.com","FEAT_DESC_VALUES":[{"FEAT_DESC":"bsmith@work.com","LIB_FEAT_ID":100005}],"LIB_FEAT_ID":100005}],"NAME":[{"FEAT_DESC":"B Smith","FEAT_DESC_VALUES":[{"FEAT_DESC":"B Smith","LIB_FEAT_ID":100038}],"LIB_FEAT_ID":100038,"USAGE_TYPE":"PRIMARY"},{"FEAT_DESC":"Robert Smith","FEAT_DESC_VALUES":[{"FEAT_DESC":"Bob J Smith","LIB_FEAT_ID":100032},{"FEAT_DESC":"Bob Smith","LIB_FEAT_ID":100018},{"FEAT_DESC":"Robbie Smith","LIB_FEAT_ID":100048},{"FEAT_DESC":"Robert Smith","LIB_FEAT_ID":100001}],"LIB_FEAT_ID":100001,"USAGE_TYPE":"PRIMARY"}],"PHONE":[{"FEAT_DESC":"702-919-1300","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":100004}],"LIB_FEAT_ID":100004,"USAGE_TYPE":"HOME"},{"FEAT_DESC":"702-919-1300","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":100004}],"LIB_FEAT_ID":100004,"USAGE_TYPE":"MOBILE"}],"RECORD_TYPE":[{"FEAT_DESC":"PERSON","FEAT_DESC_VALUES":[{"FEAT_DESC":"PERSON","LIB_FEAT_ID":100016}],"LIB_FEAT_ID":100016}]},"LAST_SEEN_DT":null,"RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","FIRST_SEEN_DT":null,"LAST_SEEN_DT":null,"RECORD_COUNT":5}]}},"MATCH_INFO":{"ERRULE_CODE":"SF1","FEATURE_SCORES":{"EMAIL":[{"CANDIDATE_FEAT":"bsmith@work.com","FULL_SCORE":100,"INBOUND_FEAT":"bsmith@work.com"}],"NAME":[{"CANDIDATE_FEAT":"Bob J Smith","GENERATION_MATCH":-1,"GNR_FN":83,"GNR_GN":40,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Smith"},{"CANDIDATE_FEAT":"Robbie Smith","GENERATION_MATCH":-1,"GNR_FN":88,"GNR_GN":40,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Smith"},{"CANDIDATE_FEAT":"Robert Smith","GENERATION_MATCH":-1,"GNR_FN":88,"GNR_GN":40,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Smith"}]},"MATCH_KEY":"+PNAME+EMAIL","MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED"}}]}
}

func ExampleG2engine_SearchByAttributes_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	jsonData := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	flags := int64(0)
	result, err := g2engine.SearchByAttributes_V2(ctx, jsonData, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.NormalizeAndSort(result)))
	// Output: {"RESOLVED_ENTITIES":[{"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":100001}},"MATCH_INFO":{"ERRULE_CODE":"SF1","MATCH_KEY":"+PNAME+EMAIL","MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED"}}]}
}

func ExampleG2engine_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_Stats() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.Stats(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 16))
	// Output: { "workload":...
}

func ExampleG2engine_WhyEntities() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := &G2engine{}
	var entityID1 int64 = 1
	var entityID2 int64 = 100006
	result, err := g2engine.WhyEntities(ctx, entityID1, entityID2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 74))
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":100006,"MATCH_INFO":{"WHY_...
}

func ExampleG2engine_WhyEntities_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := &G2engine{}
	var entityID1 int64 = 1
	var entityID2 int64 = 100006
	flags := int64(0)
	result, err := g2engine.WhyEntities_V2(ctx, entityID1, entityID2, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":100006,"MATCH_INFO":{"WHY_KEY":"","WHY_ERRULE_CODE":"","MATCH_LEVEL_CODE":""}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}},{"RESOLVED_ENTITY":{"ENTITY_ID":100006}}]}
}

func ExampleG2engine_WhyRecords() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	result, err := g2engine.WhyRecords(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 116))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":200009,"ENTITY_ID":100001,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":...
}

func ExampleG2engine_WhyRecords_V2() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	flags := int64(0)
	result, err := g2engine.WhyRecords_V2(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":200009,"ENTITY_ID":100001,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],"INTERNAL_ID_2":100002,"ENTITY_ID_2":100001,"FOCUS_RECORDS_2":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}],"MATCH_INFO":{"WHY_KEY":"+NAME+DOB+PHONE","WHY_ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_LEVEL_CODE":"RESOLVED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

// func ExampleG2engine_ProcessRedoRecord() {
// 	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
// 	ctx := context.TODO()
// 	g2engine := getG2Engine(ctx)
// 	result, err := g2engine.ProcessRedoRecord(ctx)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(result)
// 	// Output:
// }

// func ExampleG2engine_ProcessRedoRecordWithInfo() {
// 	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
// 	ctx := context.TODO()
// 	g2engine := getG2Engine(ctx)
// 	flags := int64(0)
// 	_, result, err := g2engine.ProcessRedoRecordWithInfo(ctx, flags)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(result)
// 	// Output:
// }

func ExampleG2engine_ReevaluateEntity() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	err := g2engine.ReevaluateEntity(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
func ExampleG2engine_ReevaluateEntityWithInfo() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.ReevaluateEntityWithInfo(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":100002}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_ReevaluateRecord() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := int64(0)
	err := g2engine.ReevaluateRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_ReevaluateRecordWithInfo() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := int64(0)
	result, err := g2engine.ReevaluateRecordWithInfo(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":100002}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_ReplaceRecord() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	loadID := "G2Engine_test"
	err := g2engine.ReplaceRecord(ctx, dataSourceCode, recordID, jsonData, loadID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_ReplaceRecordWithInfo() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	loadID := "G2Engine_test"
	flags := int64(0)
	result, err := g2engine.ReplaceRecordWithInfo(ctx, dataSourceCode, recordID, jsonData, loadID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_DeleteRecord() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	loadID := "G2Engine_test"
	err := g2engine.DeleteRecord(ctx, dataSourceCode, recordID, loadID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_DeleteRecordWithInfo() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	loadID := "G2Engine_test"
	flags := int64(0)
	result, err := g2engine.DeleteRecordWithInfo(ctx, dataSourceCode, recordID, loadID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_Init() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	moduleName := "Test module name"
	iniParams, err := getIniParams()
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := int64(0)
	err = g2engine.Init(ctx, moduleName, iniParams, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_InitWithConfigID() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	moduleName := "Test module name"
	iniParams, err := getIniParams()
	if err != nil {
		fmt.Println(err)
	}
	initConfigID := senzingConfigId
	verboseLogging := int64(0)
	err = g2engine.InitWithConfigID(ctx, moduleName, iniParams, initConfigID, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_Reinit() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	initConfigID, _ := g2engine.GetActiveConfigID(ctx) // Example initConfigID.
	err := g2engine.Reinit(ctx, initConfigID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_Destroy() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
