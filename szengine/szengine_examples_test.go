//go:build linux

package szengine

import (
	"context"
	"fmt"

	jutil "github.com/senzing-garage/go-common/jsonutil"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go/sz"
)

// ----------------------------------------------------------------------------
// Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzengine_SetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2engine.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzengine_GetObserverOrigin() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/g2config/g2engine_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	g2engine.SetObserverOrigin(ctx, origin)
	result := g2engine.GetObserverOrigin(ctx)
	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}

func ExampleSzengine_AddRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
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

func ExampleSzengine_AddRecord_secondRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
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

func ExampleSzengine_AddRecord_withInfo() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "TEST"
	recordId := "ABC123"
	jsonData := `{"DATA_SOURCE": "TEST", "RECORD_ID": "ABC123", "NAME_FULL": "JOE SCHMOE", "DATE_OF_BIRTH": "12/11/1978", "EMAIL_ADDRESS": "joeschmoe@nowhere.com"}`
	flags := sz.SZ_WITH_INFO
	result, err := g2engine.AddRecord(ctx, dataSourceCode, recordId, jsonData, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.Redact(result, "ENTITY_ID")))
	// Output: {"AFFECTED_ENTITIES":[{"ENTITY_ID":null}],"DATA_SOURCE":"TEST","INTERESTING_ENTITIES":{"ENTITIES":[]},"RECORD_ID":"ABC123"}
}

func ExampleSzengine_CloseExport() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	flags := int64(0)
	responseHandle, err := g2engine.ExportJsonEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	g2engine.CloseExport(ctx, responseHandle)
	// Output:
}

func ExampleSzengine_CountRedoRecords() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.CountRedoRecords(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: 1
}

func ExampleSzengine_ExportCsvEntityReport() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	csvColumnList := ""
	flags := int64(0)
	responseHandle, err := g2engine.ExportCsvEntityReport(ctx, csvColumnList, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzengine_ExportCsvEntityReportIterator() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	csvColumnList := ""
	flags := int64(0)
	for result := range g2engine.ExportCsvEntityReportIterator(ctx, csvColumnList, flags) {
		if result.Error != nil {
			fmt.Println(result.Error)
			break
		}
		fmt.Println(result.Value)
	}
	// Output: RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL,MATCH_KEY,DATA_SOURCE,RECORD_ID
}

func ExampleSzengine_ExportJsonEntityReport() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	flags := int64(0)
	responseHandle, err := g2engine.ExportJsonEntityReport(ctx, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(responseHandle > 0) // Dummy output.
	// Output: true
}

func ExampleSzengine_ExportJsonEntityReportIterator() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
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

func ExampleSzengine_FetchNext() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
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

func ExampleSzengine_FindNetworkByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1001") + `}, {"ENTITY_ID": ` + getEntityIdStringForRecord("CUSTOMERS", "1002") + `}]}`
	maxDegrees := int64(2)
	buildOutDegree := int64(1)
	maxEntities := int64(10)
	flags := int64(0)
	result, err := g2engine.FindNetworkByEntityId(ctx, entityList, maxDegrees, buildOutDegree, maxEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.NormalizeAndSort(jutil.Flatten(jutil.Redact(result, "FIRST_SEEN_DT", "LAST_SEEN_DT")))))
	// Output: {"ENTITIES":[{"RELATED_ENTITIES":[],"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","LAST_SEEN_DT":null,"RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","FIRST_SEEN_DT":null,"LAST_SEEN_DT":null,"RECORD_COUNT":3}]}}],"ENTITY_PATHS":[]}
}

func ExampleSzengine_FindNetworkByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001"}, {"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002"}]}`
	maxDegrees := int64(1)
	buildOutDegree := int64(2)
	maxEntities := int64(10)
	flags := int64(0)
	result, err := g2engine.FindNetworkByRecordId(ctx, recordList, maxDegrees, buildOutDegree, maxEntities, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.NormalizeAndSort(jutil.Flatten(jutil.Redact(result, "FIRST_SEEN_DT", "LAST_SEEN_DT")))))
	// Output: {"ENTITIES":[{"RELATED_ENTITIES":[],"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","LAST_SEEN_DT":null,"RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","FIRST_SEEN_DT":null,"LAST_SEEN_DT":null,"RECORD_COUNT":3}]}}],"ENTITY_PATHS":[]}
}

