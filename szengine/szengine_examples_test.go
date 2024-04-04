//go:build linux

package szengine

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-common/jsonutil"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go/sz"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzEngine_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2engine.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzEngine_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2config/g2engine_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2engine.SetObserverOrigin(ctx, origin)
	result := g2engine.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleSzEngine_AddRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	flags := int64(0)
	_, err := g2engine.AddRecord(ctx, dataSourceCode, recordId, jsonData, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_AddRecord_secondRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1002"
	jsonData := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`
	flags := int64(0)
	_, err := g2engine.AddRecord(ctx, dataSourceCode, recordId, jsonData, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_AddRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	dataSourceCode := "TEST"
	recordId := "ABC123"
	jsonData := `{"DATA_SOURCE": "TEST", "RECORD_ID": "ABC123", "NAME_FULL": "JOE SCHMOE", "DATE_OF_BIRTH": "12/11/1978", "EMAIL_ADDRESS": "joeschmoe@nowhere.com"}`
	flags := sz.SZ_WITH_INFO
	result, err := g2engine.AddRecord(ctx, dataSourceCode, recordId, jsonData, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jsonutil.Flatten(jsonutil.Redact(result, "ENTITY_ID")))
	// Output: {"AFFECTED_ENTITIES":[{"ENTITY_ID":null}],"DATA_SOURCE":"TEST","INTERESTING_ENTITIES":{"ENTITIES":[]},"RECORD_ID":"ABC123"}
}

func ExampleSzEngine_CloseExport() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	flags := int64(0)
	responseHandle, err := g2engine.ExportJsonEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	g2engine.CloseExport(ctx, responseHandle)
	// Output:
}

func ExampleSzEngine_CountRedoRecords() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	result, err := g2engine.CountRedoRecords(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: 1
}

func ExampleSzEngine_ExportCsvEntityReport() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	csvColumnList := ""
	flags := int64(0)
	responseHandle, err := g2engine.ExportCsvEntityReport(ctx, csvColumnList, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzEngine_ExportCsvEntityReportIterator() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	csvColumnList := ""
	flags := int64(0)
	for result := range g2engine.ExportCsvEntityReportIterator(ctx, csvColumnList, flags) {
		if result.Error != nil {
			fmt.Println(result.Error)
			break
		}
		fmt.Println(result.Value)
	}
	// Output: RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID
}

func ExampleSzEngine_ExportJsonEntityReport() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	flags := int64(0)
	responseHandle, err := g2engine.ExportJsonEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzEngine_ExportJsonEntityReportIterator() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	flags := int64(0)
	for result := range g2engine.ExportJsonEntityReportIterator(ctx, flags) {
		if result.Error != nil {
			fmt.Println(result.Error)
			break
		}
		fmt.Println(result.Value)
	}
	// Output:
}

func ExampleSzEngine_FetchNext() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	flags := int64(0)
	responseHandle, err := g2engine.ExportJsonEntityReport(ctx, flags)
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

func ExampleSzEngine_FindNetworkByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1001") + `}, {"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1002") + `}]}`
	maxDegrees := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	flags := int64(0)
	result, err := g2engine.FindNetworkByEntityId(ctx, entityList, maxDegrees, buildOutDegree, maxEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleSzEngine_FindNetworkByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`
	maxDegrees := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	flags := int64(0)
	result, err := g2engine.FindNetworkByRecordId(ctx, recordList, maxDegrees, buildOutDegree, maxEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"ENTITY_PATHS":[],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1}}]}
}

func ExampleSzEngine_FindPathByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	entityId1 := getEntityIdForRecord("CUSTOMERS", "1001")
	entityId2 := getEntityIdForRecord("CUSTOMERS", "1002")
	maxDegrees := int64(1)
	exclusions := ""
	requiredDataSources := ""
	flags := int64(0)
	result, err := g2engine.FindPathByEntityId(ctx, entityId1, entityId2, maxDegrees, exclusions, requiredDataSources, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 107))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":[{"RESOLVED_ENTITY":...
}

func ExampleSzEngine_FindPathByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordId1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordId2 := "1002"
	maxDegrees := int64(1)
	exclusions := ""
	requiredDataSources := ""
	flags := int64(0)
	result, err := g2engine.FindPathByRecordId(ctx, dataSourceCode1, recordId1, dataSourceCode2, recordId2, maxDegrees, exclusions, requiredDataSources, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 87))
	// Output: {"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":1,"ENTITIES":[1]}],"ENTITIES":...
}

func ExampleSzEngine_GetActiveConfigId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	result, err := g2engine.GetActiveConfigId(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleSzEngine_GetEntityByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	entityId := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.GetEntityByEntityId(ctx, entityId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleSzEngine_GetEntityByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	flags := int64(0)
	result, err := g2engine.GetEntityByRecordId(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleSzEngine_GetRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	flags := int64(0)
	result, err := g2engine.GetRecord(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jsonutil.Flatten(jsonutil.Normalize(result)))
	// Output: {"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}
}

func ExampleSzEngine_GetRedoRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	result, err := g2engine.GetRedoRecord(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"REASON":"Replaced observed entity","DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002","DSRC_ACTION":"X"}
}

func ExampleSzEngine_GetRepositoryLastModifiedTime() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	result, err := g2engine.GetRepositoryLastModifiedTime(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleSzEngine_GetStats() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	result, err := g2engine.GetStats(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 16))
	// Output: { "workload":...
}

func ExampleSzEngine_GetVirtualEntityByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`
	flags := int64(0)
	result, err := g2engine.GetVirtualEntityByRecordId(ctx, recordList, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1}}
}

