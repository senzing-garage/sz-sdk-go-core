//go:build linux

package g2engine

import (
	"context"
	"fmt"

	jutil "github.com/senzing/go-common/jsonutil"
	"github.com/senzing/go-logging/logging"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleG2engine_SetObserverOrigin() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2engine.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleG2engine_GetObserverOrigin() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2config/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2engine.SetObserverOrigin(ctx, origin)
	result := g2engine.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleG2engine_AddRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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

// func ExampleG2engine_AddRecordWithInfoWithReturnedRecordID() {
// 	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
// 	ctx := context.TODO()
// 	g2engine := getG2Engine(ctx)
// 	dataSourceCode := "TEST"
// 	jsonData := `{"DATA_SOURCE": "TEST", "NAME_FULL": "SUSAN SOMEBODY", "EMAIL_ADDRESS": "somesusan@somewhere.com"}`
// 	loadID := "G2Engine_test"
// 	flags := int64(0)
// 	result, _, err := g2engine.AddRecordWithInfoWithReturnedRecordID(ctx, dataSourceCode, jsonData, loadID, flags)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(jutil.Flatten(jutil.Redact(result, "ENTITY_ID", "RECORD_ID")))
// 	// Output: {"AFFECTED_ENTITIES":[{"ENTITY_ID":null}],"DATA_SOURCE":"TEST","INTERESTING_ENTITIES":{"ENTITIES":[]},"RECORD_ID":null}
// }

// func ExampleG2engine_AddRecordWithReturnedRecordID() {
// 	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
// 	ctx := context.TODO()
// 	g2engine := getG2Engine(ctx)
// 	dataSourceCode := "TEST"
// 	jsonData := `{"DATA_SOURCE": "TEST", "NAME_FULL": "JOHN DOEMAN", "EMAIL_ADDRESS": "jdoeman@anywhere.com"}`
// 	loadID := "G2Engine_test"
// 	result, err := g2engine.AddRecordWithReturnedRecordID(ctx, dataSourceCode, jsonData, loadID)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Printf("Length of record identifier is %d hexadecimal characters.\n", len(result))
// 	// Output: Length of record identifier is 40 hexadecimal characters.
// }

// func ExampleG2engine_CheckRecord() {
// 	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
// 	ctx := context.TODO()
// 	g2engine := getG2Engine(ctx)
// 	record := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
// 	recordQueryList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "123456789"}]}`
// 	result, err := g2engine.CheckRecord(ctx, record, recordQueryList)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(result)
// 	// Output: {"CHECK_RECORD_RESPONSE":[{"DSRC_CODE":"CUSTOMERS","RECORD_ID":"1001","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","MATCH_KEY":"","ERRULE_CODE":"","ERRULE_ID":0,"CANDIDATE_MATCH":"N","NON_GENERIC_CANDIDATE_MATCH":"N"}]}
// }

func ExampleG2engine_CloseExport() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.CountRedoRecords(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: 0
}

func ExampleG2engine_ExportCSVEntityReport() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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

func ExampleG2engine_ExportConfig() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	_, configId, err := g2engine.ExportConfigAndConfigID(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configId > 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_ExportJSONEntityReport() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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

func ExampleG2engine_FetchNext() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	flags := int64(0)
	responseHandle, err := g2engine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	anEntity, _ := g2engine.FetchNext(ctx, responseHandle)
	fmt.Println(len(anEntity) >= 0) // Dummy output.
	// Output: true
}

func ExampleG2engine_FindInterestingEntitiesByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITIES":[{"RELATED_ENTITIES":[],"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","LAST_SEEN_DT":null,"RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","FIRST_SEEN_DT":null,"LAST_SEEN_DT":null,"RECORD_COUNT":3}]}}],"ENTITY_PATHS":[]}
}

func ExampleG2engine_FindNetworkByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindNetworkByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITIES":[{"RELATED_ENTITIES":[],"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","LAST_SEEN_DT":null,"RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","FIRST_SEEN_DT":null,"LAST_SEEN_DT":null,"RECORD_COUNT":3}]}}],"ENTITY_PATHS":[]}
}

func ExampleG2engine_FindNetworkByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2engine_FindPathByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":...
}

