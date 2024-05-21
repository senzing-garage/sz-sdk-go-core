//go:build linux

package szengine

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go/senzing"
)

// ----------------------------------------------------------------------------
// Interface functions - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzengine() {}

func ExampleSzengine_AddRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	flags := senzing.SzWithoutInfo
	result, err := szEngine.AddRecord(ctx, dataSourceCode, recordID, recordDefinition, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {}
}

func ExampleSzengine_AddRecord_secondRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1002"
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`
	flags := senzing.SzWithoutInfo
	result, err := szEngine.AddRecord(ctx, dataSourceCode, recordID, recordDefinition, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {}
}

func ExampleSzengine_AddRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "PRIMARY_NAME_MIDDLE": "J", "DATE_OF_BIRTH": "12/11/1978", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "4/9/16", "STATUS": "Inactive", "AMOUNT": "300"}`
	flags := senzing.SzWithInfo
	result, err := szEngine.AddRecord(ctx, dataSourceCode, recordID, recordDefinition, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[{"ENTITY_ID":100001}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzengine_CloseExport() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	flags := senzing.SzNoFlags
	exportHandle, err := szEngine.ExportJsonEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	err = szEngine.CloseExport(ctx, exportHandle)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzengine_CountRedoRecords() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	result, err := szEngine.CountRedoRecords(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: 2
}

func ExampleSzengine_DeleteRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	flags := senzing.SzWithoutInfo
	result, err := szEngine.DeleteRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {}
}

func ExampleSzengine_DeleteRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	flags := senzing.SzWithInfo
	result, err := szEngine.DeleteRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","AFFECTED_ENTITIES":[],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzengine_ExportCsvEntityReport() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	csvColumnList := ""
	flags := senzing.SzNoFlags
	exportHandle, err := szEngine.ExportCsvEntityReport(ctx, csvColumnList, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(exportHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzengine_ExportCsvEntityReportIterator() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	csvColumnList := ""
	flags := senzing.SzNoFlags
	for result := range szEngine.ExportCsvEntityReportIterator(ctx, csvColumnList, flags) {
		if result.Error != nil {
			fmt.Println(result.Error)
			break
		}
		fmt.Println(result.Value)
	}
	// Output: RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID
}

func ExampleSzengine_ExportJsonEntityReport() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	flags := senzing.SzNoFlags
	exportHandle, err := szEngine.ExportJsonEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(exportHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzengine_ExportJsonEntityReportIterator() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	flags := senzing.SzNoFlags
	for result := range szEngine.ExportJsonEntityReportIterator(ctx, flags) {
		if result.Error != nil {
			fmt.Println(result.Error)
			break
		}
		fmt.Println(result.Value)
	}
	// Output:
}

func ExampleSzengine_FetchNext() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	flags := senzing.SzNoFlags
	exportHandle, err := szEngine.ExportJsonEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		err = szEngine.CloseExport(ctx, exportHandle)
	}()
	jsonEntityReport := ""
	for {
		jsonEntityReportFragment, err := szEngine.FetchNext(ctx, exportHandle)
		if err != nil {
			fmt.Println(err)
		}
		if len(jsonEntityReportFragment) == 0 {
			break
		}
		jsonEntityReport += jsonEntityReportFragment
	}
}