func ExampleSzengine_FindPathByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
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

func ExampleSzengine_FindPathByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
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

func ExampleSzengine_GetActiveConfigId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.GetActiveConfigId(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleSzengine_GetEntityByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityId := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.GetEntityByEntityId(ctx, entityId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":...
}

func ExampleSzengine_GetEntityByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	flags := int64(0)
	result, err := g2engine.GetEntityByRecordId(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 35))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":...
}

func ExampleSzengine_GetRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	flags := int64(0)
	result, err := g2engine.GetRecord(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.Normalize(result)))
	// Output: {"DATA_SOURCE":"CUSTOMERS","JSON_DATA":{"ADDR_LINE1":"123 Main Street, Las Vegas NV 89132","ADDR_TYPE":"MAILING","AMOUNT":"100","DATA_SOURCE":"CUSTOMERS","DATE":"1/2/18","DATE_OF_BIRTH":"12/11/1978","EMAIL_ADDRESS":"bsmith@work.com","PHONE_NUMBER":"702-919-1300","PHONE_TYPE":"HOME","PRIMARY_NAME_FIRST":"Robert","PRIMARY_NAME_LAST":"Smith","RECORD_ID":"1001","RECORD_TYPE":"PERSON","STATUS":"Active"},"RECORD_ID":"1001"}
}

func ExampleSzengine_GetRedoRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.GetRedoRecord(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output: {"REASON":"deferred delete","DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1002","ENTITY_TYPE":"GENERIC","DSRC_ACTION":"X"}
}

func ExampleSzengine_GetRepositoryLastModifiedTime() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.GetRepositoryLastModifiedTime(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleSzengine_GetStats() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	result, err := g2engine.GetStats(ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 16))
	// Output: { "workload":...
}