func ExampleG2engine_FindPathByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathExcludingByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2engine_FindPathExcludingByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathExcludingByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2engine_FindPathExcludingByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathIncludingSourceByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleG2engine_FindPathIncludingSourceByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_FindPathIncludingSourceByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleG2engine_FindPathIncludingSourceByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_GetActiveConfigID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	result, err := g2engine.GetEntityByEntityID(ctx, entityID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleG2engine_GetEntityByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.GetEntityByEntityID_V2(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2engine_GetEntityByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2engine_GetRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.GetRedoRecord(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleG2engine_GetRepositoryLastModifiedTime() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`
	result, err := g2engine.GetVirtualEntityByRecordID(ctx, recordList)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleG2engine_GetVirtualEntityByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`
	flags := int64(0)
	result, err := g2engine.GetVirtualEntityByRecordID_V2(ctx, recordList, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleG2engine_HowEntityByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	result, err := g2engine.HowEntityByEntityID(ctx, entityID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.NormalizeAndSort(jutil.Flatten(jutil.Redact(result, "RECORD_ID", "INBOUND_FEAT_USAGE_TYPE")))))
	// Output: {"HOW_RESULTS":{"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":3,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V1-S2"}]},"RESOLUTION_STEPS":[{"INBOUND_VIRTUAL_ENTITY_ID":"V1-S1","MATCH_INFO":{"ERRULE_CODE":"SF1_PNAME_CSTAB","FEATURE_SCORES":{"DOB":[{"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_ID":2,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"INBOUND_FEAT":"12/11/1978","INBOUND_FEAT_ID":2,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FMES","SCORE_BUCKET":"SAME"}],"EMAIL":[{"CANDIDATE_FEAT":"bsmith@work.com","CANDIDATE_FEAT_ID":5,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"INBOUND_FEAT":"bsmith@work.com","INBOUND_FEAT_ID":5,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"F1","SCORE_BUCKET":"SAME"}],"NAME":[{"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_ID":32,"CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GENERATION_MATCH":-1,"GNR_FN":90,"GNR_GN":88,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Robert Smith","INBOUND_FEAT_ID":1,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"NAME","SCORE_BUCKET":"CLOSE"},{"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_ID":32,"CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GENERATION_MATCH":-1,"GNR_FN":93,"GNR_GN":93,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_ID":18,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"NAME","SCORE_BUCKET":"CLOSE"}],"RECORD_TYPE":[{"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_ID":16,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_ID":16,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FVME","SCORE_BUCKET":"SAME"}]},"MATCH_KEY":"+NAME+DOB+EMAIL"},"RESULT_VIRTUAL_ENTITY_ID":"V1-S2","STEP":2,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V1-S1"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":3,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V3"}},{"INBOUND_VIRTUAL_ENTITY_ID":"V2","MATCH_INFO":{"ERRULE_CODE":"CNAME_CFF_CEXCL","FEATURE_SCORES":{"ADDRESS":[{"CANDIDATE_FEAT":"123 Main Street, Las Vegas NV 89132","CANDIDATE_FEAT_ID":3,"CANDIDATE_FEAT_USAGE_TYPE":"MAILING","FULL_SCORE":42,"INBOUND_FEAT":"1515 Adela Lane Las Vegas NV 89111","INBOUND_FEAT_ID":20,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FF","SCORE_BUCKET":"NO_CHANCE"}],"DOB":[{"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_ID":2,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":95,"INBOUND_FEAT":"11/12/1978","INBOUND_FEAT_ID":19,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FMES","SCORE_BUCKET":"CLOSE"}],"NAME":[{"CANDIDATE_FEAT":"Robert Smith","CANDIDATE_FEAT_ID":1,"CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GENERATION_MATCH":-1,"GNR_FN":97,"GNR_GN":95,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_ID":18,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"NAME","SCORE_BUCKET":"CLOSE"}],"PHONE":[{"CANDIDATE_FEAT":"702-919-1300","CANDIDATE_FEAT_ID":4,"CANDIDATE_FEAT_USAGE_TYPE":"HOME","FULL_SCORE":100,"INBOUND_FEAT":"702-919-1300","INBOUND_FEAT_ID":4,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FF","SCORE_BUCKET":"SAME"}],"RECORD_TYPE":[{"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_ID":16,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_ID":16,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FVME","SCORE_BUCKET":"SAME"}]},"MATCH_KEY":"+NAME+DOB+PHONE"},"RESULT_VIRTUAL_ENTITY_ID":"V1-S1","STEP":1,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V1"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V2"}}]}}
}

func ExampleG2engine_HowEntityByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.HowEntityByEntityID_V2(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.NormalizeAndSort(jutil.Flatten(jutil.Redact(result, "RECORD_ID")))))
	// Output: {"HOW_RESULTS":{"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":3,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V1-S2"}]},"RESOLUTION_STEPS":[{"INBOUND_VIRTUAL_ENTITY_ID":"V1-S1","MATCH_INFO":{"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL"},"RESULT_VIRTUAL_ENTITY_ID":"V1-S2","STEP":2,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V1-S1"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":3,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V3"}},{"INBOUND_VIRTUAL_ENTITY_ID":"V2","MATCH_INFO":{"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE"},"RESULT_VIRTUAL_ENTITY_ID":"V1-S1","STEP":1,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V1"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V2"}}]}}
}