func ExampleSzengine_FindInterestingEntitiesByEntityId() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := szEngine.FindInterestingEntitiesByEntityId(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzengine_FindInterestingEntitiesByRecordId() {
	// For more information, visit https://github.com/senzing-garage/g2-sdk-go-base/blob/main/g2engine/g2engine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := int64(0)
	result, err := szEngine.FindInterestingEntitiesByRecordId(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzengine_FindNetworkByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1001") + `}, {"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1002") + `}]}`
	maxDegrees := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	flags := senzing.SzNoFlags
	result, err := szEngine.FindNetworkByEntityId(ctx, entityList, maxDegrees, buildOutDegree, maxEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleSzengine_FindNetworkByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`
	maxDegrees := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	flags := senzing.SzNoFlags
	result, err := szEngine.FindNetworkByRecordId(ctx, recordList, maxDegrees, buildOutDegree, maxEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleSzengine_FindPathByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	startEntityID := getEntityIdForRecord("CUSTOMERS", "1001")
	endEntityID := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegrees := int64(1)
	exclusions := ""
	requiredDataSources := ""
	flags := senzing.SzNoFlags
	result, err := szEngine.FindPathByEntityId(ctx, startEntityID, endEntityID, maxDegrees, exclusions, requiredDataSources, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[100001]}],"ENTITIES":[{"RE...
}

func ExampleSzengine_FindPathByEntityId_excluding() {
	// TODO: Implement ExampleSzEngine_FindPathByEntityId_excluding
	// // For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	// ctx := context.TODO()
	// szEngine := getSzEngine(ctx)
	// startEntityId := getEntityIdForRecord("CUSTOMERS", "1001")
	// endEntityId := getEntityIdForRecord("CUSTOMERS", "1002")
	// maxDegrees := int64(1)
	// exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	// requiredDataSources := ""
	// flags := senzing.SzNoFlags
	// result, err := szEngine.FindPathByEntityId(ctx, startEntityId, endEntityId, maxDegrees, exclusions, requiredDataSources, flags)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(truncate(result, 107))
	// // Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleSzEngine_FindPathByEntityId_excludingAndIncluding() {
	// TODO: Implement ExampleSzEngine_FindPathByEntityId_excludingAndIncluding
}

func ExampleSzengine_FindPathByEntityId_including() {
	// TODO: Implement ExampleSzEngine_FindPathByEntityId_including
	// // For more information, visit https://github.com/senzing-garage/sz-sdk-go-grpc/blob/main/szengine/szengine_examples_test.go
	// ctx := context.TODO()
	// szEngine := getSzEngine(ctx)
	// startEntityId := getEntityIdForRecord("CUSTOMERS", "1001")
	// endEntityId := getEntityIdForRecord("CUSTOMERS", "1002")
	// maxDegree := int64(1)
	// exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	// requiredDataSources := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	// flags := senzing.SzNoFlags
	// result, err := szEngine.FindPathByEntityId(ctx, startEntityId, endEntityId, maxDegree, exclusions, requiredDataSources, flags)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(truncate(result, 106))
	// // Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleSzengine_FindPathByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	startDataSourceCode := "CUSTOMERS"
	startRecordID := "1001"
	endDataSourceCode := "CUSTOMERS"
	endRecordID := "1002"
	maxDegrees := int64(1)
	exclusions := ""
	requiredDataSources := ""
	flags := senzing.SzNoFlags
	result, err := szEngine.FindPathByRecordId(ctx, startDataSourceCode, startRecordID, endDataSourceCode, endRecordID, maxDegrees, exclusions, requiredDataSources, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 87))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":100001,"END_ENTITY_ID":100001,"ENTITIES":[100001...
}

func ExampleSzengine_FindPathByRecordId_excluding() {
	// TODO: Implement ExampleSzEngine_FindPathByRecordId_excluding
	// // For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	// ctx := context.TODO()
	// szEngine := getSzEngine(ctx)
	// startDataSourceCode := "CUSTOMERS"
	// startRecordId := "1001"
	// endDataSourceCode := "CUSTOMERS"
	// endRecordId := "1002"
	// maxDegree := int64(1)
	// exclusions := `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003"}]}`
	// requiredDataSources := ""
	// flags := senzing.SzNoFlags
	// result, err := szEngine.FindPathByRecordId(ctx, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegree, exclusions, requiredDataSources, flags)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(truncate(result, 107))
	// // Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleSzEngine_FindPathByRecordId_excludingAndIncluding() {
	// TODO: Implement ExampleSzEngine_FindPathByRecordId_excludingAndIncluding
}

func ExampleSzengine_FindPathByRecordId_including() {
	// TODO: Implement ExampleSzEngine_FindPathByRecordId_including
	// // For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	// ctx := context.TODO()
	// szEngine := getSzEngine(ctx)
	// startDataSourceCode := "CUSTOMERS"
	// startRecordId := "1001"
	// endDataSourceCode := "CUSTOMERS"
	// endRecordId := "1002"
	// maxDegrees := int64(1)
	// exclusions := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1003") + `}]}`
	// requiredDataSources := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	// flags := senzing.SzNoFlags
	// result, err := szEngine.FindPathByRecordId(ctx, startDataSourceCode, startRecordId, endDataSourceCode, endRecordId, maxDegrees, exclusions, requiredDataSources, flags)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(truncate(result, 119))
	// // Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleSzengine_GetActiveConfigId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	result, err := szEngine.GetActiveConfigId(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleSzengine_GetEntityByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := senzing.SzNoFlags
	result, err := szEngine.GetEntityByEntityId(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":100001}}
}

func ExampleSzengine_GetEntityByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	flags := senzing.SzNoFlags
	result, err := szEngine.GetEntityByRecordId(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":100001}}
}

func ExampleSzengine_GetRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	flags := senzing.SzNoFlags
	result, err := szEngine.GetRecord(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jsonutil.Flatten(jsonutil.Normalize(result)))
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}
}

func ExampleSzengine_GetRedoRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	result, err := szEngine.GetRedoRecord(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"REASON":"deferred delete","DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1003","DSRC_ACTION":"X"}
}

func ExampleSzengine_GetStats() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	result, err := szEngine.GetStats(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 16))
	// Output: { "workload":...
}