func ExampleSzEngine_HowEntityByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	entityId := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.HowEntityByEntityId(ctx, entityId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jsonutil.Flatten(jsonutil.NormalizeAndSort(jsonutil.Flatten(jsonutil.Redact(result, "RECORD_ID", "INBOUND_FEAT_USAGE_TYPE")))))
	// Output: {"HOW_RESULTS":{"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[{"MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":3,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V1-S2"}]},"RESOLUTION_STEPS":[{"INBOUND_VIRTUAL_ENTITY_ID":"V1-S1","MATCH_INFO":{"ERRULE_CODE":"SF1_PNAME_CSTAB","MATCH_KEY":"+NAME+DOB+EMAIL"},"RESULT_VIRTUAL_ENTITY_ID":"V1-S2","STEP":2,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V1-S1"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":3,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V3"}},{"INBOUND_VIRTUAL_ENTITY_ID":"V2","MATCH_INFO":{"ERRULE_CODE":"CNAME_CFF_CEXCL","MATCH_KEY":"+NAME+DOB+PHONE"},"RESULT_VIRTUAL_ENTITY_ID":"V1-S1","STEP":1,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":1,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V1"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V2"}}]}}
}

func ExampleSzEngine_PrimeEngine() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	err := g2engine.PrimeEngine(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_SearchByAttributes() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	searchProfile := ""
	flags := int64(0)
	result, err := g2engine.SearchByAttributes(ctx, attributes, searchProfile, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jsonutil.Flatten(jsonutil.Redact(jsonutil.Flatten(jsonutil.NormalizeAndSort(result)), "FIRST_SEEN_DT", "LAST_SEEN_DT")))
	// Output: {"RESOLVED_ENTITIES":[{"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1}},"MATCH_INFO":{"ERRULE_CODE":"SF1","MATCH_KEY":"+PNAME+EMAIL","MATCH_LEVEL_CODE":"POSSIBLY_RELATED"}}]}
}

func ExampleSzEngine_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	err := g2engine.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_WhyEntities() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := &Szengine{}
	var entityId1 int64 = 1
	var entityId2 int64 = 100003
	var flags int64 = 0
	result, err := g2engine.WhyEntities(ctx, entityId1, entityId2, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 74))
	// Output: {"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":100003,"MATCH_INFO":{"WHY_...
}

func ExampleSzEngine_WhyRecords() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	dataSourceCode1 := "CUSTOMERS"
	recordId1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordId2 := "1002"
	flags := int64(0)
	result, err := g2engine.WhyRecords(ctx, dataSourceCode1, recordId1, dataSourceCode2, recordId2, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 116))
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":1,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],"...
}

func ExampleSzEngine_ProcessRedoRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	redoRecord := ""
	flags := int64(0)
	result, err := g2engine.ProcessRedoRecord(ctx, redoRecord, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {}
}

func ExampleSzEngine_ReevaluateEntity() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	entityId := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	_, err := g2engine.ReevaluateEntity(ctx, entityId, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_ReevaluateRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	flags := int64(0)
	_, err := g2engine.ReevaluateRecord(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_ReplaceRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	flags := int64(0)
	_, err := g2engine.ReplaceRecord(ctx, dataSourceCode, recordId, recordDefinition, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_DeleteRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1003"
	flags := int64(0)
	_, err := g2engine.DeleteRecord(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_Initialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	instanceName := "Test module name"
	settings, err := getSettings()
	if err != nil {
		fmt.Println(err)
	}
	verboseLogging := int64(0)
	configId := int64(0)
	err = g2engine.Initialize(ctx, instanceName, settings, verboseLogging, configId)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_Reinitialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	initConfigId, _ := g2engine.GetActiveConfigId(ctx) // Example initConfigId.
	err := g2engine.Reinitialize(ctx, initConfigId)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzEngine_Destroy() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getSzEngine(ctx)
	err := g2engine.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