func ExampleSzengine_GetVirtualEntityByRecordId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`
	flags := int64(0)
	result, err := g2engine.GetVirtualEntityByRecordId(ctx, recordList, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(truncate(result, 51))
	// Output: {"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":...
}

func ExampleSzengine_HowEntityByEntityId() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityId := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	result, err := g2engine.HowEntityByEntityId(ctx, entityId, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.NormalizeAndSort(jutil.Flatten(jutil.Redact(result, "RECORD_ID", "INBOUND_FEAT_USAGE_TYPE")))))
	// Output: {"HOW_RESULTS":{"FINAL_STATE":{"NEED_REEVALUATION":1,"VIRTUAL_ENTITIES":[{"MEMBER_RECORDS":[{"INTERNAL_ID":12,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":3,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V2-S2"}]},"RESOLUTION_STEPS":[{"INBOUND_VIRTUAL_ENTITY_ID":"V2","MATCH_INFO":{"ERRULE_CODE":"CNAME_CFF_CEXCL","FEATURE_SCORES":{"ADDRESS":[{"CANDIDATE_FEAT":"123 Main Street, Las Vegas NV 89132","CANDIDATE_FEAT_ID":3,"CANDIDATE_FEAT_USAGE_TYPE":"MAILING","FULL_SCORE":42,"INBOUND_FEAT":"1515 Adela Lane Las Vegas NV 89111","INBOUND_FEAT_ID":20,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FF","SCORE_BUCKET":"NO_CHANCE"}],"DOB":[{"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_ID":2,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":95,"INBOUND_FEAT":"11/12/1978","INBOUND_FEAT_ID":19,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FMES","SCORE_BUCKET":"CLOSE"}],"NAME":[{"CANDIDATE_FEAT":"Robert Smith","CANDIDATE_FEAT_ID":1,"CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GENERATION_MATCH":-1,"GNR_FN":97,"GNR_GN":95,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_ID":18,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"NAME","SCORE_BUCKET":"CLOSE"}],"PHONE":[{"CANDIDATE_FEAT":"702-919-1300","CANDIDATE_FEAT_ID":4,"CANDIDATE_FEAT_USAGE_TYPE":"HOME","FULL_SCORE":100,"INBOUND_FEAT":"702-919-1300","INBOUND_FEAT_ID":4,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FF","SCORE_BUCKET":"SAME"}],"RECORD_TYPE":[{"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_ID":16,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_ID":16,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FVME","SCORE_BUCKET":"SAME"}]},"MATCH_KEY":"+NAME+DOB+PHONE"},"RESULT_VIRTUAL_ENTITY_ID":"V2-S1","STEP":1,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V2"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":12,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V12"}},{"INBOUND_VIRTUAL_ENTITY_ID":"V2-S1","MATCH_INFO":{"ERRULE_CODE":"SF1_PNAME_CSTAB","FEATURE_SCORES":{"DOB":[{"CANDIDATE_FEAT":"12/11/1978","CANDIDATE_FEAT_ID":2,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"INBOUND_FEAT":"12/11/1978","INBOUND_FEAT_ID":2,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FMES","SCORE_BUCKET":"SAME"}],"EMAIL":[{"CANDIDATE_FEAT":"bsmith@work.com","CANDIDATE_FEAT_ID":5,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"INBOUND_FEAT":"bsmith@work.com","INBOUND_FEAT_ID":5,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"F1","SCORE_BUCKET":"SAME"}],"NAME":[{"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_ID":32,"CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GENERATION_MATCH":-1,"GNR_FN":90,"GNR_GN":88,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Robert Smith","INBOUND_FEAT_ID":1,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"NAME","SCORE_BUCKET":"CLOSE"},{"CANDIDATE_FEAT":"Bob J Smith","CANDIDATE_FEAT_ID":32,"CANDIDATE_FEAT_USAGE_TYPE":"PRIMARY","GENERATION_MATCH":-1,"GNR_FN":93,"GNR_GN":93,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Bob Smith","INBOUND_FEAT_ID":18,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"NAME","SCORE_BUCKET":"CLOSE"}],"RECORD_TYPE":[{"CANDIDATE_FEAT":"PERSON","CANDIDATE_FEAT_ID":16,"CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"INBOUND_FEAT":"PERSON","INBOUND_FEAT_ID":16,"INBOUND_FEAT_USAGE_TYPE":null,"SCORE_BEHAVIOR":"FVME","SCORE_BUCKET":"SAME"}]},"MATCH_KEY":"+NAME+DOB+EMAIL"},"RESULT_VIRTUAL_ENTITY_ID":"V2-S2","STEP":2,"VIRTUAL_ENTITY_1":{"MEMBER_RECORDS":[{"INTERNAL_ID":12,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]},{"INTERNAL_ID":2,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V2-S1"},"VIRTUAL_ENTITY_2":{"MEMBER_RECORDS":[{"INTERNAL_ID":3,"RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":null}]}],"VIRTUAL_ENTITY_ID":"V3"}}]}}
}

func ExampleSzengine_PrimeEngine() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.PrimeEngine(ctx)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzengine_SearchByAttributes() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	searchProfile := ""
	flags := int64(0)
	result, err := g2engine.SearchByAttributes(ctx, attributes, searchProfile, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(jutil.Flatten(jutil.Redact(jutil.Flatten(jutil.NormalizeAndSort(result)), "FIRST_SEEN_DT", "LAST_SEEN_DT")))
	// Output: {"RESOLVED_ENTITIES":[{"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"Robert Smith","FEATURES":{"ADDRESS":[{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","FEAT_DESC_VALUES":[{"FEAT_DESC":"123 Main Street, Las Vegas NV 89132","LIB_FEAT_ID":3}],"LIB_FEAT_ID":3,"USAGE_TYPE":"MAILING"},{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","FEAT_DESC_VALUES":[{"FEAT_DESC":"1515 Adela Lane Las Vegas NV 89111","LIB_FEAT_ID":20}],"LIB_FEAT_ID":20,"USAGE_TYPE":"HOME"}],"DOB":[{"FEAT_DESC":"12/11/1978","FEAT_DESC_VALUES":[{"FEAT_DESC":"11/12/1978","LIB_FEAT_ID":19},{"FEAT_DESC":"12/11/1978","LIB_FEAT_ID":2}],"LIB_FEAT_ID":2}],"EMAIL":[{"FEAT_DESC":"bsmith@work.com","FEAT_DESC_VALUES":[{"FEAT_DESC":"bsmith@work.com","LIB_FEAT_ID":5}],"LIB_FEAT_ID":5}],"NAME":[{"FEAT_DESC":"Robert Smith","FEAT_DESC_VALUES":[{"FEAT_DESC":"Bob J Smith","LIB_FEAT_ID":32},{"FEAT_DESC":"Bob Smith","LIB_FEAT_ID":18},{"FEAT_DESC":"Robert Smith","LIB_FEAT_ID":1}],"LIB_FEAT_ID":1,"USAGE_TYPE":"PRIMARY"}],"PHONE":[{"FEAT_DESC":"702-919-1300","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4}],"LIB_FEAT_ID":4,"USAGE_TYPE":"HOME"},{"FEAT_DESC":"702-919-1300","FEAT_DESC_VALUES":[{"FEAT_DESC":"702-919-1300","LIB_FEAT_ID":4}],"LIB_FEAT_ID":4,"USAGE_TYPE":"MOBILE"}],"RECORD_TYPE":[{"FEAT_DESC":"PERSON","FEAT_DESC_VALUES":[{"FEAT_DESC":"PERSON","LIB_FEAT_ID":16}],"LIB_FEAT_ID":16}]},"LAST_SEEN_DT":null,"RECORD_SUMMARY":[{"DATA_SOURCE":"CUSTOMERS","FIRST_SEEN_DT":null,"LAST_SEEN_DT":null,"RECORD_COUNT":3}]}},"MATCH_INFO":{"ERRULE_CODE":"SF1","FEATURE_SCORES":{"EMAIL":[{"CANDIDATE_FEAT":"bsmith@work.com","FULL_SCORE":100,"INBOUND_FEAT":"bsmith@work.com"}],"NAME":[{"CANDIDATE_FEAT":"Bob J Smith","GENERATION_MATCH":-1,"GNR_FN":83,"GNR_GN":40,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Smith"},{"CANDIDATE_FEAT":"Robert Smith","GENERATION_MATCH":-1,"GNR_FN":88,"GNR_GN":40,"GNR_ON":-1,"GNR_SN":100,"INBOUND_FEAT":"Smith"}]},"MATCH_KEY":"+PNAME+EMAIL","MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED"}}]}
}

func ExampleSzengine_SetLogLevel() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzengine_WhyEntities() {
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

func ExampleSzengine_WhyRecords() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
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
	// Output: {"WHY_RESULTS":[{"INTERNAL_ID":12,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"CUSTOMERS","RECORD_ID":"1001"}],...
}

func ExampleSzengine_ProcessRedoRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	redoRecord := ""
	flags := int64(0)
	result, err := g2engine.ProcessRedoRecord(ctx, redoRecord, flags)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
	// Output:
}

func ExampleSzengine_ReevaluateEntity() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	entityId := getEntityIdForRecord("CUSTOMERS", "1001")
	flags := int64(0)
	_, err := g2engine.ReevaluateEntity(ctx, entityId, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzengine_ReevaluateRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1001"
	flags := int64(0)
	_, err := g2engine.ReevaluateRecord(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzengine_ReplaceRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
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

func ExampleSzengine_DeleteRecord() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	dataSourceCode := "CUSTOMERS"
	recordId := "1003"
	flags := int64(0)
	_, err := g2engine.DeleteRecord(ctx, dataSourceCode, recordId, flags)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzengine_Initialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	instanceName := "Test module name"
	settings, err := getIniParams()
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

func ExampleSzengine_Reinitialize() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	initConfigId, _ := g2engine.GetActiveConfigId(ctx) // Example initConfigId.
	err := g2engine.Reinitialize(ctx, initConfigId)
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}

func ExampleSzengine_Destroy() {
	// For more information, visit https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	g2engine := getG2Engine(ctx)
	err := g2engine.Destroy(ctx)
	if err != nil {
		fmt.Println(err)
	}
}