func ExampleSzengine_GetVirtualEntityByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`
	flags := senzing.SzNoFlags
	result, err := szEngine.GetVirtualEntityByRecordId(ctx, recordList, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":100001}}
}

func ExampleSzengine_HowEntityByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := senzing.SzNoFlags
	result, err := szEngine.HowEntityByEntityId(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jsonutil.Flatten(jsonutil.NormalizeAndSort(jsonutil.Flatten(jsonutil.Redact(result, "RECORD_ID", "INBOUND_FEAT_USAGE_TYPE")))))
	// Output: {"HOW_RESULTS":{"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"MEMBER_RECORDS":[{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":100002,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100001-S1"}]},"RESOLUTION_STEPS":[{"INBOUND_VIRTUAL_ENTITY_ID":"V100002","MATCH_INFO":{"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE"},"RESULT_VIRTUAL_ENTITY_ID":"V100001-S1","STEP":1,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":100001,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100001"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":100002,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V100002"}}]}}
}

func ExampleSzengine_PrimeEngine() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	err := szEngine.PrimeEngine(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_ProcessRedoRecord() {
	// TODO: Uncomment after it has been implemented.
	// // For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	// ctx := context.TODO()
	// szEngine := getSzEngine(ctx)
	// redoRecord, err := szEngine.GetRedoRecord(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// flags := senzing.SzWithoutInfo
	// result, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(result)
	// // Output: {}
}

func ExampleSzEngine_ProcessRedoRecord_withInfo() {
	// TODO: Uncomment after it has been implemented.
	// // For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	// ctx := context.TODO()
	// szEngine := getSzEngine(ctx)
	// redoRecord, err := szEngine.GetRedoRecord(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// flags := senzing.SzWithInfo
	// result, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(result)
	// // Output: {}
}

func ExampleSzengine_ReevaluateEntity() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	entityID := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := senzing.SzWithoutInfo
	result, err := szEngine.ReevaluateEntity(ctx, entityID, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {}
}

func ExampleSzengine_ReevaluateEntity_withInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	entityId := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := senzing.SzWithInfo
	result, err := szEngine.ReevaluateEntity(ctx, entityId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":100001}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzengine_ReevaluateRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	flags := senzing.SzWithoutInfo
	result, err := szEngine.ReevaluateRecord(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {}
}

func ExampleSzengine_ReevaluateRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	flags := senzing.SzWithInfo
	result, err := szEngine.ReevaluateRecord(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001","AFFECTED_ENTITIES":[{"ENTITY_ID":100001}],"INTERESTING_ENTITIES":{"ENTITIES":[]}}
}

func ExampleSzengine_SearchByAttributes() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	searchProfile := ""
	flags := senzing.SzNoFlags
	result, err := szEngine.SearchByAttributes(ctx, attributes, searchProfile, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jsonutil.Flatten(jsonutil.Redact(jsonutil.Flatten(jsonutil.NormalizeAndSort(result)), "FIRST_SEEN_DT", "LAST_SEEN_DT")))
	// Output: {"RESOLVED_ENTITIES":[{"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":100001}},"MATCH_INFO":{"ERRULE_CODE":"SF1","MATCH_KEY":"+PNAME+EMAIL","MATCH_LEVEL_CODE":"POSSIBLY_RELATED"}}]}
}

func ExampleSzEngine_SearchByAttributes_searchProfile() {
	// TODO: Implement ExampleSzEngine_SearchByAttributes_searchProfile
}

func ExampleSzengine_WhyEntities() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	entityId1 := getEntityId(truthset.CustomerRecords["1001"])
	entityId2 := getEntityId(truthset.CustomerRecords["1002"])
	flags := senzing.SzNoFlags
	result, err := szEngine.WhyEntities(ctx, entityId1, entityId2, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 74))
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":100001,"ENTITY_ID_2":100001,"MATCH_INFO":{...
}

func ExampleSzengine_WhyRecordInEntity() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	flags := senzing.SzNoFlags
	result, err := szEngine.WhyRecordInEntity(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":100001,"ENTITY_ID":100001,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],"MATCH_INFO":{"WHY_KEY":"+NAME+DOB+PHONE","WHY_ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_LEVEL_CODE":"RESOLVED"}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}]}
}

func ExampleSzengine_WhyRecords() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordId1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordId2 := "1002"
	flags := senzing.SzNoFlags
	result, err := szEngine.WhyRecords(ctx, dataSourceCode1, recordId1, dataSourceCode2, recordId2, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 115))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":100001,"ENTITY_ID":100001,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID"...
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzengine_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	err := szEngine.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzengine_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szEngine.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzengine_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szEngine.SetObserverOrigin(ctx, origin)
	result := szEngine.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

// ----------------------------------------------------------------------------
// Object creation / destruction
// ----------------------------------------------------------------------------

func ExampleSzengine_Initialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	instanceName := "Test module name"
	settings, err := getSettings()
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := senzing.SzNoLogging
	configId := senzing.SzInitializeWithDefaultConfiguration
	err = szEngine.Initialize(ctx, instanceName, settings, configId, verboseLogging)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzengine_Initialize_withConfigId() {
	// TODO: Implement ExampleSzEngine_Initialize_withConfigId
}

func ExampleSzengine_Reinitialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	configId, _ := szEngine.GetActiveConfigId(ctx)
	err := szEngine.Reinitialize(ctx, configId)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzengine_Destroy() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	err := szEngine.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