func ExampleG2engine_PrimeEngine() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.PrimeEngine(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_SearchByAttributes() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	jsonData := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	result, err := g2engine.SearchByAttributes(ctx, jsonData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.Redact(jutil.Flatten(jutil.NormalizeAndSort(result)), "FIRST_SEEN_DT", "LAST_SEEN_DT")))
	// Output: {"RESOLVED_ENTITIES":[{"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","FEATURES":{"ADDRESS":[{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","FEAT_DESC_VALUES":[{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","LIB_FEAT_ID":3}],"LIB_FEAT_ID":3,"USAGE_TYPE":"MAILING"},{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","FEAT_DESC_VALUES":[{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","LIB_FEAT_ID":20}],"LIB_FEAT_ID":20,"USAGE_TYPE":"HOME"}],"DOB":[{"FEAT_DESC":"12/11/1978","FEAT_DESC_VALUES":[{"FEAT_DESC":"11/12/1978","LIB_FEAT_ID":19},{"FEAT_DESC":"12/11/1978","LIB_FEAT_ID":2}],"LIB_FEAT_ID":2}],"EMAIL":[{"FEAT_DESC":"bsmith@work.com","FEAT_DESC_VALUES":[{"FEAT_DESC":"bsmith@work.com","LIB_FEAT_ID":5}],"LIB_FEAT_ID":5}],"NAME":[{"FEAT_DESC":"Robert Smith","FEAT_DESC_VALUES":[{"FEAT_DESC":"Bob J Smith","LIB_FEAT_ID":32},{"FEAT_DESC":"Bob Smith","LIB_FEAT_ID":18},{"FEAT_DESC":"Robert Smith","LIB_FEAT_ID":1}],"LIB_FEAT_ID":1,"USAGE_TYPE":"PRIMARY"}],"PHONE":[{"FEAT_DESC":"702-919-1300","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4}],"LIB_FEAT_ID":4,"USAGE_TYPE":"HOME"},{"FEAT_DESC":"702-919-1300","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4}],"LIB_FEAT_ID":4,"USAGE_TYPE":"MOBILE"}],"RECORD_TYPE":[{"FEAT_DESC":"PERSON","FEAT_DESC_VALUES":[{"FEAT_DESC":"PERSON","LIB_FEAT_ID":16}],"LIB_FEAT_ID":16}]},"LAST_SEEN_DT":null,"RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","FIRST_SEEN_DT":null,"LAST_SEEN_DT":null,"RECORD_COUNT":3}]}},"MATCH_INFO":{"ERRULE_CODE":"SF1","FEATURE_SCORES":{"EMAIL":[{"CANDIDATE_FEAT":"bsmith@work.com","FULL_SCORE":100,"INBOUND_FEAT":"bsmith@work.com"}],"NAME":[{"CANDIDATE_FEAT":"Bob J Smith","GENERATION_MATCH":-1,"GNR_FN":83,"GNR_GN":40,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Smith"},{"CANDIDATE_FEAT":"Robert Smith","GENERATION_MATCH":-1,"GNR_FN":88,"GNR_GN":40,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Smith"}]},"MATCH_KEY":"+PNAME+EMAIL","MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED"}}]}
}

func ExampleG2engine_SearchByAttributes_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	jsonData := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	flags := int64(0)
	result, err := g2engine.SearchByAttributes_V2(ctx, jsonData, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.NormalizeAndSort(result)))
	// Output: {"RESOLVED_ENTITIES":[{"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1}},"MATCH_INFO":{"ERRULE_CODE":"SF1","MATCH_KEY":"+PNAME+EMAIL","MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED"}}]}
}

func ExampleG2engine_SetLogLevel() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_Stats() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.Stats(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 16))
	// Output: { "workload":...
}

// FIXME: Remove after GDEV-3576 is fixed
// func ExampleG2engine_WhyEntities() {
// 	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
// 	ctx := context.TODO()
// 	g2engine := &G2engine{}
// 	var entityID1 int64 = 1
// 	var entityID2 int64 = 4
// 	result, err := g2engine.WhyEntities(ctx, entityID1, entityID2)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(truncate(result, 74))
// 	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":4,"MATCH_INFO":{"WHY_KEY":...
// }

// FIXME: Remove after GDEV-3576 is fixed
// func ExampleG2engine_WhyEntities_V2() {
// 	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
// 	ctx := context.TODO()
// 	g2engine := &G2engine{}
// 	var entityID1 int64 = 1
// 	var entityID2 int64 = 4
// 	flags := int64(0)
// 	result, err := g2engine.WhyEntities_V2(ctx, entityID1, entityID2, flags)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(result)
// 	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":4,"MATCH_INFO":{"WHY_KEY":"","WHY_ERRULE_CODE":"","MATCH_LEVEL_CODE":""}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":4}}]}
// }

func ExampleG2engine_WhyEntityByEntityID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	result, err := g2engine.WhyEntityByEntityID(ctx, entityID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 106))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":...
}

func ExampleG2engine_WhyEntityByEntityID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.WhyEntityByEntityID_V2(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 106))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":...
}

func ExampleG2engine_WhyEntityByRecordID() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	result, err := g2engine.WhyEntityByRecordID(ctx, dataSourceCode, recordID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 106))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":...
}

func ExampleG2engine_WhyEntityByRecordID_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := int64(0)
	result, err := g2engine.WhyEntityByRecordID_V2(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 106))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":...
}

func ExampleG2engine_WhyRecords() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	fmt.Println(truncate(result, 115))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],...
}

func ExampleG2engine_WhyRecords_V2() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],"INTERNAL_ID_2":2,"ENTITY_ID_2":1,"FOCUS_RECORDS_2":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002"}],"MATCH_INFO":{"WHY_KEY":"+NAME+DOB+PHONE","WHY_ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_LEVEL_CODE":"RESOLVED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleG2engine_Process() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	record := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	err := g2engine.Process(ctx, record)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

// func ExampleG2engine_ProcessRedoRecord() {
// 	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
// 	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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

func ExampleG2engine_ProcessWithInfo() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	record := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	flags := int64(0)
	result, err := g2engine.ProcessWithInfo(ctx, record, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

// func ExampleG2engine_ProcessWithResponse() {
// 	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
// 	ctx := context.TODO()
// 	g2engine := getG2Engine(ctx)
// 	record := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
// 	result, err := g2engine.ProcessWithResponse(ctx, record)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(result)
// 	// Output: {"MESSAGE": "ER SKIPPED - DUPLICATE RECORD IN G2"}
// }

// func ExampleG2engine_ProcessWithResponseResize() {
// 	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
// 	ctx := context.TODO()
// 	g2engine := getG2Engine(ctx)
// 	record := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
// 	result, err := g2engine.ProcessWithResponseResize(ctx, record)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(result)
// 	// Output: {"MESSAGE": "ER SKIPPED - DUPLICATE RECORD IN G2"}
// }

func ExampleG2engine_ReevaluateEntity() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.ReevaluateEntityWithInfo(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_ReevaluateRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":1}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleG2engine_ReplaceRecord() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
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
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	initConfigID, _ := g2engine.GetActiveConfigID(ctx) // Example initConfigID.
	err := g2engine.Reinit(ctx, initConfigID)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_PurgeRepository() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.PurgeRepository(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleG2engine_Destroy() {
	// For more information, visit https://github.com/Senzing/g2-sdk-go-base/blob/main/g2engine/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
