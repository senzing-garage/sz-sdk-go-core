//go:build linux

package szengine_test

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/jsonutil"
	"github.com/senzing-garage/go-helpers/truthset"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/sz-sdk-go/senzing"
)

// ----------------------------------------------------------------------------
// Interface methods - Examples for godoc documentation
// ----------------------------------------------------------------------------

func ExampleSzengine_AddRecord() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	flags := senzing.SzWithoutInfo

	result, err := szEngine.AddRecord(ctx, dataSourceCode, recordID, recordDefinition, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(result)

	// Output:
}

func ExampleSzengine_AddRecord_secondRecord() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	dataSourceCode := "CUSTOMERS"
	recordID := "1002"
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1002", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "DATE_OF_BIRTH": "11/12/1978", "ADDR_TYPE": "HOME", "ADDR_LINE1": "1515 Adela Lane", "ADDR_CITY": "Las Vegas", "ADDR_STATE": "NV", "ADDR_POSTAL_CODE": "89111", "PHONE_TYPE": "MOBILE", "PHONE_NUMBER": "702-919-1300", "DATE": "3/10/17", "STATUS": "Inactive", "AMOUNT": "200"}`
	flags := senzing.SzWithoutInfo

	result, err := szEngine.AddRecord(ctx, dataSourceCode, recordID, recordDefinition, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(result)
	// Output:
}

func ExampleSzengine_AddRecord_withInfo() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Bob", "PRIMARY_NAME_MIDDLE": "J", "DATE_OF_BIRTH": "12/11/1978", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "4/9/16", "STATUS": "Inactive", "AMOUNT": "300"}`
	flags := senzing.SzWithInfo

	result, err := szEngine.AddRecord(ctx, dataSourceCode, recordID, recordDefinition, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "DATA_SOURCE": "CUSTOMERS",
	//     "RECORD_ID": "1003",
	//     "AFFECTED_ENTITIES": [
	//         {
	//             "ENTITY_ID": 100001
	//         }
	//     ]
	// }
}

func ExampleSzengine_CloseExportReport() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	flags := senzing.SzNoFlags

	exportHandle, err := szEngine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		handleError(err)
		return
	}

	err = szEngine.CloseExportReport(ctx, exportHandle)
	if err != nil {
		handleError(err)
		return
	}
	// Output:
}

func ExampleSzengine_CountRedoRecords() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	result, err := szEngine.CountRedoRecords(ctx)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(result)
	// Output: 98
}

func ExampleSzengine_ExportCsvEntityReport() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	csvColumnList := ""
	flags := senzing.SzNoFlags

	exportHandle, err := szEngine.ExportCsvEntityReport(ctx, csvColumnList, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(exportHandle > 0) // Dummy output.
	// Output: true
}

// func ExampleSzengine_ExportCsvEntityReportIterator() {
// 	// For more information, visit
// 	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
// 	ctx := context.TODO()
// 	szAbstractFactory := getSzAbstractFactory(ctx)

// 	szEngine, err := szAbstractFactory.CreateEngine(ctx)
// 	if err != nil {
// 		handleError(err)
// 	}

// 	csvColumnList := ""
// 	flags := senzing.SzNoFlags

// 	for result := range szEngine.ExportCsvEntityReportIterator(ctx, csvColumnList, flags) {
// 		if result.Error != nil {
// 			handleError(err)

// 			break
// 		}

// 		fmt.Println(result.Value)
// 	}
// 	// Output: RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID
// }

func ExampleSzengine_ExportJSONEntityReport() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	flags := senzing.SzNoFlags

	exportHandle, err := szEngine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(exportHandle > 0) // Dummy output.
	// Output: true
}

// func ExampleSzengine_ExportJSONEntityReportIterator() {
// 	// For more information, visit
// 	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
// 	ctx := context.TODO()
// 	szAbstractFactory := getSzAbstractFactory(ctx)

// 	szEngine, err := szAbstractFactory.CreateEngine(ctx)
// 	if err != nil {
// 		handleError(err)
// 	}

// 	flags := senzing.SzNoFlags

// 	for result := range szEngine.ExportJSONEntityReportIterator(ctx, flags) {
// 		if result.Error != nil {
// 			handleError(err)

// 			break
// 		}

// 		fmt.Println(result.Value)
// 	}
// 	// Output:
// }

func ExampleSzengine_FetchNext() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	flags := senzing.SzNoFlags

	exportHandle, err := szEngine.ExportJSONEntityReport(ctx, flags)
	if err != nil {
		handleError(err)
		return
	}

	defer func() {
		handleError(szEngine.CloseExportReport(ctx, exportHandle))
	}()

	jsonEntityReport := ""

	for {
		jsonEntityReportFragment, err := szEngine.FetchNext(ctx, exportHandle)
		if err != nil {
			handleError(err)
			return
		}

		if len(jsonEntityReportFragment) == 0 {
			break
		}

		jsonEntityReport += jsonEntityReportFragment
	}
	// Output:
}

func ExampleSzengine_FindInterestingEntitiesByEntityID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	entityID := getEntityIDForRecord(ctx, szEngine, "CUSTOMERS", "1001")

	flags := senzing.SzNoFlags

	result, err := szEngine.FindInterestingEntitiesByEntityID(ctx, entityID, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "INTERESTING_ENTITIES": {
	//         "ENTITIES": []
	//     }
	// }
}

func ExampleSzengine_FindInterestingEntitiesByRecordID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := senzing.SzNoFlags

	result, err := szEngine.FindInterestingEntitiesByRecordID(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "INTERESTING_ENTITIES": {
	//         "ENTITIES": []
	//     }
	// }
}

// func ExampleSzengine_FindNetworkByEntityID() {
// 	// For more information, visit
// 	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
// 	ctx := context.TODO()
// 	szAbstractFactory := createSzAbstractFactory(ctx)
// 	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

// 	szEngine, err := szAbstractFactory.CreateEngine(ctx)
// 	if err != nil {
// 		handleError(err)
// 	}
// 	defer func() { handleError(szEngine.Destroy(ctx)) }()

// 	entityID1 := getEntityIDStringForRecord(ctx, "CUSTOMERS", "1001")
// 	entityID2 := getEntityIDStringForRecord(ctx, "CUSTOMERS", "1002")
// 	entityList := `{"ENTITIES": [{"ENTITY_ID": ` + entityID1 + `}, {"ENTITY_ID": ` + entityID2 + `}]}`
// 	maxDegrees := int64(2)
// 	buildOutDegrees := int64(1)
// 	maxEntities := int64(10)
// 	flags := senzing.SzNoFlags

// 	result, err := szEngine.FindNetworkByEntityID(ctx, entityList, maxDegrees, buildOutDegrees, maxEntities, flags)
// 	if err != nil {
// 		handleError(err)
// 	}

// 	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
// 	// Output:
// 	// {
// 	//     "ENTITY_PATHS": [],
// 	//     "ENTITIES": [
// 	//         {
// 	//             "RESOLVED_ENTITY": {
// 	//                 "ENTITY_ID": 100001
// 	//             }
// 	//         }
// 	//     ]
// 	// }
// }

func ExampleSzengine_FindNetworkByEntityID_output() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	// The following code pretty-prints example output JSON.
	exampleOutput := `{"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"SEAMAN","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-11-29 22:25:18.997","LAST_SEEN_DT":"2022-11-29 22:25:19.005"}],"LAST_SEEN_DT":"2022-11-29 22:25:19.005"},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-11-29 22:25:19.009","LAST_SEEN_DT":"2022-11-29 22:25:19.009"}],"LAST_SEEN_DT":"2022-11-29 22:25:19.009"},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]}]}`
	fmt.Println(jsonutil.PrettyPrint(exampleOutput, jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 1,
	//             "END_ENTITY_ID": 2,
	//             "ENTITIES": [
	//                 1,
	//                 2
	//             ]
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 1,
	//                 "ENTITY_NAME": "SEAMAN",
	//                 "RECORD_SUMMARY": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_COUNT": 2,
	//                         "FIRST_SEEN_DT": "2022-11-29 22:25:18.997",
	//                         "LAST_SEEN_DT": "2022-11-29 22:25:19.005"
	//                     }
	//                 ],
	//                 "LAST_SEEN_DT": "2022-11-29 22:25:19.005"
	//             },
	//             "RELATED_ENTITIES": [
	//                 {
	//                     "ENTITY_ID": 2,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-DOB-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 }
	//             ]
	//         },
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 2,
	//                 "ENTITY_NAME": "Smith",
	//                 "RECORD_SUMMARY": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_COUNT": 1,
	//                         "FIRST_SEEN_DT": "2022-11-29 22:25:19.009",
	//                         "LAST_SEEN_DT": "2022-11-29 22:25:19.009"
	//                     }
	//                 ],
	//                 "LAST_SEEN_DT": "2022-11-29 22:25:19.009"
	//             },
	//             "RELATED_ENTITIES": [
	//                 {
	//                     "ENTITY_ID": 1,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-DOB-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 }
	//             ]
	//         }
	//     ]
	// }
}

func ExampleSzengine_FindNetworkByRecordID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	recordList := `
	{
		"RECORDS": [
			{
				"DATA_SOURCE": "CUSTOMERS",
				"RECORD_ID": "1001"
			},
			{
				"DATA_SOURCE": "CUSTOMERS",
				"RECORD_ID": "1002"
			}
		]
	}`
	maxDegrees := int64(1)
	buildOutDegrees := int64(2)
	maxEntities := int64(10)
	flags := senzing.SzNoFlags

	result, err := szEngine.FindNetworkByRecordID(ctx, recordList, maxDegrees, buildOutDegrees, maxEntities, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100001
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzengine_FindNetworkByRecordID_output() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	// The following code pretty-prints example output JSON.
	exampleOutput := `{"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 14:40:34.285","LAST_SEEN_DT":"2022-12-06 14:40:34.420"}],"LAST_SEEN_DT":"2022-12-06 14:40:34.420"},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 14:40:34.359","LAST_SEEN_DT":"2022-12-06 14:40:34.359"}],"LAST_SEEN_DT":"2022-12-06 14:40:34.359"},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]},{"RESOLVED_ENTITY":{"ENTITY_ID":3,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 14:40:34.424","LAST_SEEN_DT":"2022-12-06 14:40:34.424"}],"LAST_SEEN_DT":"2022-12-06 14:40:34.424"},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]}]}`
	fmt.Println(jsonutil.PrettyPrint(exampleOutput, jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 1,
	//             "END_ENTITY_ID": 2,
	//             "ENTITIES": [
	//                 1,
	//                 2
	//             ]
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 1,
	//                 "ENTITY_NAME": "JOHNSON",
	//                 "RECORD_SUMMARY": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_COUNT": 2,
	//                         "FIRST_SEEN_DT": "2022-12-06 14:40:34.285",
	//                         "LAST_SEEN_DT": "2022-12-06 14:40:34.420"
	//                     }
	//                 ],
	//                 "LAST_SEEN_DT": "2022-12-06 14:40:34.420"
	//             },
	//             "RELATED_ENTITIES": [
	//                 {
	//                     "ENTITY_ID": 2,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 },
	//                 {
	//                     "ENTITY_ID": 3,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-DOB-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 }
	//             ]
	//         },
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 2,
	//                 "ENTITY_NAME": "OCEANGUY",
	//                 "RECORD_SUMMARY": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_COUNT": 1,
	//                         "FIRST_SEEN_DT": "2022-12-06 14:40:34.359",
	//                         "LAST_SEEN_DT": "2022-12-06 14:40:34.359"
	//                     }
	//                 ],
	//                 "LAST_SEEN_DT": "2022-12-06 14:40:34.359"
	//             },
	//             "RELATED_ENTITIES": [
	//                 {
	//                     "ENTITY_ID": 1,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 },
	//                 {
	//                     "ENTITY_ID": 3,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+ADDRESS+PHONE+ACCT_NUM-DOB-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 }
	//             ]
	//         },
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 3,
	//                 "ENTITY_NAME": "Smith",
	//                 "RECORD_SUMMARY": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_COUNT": 1,
	//                         "FIRST_SEEN_DT": "2022-12-06 14:40:34.424",
	//                         "LAST_SEEN_DT": "2022-12-06 14:40:34.424"
	//                     }
	//                 ],
	//                 "LAST_SEEN_DT": "2022-12-06 14:40:34.424"
	//             },
	//             "RELATED_ENTITIES": [
	//                 {
	//                     "ENTITY_ID": 1,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-DOB-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 },
	//                 {
	//                     "ENTITY_ID": 2,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+ADDRESS+PHONE+ACCT_NUM-DOB-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 }
	//             ]
	//         }
	//     ]
	// }
}

func ExampleSzengine_FindPathByEntityID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	startEntityID := getEntityIDForRecord(ctx, szEngine, "CUSTOMERS", "1001")
	endEntityID := getEntityIDForRecord(ctx, szEngine, "CUSTOMERS", "1002")

	maxDegrees := int64(1)
	avoidEntityIDs := ""
	requiredDataSources := ""
	flags := senzing.SzNoFlags

	result, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		maxDegrees,
		avoidEntityIDs,
		requiredDataSources,
		flags,
	)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 100001,
	//             "END_ENTITY_ID": 100001,
	//             "ENTITIES": [
	//                 100001
	//             ]
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100001
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzengine_FindPathByEntityID_avoiding() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	startEntityID := getEntityIDForRecord(ctx, szEngine, "CUSTOMERS", "1001")
	endEntityID := getEntityIDForRecord(ctx, szEngine, "CUSTOMERS", "1002")

	maxDegrees := int64(1)
	avoidEntityIDs := `{"ENTITIES": [{"ENTITY_ID": 9999}]}`
	requiredDataSources := ""
	flags := senzing.SzNoFlags

	result, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		maxDegrees,
		avoidEntityIDs,
		requiredDataSources,
		flags,
	)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 100001,
	//             "END_ENTITY_ID": 100001,
	//             "ENTITIES": [
	//                 100001
	//             ]
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100001
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzEngine_FindPathByEntityID_avoidingAndIncluding() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	startEntityID := getEntityIDForRecord(ctx, szEngine, "CUSTOMERS", "1001")
	endEntityID := getEntityIDForRecord(ctx, szEngine, "CUSTOMERS", "1002")

	maxDegree := int64(1)
	avoidEntityIDs := `{"ENTITIES": [{"ENTITY_ID": 9999}]}`
	requiredDataSources := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	flags := senzing.SzNoFlags

	result, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		maxDegree,
		avoidEntityIDs,
		requiredDataSources,
		flags,
	)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 100001,
	//             "END_ENTITY_ID": 100001,
	//             "ENTITIES": []
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100001
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzengine_FindPathByEntityID_including() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	startEntityID := getEntityIDForRecord(ctx, szEngine, "CUSTOMERS", "1001")
	endEntityID := getEntityIDForRecord(ctx, szEngine, "CUSTOMERS", "1002")

	maxDegree := int64(1)
	avoidEntityIDs := ""
	requiredDataSources := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	flags := senzing.SzNoFlags

	result, err := szEngine.FindPathByEntityID(
		ctx,
		startEntityID,
		endEntityID,
		maxDegree,
		avoidEntityIDs,
		requiredDataSources,
		flags,
	)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 100001,
	//             "END_ENTITY_ID": 100001,
	//             "ENTITIES": []
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100001
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzengine_FindPathByEntityID_output() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	// The following code pretty-prints example output JSON.
	exampleOutput := `{"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 14:43:49.024","LAST_SEEN_DT":"2022-12-06 14:43:49.164"}],"LAST_SEEN_DT":"2022-12-06 14:43:49.164"},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 14:43:49.104","LAST_SEEN_DT":"2022-12-06 14:43:49.104"}],"LAST_SEEN_DT":"2022-12-06 14:43:49.104"},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]}]}`
	fmt.Println(jsonutil.PrettyPrint(exampleOutput, jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 1,
	//             "END_ENTITY_ID": 2,
	//             "ENTITIES": [
	//                 1,
	//                 2
	//             ]
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 1,
	//                 "ENTITY_NAME": "JOHNSON",
	//                 "RECORD_SUMMARY": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_COUNT": 2,
	//                         "FIRST_SEEN_DT": "2022-12-06 14:43:49.024",
	//                         "LAST_SEEN_DT": "2022-12-06 14:43:49.164"
	//                     }
	//                 ],
	//                 "LAST_SEEN_DT": "2022-12-06 14:43:49.164"
	//             },
	//             "RELATED_ENTITIES": [
	//                 {
	//                     "ENTITY_ID": 2,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 },
	//                 {
	//                     "ENTITY_ID": 3,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-DOB-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 }
	//             ]
	//         },
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 2,
	//                 "ENTITY_NAME": "OCEANGUY",
	//                 "RECORD_SUMMARY": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_COUNT": 1,
	//                         "FIRST_SEEN_DT": "2022-12-06 14:43:49.104",
	//                         "LAST_SEEN_DT": "2022-12-06 14:43:49.104"
	//                     }
	//                 ],
	//                 "LAST_SEEN_DT": "2022-12-06 14:43:49.104"
	//             },
	//             "RELATED_ENTITIES": [
	//                 {
	//                     "ENTITY_ID": 1,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 },
	//                 {
	//                     "ENTITY_ID": 3,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+ADDRESS+PHONE+ACCT_NUM-DOB-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 }
	//             ]
	//         }
	//     ]
	// }
}

func ExampleSzengine_FindPathByRecordID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	startDataSourceCode := "CUSTOMERS"
	startRecordID := "1001"
	endDataSourceCode := "CUSTOMERS"
	endRecordID := "1002"
	maxDegrees := int64(1)
	avoidRecordKeys := ""
	requiredDataSources := ""
	flags := senzing.SzNoFlags

	result, err := szEngine.FindPathByRecordID(
		ctx,
		startDataSourceCode,
		startRecordID,
		endDataSourceCode,
		endRecordID,
		maxDegrees,
		avoidRecordKeys,
		requiredDataSources,
		flags,
	)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 100001,
	//             "END_ENTITY_ID": 100001,
	//             "ENTITIES": [
	//                 100001
	//             ]
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100001
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzengine_FindPathByRecordID_avoiding() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	startDataSourceCode := "CUSTOMERS"
	startRecordID := "1001"
	endDataSourceCode := "CUSTOMERS"
	endRecordID := "1002"
	maxDegree := int64(1)
	avoidRecordKeys := `{"RECORDS": [{ "DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1003"}]}`
	requiredDataSources := ""
	flags := senzing.SzNoFlags

	result, err := szEngine.FindPathByRecordID(
		ctx,
		startDataSourceCode,
		startRecordID,
		endDataSourceCode,
		endRecordID,
		maxDegree,
		avoidRecordKeys,
		requiredDataSources,
		flags,
	)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 100001,
	//             "END_ENTITY_ID": 100001,
	//             "ENTITIES": [
	//                 100001
	//             ]
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100001
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzEngine_FindPathByRecordID_avoidingAndIncluding() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	startDataSourceCode := "CUSTOMERS"
	startRecordID := "1001"
	endDataSourceCode := "CUSTOMERS"
	endRecordID := "1002"
	maxDegrees := int64(1)
	avoidRecordKeys := `{"ENTITIES": [{"ENTITY_ID": 9999}]}`
	requiredDataSources := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	flags := senzing.SzNoFlags

	result, err := szEngine.FindPathByRecordID(
		ctx,
		startDataSourceCode,
		startRecordID,
		endDataSourceCode,
		endRecordID,
		maxDegrees,
		avoidRecordKeys,
		requiredDataSources,
		flags,
	)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 100001,
	//             "END_ENTITY_ID": 100001,
	//             "ENTITIES": []
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100001
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzengine_FindPathByRecordID_including() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	startDataSourceCode := "CUSTOMERS"
	startRecordID := "1001"
	endDataSourceCode := "CUSTOMERS"
	endRecordID := "1002"
	maxDegrees := int64(1)
	avoidRecordKeys := ""
	requiredDataSources := `{"DATA_SOURCES": ["CUSTOMERS"]}`
	flags := senzing.SzNoFlags

	result, err := szEngine.FindPathByRecordID(
		ctx,
		startDataSourceCode,
		startRecordID,
		endDataSourceCode,
		endRecordID,
		maxDegrees,
		avoidRecordKeys,
		requiredDataSources,
		flags,
	)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 100001,
	//             "END_ENTITY_ID": 100001,
	//             "ENTITIES": []
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100001
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzengine_FindPathByRecordID_output() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	// The following code pretty-prints example output JSON.
	exampleOutput := `{"ENTITY_PATHS":[{"START_ENTITY_ID":1,"END_ENTITY_ID":2,"ENTITIES":[1,2]}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 14:48:19.522","LAST_SEEN_DT":"2022-12-06 14:48:19.667"}],"LAST_SEEN_DT":"2022-12-06 14:48:19.667"},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 14:48:19.593","LAST_SEEN_DT":"2022-12-06 14:48:19.593"}],"LAST_SEEN_DT":"2022-12-06 14:48:19.593"},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0}]}]}`
	fmt.Println(jsonutil.PrettyPrint(exampleOutput, jsonIndentation))
	// Output:
	// {
	//     "ENTITY_PATHS": [
	//         {
	//             "START_ENTITY_ID": 1,
	//             "END_ENTITY_ID": 2,
	//             "ENTITIES": [
	//                 1,
	//                 2
	//             ]
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 1,
	//                 "ENTITY_NAME": "JOHNSON",
	//                 "RECORD_SUMMARY": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_COUNT": 2,
	//                         "FIRST_SEEN_DT": "2022-12-06 14:48:19.522",
	//                         "LAST_SEEN_DT": "2022-12-06 14:48:19.667"
	//                     }
	//                 ],
	//                 "LAST_SEEN_DT": "2022-12-06 14:48:19.667"
	//             },
	//             "RELATED_ENTITIES": [
	//                 {
	//                     "ENTITY_ID": 2,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 },
	//                 {
	//                     "ENTITY_ID": 3,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-DOB-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 }
	//             ]
	//         },
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 2,
	//                 "ENTITY_NAME": "OCEANGUY",
	//                 "RECORD_SUMMARY": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_COUNT": 1,
	//                         "FIRST_SEEN_DT": "2022-12-06 14:48:19.593",
	//                         "LAST_SEEN_DT": "2022-12-06 14:48:19.593"
	//                     }
	//                 ],
	//                 "LAST_SEEN_DT": "2022-12-06 14:48:19.593"
	//             },
	//             "RELATED_ENTITIES": [
	//                 {
	//                     "ENTITY_ID": 1,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 },
	//                 {
	//                     "ENTITY_ID": 3,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+ADDRESS+PHONE+ACCT_NUM-DOB-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0
	//                 }
	//             ]
	//         }
	//     ]
	// }
}

func ExampleSzengine_GetActiveConfigID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	result, err := szEngine.GetActiveConfigID(ctx)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(result > 0) // Dummy output.
	// Output: true
}

func ExampleSzengine_GetEntityByEntityID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	entityID := getEntityIDForRecord(ctx, szEngine, "CUSTOMERS", "1001")

	flags := senzing.SzNoFlags

	result, err := szEngine.GetEntityByEntityID(ctx, entityID, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "RESOLVED_ENTITY": {
	//         "ENTITY_ID": 100001
	//     }
	// }
}

func ExampleSzengine_GetEntityByEntityID_output() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	// The following code pretty-prints example output JSON.
	exampleOutput := `{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4}]}],"DOB":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7}]}],"NAME":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5}]}],"SSN":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"FEAT_DESC_VALUES":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 15:09:48.577","LAST_SEEN_DT":"2022-12-06 15:09:48.705"}],"LAST_SEEN_DT":"2022-12-06 15:09:48.705","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 15:09:48.577"},{"DATA_SOURCE":"TEST","RECORD_ID":"FCCE9793DAAD23159DBCCEB97FF2745B92CE7919","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+EXACTLY_SAME","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 15:09:48.705"}]},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:09:48.647","LAST_SEEN_DT":"2022-12-06 15:09:48.647"}],"LAST_SEEN_DT":"2022-12-06 15:09:48.647"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:09:48.709","LAST_SEEN_DT":"2022-12-06 15:09:48.709"}],"LAST_SEEN_DT":"2022-12-06 15:09:48.709"}]}`
	fmt.Println(jsonutil.PrettyPrint(exampleOutput, jsonIndentation))
	// Output:
	// {
	//     "RESOLVED_ENTITY": {
	//         "ENTITY_ID": 1,
	//         "ENTITY_NAME": "JOHNSON",
	//         "FEATURES": {
	//             "ACCT_NUM": [
	//                 {
	//                     "FEAT_DESC": "5534202208773608",
	//                     "LIB_FEAT_ID": 8,
	//                     "USAGE_TYPE": "CC",
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "5534202208773608",
	//                             "LIB_FEAT_ID": 8
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "ADDRESS": [
	//                 {
	//                     "FEAT_DESC": "772 Armstrong RD Delhi LA 71232",
	//                     "LIB_FEAT_ID": 4,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "772 Armstrong RD Delhi LA 71232",
	//                             "LIB_FEAT_ID": 4
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "DOB": [
	//                 {
	//                     "FEAT_DESC": "4/8/1983",
	//                     "LIB_FEAT_ID": 2,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "4/8/1983",
	//                             "LIB_FEAT_ID": 2
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "GENDER": [
	//                 {
	//                     "FEAT_DESC": "F",
	//                     "LIB_FEAT_ID": 3,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "F",
	//                             "LIB_FEAT_ID": 3
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "LOGIN_ID": [
	//                 {
	//                     "FEAT_DESC": "flavorh",
	//                     "LIB_FEAT_ID": 7,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "flavorh",
	//                             "LIB_FEAT_ID": 7
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "NAME": [
	//                 {
	//                     "FEAT_DESC": "JOHNSON",
	//                     "LIB_FEAT_ID": 1,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "JOHNSON",
	//                             "LIB_FEAT_ID": 1
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "PHONE": [
	//                 {
	//                     "FEAT_DESC": "225-671-0796",
	//                     "LIB_FEAT_ID": 5,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "225-671-0796",
	//                             "LIB_FEAT_ID": 5
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "SSN": [
	//                 {
	//                     "FEAT_DESC": "053-39-3251",
	//                     "LIB_FEAT_ID": 6,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "053-39-3251",
	//                             "LIB_FEAT_ID": 6
	//                         }
	//                     ]
	//                 }
	//             ]
	//         },
	//         "RECORD_SUMMARY": [
	//             {
	//                 "DATA_SOURCE": "TEST",
	//                 "RECORD_COUNT": 2,
	//                 "FIRST_SEEN_DT": "2022-12-06 15:09:48.577",
	//                 "LAST_SEEN_DT": "2022-12-06 15:09:48.705"
	//             }
	//         ],
	//         "LAST_SEEN_DT": "2022-12-06 15:09:48.705",
	//         "RECORDS": [
	//             {
	//                 "DATA_SOURCE": "TEST",
	//                 "RECORD_ID": "111",
	//                 "ENTITY_TYPE": "TEST",
	//                 "INTERNAL_ID": 1,
	//                 "ENTITY_KEY": "C6063D4396612FBA7324DB0739273BA1FE815C43",
	//                 "ENTITY_DESC": "JOHNSON",
	//                 "MATCH_KEY": "",
	//                 "MATCH_LEVEL": 0,
	//                 "MATCH_LEVEL_CODE": "",
	//                 "ERRULE_CODE": "",
	//                 "LAST_SEEN_DT": "2022-12-06 15:09:48.577"
	//             },
	//             {
	//                 "DATA_SOURCE": "TEST",
	//                 "RECORD_ID": "FCCE9793DAAD23159DBCCEB97FF2745B92CE7919",
	//                 "ENTITY_TYPE": "TEST",
	//                 "INTERNAL_ID": 1,
	//                 "ENTITY_KEY": "C6063D4396612FBA7324DB0739273BA1FE815C43",
	//                 "ENTITY_DESC": "JOHNSON",
	//                 "MATCH_KEY": "+EXACTLY_SAME",
	//                 "MATCH_LEVEL": 0,
	//                 "MATCH_LEVEL_CODE": "",
	//                 "ERRULE_CODE": "",
	//                 "LAST_SEEN_DT": "2022-12-06 15:09:48.705"
	//             }
	//         ]
	//     },
	//     "RELATED_ENTITIES": [
	//         {
	//             "ENTITY_ID": 2,
	//             "MATCH_LEVEL": 3,
	//             "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//             "MATCH_KEY": "+PHONE+ACCT_NUM-SSN",
	//             "ERRULE_CODE": "SF1",
	//             "IS_DISCLOSED": 0,
	//             "IS_AMBIGUOUS": 0,
	//             "ENTITY_NAME": "OCEANGUY",
	//             "RECORD_SUMMARY": [
	//                 {
	//                     "DATA_SOURCE": "TEST",
	//                     "RECORD_COUNT": 1,
	//                     "FIRST_SEEN_DT": "2022-12-06 15:09:48.647",
	//                     "LAST_SEEN_DT": "2022-12-06 15:09:48.647"
	//                 }
	//             ],
	//             "LAST_SEEN_DT": "2022-12-06 15:09:48.647"
	//         },
	//         {
	//             "ENTITY_ID": 3,
	//             "MATCH_LEVEL": 3,
	//             "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//             "MATCH_KEY": "+PHONE+ACCT_NUM-DOB-SSN",
	//             "ERRULE_CODE": "SF1",
	//             "IS_DISCLOSED": 0,
	//             "IS_AMBIGUOUS": 0,
	//             "ENTITY_NAME": "Smith",
	//             "RECORD_SUMMARY": [
	//                 {
	//                     "DATA_SOURCE": "TEST",
	//                     "RECORD_COUNT": 1,
	//                     "FIRST_SEEN_DT": "2022-12-06 15:09:48.709",
	//                     "LAST_SEEN_DT": "2022-12-06 15:09:48.709"
	//                 }
	//             ],
	//             "LAST_SEEN_DT": "2022-12-06 15:09:48.709"
	//         }
	//     ]
	// }
}

func ExampleSzengine_GetEntityByRecordID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := senzing.SzNoFlags

	result, err := szEngine.GetEntityByRecordID(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "RESOLVED_ENTITY": {
	//         "ENTITY_ID": 100001
	//     }
	// }
}

func ExampleSzengine_GetEntityByRecordID_output() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	// The following code pretty-prints example output JSON.
	exampleOutput := `{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4}]}],"DOB":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7}]}],"NAME":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5}]}],"SSN":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"FEAT_DESC_VALUES":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 15:12:25.464","LAST_SEEN_DT":"2022-12-06 15:12:25.597"}],"LAST_SEEN_DT":"2022-12-06 15:12:25.597","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 15:12:25.464"},{"DATA_SOURCE":"TEST","RECORD_ID":"FCCE9793DAAD23159DBCCEB97FF2745B92CE7919","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+EXACTLY_SAME","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 15:12:25.597"}]},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:12:25.536","LAST_SEEN_DT":"2022-12-06 15:12:25.536"}],"LAST_SEEN_DT":"2022-12-06 15:12:25.536"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:12:25.603","LAST_SEEN_DT":"2022-12-06 15:12:25.603"}],"LAST_SEEN_DT":"2022-12-06 15:12:25.603"}]}`
	fmt.Println(jsonutil.PrettyPrint(exampleOutput, jsonIndentation))
	// Output:
	// {
	//     "RESOLVED_ENTITY": {
	//         "ENTITY_ID": 1,
	//         "ENTITY_NAME": "JOHNSON",
	//         "FEATURES": {
	//             "ACCT_NUM": [
	//                 {
	//                     "FEAT_DESC": "5534202208773608",
	//                     "LIB_FEAT_ID": 8,
	//                     "USAGE_TYPE": "CC",
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "5534202208773608",
	//                             "LIB_FEAT_ID": 8
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "ADDRESS": [
	//                 {
	//                     "FEAT_DESC": "772 Armstrong RD Delhi LA 71232",
	//                     "LIB_FEAT_ID": 4,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "772 Armstrong RD Delhi LA 71232",
	//                             "LIB_FEAT_ID": 4
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "DOB": [
	//                 {
	//                     "FEAT_DESC": "4/8/1983",
	//                     "LIB_FEAT_ID": 2,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "4/8/1983",
	//                             "LIB_FEAT_ID": 2
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "GENDER": [
	//                 {
	//                     "FEAT_DESC": "F",
	//                     "LIB_FEAT_ID": 3,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "F",
	//                             "LIB_FEAT_ID": 3
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "LOGIN_ID": [
	//                 {
	//                     "FEAT_DESC": "flavorh",
	//                     "LIB_FEAT_ID": 7,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "flavorh",
	//                             "LIB_FEAT_ID": 7
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "NAME": [
	//                 {
	//                     "FEAT_DESC": "JOHNSON",
	//                     "LIB_FEAT_ID": 1,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "JOHNSON",
	//                             "LIB_FEAT_ID": 1
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "PHONE": [
	//                 {
	//                     "FEAT_DESC": "225-671-0796",
	//                     "LIB_FEAT_ID": 5,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "225-671-0796",
	//                             "LIB_FEAT_ID": 5
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "SSN": [
	//                 {
	//                     "FEAT_DESC": "053-39-3251",
	//                     "LIB_FEAT_ID": 6,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "053-39-3251",
	//                             "LIB_FEAT_ID": 6
	//                         }
	//                     ]
	//                 }
	//             ]
	//         },
	//         "RECORD_SUMMARY": [
	//             {
	//                 "DATA_SOURCE": "TEST",
	//                 "RECORD_COUNT": 2,
	//                 "FIRST_SEEN_DT": "2022-12-06 15:12:25.464",
	//                 "LAST_SEEN_DT": "2022-12-06 15:12:25.597"
	//             }
	//         ],
	//         "LAST_SEEN_DT": "2022-12-06 15:12:25.597",
	//         "RECORDS": [
	//             {
	//                 "DATA_SOURCE": "TEST",
	//                 "RECORD_ID": "111",
	//                 "ENTITY_TYPE": "TEST",
	//                 "INTERNAL_ID": 1,
	//                 "ENTITY_KEY": "C6063D4396612FBA7324DB0739273BA1FE815C43",
	//                 "ENTITY_DESC": "JOHNSON",
	//                 "MATCH_KEY": "",
	//                 "MATCH_LEVEL": 0,
	//                 "MATCH_LEVEL_CODE": "",
	//                 "ERRULE_CODE": "",
	//                 "LAST_SEEN_DT": "2022-12-06 15:12:25.464"
	//             },
	//             {
	//                 "DATA_SOURCE": "TEST",
	//                 "RECORD_ID": "FCCE9793DAAD23159DBCCEB97FF2745B92CE7919",
	//                 "ENTITY_TYPE": "TEST",
	//                 "INTERNAL_ID": 1,
	//                 "ENTITY_KEY": "C6063D4396612FBA7324DB0739273BA1FE815C43",
	//                 "ENTITY_DESC": "JOHNSON",
	//                 "MATCH_KEY": "+EXACTLY_SAME",
	//                 "MATCH_LEVEL": 0,
	//                 "MATCH_LEVEL_CODE": "",
	//                 "ERRULE_CODE": "",
	//                 "LAST_SEEN_DT": "2022-12-06 15:12:25.597"
	//             }
	//         ]
	//     },
	//     "RELATED_ENTITIES": [
	//         {
	//             "ENTITY_ID": 2,
	//             "MATCH_LEVEL": 3,
	//             "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//             "MATCH_KEY": "+PHONE+ACCT_NUM-SSN",
	//             "ERRULE_CODE": "SF1",
	//             "IS_DISCLOSED": 0,
	//             "IS_AMBIGUOUS": 0,
	//             "ENTITY_NAME": "OCEANGUY",
	//             "RECORD_SUMMARY": [
	//                 {
	//                     "DATA_SOURCE": "TEST",
	//                     "RECORD_COUNT": 1,
	//                     "FIRST_SEEN_DT": "2022-12-06 15:12:25.536",
	//                     "LAST_SEEN_DT": "2022-12-06 15:12:25.536"
	//                 }
	//             ],
	//             "LAST_SEEN_DT": "2022-12-06 15:12:25.536"
	//         },
	//         {
	//             "ENTITY_ID": 3,
	//             "MATCH_LEVEL": 3,
	//             "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//             "MATCH_KEY": "+PHONE+ACCT_NUM-DOB-SSN",
	//             "ERRULE_CODE": "SF1",
	//             "IS_DISCLOSED": 0,
	//             "IS_AMBIGUOUS": 0,
	//             "ENTITY_NAME": "Smith",
	//             "RECORD_SUMMARY": [
	//                 {
	//                     "DATA_SOURCE": "TEST",
	//                     "RECORD_COUNT": 1,
	//                     "FIRST_SEEN_DT": "2022-12-06 15:12:25.603",
	//                     "LAST_SEEN_DT": "2022-12-06 15:12:25.603"
	//                 }
	//             ],
	//             "LAST_SEEN_DT": "2022-12-06 15:12:25.603"
	//         }
	//     ]
	// }
}

func ExampleSzengine_GetRecord() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := senzing.SzNoFlags

	result, err := szEngine.GetRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(jsonutil.Flatten(jsonutil.Normalize(result)), jsonIndentation))
	// Output:
	// {
	//     "DATA_SOURCE": "CUSTOMERS",
	//     "RECORD_ID": "1001"
	// }
}

func ExampleSzengine_GetRedoRecord() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	result, err := szEngine.GetRedoRecord(ctx)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "UMF_PROC": {
	//         "NAME": "REPAIR_ENTITY",
	//         "PARAMS": [
	//             {
	//                 "PARAM": {
	//                     "NAME": "ENTITY_ID",
	//                     "VALUE": 288
	//                 }
	//             },
	//             {
	//                 "PARAM": {
	//                     "NAME": "ENTITY_CORRUPTION_TRANSIENT",
	//                     "VALUE": 1
	//                 }
	//             },
	//             {
	//                 "PARAM": {
	//                     "NAME": "REEVAL_ITERATION",
	//                     "VALUE": 1
	//                 }
	//             },
	//             {
	//                 "PARAM": {
	//                     "NAME": "REASON",
	//                     "VALUE": "deferred delete: Resolved Entity 288"
	//                 }
	//             }
	//         ]
	//     }
	// }
}

func ExampleSzengine_GetStats() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	result, err := szEngine.GetStats(ctx)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.Truncate(result, 2))
	// Output: {"workload":{...
}

func ExampleSzengine_GetStats_output() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	// The following code pretty-prints example output JSON.
	exampleOutput := `{"workload":{"apiVersion":"4.0.0.25245","datetimestamp":"2025-09-02T16:24:39Z","license":{"status":"ok","type":"non-production","dsrLimit":"ok"},"loadedRecords":3,"processing":{"addedRecords":0,"batchAddedRecords":0,"reevaluations":0,"repairedEntities":0,"deletedRecords":0,"details":{"optimizedOut":0,"optimizedOutSkipped":0,"newObsEnt":0,"obsEntHashSame":0,"obsEntHashDiff":0,"filteredObsFeat":0,"partiallyResolved":0,"changeDeletes":0,"retries":0,"candidates":0,"duration":0},"ambiguous":{"actualTest":0,"cachedTest":0}},"caches":{"libFeatCacheHit":0,"libFeatCacheMiss":0,"resFeatStatCacheHit":0,"resFeatStatCacheMiss":0,"libFeatInsert":0,"resFeatStatInsert":0,"resFeatStatUpdateAttempt":0,"resFeatStatUpdateFail":0},"lockWaits":{"refreshLocks":{"maxMS":0,"totalMS":0,"count":0}},"unresolve":{"triggers":{"normalResolve":0,"update":0,"relLink":0,"extensiveResolve":0,"ambiguousNoResolve":0,"ambiguousMultiResolve":0},"unresolveTest":0,"abortedUnresolve":0},"reresolve":{"triggers":{"skipped":0,"abortRetry":0,"unresolveMovement":0,"multipleResolvableCandidates":0,"resolveNewFeatures":0,"newFeatureFTypes":[]},"suppressedCandidateBuildersForReresolve":[],"suppressedScoredFeatureTypeForReresolve":[]},"expressedFeatures":{"calls":[],"created":[]},"scoring":{"scoredPairs":[],"cacheHit":[],"cacheMiss":[],"suppressedScoredFeatureType":[],"suppressedDisclosedRelationshipDomainCount":0},"redoTriggers":[],"contention":{"valuelatch":[],"feature":[],"resent":[]},"genericDetect":[],"candidates":{"candidateBuilders":[],"suppressedCandidateBuilders":[]},"repairDiagnosis":{"types":0},"threadState":{"active":0,"idle":5,"governorContention":0,"sqlExecuting":0,"loader":0,"resolver":0,"scoring":0,"dataLatchContention":0,"obsEntContention":0,"resEntContention":0},"systemResources":{"initResources":[{"physicalCores":16},{"logicalCores":16},{"totalMemory":"62.6GB"},{"availableMemory":"53.9GB"}],"currResources":[{"availableMemory":"49.4GB"},{"processMemory":"5.4GB"},{"activeThreads":0},{"workerThreads":5},{"systemLoad":[{"cpuUser":6.1},{"cpuSystem":0},{"cpuIdle":93.9},{"cpuWait":0}]}]}}}`
	fmt.Println(jsonutil.PrettyPrint(exampleOutput, jsonIndentation))
	// Output:
	// {
	//     "workload": {
	//         "apiVersion": "4.0.0.25245",
	//         "datetimestamp": "2025-09-02T16:24:39Z",
	//         "license": {
	//             "status": "ok",
	//             "type": "non-production",
	//             "dsrLimit": "ok"
	//         },
	//         "loadedRecords": 3,
	//         "processing": {
	//             "addedRecords": 0,
	//             "batchAddedRecords": 0,
	//             "reevaluations": 0,
	//             "repairedEntities": 0,
	//             "deletedRecords": 0,
	//             "details": {
	//                 "optimizedOut": 0,
	//                 "optimizedOutSkipped": 0,
	//                 "newObsEnt": 0,
	//                 "obsEntHashSame": 0,
	//                 "obsEntHashDiff": 0,
	//                 "filteredObsFeat": 0,
	//                 "partiallyResolved": 0,
	//                 "changeDeletes": 0,
	//                 "retries": 0,
	//                 "candidates": 0,
	//                 "duration": 0
	//             },
	//             "ambiguous": {
	//                 "actualTest": 0,
	//                 "cachedTest": 0
	//             }
	//         },
	//         "caches": {
	//             "libFeatCacheHit": 0,
	//             "libFeatCacheMiss": 0,
	//             "resFeatStatCacheHit": 0,
	//             "resFeatStatCacheMiss": 0,
	//             "libFeatInsert": 0,
	//             "resFeatStatInsert": 0,
	//             "resFeatStatUpdateAttempt": 0,
	//             "resFeatStatUpdateFail": 0
	//         },
	//         "lockWaits": {
	//             "refreshLocks": {
	//                 "maxMS": 0,
	//                 "totalMS": 0,
	//                 "count": 0
	//             }
	//         },
	//         "unresolve": {
	//             "triggers": {
	//                 "normalResolve": 0,
	//                 "update": 0,
	//                 "relLink": 0,
	//                 "extensiveResolve": 0,
	//                 "ambiguousNoResolve": 0,
	//                 "ambiguousMultiResolve": 0
	//             },
	//             "unresolveTest": 0,
	//             "abortedUnresolve": 0
	//         },
	//         "reresolve": {
	//             "triggers": {
	//                 "skipped": 0,
	//                 "abortRetry": 0,
	//                 "unresolveMovement": 0,
	//                 "multipleResolvableCandidates": 0,
	//                 "resolveNewFeatures": 0,
	//                 "newFeatureFTypes": []
	//             },
	//             "suppressedCandidateBuildersForReresolve": [],
	//             "suppressedScoredFeatureTypeForReresolve": []
	//         },
	//         "expressedFeatures": {
	//             "calls": [],
	//             "created": []
	//         },
	//         "scoring": {
	//             "scoredPairs": [],
	//             "cacheHit": [],
	//             "cacheMiss": [],
	//             "suppressedScoredFeatureType": [],
	//             "suppressedDisclosedRelationshipDomainCount": 0
	//         },
	//         "redoTriggers": [],
	//         "contention": {
	//             "valuelatch": [],
	//             "feature": [],
	//             "resent": []
	//         },
	//         "genericDetect": [],
	//         "candidates": {
	//             "candidateBuilders": [],
	//             "suppressedCandidateBuilders": []
	//         },
	//         "repairDiagnosis": {
	//             "types": 0
	//         },
	//         "threadState": {
	//             "active": 0,
	//             "idle": 5,
	//             "governorContention": 0,
	//             "sqlExecuting": 0,
	//             "loader": 0,
	//             "resolver": 0,
	//             "scoring": 0,
	//             "dataLatchContention": 0,
	//             "obsEntContention": 0,
	//             "resEntContention": 0
	//         },
	//         "systemResources": {
	//             "initResources": [
	//                 {
	//                     "physicalCores": 16
	//                 },
	//                 {
	//                     "logicalCores": 16
	//                 },
	//                 {
	//                     "totalMemory": "62.6GB"
	//                 },
	//                 {
	//                     "availableMemory": "53.9GB"
	//                 }
	//             ],
	//             "currResources": [
	//                 {
	//                     "availableMemory": "49.4GB"
	//                 },
	//                 {
	//                     "processMemory": "5.4GB"
	//                 },
	//                 {
	//                     "activeThreads": 0
	//                 },
	//                 {
	//                     "workerThreads": 5
	//                 },
	//                 {
	//                     "systemLoad": [
	//                         {
	//                             "cpuUser": 6.1
	//                         },
	//                         {
	//                             "cpuSystem": 0
	//                         },
	//                         {
	//                             "cpuIdle": 93.9
	//                         },
	//                         {
	//                             "cpuWait": 0
	//                         }
	//                     ]
	//                 }
	//             ]
	//         }
	//     }
	// }
}

func ExampleSzengine_GetVirtualEntityByRecordID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	recordList := `{"RECORDS": [{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1001"},{"DATA_SOURCE": "CUSTOMERS","RECORD_ID": "1002"}]}`
	flags := senzing.SzNoFlags

	result, err := szEngine.GetVirtualEntityByRecordID(ctx, recordList, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "RESOLVED_ENTITY": {
	//         "ENTITY_ID": 100001
	//     }
	// }
}

func ExampleSzengine_GetVirtualEntityByRecordID_output() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	// The following code pretty-prints example output JSON.
	exampleOutput := `{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDR_KEY":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"DOB":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"FEAT_DESC_VALUES":[{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ID_KEY":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"FEAT_DESC_VALUES":[{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME_KEY":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE_KEY":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"FEAT_DESC_VALUES":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SEARCH_KEY":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SSN":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"FEAT_DESC_VALUES":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"FEAT_DESC_VALUES":[{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":2,"FIRST_SEEN_DT":"2022-12-06 15:20:17.088","LAST_SEEN_DT":"2022-12-06 15:20:17.161"}],"LAST_SEEN_DT":"2022-12-06 15:20:17.161","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","LAST_SEEN_DT":"2022-12-06 15:20:17.088","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"222","ENTITY_TYPE":"TEST","INTERNAL_ID":2,"ENTITY_KEY":"740BA22D15CA88462A930AF8A7C904FF5E48226C","ENTITY_DESC":"OCEANGUY","LAST_SEEN_DT":"2022-12-06 15:20:17.161","FEATURES":[{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":24},{"LIB_FEAT_ID":25},{"LIB_FEAT_ID":26},{"LIB_FEAT_ID":27},{"LIB_FEAT_ID":28},{"LIB_FEAT_ID":29},{"LIB_FEAT_ID":30},{"LIB_FEAT_ID":31},{"LIB_FEAT_ID":32},{"LIB_FEAT_ID":33},{"LIB_FEAT_ID":34},{"LIB_FEAT_ID":35},{"LIB_FEAT_ID":36},{"LIB_FEAT_ID":37},{"LIB_FEAT_ID":38},{"LIB_FEAT_ID":39},{"LIB_FEAT_ID":40}]}]}}`
	fmt.Println(jsonutil.PrettyPrint(exampleOutput, jsonIndentation))
	// Output:
	// {
	//     "RESOLVED_ENTITY": {
	//         "ENTITY_ID": 1,
	//         "ENTITY_NAME": "JOHNSON",
	//         "FEATURES": {
	//             "ACCT_NUM": [
	//                 {
	//                     "FEAT_DESC": "5534202208773608",
	//                     "LIB_FEAT_ID": 8,
	//                     "USAGE_TYPE": "CC",
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "5534202208773608",
	//                             "LIB_FEAT_ID": 8,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "Y",
	//                             "ENTITY_COUNT": 3,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "ADDRESS": [
	//                 {
	//                     "FEAT_DESC": "772 Armstrong RD Delhi LA 71232",
	//                     "LIB_FEAT_ID": 4,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "772 Armstrong RD Delhi LA 71232",
	//                             "LIB_FEAT_ID": 4,
	//                             "USED_FOR_CAND": "N",
	//                             "USED_FOR_SCORING": "Y",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "772 Armstrong RD Delhi WI 53543",
	//                     "LIB_FEAT_ID": 26,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "772 Armstrong RD Delhi WI 53543",
	//                             "LIB_FEAT_ID": 26,
	//                             "USED_FOR_CAND": "N",
	//                             "USED_FOR_SCORING": "Y",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "ADDR_KEY": [
	//                 {
	//                     "FEAT_DESC": "772|ARMSTRNK||53543",
	//                     "LIB_FEAT_ID": 37,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "772|ARMSTRNK||53543",
	//                             "LIB_FEAT_ID": 37,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "772|ARMSTRNK||71232",
	//                     "LIB_FEAT_ID": 18,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "772|ARMSTRNK||71232",
	//                             "LIB_FEAT_ID": 18,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "772|ARMSTRNK||TL",
	//                     "LIB_FEAT_ID": 17,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "772|ARMSTRNK||TL",
	//                             "LIB_FEAT_ID": 17,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 3,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "DOB": [
	//                 {
	//                     "FEAT_DESC": "4/8/1983",
	//                     "LIB_FEAT_ID": 2,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "4/8/1983",
	//                             "LIB_FEAT_ID": 2,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "Y",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "6/9/1983",
	//                     "LIB_FEAT_ID": 25,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "6/9/1983",
	//                             "LIB_FEAT_ID": 25,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "Y",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "GENDER": [
	//                 {
	//                     "FEAT_DESC": "F",
	//                     "LIB_FEAT_ID": 3,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "F",
	//                             "LIB_FEAT_ID": 3,
	//                             "USED_FOR_CAND": "N",
	//                             "USED_FOR_SCORING": "Y",
	//                             "ENTITY_COUNT": 3,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "ID_KEY": [
	//                 {
	//                     "FEAT_DESC": "ACCT_NUM=5534202208773608",
	//                     "LIB_FEAT_ID": 19,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "ACCT_NUM=5534202208773608",
	//                             "LIB_FEAT_ID": 19,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 3,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "SSN=053-39-3251",
	//                     "LIB_FEAT_ID": 20,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "SSN=053-39-3251",
	//                             "LIB_FEAT_ID": 20,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "SSN=153-33-5185",
	//                     "LIB_FEAT_ID": 38,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "SSN=153-33-5185",
	//                             "LIB_FEAT_ID": 38,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "LOGIN_ID": [
	//                 {
	//                     "FEAT_DESC": "flavorh",
	//                     "LIB_FEAT_ID": 7,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "flavorh",
	//                             "LIB_FEAT_ID": 7,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "Y",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "flavorh2",
	//                     "LIB_FEAT_ID": 28,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "flavorh2",
	//                             "LIB_FEAT_ID": 28,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "Y",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "NAME": [
	//                 {
	//                     "FEAT_DESC": "JOHNSON",
	//                     "LIB_FEAT_ID": 1,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "JOHNSON",
	//                             "LIB_FEAT_ID": 1,
	//                             "USED_FOR_CAND": "N",
	//                             "USED_FOR_SCORING": "Y",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "OCEANGUY",
	//                     "LIB_FEAT_ID": 24,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "OCEANGUY",
	//                             "LIB_FEAT_ID": 24,
	//                             "USED_FOR_CAND": "N",
	//                             "USED_FOR_SCORING": "Y",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "NAME_KEY": [
	//                 {
	//                     "FEAT_DESC": "ASNK",
	//                     "LIB_FEAT_ID": 29,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "ASNK",
	//                             "LIB_FEAT_ID": 29,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "ASNK|ADDRESS.CITY_STD=TL",
	//                     "LIB_FEAT_ID": 34,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "ASNK|ADDRESS.CITY_STD=TL",
	//                             "LIB_FEAT_ID": 34,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "ASNK|DOB.MMDD_HASH=0906",
	//                     "LIB_FEAT_ID": 32,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "ASNK|DOB.MMDD_HASH=0906",
	//                             "LIB_FEAT_ID": 32,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "ASNK|DOB.MMYY_HASH=0683",
	//                     "LIB_FEAT_ID": 30,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "ASNK|DOB.MMYY_HASH=0683",
	//                             "LIB_FEAT_ID": 30,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "ASNK|DOB=80906",
	//                     "LIB_FEAT_ID": 31,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "ASNK|DOB=80906",
	//                             "LIB_FEAT_ID": 31,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "ASNK|PHONE.PHONE_LAST_5=10796",
	//                     "LIB_FEAT_ID": 33,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "ASNK|PHONE.PHONE_LAST_5=10796",
	//                             "LIB_FEAT_ID": 33,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "ASNK|POST=53543",
	//                     "LIB_FEAT_ID": 36,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "ASNK|POST=53543",
	//                             "LIB_FEAT_ID": 36,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "ASNK|SSN=5185",
	//                     "LIB_FEAT_ID": 35,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "ASNK|SSN=5185",
	//                             "LIB_FEAT_ID": 35,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "JNSN",
	//                     "LIB_FEAT_ID": 11,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "JNSN",
	//                             "LIB_FEAT_ID": 11,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "JNSN|ADDRESS.CITY_STD=TL",
	//                     "LIB_FEAT_ID": 12,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "JNSN|ADDRESS.CITY_STD=TL",
	//                             "LIB_FEAT_ID": 12,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "JNSN|DOB.MMDD_HASH=0804",
	//                     "LIB_FEAT_ID": 9,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "JNSN|DOB.MMDD_HASH=0804",
	//                             "LIB_FEAT_ID": 9,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "JNSN|DOB.MMYY_HASH=0483",
	//                     "LIB_FEAT_ID": 10,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "JNSN|DOB.MMYY_HASH=0483",
	//                             "LIB_FEAT_ID": 10,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "JNSN|DOB=80804",
	//                     "LIB_FEAT_ID": 13,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "JNSN|DOB=80804",
	//                             "LIB_FEAT_ID": 13,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "JNSN|PHONE.PHONE_LAST_5=10796",
	//                     "LIB_FEAT_ID": 15,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "JNSN|PHONE.PHONE_LAST_5=10796",
	//                             "LIB_FEAT_ID": 15,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "JNSN|POST=71232",
	//                     "LIB_FEAT_ID": 14,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "JNSN|POST=71232",
	//                             "LIB_FEAT_ID": 14,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "JNSN|SSN=3251",
	//                     "LIB_FEAT_ID": 16,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "JNSN|SSN=3251",
	//                             "LIB_FEAT_ID": 16,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "PHONE": [
	//                 {
	//                     "FEAT_DESC": "225-671-0796",
	//                     "LIB_FEAT_ID": 5,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "225-671-0796",
	//                             "LIB_FEAT_ID": 5,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "Y",
	//                             "ENTITY_COUNT": 3,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "PHONE_KEY": [
	//                 {
	//                     "FEAT_DESC": "2256710796",
	//                     "LIB_FEAT_ID": 21,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "2256710796",
	//                             "LIB_FEAT_ID": 21,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 3,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "SEARCH_KEY": [
	//                 {
	//                     "FEAT_DESC": "LOGIN_ID:FLAVORH2|",
	//                     "LIB_FEAT_ID": 40,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "LOGIN_ID:FLAVORH2|",
	//                             "LIB_FEAT_ID": 40,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "LOGIN_ID:FLAVORH|",
	//                     "LIB_FEAT_ID": 22,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "LOGIN_ID:FLAVORH|",
	//                             "LIB_FEAT_ID": 22,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "SSN:3251|80804|",
	//                     "LIB_FEAT_ID": 23,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "SSN:3251|80804|",
	//                             "LIB_FEAT_ID": 23,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "SSN:5185|80906|",
	//                     "LIB_FEAT_ID": 39,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "SSN:5185|80906|",
	//                             "LIB_FEAT_ID": 39,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "N",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 }
	//             ],
	//             "SSN": [
	//                 {
	//                     "FEAT_DESC": "053-39-3251",
	//                     "LIB_FEAT_ID": 6,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "053-39-3251",
	//                             "LIB_FEAT_ID": 6,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "Y",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 },
	//                 {
	//                     "FEAT_DESC": "153-33-5185",
	//                     "LIB_FEAT_ID": 27,
	//                     "FEAT_DESC_VALUES": [
	//                         {
	//                             "FEAT_DESC": "153-33-5185",
	//                             "LIB_FEAT_ID": 27,
	//                             "USED_FOR_CAND": "Y",
	//                             "USED_FOR_SCORING": "Y",
	//                             "ENTITY_COUNT": 1,
	//                             "CANDIDATE_CAP_REACHED": "N",
	//                             "SCORING_CAP_REACHED": "N",
	//                             "SUPPRESSED": "N"
	//                         }
	//                     ]
	//                 }
	//             ]
	//         },
	//         "RECORD_SUMMARY": [
	//             {
	//                 "DATA_SOURCE": "TEST",
	//                 "RECORD_COUNT": 2,
	//                 "FIRST_SEEN_DT": "2022-12-06 15:20:17.088",
	//                 "LAST_SEEN_DT": "2022-12-06 15:20:17.161"
	//             }
	//         ],
	//         "LAST_SEEN_DT": "2022-12-06 15:20:17.161",
	//         "RECORDS": [
	//             {
	//                 "DATA_SOURCE": "TEST",
	//                 "RECORD_ID": "111",
	//                 "ENTITY_TYPE": "TEST",
	//                 "INTERNAL_ID": 1,
	//                 "ENTITY_KEY": "C6063D4396612FBA7324DB0739273BA1FE815C43",
	//                 "ENTITY_DESC": "JOHNSON",
	//                 "LAST_SEEN_DT": "2022-12-06 15:20:17.088",
	//                 "FEATURES": [
	//                     {
	//                         "LIB_FEAT_ID": 1
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 2
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 3
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 4
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 5
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 6
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 7
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 8,
	//                         "USAGE_TYPE": "CC"
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 9
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 10
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 11
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 12
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 13
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 14
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 15
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 16
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 17
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 18
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 19
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 20
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 21
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 22
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 23
	//                     }
	//                 ]
	//             },
	//             {
	//                 "DATA_SOURCE": "TEST",
	//                 "RECORD_ID": "222",
	//                 "ENTITY_TYPE": "TEST",
	//                 "INTERNAL_ID": 2,
	//                 "ENTITY_KEY": "740BA22D15CA88462A930AF8A7C904FF5E48226C",
	//                 "ENTITY_DESC": "OCEANGUY",
	//                 "LAST_SEEN_DT": "2022-12-06 15:20:17.161",
	//                 "FEATURES": [
	//                     {
	//                         "LIB_FEAT_ID": 3
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 5
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 8,
	//                         "USAGE_TYPE": "CC"
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 17
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 19
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 21
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 24
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 25
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 26
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 27
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 28
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 29
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 30
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 31
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 32
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 33
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 34
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 35
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 36
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 37
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 38
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 39
	//                     },
	//                     {
	//                         "LIB_FEAT_ID": 40
	//                     }
	//                 ]
	//             }
	//         ]
	//     }
	// }
}

func ExampleSzengine_HowEntityByEntityID() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	entityID := getEntityIDForRecord(ctx, szEngine, "CUSTOMERS", "1001")

	flags := senzing.SzNoFlags

	result, err := szEngine.HowEntityByEntityID(ctx, entityID, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.Truncate(result, 5))
	// Output: {"HOW_RESULTS":{"FINAL_STATE":{"NEED_REEVALUATION":0,"VIRTUAL_ENTITIES":[...
}

func ExampleSzengine_GetRecordPreview() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	recordDefinition := `{"DATA_SOURCE": "CUSTOMERS", "RECORD_ID": "1001", "RECORD_TYPE": "PERSON", "PRIMARY_NAME_LAST": "Smith", "PRIMARY_NAME_FIRST": "Robert", "DATE_OF_BIRTH": "12/11/1978", "ADDR_TYPE": "MAILING", "ADDR_LINE1": "123 Main Street, Las Vegas NV 89132", "PHONE_TYPE": "HOME", "PHONE_NUMBER": "702-919-1300", "EMAIL_ADDRESS": "bsmith@work.com", "DATE": "1/2/18", "STATUS": "Active", "AMOUNT": "100"}`
	flags := senzing.SzNoFlags

	result, err := szEngine.GetRecordPreview(ctx, recordDefinition, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(result)
	// Output: {}
}

func ExampleSzengine_PrimeEngine() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	err = szEngine.PrimeEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}
	// Output:
}

func ExampleSzEngine_ProcessRedoRecord() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	redoRecord, err := szEngine.GetRedoRecord(ctx)
	if err != nil {
		handleError(err)
		return
	}

	flags := senzing.SzWithoutInfo

	result, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(result)
	// Output:
}

func ExampleSzEngine_ProcessRedoRecord_withInfo() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	redoRecord, err := szEngine.GetRedoRecord(ctx)
	if err != nil {
		handleError(err)
		return
	}

	flags := senzing.SzWithInfo

	result, err := szEngine.ProcessRedoRecord(ctx, redoRecord, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "DATA_SOURCE": "",
	//     "RECORD_ID": "",
	//     "AFFECTED_ENTITIES": [
	//         {
	//             "ENTITY_ID": 90
	//         }
	//     ]
	// }
}

func ExampleSzengine_SearchByAttributes() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	searchProfile := ""
	flags := senzing.SzNoFlags

	result, err := szEngine.SearchByAttributes(ctx, attributes, searchProfile, flags)
	if err != nil {
		handleError(err)
		return
	}

	redactKeys := []string{"FIRST_SEEN_DT", "LAST_SEEN_DT"}
	fmt.Println(
		jsonutil.PrettyPrint(
			jsonutil.Flatten(
				jsonutil.Redact(
					jsonutil.Flatten(jsonutil.NormalizeAndSort(result)),
					redactKeys...,
				),
			), jsonIndentation))
	// Output:
	// {
	//     "RESOLVED_ENTITIES": [
	//         {
	//             "ENTITY": {
	//                 "RESOLVED_ENTITY": {
	//                     "ENTITY_ID": 100001
	//                 }
	//             },
	//             "MATCH_INFO": {
	//                 "ERRULE_CODE": "SF1",
	//                 "MATCH_KEY": "+PNAME+EMAIL",
	//                 "MATCH_LEVEL_CODE": "POSSIBLY_RELATED"
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzEngine_SearchByAttributes_searchProfile() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	searchProfile := "SEARCH"
	flags := senzing.SzNoFlags

	result, err := szEngine.SearchByAttributes(ctx, attributes, searchProfile, flags)
	if err != nil {
		handleError(err)
		return
	}

	redactKeys := []string{"FIRST_SEEN_DT", "LAST_SEEN_DT"}
	fmt.Println(
		jsonutil.PrettyPrint(
			jsonutil.Flatten(
				jsonutil.Redact(
					jsonutil.Flatten(jsonutil.NormalizeAndSort(result)),
					redactKeys...,
				),
			), jsonIndentation))
	// Output:
	// {
	//     "RESOLVED_ENTITIES": [
	//         {
	//             "ENTITY": {
	//                 "RESOLVED_ENTITY": {
	//                     "ENTITY_ID": 100001
	//                 }
	//             },
	//             "MATCH_INFO": {
	//                 "ERRULE_CODE": "SF1",
	//                 "MATCH_KEY": "+PNAME+EMAIL",
	//                 "MATCH_LEVEL_CODE": "POSSIBLY_RELATED"
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzengine_SearchByAttributes_output() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	// The following code pretty-prints example output JSON.
	exampleOutput := `{"RESOLVED_ENTITIES":[{"MATCH_INFO":{"MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","MATCH_KEY":"+NAME+SSN","ERRULE_CODE":"SF1_PNAME_CSTAB","FEATURE_SCORES":{"NAME":[{"INBOUND_FEAT":"JOHNSON","CANDIDATE_FEAT":"JOHNSON","GNR_FN":100,"GNR_SN":100,"GNR_GN":70,"GENERATION_MATCH":-1,"GNR_ON":-1}],"SSN":[{"INBOUND_FEAT":"053-39-3251","CANDIDATE_FEAT":"053-39-3251","FULL_SCORE":100}]}},"ENTITY":{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4}]}],"DOB":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2}]},{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7}]}],"NAME":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5}]}],"SSN":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"FEAT_DESC_VALUES":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":6,"FIRST_SEEN_DT":"2022-12-06 15:38:06.175","LAST_SEEN_DT":"2022-12-06 15:38:06.957"}],"LAST_SEEN_DT":"2022-12-06 15:38:06.957"}}}]}`
	fmt.Println(jsonutil.PrettyPrint(exampleOutput, jsonIndentation))
	// Output:
	// {
	//     "RESOLVED_ENTITIES": [
	//         {
	//             "MATCH_INFO": {
	//                 "MATCH_LEVEL": 1,
	//                 "MATCH_LEVEL_CODE": "RESOLVED",
	//                 "MATCH_KEY": "+NAME+SSN",
	//                 "ERRULE_CODE": "SF1_PNAME_CSTAB",
	//                 "FEATURE_SCORES": {
	//                     "NAME": [
	//                         {
	//                             "INBOUND_FEAT": "JOHNSON",
	//                             "CANDIDATE_FEAT": "JOHNSON",
	//                             "GNR_FN": 100,
	//                             "GNR_SN": 100,
	//                             "GNR_GN": 70,
	//                             "GENERATION_MATCH": -1,
	//                             "GNR_ON": -1
	//                         }
	//                     ],
	//                     "SSN": [
	//                         {
	//                             "INBOUND_FEAT": "053-39-3251",
	//                             "CANDIDATE_FEAT": "053-39-3251",
	//                             "FULL_SCORE": 100
	//                         }
	//                     ]
	//                 }
	//             },
	//             "ENTITY": {
	//                 "RESOLVED_ENTITY": {
	//                     "ENTITY_ID": 1,
	//                     "ENTITY_NAME": "JOHNSON",
	//                     "FEATURES": {
	//                         "ACCT_NUM": [
	//                             {
	//                                 "FEAT_DESC": "5534202208773608",
	//                                 "LIB_FEAT_ID": 8,
	//                                 "USAGE_TYPE": "CC",
	//                                 "FEAT_DESC_VALUES": [
	//                                     {
	//                                         "FEAT_DESC": "5534202208773608",
	//                                         "LIB_FEAT_ID": 8
	//                                     }
	//                                 ]
	//                             }
	//                         ],
	//                         "ADDRESS": [
	//                             {
	//                                 "FEAT_DESC": "772 Armstrong RD Delhi LA 71232",
	//                                 "LIB_FEAT_ID": 4,
	//                                 "FEAT_DESC_VALUES": [
	//                                     {
	//                                         "FEAT_DESC": "772 Armstrong RD Delhi LA 71232",
	//                                         "LIB_FEAT_ID": 4
	//                                     }
	//                                 ]
	//                             }
	//                         ],
	//                         "DOB": [
	//                             {
	//                                 "FEAT_DESC": "4/8/1983",
	//                                 "LIB_FEAT_ID": 2,
	//                                 "FEAT_DESC_VALUES": [
	//                                     {
	//                                         "FEAT_DESC": "4/8/1983",
	//                                         "LIB_FEAT_ID": 2
	//                                     }
	//                                 ]
	//                             },
	//                             {
	//                                 "FEAT_DESC": "4/8/1985",
	//                                 "LIB_FEAT_ID": 100001,
	//                                 "FEAT_DESC_VALUES": [
	//                                     {
	//                                         "FEAT_DESC": "4/8/1985",
	//                                         "LIB_FEAT_ID": 100001
	//                                     }
	//                                 ]
	//                             }
	//                         ],
	//                         "GENDER": [
	//                             {
	//                                 "FEAT_DESC": "F",
	//                                 "LIB_FEAT_ID": 3,
	//                                 "FEAT_DESC_VALUES": [
	//                                     {
	//                                         "FEAT_DESC": "F",
	//                                         "LIB_FEAT_ID": 3
	//                                     }
	//                                 ]
	//                             }
	//                         ],
	//                         "LOGIN_ID": [
	//                             {
	//                                 "FEAT_DESC": "flavorh",
	//                                 "LIB_FEAT_ID": 7,
	//                                 "FEAT_DESC_VALUES": [
	//                                     {
	//                                         "FEAT_DESC": "flavorh",
	//                                         "LIB_FEAT_ID": 7
	//                                     }
	//                                 ]
	//                             }
	//                         ],
	//                         "NAME": [
	//                             {
	//                                 "FEAT_DESC": "JOHNSON",
	//                                 "LIB_FEAT_ID": 1,
	//                                 "FEAT_DESC_VALUES": [
	//                                     {
	//                                         "FEAT_DESC": "JOHNSON",
	//                                         "LIB_FEAT_ID": 1
	//                                     }
	//                                 ]
	//                             }
	//                         ],
	//                         "PHONE": [
	//                             {
	//                                 "FEAT_DESC": "225-671-0796",
	//                                 "LIB_FEAT_ID": 5,
	//                                 "FEAT_DESC_VALUES": [
	//                                     {
	//                                         "FEAT_DESC": "225-671-0796",
	//                                         "LIB_FEAT_ID": 5
	//                                     }
	//                                 ]
	//                             }
	//                         ],
	//                         "SSN": [
	//                             {
	//                                 "FEAT_DESC": "053-39-3251",
	//                                 "LIB_FEAT_ID": 6,
	//                                 "FEAT_DESC_VALUES": [
	//                                     {
	//                                         "FEAT_DESC": "053-39-3251",
	//                                         "LIB_FEAT_ID": 6
	//                                     }
	//                                 ]
	//                             }
	//                         ]
	//                     },
	//                     "RECORD_SUMMARY": [
	//                         {
	//                             "DATA_SOURCE": "TEST",
	//                             "RECORD_COUNT": 6,
	//                             "FIRST_SEEN_DT": "2022-12-06 15:38:06.175",
	//                             "LAST_SEEN_DT": "2022-12-06 15:38:06.957"
	//                         }
	//                     ],
	//                     "LAST_SEEN_DT": "2022-12-06 15:38:06.957"
	//                 }
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzengine_WhyEntities() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	entityID1 := getEntityID(ctx, szEngine, truthset.CustomerRecords["1001"])
	entityID2 := getEntityID(ctx, szEngine, truthset.CustomerRecords["1002"])
	flags := senzing.SzNoFlags

	result, err := szEngine.WhyEntities(ctx, entityID1, entityID2, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "WHY_RESULTS": [
	//         {
	//             "ENTITY_ID": 100001,
	//             "ENTITY_ID_2": 100001,
	//             "MATCH_INFO": {
	//                 "WHY_KEY": "+NAME+DOB+ADDRESS+PHONE+EMAIL",
	//                 "WHY_ERRULE_CODE": "SF1_SNAME_CFF_CSTAB",
	//                 "MATCH_LEVEL_CODE": "RESOLVED"
	//             }
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100001
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzengine_WhyEntities_output() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	// The following code pretty-prints example output JSON.
	exampleOutput := `{"WHY_RESULTS":[{"ENTITY_ID":1,"ENTITY_ID_2":2,"MATCH_INFO":{"WHY_KEY":"+PHONE+ACCT_NUM-SSN","WHY_ERRULE_CODE":"SF1","MATCH_LEVEL_CODE":"POSSIBLY_RELATED","CANDIDATE_KEYS":{"ACCT_NUM":[{"FEAT_ID":8,"FEAT_DESC":"5534202208773608"}],"ADDR_KEY":[{"FEAT_ID":17,"FEAT_DESC":"772|ARMSTRNK||TL"}],"ID_KEY":[{"FEAT_ID":19,"FEAT_DESC":"ACCT_NUM=5534202208773608"}],"PHONE":[{"FEAT_ID":5,"FEAT_DESC":"225-671-0796"}],"PHONE_KEY":[{"FEAT_ID":21,"FEAT_DESC":"2256710796"}]},"DISCLOSED_RELATIONS":{},"FEATURE_SCORES":{"ACCT_NUM":[{"INBOUND_FEAT_ID":8,"INBOUND_FEAT":"5534202208773608","INBOUND_FEAT_USAGE_TYPE":"CC","CANDIDATE_FEAT_ID":8,"CANDIDATE_FEAT":"5534202208773608","CANDIDATE_FEAT_USAGE_TYPE":"CC","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"F1"}],"ADDRESS":[{"INBOUND_FEAT_ID":4,"INBOUND_FEAT":"772 Armstrong RD Delhi LA 71232","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":26,"CANDIDATE_FEAT":"772 Armstrong RD Delhi WI 53543","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":81,"SCORE_BUCKET":"LIKELY","SCORE_BEHAVIOR":"FF"}],"DOB":[{"INBOUND_FEAT_ID":100001,"INBOUND_FEAT":"4/8/1985","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":25,"CANDIDATE_FEAT":"6/9/1983","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":79,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"FMES"},{"INBOUND_FEAT_ID":2,"INBOUND_FEAT":"4/8/1983","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":25,"CANDIDATE_FEAT":"6/9/1983","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":86,"SCORE_BUCKET":"PLAUSIBLE","SCORE_BEHAVIOR":"FMES"}],"GENDER":[{"INBOUND_FEAT_ID":3,"INBOUND_FEAT":"F","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":3,"CANDIDATE_FEAT":"F","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}],"LOGIN_ID":[{"INBOUND_FEAT_ID":7,"INBOUND_FEAT":"flavorh","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":28,"CANDIDATE_FEAT":"flavorh2","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":0,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"F1"}],"NAME":[{"INBOUND_FEAT_ID":1,"INBOUND_FEAT":"JOHNSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":24,"CANDIDATE_FEAT":"OCEANGUY","CANDIDATE_FEAT_USAGE_TYPE":"","GNR_FN":33,"GNR_SN":32,"GNR_GN":70,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"NAME"}],"PHONE":[{"INBOUND_FEAT_ID":5,"INBOUND_FEAT":"225-671-0796","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":5,"CANDIDATE_FEAT":"225-671-0796","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FF"}],"SSN":[{"INBOUND_FEAT_ID":6,"INBOUND_FEAT":"053-39-3251","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":27,"CANDIDATE_FEAT":"153-33-5185","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":0,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"F1ES"}]}}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDR_KEY":[{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"DOB":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ID_KEY":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME_KEY":[{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0485","LIB_FEAT_ID":100002,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0485","LIB_FEAT_ID":100002,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE_KEY":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"FEAT_DESC_VALUES":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SEARCH_KEY":[{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SSN":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"FEAT_DESC_VALUES":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":6,"FIRST_SEEN_DT":"2022-12-06 15:58:57.129","LAST_SEEN_DT":"2022-12-06 15:58:57.906"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.906","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111","ENTITY_TYPE":"TEST","INTERNAL_ID":100001,"ENTITY_KEY":"A6C927986DF7329D1D2CDE0E8F34328AE640FB7E","ENTITY_DESC":"JOHNSON","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 15:58:57.906","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23},{"LIB_FEAT_ID":100001},{"LIB_FEAT_ID":100002}]},{"DATA_SOURCE":"TEST","RECORD_ID":"444","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 15:58:57.400","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"555","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 15:58:57.404","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"666","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 15:58:57.407","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"777","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 15:58:57.410","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"FCCE9793DAAD23159DBCCEB97FF2745B92CE7919","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 15:58:57.259","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]}]},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:58:57.201","LAST_SEEN_DT":"2022-12-06 15:58:57.201"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.201"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:58:57.263","LAST_SEEN_DT":"2022-12-06 15:58:57.263"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.263"}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDR_KEY":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"DOB":[{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"FEAT_DESC_VALUES":[{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ID_KEY":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME":[{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"FEAT_DESC_VALUES":[{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME_KEY":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE_KEY":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"FEAT_DESC_VALUES":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SEARCH_KEY":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SSN":[{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"FEAT_DESC_VALUES":[{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:58:57.201","LAST_SEEN_DT":"2022-12-06 15:58:57.201"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.201","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"222","ENTITY_TYPE":"TEST","INTERNAL_ID":2,"ENTITY_KEY":"740BA22D15CA88462A930AF8A7C904FF5E48226C","ENTITY_DESC":"OCEANGUY","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 15:58:57.201","FEATURES":[{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":24},{"LIB_FEAT_ID":25},{"LIB_FEAT_ID":26},{"LIB_FEAT_ID":27},{"LIB_FEAT_ID":28},{"LIB_FEAT_ID":29},{"LIB_FEAT_ID":30},{"LIB_FEAT_ID":31},{"LIB_FEAT_ID":32},{"LIB_FEAT_ID":33},{"LIB_FEAT_ID":34},{"LIB_FEAT_ID":35},{"LIB_FEAT_ID":36},{"LIB_FEAT_ID":37},{"LIB_FEAT_ID":38},{"LIB_FEAT_ID":39},{"LIB_FEAT_ID":40}]}]},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":6,"FIRST_SEEN_DT":"2022-12-06 15:58:57.129","LAST_SEEN_DT":"2022-12-06 15:58:57.906"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.906"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 15:58:57.263","LAST_SEEN_DT":"2022-12-06 15:58:57.263"}],"LAST_SEEN_DT":"2022-12-06 15:58:57.263"}]}]}`
	fmt.Println(jsonutil.PrettyPrint(exampleOutput, jsonIndentation))
	// Output:
	// {
	//     "WHY_RESULTS": [
	//         {
	//             "ENTITY_ID": 1,
	//             "ENTITY_ID_2": 2,
	//             "MATCH_INFO": {
	//                 "WHY_KEY": "+PHONE+ACCT_NUM-SSN",
	//                 "WHY_ERRULE_CODE": "SF1",
	//                 "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                 "CANDIDATE_KEYS": {
	//                     "ACCT_NUM": [
	//                         {
	//                             "FEAT_ID": 8,
	//                             "FEAT_DESC": "5534202208773608"
	//                         }
	//                     ],
	//                     "ADDR_KEY": [
	//                         {
	//                             "FEAT_ID": 17,
	//                             "FEAT_DESC": "772|ARMSTRNK||TL"
	//                         }
	//                     ],
	//                     "ID_KEY": [
	//                         {
	//                             "FEAT_ID": 19,
	//                             "FEAT_DESC": "ACCT_NUM=5534202208773608"
	//                         }
	//                     ],
	//                     "PHONE": [
	//                         {
	//                             "FEAT_ID": 5,
	//                             "FEAT_DESC": "225-671-0796"
	//                         }
	//                     ],
	//                     "PHONE_KEY": [
	//                         {
	//                             "FEAT_ID": 21,
	//                             "FEAT_DESC": "2256710796"
	//                         }
	//                     ]
	//                 },
	//                 "DISCLOSED_RELATIONS": {},
	//                 "FEATURE_SCORES": {
	//                     "ACCT_NUM": [
	//                         {
	//                             "INBOUND_FEAT_ID": 8,
	//                             "INBOUND_FEAT": "5534202208773608",
	//                             "INBOUND_FEAT_USAGE_TYPE": "CC",
	//                             "CANDIDATE_FEAT_ID": 8,
	//                             "CANDIDATE_FEAT": "5534202208773608",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "CC",
	//                             "FULL_SCORE": 100,
	//                             "SCORE_BUCKET": "SAME",
	//                             "SCORE_BEHAVIOR": "F1"
	//                         }
	//                     ],
	//                     "ADDRESS": [
	//                         {
	//                             "INBOUND_FEAT_ID": 4,
	//                             "INBOUND_FEAT": "772 Armstrong RD Delhi LA 71232",
	//                             "INBOUND_FEAT_USAGE_TYPE": "",
	//                             "CANDIDATE_FEAT_ID": 26,
	//                             "CANDIDATE_FEAT": "772 Armstrong RD Delhi WI 53543",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "",
	//                             "FULL_SCORE": 81,
	//                             "SCORE_BUCKET": "LIKELY",
	//                             "SCORE_BEHAVIOR": "FF"
	//                         }
	//                     ],
	//                     "DOB": [
	//                         {
	//                             "INBOUND_FEAT_ID": 100001,
	//                             "INBOUND_FEAT": "4/8/1985",
	//                             "INBOUND_FEAT_USAGE_TYPE": "",
	//                             "CANDIDATE_FEAT_ID": 25,
	//                             "CANDIDATE_FEAT": "6/9/1983",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "",
	//                             "FULL_SCORE": 79,
	//                             "SCORE_BUCKET": "NO_CHANCE",
	//                             "SCORE_BEHAVIOR": "FMES"
	//                         },
	//                         {
	//                             "INBOUND_FEAT_ID": 2,
	//                             "INBOUND_FEAT": "4/8/1983",
	//                             "INBOUND_FEAT_USAGE_TYPE": "",
	//                             "CANDIDATE_FEAT_ID": 25,
	//                             "CANDIDATE_FEAT": "6/9/1983",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "",
	//                             "FULL_SCORE": 86,
	//                             "SCORE_BUCKET": "PLAUSIBLE",
	//                             "SCORE_BEHAVIOR": "FMES"
	//                         }
	//                     ],
	//                     "GENDER": [
	//                         {
	//                             "INBOUND_FEAT_ID": 3,
	//                             "INBOUND_FEAT": "F",
	//                             "INBOUND_FEAT_USAGE_TYPE": "",
	//                             "CANDIDATE_FEAT_ID": 3,
	//                             "CANDIDATE_FEAT": "F",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "",
	//                             "FULL_SCORE": 100,
	//                             "SCORE_BUCKET": "SAME",
	//                             "SCORE_BEHAVIOR": "FVME"
	//                         }
	//                     ],
	//                     "LOGIN_ID": [
	//                         {
	//                             "INBOUND_FEAT_ID": 7,
	//                             "INBOUND_FEAT": "flavorh",
	//                             "INBOUND_FEAT_USAGE_TYPE": "",
	//                             "CANDIDATE_FEAT_ID": 28,
	//                             "CANDIDATE_FEAT": "flavorh2",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "",
	//                             "FULL_SCORE": 0,
	//                             "SCORE_BUCKET": "NO_CHANCE",
	//                             "SCORE_BEHAVIOR": "F1"
	//                         }
	//                     ],
	//                     "NAME": [
	//                         {
	//                             "INBOUND_FEAT_ID": 1,
	//                             "INBOUND_FEAT": "JOHNSON",
	//                             "INBOUND_FEAT_USAGE_TYPE": "",
	//                             "CANDIDATE_FEAT_ID": 24,
	//                             "CANDIDATE_FEAT": "OCEANGUY",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "",
	//                             "GNR_FN": 33,
	//                             "GNR_SN": 32,
	//                             "GNR_GN": 70,
	//                             "GENERATION_MATCH": -1,
	//                             "GNR_ON": -1,
	//                             "SCORE_BUCKET": "NO_CHANCE",
	//                             "SCORE_BEHAVIOR": "NAME"
	//                         }
	//                     ],
	//                     "PHONE": [
	//                         {
	//                             "INBOUND_FEAT_ID": 5,
	//                             "INBOUND_FEAT": "225-671-0796",
	//                             "INBOUND_FEAT_USAGE_TYPE": "",
	//                             "CANDIDATE_FEAT_ID": 5,
	//                             "CANDIDATE_FEAT": "225-671-0796",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "",
	//                             "FULL_SCORE": 100,
	//                             "SCORE_BUCKET": "SAME",
	//                             "SCORE_BEHAVIOR": "FF"
	//                         }
	//                     ],
	//                     "SSN": [
	//                         {
	//                             "INBOUND_FEAT_ID": 6,
	//                             "INBOUND_FEAT": "053-39-3251",
	//                             "INBOUND_FEAT_USAGE_TYPE": "",
	//                             "CANDIDATE_FEAT_ID": 27,
	//                             "CANDIDATE_FEAT": "153-33-5185",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "",
	//                             "FULL_SCORE": 0,
	//                             "SCORE_BUCKET": "NO_CHANCE",
	//                             "SCORE_BEHAVIOR": "F1ES"
	//                         }
	//                     ]
	//                 }
	//             }
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 1,
	//                 "ENTITY_NAME": "JOHNSON",
	//                 "FEATURES": {
	//                     "ACCT_NUM": [
	//                         {
	//                             "FEAT_DESC": "5534202208773608",
	//                             "LIB_FEAT_ID": 8,
	//                             "USAGE_TYPE": "CC",
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "5534202208773608",
	//                                     "LIB_FEAT_ID": 8,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "ADDRESS": [
	//                         {
	//                             "FEAT_DESC": "772 Armstrong RD Delhi LA 71232",
	//                             "LIB_FEAT_ID": 4,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "772 Armstrong RD Delhi LA 71232",
	//                                     "LIB_FEAT_ID": 4,
	//                                     "USED_FOR_CAND": "N",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "ADDR_KEY": [
	//                         {
	//                             "FEAT_DESC": "772|ARMSTRNK||71232",
	//                             "LIB_FEAT_ID": 18,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "772|ARMSTRNK||71232",
	//                                     "LIB_FEAT_ID": 18,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "772|ARMSTRNK||TL",
	//                             "LIB_FEAT_ID": 17,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "772|ARMSTRNK||TL",
	//                                     "LIB_FEAT_ID": 17,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "DOB": [
	//                         {
	//                             "FEAT_DESC": "4/8/1983",
	//                             "LIB_FEAT_ID": 2,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "4/8/1983",
	//                                     "LIB_FEAT_ID": 2,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "4/8/1985",
	//                             "LIB_FEAT_ID": 100001,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "4/8/1985",
	//                                     "LIB_FEAT_ID": 100001,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "GENDER": [
	//                         {
	//                             "FEAT_DESC": "F",
	//                             "LIB_FEAT_ID": 3,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "F",
	//                                     "LIB_FEAT_ID": 3,
	//                                     "USED_FOR_CAND": "N",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "ID_KEY": [
	//                         {
	//                             "FEAT_DESC": "ACCT_NUM=5534202208773608",
	//                             "LIB_FEAT_ID": 19,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ACCT_NUM=5534202208773608",
	//                                     "LIB_FEAT_ID": 19,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "SSN=053-39-3251",
	//                             "LIB_FEAT_ID": 20,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "SSN=053-39-3251",
	//                                     "LIB_FEAT_ID": 20,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "LOGIN_ID": [
	//                         {
	//                             "FEAT_DESC": "flavorh",
	//                             "LIB_FEAT_ID": 7,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "flavorh",
	//                                     "LIB_FEAT_ID": 7,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "NAME": [
	//                         {
	//                             "FEAT_DESC": "JOHNSON",
	//                             "LIB_FEAT_ID": 1,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JOHNSON",
	//                                     "LIB_FEAT_ID": 1,
	//                                     "USED_FOR_CAND": "N",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "NAME_KEY": [
	//                         {
	//                             "FEAT_DESC": "JNSN",
	//                             "LIB_FEAT_ID": 11,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN",
	//                                     "LIB_FEAT_ID": 11,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|ADDRESS.CITY_STD=TL",
	//                             "LIB_FEAT_ID": 12,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|ADDRESS.CITY_STD=TL",
	//                                     "LIB_FEAT_ID": 12,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|DOB.MMDD_HASH=0804",
	//                             "LIB_FEAT_ID": 9,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|DOB.MMDD_HASH=0804",
	//                                     "LIB_FEAT_ID": 9,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|DOB.MMYY_HASH=0483",
	//                             "LIB_FEAT_ID": 10,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|DOB.MMYY_HASH=0483",
	//                                     "LIB_FEAT_ID": 10,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|DOB.MMYY_HASH=0485",
	//                             "LIB_FEAT_ID": 100002,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|DOB.MMYY_HASH=0485",
	//                                     "LIB_FEAT_ID": 100002,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|DOB=80804",
	//                             "LIB_FEAT_ID": 13,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|DOB=80804",
	//                                     "LIB_FEAT_ID": 13,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|PHONE.PHONE_LAST_5=10796",
	//                             "LIB_FEAT_ID": 15,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|PHONE.PHONE_LAST_5=10796",
	//                                     "LIB_FEAT_ID": 15,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|POST=71232",
	//                             "LIB_FEAT_ID": 14,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|POST=71232",
	//                                     "LIB_FEAT_ID": 14,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|SSN=3251",
	//                             "LIB_FEAT_ID": 16,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|SSN=3251",
	//                                     "LIB_FEAT_ID": 16,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "PHONE": [
	//                         {
	//                             "FEAT_DESC": "225-671-0796",
	//                             "LIB_FEAT_ID": 5,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "225-671-0796",
	//                                     "LIB_FEAT_ID": 5,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "PHONE_KEY": [
	//                         {
	//                             "FEAT_DESC": "2256710796",
	//                             "LIB_FEAT_ID": 21,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "2256710796",
	//                                     "LIB_FEAT_ID": 21,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "SEARCH_KEY": [
	//                         {
	//                             "FEAT_DESC": "LOGIN_ID:FLAVORH|",
	//                             "LIB_FEAT_ID": 22,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "LOGIN_ID:FLAVORH|",
	//                                     "LIB_FEAT_ID": 22,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "SSN:3251|80804|",
	//                             "LIB_FEAT_ID": 23,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "SSN:3251|80804|",
	//                                     "LIB_FEAT_ID": 23,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "SSN": [
	//                         {
	//                             "FEAT_DESC": "053-39-3251",
	//                             "LIB_FEAT_ID": 6,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "053-39-3251",
	//                                     "LIB_FEAT_ID": 6,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ]
	//                 },
	//                 "RECORD_SUMMARY": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_COUNT": 6,
	//                         "FIRST_SEEN_DT": "2022-12-06 15:58:57.129",
	//                         "LAST_SEEN_DT": "2022-12-06 15:58:57.906"
	//                     }
	//                 ],
	//                 "LAST_SEEN_DT": "2022-12-06 15:58:57.906",
	//                 "RECORDS": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_ID": "111",
	//                         "ENTITY_TYPE": "TEST",
	//                         "INTERNAL_ID": 100001,
	//                         "ENTITY_KEY": "A6C927986DF7329D1D2CDE0E8F34328AE640FB7E",
	//                         "ENTITY_DESC": "JOHNSON",
	//                         "MATCH_KEY": "",
	//                         "MATCH_LEVEL": 0,
	//                         "MATCH_LEVEL_CODE": "",
	//                         "ERRULE_CODE": "",
	//                         "LAST_SEEN_DT": "2022-12-06 15:58:57.906",
	//                         "FEATURES": [
	//                             {
	//                                 "LIB_FEAT_ID": 1
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 3
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 4
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 5
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 6
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 7
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 8,
	//                                 "USAGE_TYPE": "CC"
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 9
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 11
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 12
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 13
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 14
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 15
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 16
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 17
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 18
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 19
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 20
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 21
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 22
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 23
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 100001
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 100002
	//                             }
	//                         ]
	//                     },
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_ID": "444",
	//                         "ENTITY_TYPE": "TEST",
	//                         "INTERNAL_ID": 1,
	//                         "ENTITY_KEY": "C6063D4396612FBA7324DB0739273BA1FE815C43",
	//                         "ENTITY_DESC": "JOHNSON",
	//                         "MATCH_KEY": "+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM",
	//                         "MATCH_LEVEL": 1,
	//                         "MATCH_LEVEL_CODE": "RESOLVED",
	//                         "ERRULE_CODE": "SF1_PNAME_CFF_CSTAB",
	//                         "LAST_SEEN_DT": "2022-12-06 15:58:57.400",
	//                         "FEATURES": [
	//                             {
	//                                 "LIB_FEAT_ID": 1
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 2
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 3
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 4
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 5
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 6
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 7
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 8,
	//                                 "USAGE_TYPE": "CC"
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 9
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 10
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 11
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 12
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 13
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 14
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 15
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 16
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 17
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 18
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 19
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 20
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 21
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 22
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 23
	//                             }
	//                         ]
	//                     },
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_ID": "555",
	//                         "ENTITY_TYPE": "TEST",
	//                         "INTERNAL_ID": 1,
	//                         "ENTITY_KEY": "C6063D4396612FBA7324DB0739273BA1FE815C43",
	//                         "ENTITY_DESC": "JOHNSON",
	//                         "MATCH_KEY": "+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM",
	//                         "MATCH_LEVEL": 1,
	//                         "MATCH_LEVEL_CODE": "RESOLVED",
	//                         "ERRULE_CODE": "SF1_PNAME_CFF_CSTAB",
	//                         "LAST_SEEN_DT": "2022-12-06 15:58:57.404",
	//                         "FEATURES": [
	//                             {
	//                                 "LIB_FEAT_ID": 1
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 2
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 3
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 4
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 5
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 6
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 7
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 8,
	//                                 "USAGE_TYPE": "CC"
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 9
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 10
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 11
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 12
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 13
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 14
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 15
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 16
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 17
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 18
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 19
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 20
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 21
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 22
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 23
	//                             }
	//                         ]
	//                     },
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_ID": "666",
	//                         "ENTITY_TYPE": "TEST",
	//                         "INTERNAL_ID": 1,
	//                         "ENTITY_KEY": "C6063D4396612FBA7324DB0739273BA1FE815C43",
	//                         "ENTITY_DESC": "JOHNSON",
	//                         "MATCH_KEY": "+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM",
	//                         "MATCH_LEVEL": 1,
	//                         "MATCH_LEVEL_CODE": "RESOLVED",
	//                         "ERRULE_CODE": "SF1_PNAME_CFF_CSTAB",
	//                         "LAST_SEEN_DT": "2022-12-06 15:58:57.407",
	//                         "FEATURES": [
	//                             {
	//                                 "LIB_FEAT_ID": 1
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 2
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 3
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 4
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 5
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 6
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 7
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 8,
	//                                 "USAGE_TYPE": "CC"
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 9
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 10
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 11
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 12
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 13
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 14
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 15
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 16
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 17
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 18
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 19
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 20
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 21
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 22
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 23
	//                             }
	//                         ]
	//                     },
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_ID": "777",
	//                         "ENTITY_TYPE": "TEST",
	//                         "INTERNAL_ID": 1,
	//                         "ENTITY_KEY": "C6063D4396612FBA7324DB0739273BA1FE815C43",
	//                         "ENTITY_DESC": "JOHNSON",
	//                         "MATCH_KEY": "+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM",
	//                         "MATCH_LEVEL": 1,
	//                         "MATCH_LEVEL_CODE": "RESOLVED",
	//                         "ERRULE_CODE": "SF1_PNAME_CFF_CSTAB",
	//                         "LAST_SEEN_DT": "2022-12-06 15:58:57.410",
	//                         "FEATURES": [
	//                             {
	//                                 "LIB_FEAT_ID": 1
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 2
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 3
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 4
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 5
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 6
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 7
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 8,
	//                                 "USAGE_TYPE": "CC"
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 9
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 10
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 11
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 12
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 13
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 14
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 15
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 16
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 17
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 18
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 19
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 20
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 21
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 22
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 23
	//                             }
	//                         ]
	//                     },
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_ID": "FCCE9793DAAD23159DBCCEB97FF2745B92CE7919",
	//                         "ENTITY_TYPE": "TEST",
	//                         "INTERNAL_ID": 1,
	//                         "ENTITY_KEY": "C6063D4396612FBA7324DB0739273BA1FE815C43",
	//                         "ENTITY_DESC": "JOHNSON",
	//                         "MATCH_KEY": "+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM",
	//                         "MATCH_LEVEL": 1,
	//                         "MATCH_LEVEL_CODE": "RESOLVED",
	//                         "ERRULE_CODE": "SF1_PNAME_CFF_CSTAB",
	//                         "LAST_SEEN_DT": "2022-12-06 15:58:57.259",
	//                         "FEATURES": [
	//                             {
	//                                 "LIB_FEAT_ID": 1
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 2
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 3
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 4
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 5
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 6
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 7
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 8,
	//                                 "USAGE_TYPE": "CC"
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 9
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 10
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 11
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 12
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 13
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 14
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 15
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 16
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 17
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 18
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 19
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 20
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 21
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 22
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 23
	//                             }
	//                         ]
	//                     }
	//                 ]
	//             },
	//             "RELATED_ENTITIES": [
	//                 {
	//                     "ENTITY_ID": 2,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0,
	//                     "ENTITY_NAME": "OCEANGUY",
	//                     "RECORD_SUMMARY": [
	//                         {
	//                             "DATA_SOURCE": "TEST",
	//                             "RECORD_COUNT": 1,
	//                             "FIRST_SEEN_DT": "2022-12-06 15:58:57.201",
	//                             "LAST_SEEN_DT": "2022-12-06 15:58:57.201"
	//                         }
	//                     ],
	//                     "LAST_SEEN_DT": "2022-12-06 15:58:57.201"
	//                 },
	//                 {
	//                     "ENTITY_ID": 3,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-DOB-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0,
	//                     "ENTITY_NAME": "Smith",
	//                     "RECORD_SUMMARY": [
	//                         {
	//                             "DATA_SOURCE": "TEST",
	//                             "RECORD_COUNT": 1,
	//                             "FIRST_SEEN_DT": "2022-12-06 15:58:57.263",
	//                             "LAST_SEEN_DT": "2022-12-06 15:58:57.263"
	//                         }
	//                     ],
	//                     "LAST_SEEN_DT": "2022-12-06 15:58:57.263"
	//                 }
	//             ]
	//         },
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 2,
	//                 "ENTITY_NAME": "OCEANGUY",
	//                 "FEATURES": {
	//                     "ACCT_NUM": [
	//                         {
	//                             "FEAT_DESC": "5534202208773608",
	//                             "LIB_FEAT_ID": 8,
	//                             "USAGE_TYPE": "CC",
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "5534202208773608",
	//                                     "LIB_FEAT_ID": 8,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "ADDRESS": [
	//                         {
	//                             "FEAT_DESC": "772 Armstrong RD Delhi WI 53543",
	//                             "LIB_FEAT_ID": 26,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "772 Armstrong RD Delhi WI 53543",
	//                                     "LIB_FEAT_ID": 26,
	//                                     "USED_FOR_CAND": "N",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "ADDR_KEY": [
	//                         {
	//                             "FEAT_DESC": "772|ARMSTRNK||53543",
	//                             "LIB_FEAT_ID": 37,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "772|ARMSTRNK||53543",
	//                                     "LIB_FEAT_ID": 37,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "772|ARMSTRNK||TL",
	//                             "LIB_FEAT_ID": 17,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "772|ARMSTRNK||TL",
	//                                     "LIB_FEAT_ID": 17,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "DOB": [
	//                         {
	//                             "FEAT_DESC": "6/9/1983",
	//                             "LIB_FEAT_ID": 25,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "6/9/1983",
	//                                     "LIB_FEAT_ID": 25,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "GENDER": [
	//                         {
	//                             "FEAT_DESC": "F",
	//                             "LIB_FEAT_ID": 3,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "F",
	//                                     "LIB_FEAT_ID": 3,
	//                                     "USED_FOR_CAND": "N",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "ID_KEY": [
	//                         {
	//                             "FEAT_DESC": "ACCT_NUM=5534202208773608",
	//                             "LIB_FEAT_ID": 19,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ACCT_NUM=5534202208773608",
	//                                     "LIB_FEAT_ID": 19,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "SSN=153-33-5185",
	//                             "LIB_FEAT_ID": 38,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "SSN=153-33-5185",
	//                                     "LIB_FEAT_ID": 38,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "LOGIN_ID": [
	//                         {
	//                             "FEAT_DESC": "flavorh2",
	//                             "LIB_FEAT_ID": 28,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "flavorh2",
	//                                     "LIB_FEAT_ID": 28,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "NAME": [
	//                         {
	//                             "FEAT_DESC": "OCEANGUY",
	//                             "LIB_FEAT_ID": 24,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "OCEANGUY",
	//                                     "LIB_FEAT_ID": 24,
	//                                     "USED_FOR_CAND": "N",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "NAME_KEY": [
	//                         {
	//                             "FEAT_DESC": "ASNK",
	//                             "LIB_FEAT_ID": 29,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK",
	//                                     "LIB_FEAT_ID": 29,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "ASNK|ADDRESS.CITY_STD=TL",
	//                             "LIB_FEAT_ID": 34,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK|ADDRESS.CITY_STD=TL",
	//                                     "LIB_FEAT_ID": 34,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "ASNK|DOB.MMDD_HASH=0906",
	//                             "LIB_FEAT_ID": 32,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK|DOB.MMDD_HASH=0906",
	//                                     "LIB_FEAT_ID": 32,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "ASNK|DOB.MMYY_HASH=0683",
	//                             "LIB_FEAT_ID": 30,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK|DOB.MMYY_HASH=0683",
	//                                     "LIB_FEAT_ID": 30,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "ASNK|DOB=80906",
	//                             "LIB_FEAT_ID": 31,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK|DOB=80906",
	//                                     "LIB_FEAT_ID": 31,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "ASNK|PHONE.PHONE_LAST_5=10796",
	//                             "LIB_FEAT_ID": 33,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK|PHONE.PHONE_LAST_5=10796",
	//                                     "LIB_FEAT_ID": 33,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "ASNK|POST=53543",
	//                             "LIB_FEAT_ID": 36,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK|POST=53543",
	//                                     "LIB_FEAT_ID": 36,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "ASNK|SSN=5185",
	//                             "LIB_FEAT_ID": 35,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK|SSN=5185",
	//                                     "LIB_FEAT_ID": 35,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "PHONE": [
	//                         {
	//                             "FEAT_DESC": "225-671-0796",
	//                             "LIB_FEAT_ID": 5,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "225-671-0796",
	//                                     "LIB_FEAT_ID": 5,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "PHONE_KEY": [
	//                         {
	//                             "FEAT_DESC": "2256710796",
	//                             "LIB_FEAT_ID": 21,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "2256710796",
	//                                     "LIB_FEAT_ID": 21,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "SEARCH_KEY": [
	//                         {
	//                             "FEAT_DESC": "LOGIN_ID:FLAVORH2|",
	//                             "LIB_FEAT_ID": 40,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "LOGIN_ID:FLAVORH2|",
	//                                     "LIB_FEAT_ID": 40,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "SSN:5185|80906|",
	//                             "LIB_FEAT_ID": 39,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "SSN:5185|80906|",
	//                                     "LIB_FEAT_ID": 39,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "SSN": [
	//                         {
	//                             "FEAT_DESC": "153-33-5185",
	//                             "LIB_FEAT_ID": 27,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "153-33-5185",
	//                                     "LIB_FEAT_ID": 27,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ]
	//                 },
	//                 "RECORD_SUMMARY": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_COUNT": 1,
	//                         "FIRST_SEEN_DT": "2022-12-06 15:58:57.201",
	//                         "LAST_SEEN_DT": "2022-12-06 15:58:57.201"
	//                     }
	//                 ],
	//                 "LAST_SEEN_DT": "2022-12-06 15:58:57.201",
	//                 "RECORDS": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_ID": "222",
	//                         "ENTITY_TYPE": "TEST",
	//                         "INTERNAL_ID": 2,
	//                         "ENTITY_KEY": "740BA22D15CA88462A930AF8A7C904FF5E48226C",
	//                         "ENTITY_DESC": "OCEANGUY",
	//                         "MATCH_KEY": "",
	//                         "MATCH_LEVEL": 0,
	//                         "MATCH_LEVEL_CODE": "",
	//                         "ERRULE_CODE": "",
	//                         "LAST_SEEN_DT": "2022-12-06 15:58:57.201",
	//                         "FEATURES": [
	//                             {
	//                                 "LIB_FEAT_ID": 3
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 5
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 8,
	//                                 "USAGE_TYPE": "CC"
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 17
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 19
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 21
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 24
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 25
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 26
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 27
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 28
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 29
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 30
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 31
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 32
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 33
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 34
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 35
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 36
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 37
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 38
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 39
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 40
	//                             }
	//                         ]
	//                     }
	//                 ]
	//             },
	//             "RELATED_ENTITIES": [
	//                 {
	//                     "ENTITY_ID": 1,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0,
	//                     "ENTITY_NAME": "JOHNSON",
	//                     "RECORD_SUMMARY": [
	//                         {
	//                             "DATA_SOURCE": "TEST",
	//                             "RECORD_COUNT": 6,
	//                             "FIRST_SEEN_DT": "2022-12-06 15:58:57.129",
	//                             "LAST_SEEN_DT": "2022-12-06 15:58:57.906"
	//                         }
	//                     ],
	//                     "LAST_SEEN_DT": "2022-12-06 15:58:57.906"
	//                 },
	//                 {
	//                     "ENTITY_ID": 3,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+ADDRESS+PHONE+ACCT_NUM-DOB-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0,
	//                     "ENTITY_NAME": "Smith",
	//                     "RECORD_SUMMARY": [
	//                         {
	//                             "DATA_SOURCE": "TEST",
	//                             "RECORD_COUNT": 1,
	//                             "FIRST_SEEN_DT": "2022-12-06 15:58:57.263",
	//                             "LAST_SEEN_DT": "2022-12-06 15:58:57.263"
	//                         }
	//                     ],
	//                     "LAST_SEEN_DT": "2022-12-06 15:58:57.263"
	//                 }
	//             ]
	//         }
	//     ]
	// }
}

func ExampleSzengine_WhyRecordInEntity() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := senzing.SzNoFlags

	result, err := szEngine.WhyRecordInEntity(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "WHY_RESULTS": [
	//         {
	//             "INTERNAL_ID": 100001,
	//             "ENTITY_ID": 100001,
	//             "FOCUS_RECORDS": [
	//                 {
	//                     "DATA_SOURCE": "CUSTOMERS",
	//                     "RECORD_ID": "1001"
	//                 }
	//             ],
	//             "MATCH_INFO": {
	//                 "WHY_KEY": "+NAME+DOB+PHONE+EMAIL",
	//                 "WHY_ERRULE_CODE": "SF1_SNAME_CFF_CSTAB",
	//                 "MATCH_LEVEL_CODE": "RESOLVED"
	//             }
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100001
	//             }
	//         }
	//     ]
	// }
}

func ExampleSzengine_WhyRecords() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	dataSourceCode1 := "CUSTOMERS"
	recordID1 := "1001"
	dataSourceCode2 := "CUSTOMERS"
	recordID2 := "1002"
	flags := senzing.SzNoFlags

	result, err := szEngine.WhyRecords(ctx, dataSourceCode1, recordID1, dataSourceCode2, recordID2, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.Truncate(result, 7))
	// Output: {"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":100001}}...
}

func ExampleSzengine_WhyRecords_output() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	// The following code pretty-prints example output JSON.
	exampleOutput := `{"WHY_RESULTS":[{"INTERNAL_ID":100001,"ENTITY_ID":1,"FOCUS_RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111"}],"INTERNAL_ID_2":2,"ENTITY_ID_2":2,"FOCUS_RECORDS_2":[{"DATA_SOURCE":"TEST","RECORD_ID":"222"}],"MATCH_INFO":{"WHY_KEY":"+PHONE+ACCT_NUM-DOB-SSN","WHY_ERRULE_CODE":"SF1","MATCH_LEVEL_CODE":"POSSIBLY_RELATED","CANDIDATE_KEYS":{"ACCT_NUM":[{"FEAT_ID":8,"FEAT_DESC":"5534202208773608"}],"ADDR_KEY":[{"FEAT_ID":17,"FEAT_DESC":"772|ARMSTRNK||TL"}],"ID_KEY":[{"FEAT_ID":19,"FEAT_DESC":"ACCT_NUM=5534202208773608"}],"PHONE":[{"FEAT_ID":5,"FEAT_DESC":"225-671-0796"}],"PHONE_KEY":[{"FEAT_ID":21,"FEAT_DESC":"2256710796"}]},"DISCLOSED_RELATIONS":{},"FEATURE_SCORES":{"ACCT_NUM":[{"INBOUND_FEAT_ID":8,"INBOUND_FEAT":"5534202208773608","INBOUND_FEAT_USAGE_TYPE":"CC","CANDIDATE_FEAT_ID":8,"CANDIDATE_FEAT":"5534202208773608","CANDIDATE_FEAT_USAGE_TYPE":"CC","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"F1"}],"ADDRESS":[{"INBOUND_FEAT_ID":4,"INBOUND_FEAT":"772 Armstrong RD Delhi LA 71232","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":26,"CANDIDATE_FEAT":"772 Armstrong RD Delhi WI 53543","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":81,"SCORE_BUCKET":"LIKELY","SCORE_BEHAVIOR":"FF"}],"DOB":[{"INBOUND_FEAT_ID":100001,"INBOUND_FEAT":"4/8/1985","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":25,"CANDIDATE_FEAT":"6/9/1983","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":79,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"FMES"}],"GENDER":[{"INBOUND_FEAT_ID":3,"INBOUND_FEAT":"F","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":3,"CANDIDATE_FEAT":"F","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FVME"}],"LOGIN_ID":[{"INBOUND_FEAT_ID":7,"INBOUND_FEAT":"flavorh","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":28,"CANDIDATE_FEAT":"flavorh2","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":0,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"F1"}],"NAME":[{"INBOUND_FEAT_ID":1,"INBOUND_FEAT":"JOHNSON","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":24,"CANDIDATE_FEAT":"OCEANGUY","CANDIDATE_FEAT_USAGE_TYPE":"","GNR_FN":33,"GNR_SN":32,"GNR_GN":70,"GENERATION_MATCH":-1,"GNR_ON":-1,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"NAME"}],"PHONE":[{"INBOUND_FEAT_ID":5,"INBOUND_FEAT":"225-671-0796","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":5,"CANDIDATE_FEAT":"225-671-0796","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":100,"SCORE_BUCKET":"SAME","SCORE_BEHAVIOR":"FF"}],"SSN":[{"INBOUND_FEAT_ID":6,"INBOUND_FEAT":"053-39-3251","INBOUND_FEAT_USAGE_TYPE":"","CANDIDATE_FEAT_ID":27,"CANDIDATE_FEAT":"153-33-5185","CANDIDATE_FEAT_USAGE_TYPE":"","FULL_SCORE":0,"SCORE_BUCKET":"NO_CHANCE","SCORE_BEHAVIOR":"F1ES"}]}}}],"ENTITIES":[{"RESOLVED_ENTITY":{"ENTITY_ID":1,"ENTITY_NAME":"JOHNSON","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi LA 71232","LIB_FEAT_ID":4,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDR_KEY":[{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||71232","LIB_FEAT_ID":18,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"DOB":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1983","LIB_FEAT_ID":2,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001,"FEAT_DESC_VALUES":[{"FEAT_DESC":"4/8/1985","LIB_FEAT_ID":100001,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ID_KEY":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=053-39-3251","LIB_FEAT_ID":20,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh","LIB_FEAT_ID":7,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JOHNSON","LIB_FEAT_ID":1,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME_KEY":[{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN","LIB_FEAT_ID":11,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":12,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMDD_HASH=0804","LIB_FEAT_ID":9,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0483","LIB_FEAT_ID":10,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0485","LIB_FEAT_ID":100002,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB.MMYY_HASH=0485","LIB_FEAT_ID":100002,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|DOB=80804","LIB_FEAT_ID":13,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":15,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|POST=71232","LIB_FEAT_ID":14,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"FEAT_DESC_VALUES":[{"FEAT_DESC":"JNSN|SSN=3251","LIB_FEAT_ID":16,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE_KEY":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"FEAT_DESC_VALUES":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SEARCH_KEY":[{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH|","LIB_FEAT_ID":22,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:3251|80804|","LIB_FEAT_ID":23,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SSN":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"FEAT_DESC_VALUES":[{"FEAT_DESC":"053-39-3251","LIB_FEAT_ID":6,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":6,"FIRST_SEEN_DT":"2022-12-06 16:13:27.135","LAST_SEEN_DT":"2022-12-06 16:13:27.916"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.916","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"111","ENTITY_TYPE":"TEST","INTERNAL_ID":100001,"ENTITY_KEY":"A6C927986DF7329D1D2CDE0E8F34328AE640FB7E","ENTITY_DESC":"JOHNSON","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 16:13:27.916","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23},{"LIB_FEAT_ID":100001},{"LIB_FEAT_ID":100002}]},{"DATA_SOURCE":"TEST","RECORD_ID":"444","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 16:13:27.405","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"555","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 16:13:27.408","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"666","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 16:13:27.411","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"777","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 16:13:27.418","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]},{"DATA_SOURCE":"TEST","RECORD_ID":"FCCE9793DAAD23159DBCCEB97FF2745B92CE7919","ENTITY_TYPE":"TEST","INTERNAL_ID":1,"ENTITY_KEY":"C6063D4396612FBA7324DB0739273BA1FE815C43","ENTITY_DESC":"JOHNSON","MATCH_KEY":"+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM","MATCH_LEVEL":1,"MATCH_LEVEL_CODE":"RESOLVED","ERRULE_CODE":"SF1_PNAME_CFF_CSTAB","LAST_SEEN_DT":"2022-12-06 16:13:27.265","FEATURES":[{"LIB_FEAT_ID":1},{"LIB_FEAT_ID":2},{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":4},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":6},{"LIB_FEAT_ID":7},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":9},{"LIB_FEAT_ID":10},{"LIB_FEAT_ID":11},{"LIB_FEAT_ID":12},{"LIB_FEAT_ID":13},{"LIB_FEAT_ID":14},{"LIB_FEAT_ID":15},{"LIB_FEAT_ID":16},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":18},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":20},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":22},{"LIB_FEAT_ID":23}]}]},"RELATED_ENTITIES":[{"ENTITY_ID":2,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"OCEANGUY","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 16:13:27.208","LAST_SEEN_DT":"2022-12-06 16:13:27.208"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.208"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 16:13:27.272","LAST_SEEN_DT":"2022-12-06 16:13:27.272"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.272"}]},{"RESOLVED_ENTITY":{"ENTITY_ID":2,"ENTITY_NAME":"OCEANGUY","FEATURES":{"ACCT_NUM":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USAGE_TYPE":"CC","FEAT_DESC_VALUES":[{"FEAT_DESC":"5534202208773608","LIB_FEAT_ID":8,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDRESS":[{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772 Armstrong RD Delhi WI 53543","LIB_FEAT_ID":26,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ADDR_KEY":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||53543","LIB_FEAT_ID":37,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"FEAT_DESC_VALUES":[{"FEAT_DESC":"772|ARMSTRNK||TL","LIB_FEAT_ID":17,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"DOB":[{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"FEAT_DESC_VALUES":[{"FEAT_DESC":"6/9/1983","LIB_FEAT_ID":25,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"GENDER":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"FEAT_DESC_VALUES":[{"FEAT_DESC":"F","LIB_FEAT_ID":3,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"ID_KEY":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ACCT_NUM=5534202208773608","LIB_FEAT_ID":19,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN=153-33-5185","LIB_FEAT_ID":38,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"LOGIN_ID":[{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"FEAT_DESC_VALUES":[{"FEAT_DESC":"flavorh2","LIB_FEAT_ID":28,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME":[{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"FEAT_DESC_VALUES":[{"FEAT_DESC":"OCEANGUY","LIB_FEAT_ID":24,"USED_FOR_CAND":"N","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"NAME_KEY":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK","LIB_FEAT_ID":29,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|ADDRESS.CITY_STD=TL","LIB_FEAT_ID":34,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMDD_HASH=0906","LIB_FEAT_ID":32,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB.MMYY_HASH=0683","LIB_FEAT_ID":30,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|DOB=80906","LIB_FEAT_ID":31,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|PHONE.PHONE_LAST_5=10796","LIB_FEAT_ID":33,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|POST=53543","LIB_FEAT_ID":36,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"FEAT_DESC_VALUES":[{"FEAT_DESC":"ASNK|SSN=5185","LIB_FEAT_ID":35,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"FEAT_DESC_VALUES":[{"FEAT_DESC":"225-671-0796","LIB_FEAT_ID":5,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"PHONE_KEY":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"FEAT_DESC_VALUES":[{"FEAT_DESC":"2256710796","LIB_FEAT_ID":21,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":3,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SEARCH_KEY":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"FEAT_DESC_VALUES":[{"FEAT_DESC":"LOGIN_ID:FLAVORH2|","LIB_FEAT_ID":40,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]},{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"FEAT_DESC_VALUES":[{"FEAT_DESC":"SSN:5185|80906|","LIB_FEAT_ID":39,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"N","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}],"SSN":[{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"FEAT_DESC_VALUES":[{"FEAT_DESC":"153-33-5185","LIB_FEAT_ID":27,"USED_FOR_CAND":"Y","USED_FOR_SCORING":"Y","ENTITY_COUNT":1,"CANDIDATE_CAP_REACHED":"N","SCORING_CAP_REACHED":"N","SUPPRESSED":"N"}]}]},"RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 16:13:27.208","LAST_SEEN_DT":"2022-12-06 16:13:27.208"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.208","RECORDS":[{"DATA_SOURCE":"TEST","RECORD_ID":"222","ENTITY_TYPE":"TEST","INTERNAL_ID":2,"ENTITY_KEY":"740BA22D15CA88462A930AF8A7C904FF5E48226C","ENTITY_DESC":"OCEANGUY","MATCH_KEY":"","MATCH_LEVEL":0,"MATCH_LEVEL_CODE":"","ERRULE_CODE":"","LAST_SEEN_DT":"2022-12-06 16:13:27.208","FEATURES":[{"LIB_FEAT_ID":3},{"LIB_FEAT_ID":5},{"LIB_FEAT_ID":8,"USAGE_TYPE":"CC"},{"LIB_FEAT_ID":17},{"LIB_FEAT_ID":19},{"LIB_FEAT_ID":21},{"LIB_FEAT_ID":24},{"LIB_FEAT_ID":25},{"LIB_FEAT_ID":26},{"LIB_FEAT_ID":27},{"LIB_FEAT_ID":28},{"LIB_FEAT_ID":29},{"LIB_FEAT_ID":30},{"LIB_FEAT_ID":31},{"LIB_FEAT_ID":32},{"LIB_FEAT_ID":33},{"LIB_FEAT_ID":34},{"LIB_FEAT_ID":35},{"LIB_FEAT_ID":36},{"LIB_FEAT_ID":37},{"LIB_FEAT_ID":38},{"LIB_FEAT_ID":39},{"LIB_FEAT_ID":40}]}]},"RELATED_ENTITIES":[{"ENTITY_ID":1,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+PHONE+ACCT_NUM-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"JOHNSON","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":6,"FIRST_SEEN_DT":"2022-12-06 16:13:27.135","LAST_SEEN_DT":"2022-12-06 16:13:27.916"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.916"},{"ENTITY_ID":3,"MATCH_LEVEL":3,"MATCH_LEVEL_CODE":"POSSIBLY_RELATED","MATCH_KEY":"+ADDRESS+PHONE+ACCT_NUM-DOB-SSN","ERRULE_CODE":"SF1","IS_DISCLOSED":0,"IS_AMBIGUOUS":0,"ENTITY_NAME":"Smith","RECORD_SUMMARY":[{"DATA_SOURCE":"TEST","RECORD_COUNT":1,"FIRST_SEEN_DT":"2022-12-06 16:13:27.272","LAST_SEEN_DT":"2022-12-06 16:13:27.272"}],"LAST_SEEN_DT":"2022-12-06 16:13:27.272"}]}]}`
	fmt.Println(jsonutil.PrettyPrint(exampleOutput, jsonIndentation))
	// Output:
	// {
	//     "WHY_RESULTS": [
	//         {
	//             "INTERNAL_ID": 100001,
	//             "ENTITY_ID": 1,
	//             "FOCUS_RECORDS": [
	//                 {
	//                     "DATA_SOURCE": "TEST",
	//                     "RECORD_ID": "111"
	//                 }
	//             ],
	//             "INTERNAL_ID_2": 2,
	//             "ENTITY_ID_2": 2,
	//             "FOCUS_RECORDS_2": [
	//                 {
	//                     "DATA_SOURCE": "TEST",
	//                     "RECORD_ID": "222"
	//                 }
	//             ],
	//             "MATCH_INFO": {
	//                 "WHY_KEY": "+PHONE+ACCT_NUM-DOB-SSN",
	//                 "WHY_ERRULE_CODE": "SF1",
	//                 "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                 "CANDIDATE_KEYS": {
	//                     "ACCT_NUM": [
	//                         {
	//                             "FEAT_ID": 8,
	//                             "FEAT_DESC": "5534202208773608"
	//                         }
	//                     ],
	//                     "ADDR_KEY": [
	//                         {
	//                             "FEAT_ID": 17,
	//                             "FEAT_DESC": "772|ARMSTRNK||TL"
	//                         }
	//                     ],
	//                     "ID_KEY": [
	//                         {
	//                             "FEAT_ID": 19,
	//                             "FEAT_DESC": "ACCT_NUM=5534202208773608"
	//                         }
	//                     ],
	//                     "PHONE": [
	//                         {
	//                             "FEAT_ID": 5,
	//                             "FEAT_DESC": "225-671-0796"
	//                         }
	//                     ],
	//                     "PHONE_KEY": [
	//                         {
	//                             "FEAT_ID": 21,
	//                             "FEAT_DESC": "2256710796"
	//                         }
	//                     ]
	//                 },
	//                 "DISCLOSED_RELATIONS": {},
	//                 "FEATURE_SCORES": {
	//                     "ACCT_NUM": [
	//                         {
	//                             "INBOUND_FEAT_ID": 8,
	//                             "INBOUND_FEAT": "5534202208773608",
	//                             "INBOUND_FEAT_USAGE_TYPE": "CC",
	//                             "CANDIDATE_FEAT_ID": 8,
	//                             "CANDIDATE_FEAT": "5534202208773608",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "CC",
	//                             "FULL_SCORE": 100,
	//                             "SCORE_BUCKET": "SAME",
	//                             "SCORE_BEHAVIOR": "F1"
	//                         }
	//                     ],
	//                     "ADDRESS": [
	//                         {
	//                             "INBOUND_FEAT_ID": 4,
	//                             "INBOUND_FEAT": "772 Armstrong RD Delhi LA 71232",
	//                             "INBOUND_FEAT_USAGE_TYPE": "",
	//                             "CANDIDATE_FEAT_ID": 26,
	//                             "CANDIDATE_FEAT": "772 Armstrong RD Delhi WI 53543",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "",
	//                             "FULL_SCORE": 81,
	//                             "SCORE_BUCKET": "LIKELY",
	//                             "SCORE_BEHAVIOR": "FF"
	//                         }
	//                     ],
	//                     "DOB": [
	//                         {
	//                             "INBOUND_FEAT_ID": 100001,
	//                             "INBOUND_FEAT": "4/8/1985",
	//                             "INBOUND_FEAT_USAGE_TYPE": "",
	//                             "CANDIDATE_FEAT_ID": 25,
	//                             "CANDIDATE_FEAT": "6/9/1983",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "",
	//                             "FULL_SCORE": 79,
	//                             "SCORE_BUCKET": "NO_CHANCE",
	//                             "SCORE_BEHAVIOR": "FMES"
	//                         }
	//                     ],
	//                     "GENDER": [
	//                         {
	//                             "INBOUND_FEAT_ID": 3,
	//                             "INBOUND_FEAT": "F",
	//                             "INBOUND_FEAT_USAGE_TYPE": "",
	//                             "CANDIDATE_FEAT_ID": 3,
	//                             "CANDIDATE_FEAT": "F",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "",
	//                             "FULL_SCORE": 100,
	//                             "SCORE_BUCKET": "SAME",
	//                             "SCORE_BEHAVIOR": "FVME"
	//                         }
	//                     ],
	//                     "LOGIN_ID": [
	//                         {
	//                             "INBOUND_FEAT_ID": 7,
	//                             "INBOUND_FEAT": "flavorh",
	//                             "INBOUND_FEAT_USAGE_TYPE": "",
	//                             "CANDIDATE_FEAT_ID": 28,
	//                             "CANDIDATE_FEAT": "flavorh2",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "",
	//                             "FULL_SCORE": 0,
	//                             "SCORE_BUCKET": "NO_CHANCE",
	//                             "SCORE_BEHAVIOR": "F1"
	//                         }
	//                     ],
	//                     "NAME": [
	//                         {
	//                             "INBOUND_FEAT_ID": 1,
	//                             "INBOUND_FEAT": "JOHNSON",
	//                             "INBOUND_FEAT_USAGE_TYPE": "",
	//                             "CANDIDATE_FEAT_ID": 24,
	//                             "CANDIDATE_FEAT": "OCEANGUY",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "",
	//                             "GNR_FN": 33,
	//                             "GNR_SN": 32,
	//                             "GNR_GN": 70,
	//                             "GENERATION_MATCH": -1,
	//                             "GNR_ON": -1,
	//                             "SCORE_BUCKET": "NO_CHANCE",
	//                             "SCORE_BEHAVIOR": "NAME"
	//                         }
	//                     ],
	//                     "PHONE": [
	//                         {
	//                             "INBOUND_FEAT_ID": 5,
	//                             "INBOUND_FEAT": "225-671-0796",
	//                             "INBOUND_FEAT_USAGE_TYPE": "",
	//                             "CANDIDATE_FEAT_ID": 5,
	//                             "CANDIDATE_FEAT": "225-671-0796",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "",
	//                             "FULL_SCORE": 100,
	//                             "SCORE_BUCKET": "SAME",
	//                             "SCORE_BEHAVIOR": "FF"
	//                         }
	//                     ],
	//                     "SSN": [
	//                         {
	//                             "INBOUND_FEAT_ID": 6,
	//                             "INBOUND_FEAT": "053-39-3251",
	//                             "INBOUND_FEAT_USAGE_TYPE": "",
	//                             "CANDIDATE_FEAT_ID": 27,
	//                             "CANDIDATE_FEAT": "153-33-5185",
	//                             "CANDIDATE_FEAT_USAGE_TYPE": "",
	//                             "FULL_SCORE": 0,
	//                             "SCORE_BUCKET": "NO_CHANCE",
	//                             "SCORE_BEHAVIOR": "F1ES"
	//                         }
	//                     ]
	//                 }
	//             }
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 1,
	//                 "ENTITY_NAME": "JOHNSON",
	//                 "FEATURES": {
	//                     "ACCT_NUM": [
	//                         {
	//                             "FEAT_DESC": "5534202208773608",
	//                             "LIB_FEAT_ID": 8,
	//                             "USAGE_TYPE": "CC",
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "5534202208773608",
	//                                     "LIB_FEAT_ID": 8,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "ADDRESS": [
	//                         {
	//                             "FEAT_DESC": "772 Armstrong RD Delhi LA 71232",
	//                             "LIB_FEAT_ID": 4,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "772 Armstrong RD Delhi LA 71232",
	//                                     "LIB_FEAT_ID": 4,
	//                                     "USED_FOR_CAND": "N",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "ADDR_KEY": [
	//                         {
	//                             "FEAT_DESC": "772|ARMSTRNK||71232",
	//                             "LIB_FEAT_ID": 18,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "772|ARMSTRNK||71232",
	//                                     "LIB_FEAT_ID": 18,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "772|ARMSTRNK||TL",
	//                             "LIB_FEAT_ID": 17,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "772|ARMSTRNK||TL",
	//                                     "LIB_FEAT_ID": 17,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "DOB": [
	//                         {
	//                             "FEAT_DESC": "4/8/1983",
	//                             "LIB_FEAT_ID": 2,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "4/8/1983",
	//                                     "LIB_FEAT_ID": 2,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "4/8/1985",
	//                             "LIB_FEAT_ID": 100001,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "4/8/1985",
	//                                     "LIB_FEAT_ID": 100001,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "GENDER": [
	//                         {
	//                             "FEAT_DESC": "F",
	//                             "LIB_FEAT_ID": 3,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "F",
	//                                     "LIB_FEAT_ID": 3,
	//                                     "USED_FOR_CAND": "N",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "ID_KEY": [
	//                         {
	//                             "FEAT_DESC": "ACCT_NUM=5534202208773608",
	//                             "LIB_FEAT_ID": 19,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ACCT_NUM=5534202208773608",
	//                                     "LIB_FEAT_ID": 19,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "SSN=053-39-3251",
	//                             "LIB_FEAT_ID": 20,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "SSN=053-39-3251",
	//                                     "LIB_FEAT_ID": 20,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "LOGIN_ID": [
	//                         {
	//                             "FEAT_DESC": "flavorh",
	//                             "LIB_FEAT_ID": 7,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "flavorh",
	//                                     "LIB_FEAT_ID": 7,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "NAME": [
	//                         {
	//                             "FEAT_DESC": "JOHNSON",
	//                             "LIB_FEAT_ID": 1,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JOHNSON",
	//                                     "LIB_FEAT_ID": 1,
	//                                     "USED_FOR_CAND": "N",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "NAME_KEY": [
	//                         {
	//                             "FEAT_DESC": "JNSN",
	//                             "LIB_FEAT_ID": 11,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN",
	//                                     "LIB_FEAT_ID": 11,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|ADDRESS.CITY_STD=TL",
	//                             "LIB_FEAT_ID": 12,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|ADDRESS.CITY_STD=TL",
	//                                     "LIB_FEAT_ID": 12,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|DOB.MMDD_HASH=0804",
	//                             "LIB_FEAT_ID": 9,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|DOB.MMDD_HASH=0804",
	//                                     "LIB_FEAT_ID": 9,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|DOB.MMYY_HASH=0483",
	//                             "LIB_FEAT_ID": 10,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|DOB.MMYY_HASH=0483",
	//                                     "LIB_FEAT_ID": 10,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|DOB.MMYY_HASH=0485",
	//                             "LIB_FEAT_ID": 100002,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|DOB.MMYY_HASH=0485",
	//                                     "LIB_FEAT_ID": 100002,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|DOB=80804",
	//                             "LIB_FEAT_ID": 13,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|DOB=80804",
	//                                     "LIB_FEAT_ID": 13,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|PHONE.PHONE_LAST_5=10796",
	//                             "LIB_FEAT_ID": 15,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|PHONE.PHONE_LAST_5=10796",
	//                                     "LIB_FEAT_ID": 15,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|POST=71232",
	//                             "LIB_FEAT_ID": 14,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|POST=71232",
	//                                     "LIB_FEAT_ID": 14,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "JNSN|SSN=3251",
	//                             "LIB_FEAT_ID": 16,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "JNSN|SSN=3251",
	//                                     "LIB_FEAT_ID": 16,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "PHONE": [
	//                         {
	//                             "FEAT_DESC": "225-671-0796",
	//                             "LIB_FEAT_ID": 5,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "225-671-0796",
	//                                     "LIB_FEAT_ID": 5,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "PHONE_KEY": [
	//                         {
	//                             "FEAT_DESC": "2256710796",
	//                             "LIB_FEAT_ID": 21,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "2256710796",
	//                                     "LIB_FEAT_ID": 21,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "SEARCH_KEY": [
	//                         {
	//                             "FEAT_DESC": "LOGIN_ID:FLAVORH|",
	//                             "LIB_FEAT_ID": 22,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "LOGIN_ID:FLAVORH|",
	//                                     "LIB_FEAT_ID": 22,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "SSN:3251|80804|",
	//                             "LIB_FEAT_ID": 23,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "SSN:3251|80804|",
	//                                     "LIB_FEAT_ID": 23,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "SSN": [
	//                         {
	//                             "FEAT_DESC": "053-39-3251",
	//                             "LIB_FEAT_ID": 6,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "053-39-3251",
	//                                     "LIB_FEAT_ID": 6,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ]
	//                 },
	//                 "RECORD_SUMMARY": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_COUNT": 6,
	//                         "FIRST_SEEN_DT": "2022-12-06 16:13:27.135",
	//                         "LAST_SEEN_DT": "2022-12-06 16:13:27.916"
	//                     }
	//                 ],
	//                 "LAST_SEEN_DT": "2022-12-06 16:13:27.916",
	//                 "RECORDS": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_ID": "111",
	//                         "ENTITY_TYPE": "TEST",
	//                         "INTERNAL_ID": 100001,
	//                         "ENTITY_KEY": "A6C927986DF7329D1D2CDE0E8F34328AE640FB7E",
	//                         "ENTITY_DESC": "JOHNSON",
	//                         "MATCH_KEY": "",
	//                         "MATCH_LEVEL": 0,
	//                         "MATCH_LEVEL_CODE": "",
	//                         "ERRULE_CODE": "",
	//                         "LAST_SEEN_DT": "2022-12-06 16:13:27.916",
	//                         "FEATURES": [
	//                             {
	//                                 "LIB_FEAT_ID": 1
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 3
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 4
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 5
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 6
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 7
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 8,
	//                                 "USAGE_TYPE": "CC"
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 9
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 11
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 12
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 13
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 14
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 15
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 16
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 17
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 18
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 19
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 20
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 21
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 22
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 23
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 100001
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 100002
	//                             }
	//                         ]
	//                     },
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_ID": "444",
	//                         "ENTITY_TYPE": "TEST",
	//                         "INTERNAL_ID": 1,
	//                         "ENTITY_KEY": "C6063D4396612FBA7324DB0739273BA1FE815C43",
	//                         "ENTITY_DESC": "JOHNSON",
	//                         "MATCH_KEY": "+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM",
	//                         "MATCH_LEVEL": 1,
	//                         "MATCH_LEVEL_CODE": "RESOLVED",
	//                         "ERRULE_CODE": "SF1_PNAME_CFF_CSTAB",
	//                         "LAST_SEEN_DT": "2022-12-06 16:13:27.405",
	//                         "FEATURES": [
	//                             {
	//                                 "LIB_FEAT_ID": 1
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 2
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 3
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 4
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 5
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 6
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 7
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 8,
	//                                 "USAGE_TYPE": "CC"
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 9
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 10
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 11
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 12
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 13
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 14
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 15
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 16
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 17
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 18
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 19
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 20
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 21
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 22
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 23
	//                             }
	//                         ]
	//                     },
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_ID": "555",
	//                         "ENTITY_TYPE": "TEST",
	//                         "INTERNAL_ID": 1,
	//                         "ENTITY_KEY": "C6063D4396612FBA7324DB0739273BA1FE815C43",
	//                         "ENTITY_DESC": "JOHNSON",
	//                         "MATCH_KEY": "+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM",
	//                         "MATCH_LEVEL": 1,
	//                         "MATCH_LEVEL_CODE": "RESOLVED",
	//                         "ERRULE_CODE": "SF1_PNAME_CFF_CSTAB",
	//                         "LAST_SEEN_DT": "2022-12-06 16:13:27.408",
	//                         "FEATURES": [
	//                             {
	//                                 "LIB_FEAT_ID": 1
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 2
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 3
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 4
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 5
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 6
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 7
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 8,
	//                                 "USAGE_TYPE": "CC"
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 9
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 10
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 11
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 12
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 13
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 14
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 15
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 16
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 17
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 18
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 19
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 20
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 21
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 22
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 23
	//                             }
	//                         ]
	//                     },
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_ID": "666",
	//                         "ENTITY_TYPE": "TEST",
	//                         "INTERNAL_ID": 1,
	//                         "ENTITY_KEY": "C6063D4396612FBA7324DB0739273BA1FE815C43",
	//                         "ENTITY_DESC": "JOHNSON",
	//                         "MATCH_KEY": "+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM",
	//                         "MATCH_LEVEL": 1,
	//                         "MATCH_LEVEL_CODE": "RESOLVED",
	//                         "ERRULE_CODE": "SF1_PNAME_CFF_CSTAB",
	//                         "LAST_SEEN_DT": "2022-12-06 16:13:27.411",
	//                         "FEATURES": [
	//                             {
	//                                 "LIB_FEAT_ID": 1
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 2
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 3
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 4
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 5
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 6
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 7
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 8,
	//                                 "USAGE_TYPE": "CC"
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 9
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 10
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 11
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 12
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 13
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 14
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 15
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 16
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 17
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 18
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 19
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 20
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 21
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 22
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 23
	//                             }
	//                         ]
	//                     },
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_ID": "777",
	//                         "ENTITY_TYPE": "TEST",
	//                         "INTERNAL_ID": 1,
	//                         "ENTITY_KEY": "C6063D4396612FBA7324DB0739273BA1FE815C43",
	//                         "ENTITY_DESC": "JOHNSON",
	//                         "MATCH_KEY": "+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM",
	//                         "MATCH_LEVEL": 1,
	//                         "MATCH_LEVEL_CODE": "RESOLVED",
	//                         "ERRULE_CODE": "SF1_PNAME_CFF_CSTAB",
	//                         "LAST_SEEN_DT": "2022-12-06 16:13:27.418",
	//                         "FEATURES": [
	//                             {
	//                                 "LIB_FEAT_ID": 1
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 2
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 3
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 4
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 5
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 6
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 7
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 8,
	//                                 "USAGE_TYPE": "CC"
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 9
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 10
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 11
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 12
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 13
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 14
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 15
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 16
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 17
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 18
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 19
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 20
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 21
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 22
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 23
	//                             }
	//                         ]
	//                     },
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_ID": "FCCE9793DAAD23159DBCCEB97FF2745B92CE7919",
	//                         "ENTITY_TYPE": "TEST",
	//                         "INTERNAL_ID": 1,
	//                         "ENTITY_KEY": "C6063D4396612FBA7324DB0739273BA1FE815C43",
	//                         "ENTITY_DESC": "JOHNSON",
	//                         "MATCH_KEY": "+NAME+ADDRESS+PHONE+SSN+LOGIN_ID+ACCT_NUM",
	//                         "MATCH_LEVEL": 1,
	//                         "MATCH_LEVEL_CODE": "RESOLVED",
	//                         "ERRULE_CODE": "SF1_PNAME_CFF_CSTAB",
	//                         "LAST_SEEN_DT": "2022-12-06 16:13:27.265",
	//                         "FEATURES": [
	//                             {
	//                                 "LIB_FEAT_ID": 1
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 2
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 3
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 4
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 5
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 6
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 7
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 8,
	//                                 "USAGE_TYPE": "CC"
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 9
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 10
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 11
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 12
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 13
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 14
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 15
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 16
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 17
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 18
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 19
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 20
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 21
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 22
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 23
	//                             }
	//                         ]
	//                     }
	//                 ]
	//             },
	//             "RELATED_ENTITIES": [
	//                 {
	//                     "ENTITY_ID": 2,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0,
	//                     "ENTITY_NAME": "OCEANGUY",
	//                     "RECORD_SUMMARY": [
	//                         {
	//                             "DATA_SOURCE": "TEST",
	//                             "RECORD_COUNT": 1,
	//                             "FIRST_SEEN_DT": "2022-12-06 16:13:27.208",
	//                             "LAST_SEEN_DT": "2022-12-06 16:13:27.208"
	//                         }
	//                     ],
	//                     "LAST_SEEN_DT": "2022-12-06 16:13:27.208"
	//                 },
	//                 {
	//                     "ENTITY_ID": 3,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-DOB-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0,
	//                     "ENTITY_NAME": "Smith",
	//                     "RECORD_SUMMARY": [
	//                         {
	//                             "DATA_SOURCE": "TEST",
	//                             "RECORD_COUNT": 1,
	//                             "FIRST_SEEN_DT": "2022-12-06 16:13:27.272",
	//                             "LAST_SEEN_DT": "2022-12-06 16:13:27.272"
	//                         }
	//                     ],
	//                     "LAST_SEEN_DT": "2022-12-06 16:13:27.272"
	//                 }
	//             ]
	//         },
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 2,
	//                 "ENTITY_NAME": "OCEANGUY",
	//                 "FEATURES": {
	//                     "ACCT_NUM": [
	//                         {
	//                             "FEAT_DESC": "5534202208773608",
	//                             "LIB_FEAT_ID": 8,
	//                             "USAGE_TYPE": "CC",
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "5534202208773608",
	//                                     "LIB_FEAT_ID": 8,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "ADDRESS": [
	//                         {
	//                             "FEAT_DESC": "772 Armstrong RD Delhi WI 53543",
	//                             "LIB_FEAT_ID": 26,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "772 Armstrong RD Delhi WI 53543",
	//                                     "LIB_FEAT_ID": 26,
	//                                     "USED_FOR_CAND": "N",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "ADDR_KEY": [
	//                         {
	//                             "FEAT_DESC": "772|ARMSTRNK||53543",
	//                             "LIB_FEAT_ID": 37,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "772|ARMSTRNK||53543",
	//                                     "LIB_FEAT_ID": 37,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "772|ARMSTRNK||TL",
	//                             "LIB_FEAT_ID": 17,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "772|ARMSTRNK||TL",
	//                                     "LIB_FEAT_ID": 17,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "DOB": [
	//                         {
	//                             "FEAT_DESC": "6/9/1983",
	//                             "LIB_FEAT_ID": 25,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "6/9/1983",
	//                                     "LIB_FEAT_ID": 25,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "GENDER": [
	//                         {
	//                             "FEAT_DESC": "F",
	//                             "LIB_FEAT_ID": 3,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "F",
	//                                     "LIB_FEAT_ID": 3,
	//                                     "USED_FOR_CAND": "N",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "ID_KEY": [
	//                         {
	//                             "FEAT_DESC": "ACCT_NUM=5534202208773608",
	//                             "LIB_FEAT_ID": 19,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ACCT_NUM=5534202208773608",
	//                                     "LIB_FEAT_ID": 19,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "SSN=153-33-5185",
	//                             "LIB_FEAT_ID": 38,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "SSN=153-33-5185",
	//                                     "LIB_FEAT_ID": 38,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "LOGIN_ID": [
	//                         {
	//                             "FEAT_DESC": "flavorh2",
	//                             "LIB_FEAT_ID": 28,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "flavorh2",
	//                                     "LIB_FEAT_ID": 28,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "NAME": [
	//                         {
	//                             "FEAT_DESC": "OCEANGUY",
	//                             "LIB_FEAT_ID": 24,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "OCEANGUY",
	//                                     "LIB_FEAT_ID": 24,
	//                                     "USED_FOR_CAND": "N",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "NAME_KEY": [
	//                         {
	//                             "FEAT_DESC": "ASNK",
	//                             "LIB_FEAT_ID": 29,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK",
	//                                     "LIB_FEAT_ID": 29,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "ASNK|ADDRESS.CITY_STD=TL",
	//                             "LIB_FEAT_ID": 34,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK|ADDRESS.CITY_STD=TL",
	//                                     "LIB_FEAT_ID": 34,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "ASNK|DOB.MMDD_HASH=0906",
	//                             "LIB_FEAT_ID": 32,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK|DOB.MMDD_HASH=0906",
	//                                     "LIB_FEAT_ID": 32,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "ASNK|DOB.MMYY_HASH=0683",
	//                             "LIB_FEAT_ID": 30,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK|DOB.MMYY_HASH=0683",
	//                                     "LIB_FEAT_ID": 30,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "ASNK|DOB=80906",
	//                             "LIB_FEAT_ID": 31,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK|DOB=80906",
	//                                     "LIB_FEAT_ID": 31,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "ASNK|PHONE.PHONE_LAST_5=10796",
	//                             "LIB_FEAT_ID": 33,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK|PHONE.PHONE_LAST_5=10796",
	//                                     "LIB_FEAT_ID": 33,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "ASNK|POST=53543",
	//                             "LIB_FEAT_ID": 36,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK|POST=53543",
	//                                     "LIB_FEAT_ID": 36,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "ASNK|SSN=5185",
	//                             "LIB_FEAT_ID": 35,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "ASNK|SSN=5185",
	//                                     "LIB_FEAT_ID": 35,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "PHONE": [
	//                         {
	//                             "FEAT_DESC": "225-671-0796",
	//                             "LIB_FEAT_ID": 5,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "225-671-0796",
	//                                     "LIB_FEAT_ID": 5,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "PHONE_KEY": [
	//                         {
	//                             "FEAT_DESC": "2256710796",
	//                             "LIB_FEAT_ID": 21,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "2256710796",
	//                                     "LIB_FEAT_ID": 21,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 3,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "SEARCH_KEY": [
	//                         {
	//                             "FEAT_DESC": "LOGIN_ID:FLAVORH2|",
	//                             "LIB_FEAT_ID": 40,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "LOGIN_ID:FLAVORH2|",
	//                                     "LIB_FEAT_ID": 40,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         },
	//                         {
	//                             "FEAT_DESC": "SSN:5185|80906|",
	//                             "LIB_FEAT_ID": 39,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "SSN:5185|80906|",
	//                                     "LIB_FEAT_ID": 39,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "N",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ],
	//                     "SSN": [
	//                         {
	//                             "FEAT_DESC": "153-33-5185",
	//                             "LIB_FEAT_ID": 27,
	//                             "FEAT_DESC_VALUES": [
	//                                 {
	//                                     "FEAT_DESC": "153-33-5185",
	//                                     "LIB_FEAT_ID": 27,
	//                                     "USED_FOR_CAND": "Y",
	//                                     "USED_FOR_SCORING": "Y",
	//                                     "ENTITY_COUNT": 1,
	//                                     "CANDIDATE_CAP_REACHED": "N",
	//                                     "SCORING_CAP_REACHED": "N",
	//                                     "SUPPRESSED": "N"
	//                                 }
	//                             ]
	//                         }
	//                     ]
	//                 },
	//                 "RECORD_SUMMARY": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_COUNT": 1,
	//                         "FIRST_SEEN_DT": "2022-12-06 16:13:27.208",
	//                         "LAST_SEEN_DT": "2022-12-06 16:13:27.208"
	//                     }
	//                 ],
	//                 "LAST_SEEN_DT": "2022-12-06 16:13:27.208",
	//                 "RECORDS": [
	//                     {
	//                         "DATA_SOURCE": "TEST",
	//                         "RECORD_ID": "222",
	//                         "ENTITY_TYPE": "TEST",
	//                         "INTERNAL_ID": 2,
	//                         "ENTITY_KEY": "740BA22D15CA88462A930AF8A7C904FF5E48226C",
	//                         "ENTITY_DESC": "OCEANGUY",
	//                         "MATCH_KEY": "",
	//                         "MATCH_LEVEL": 0,
	//                         "MATCH_LEVEL_CODE": "",
	//                         "ERRULE_CODE": "",
	//                         "LAST_SEEN_DT": "2022-12-06 16:13:27.208",
	//                         "FEATURES": [
	//                             {
	//                                 "LIB_FEAT_ID": 3
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 5
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 8,
	//                                 "USAGE_TYPE": "CC"
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 17
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 19
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 21
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 24
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 25
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 26
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 27
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 28
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 29
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 30
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 31
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 32
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 33
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 34
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 35
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 36
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 37
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 38
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 39
	//                             },
	//                             {
	//                                 "LIB_FEAT_ID": 40
	//                             }
	//                         ]
	//                     }
	//                 ]
	//             },
	//             "RELATED_ENTITIES": [
	//                 {
	//                     "ENTITY_ID": 1,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+PHONE+ACCT_NUM-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0,
	//                     "ENTITY_NAME": "JOHNSON",
	//                     "RECORD_SUMMARY": [
	//                         {
	//                             "DATA_SOURCE": "TEST",
	//                             "RECORD_COUNT": 6,
	//                             "FIRST_SEEN_DT": "2022-12-06 16:13:27.135",
	//                             "LAST_SEEN_DT": "2022-12-06 16:13:27.916"
	//                         }
	//                     ],
	//                     "LAST_SEEN_DT": "2022-12-06 16:13:27.916"
	//                 },
	//                 {
	//                     "ENTITY_ID": 3,
	//                     "MATCH_LEVEL": 3,
	//                     "MATCH_LEVEL_CODE": "POSSIBLY_RELATED",
	//                     "MATCH_KEY": "+ADDRESS+PHONE+ACCT_NUM-DOB-SSN",
	//                     "ERRULE_CODE": "SF1",
	//                     "IS_DISCLOSED": 0,
	//                     "IS_AMBIGUOUS": 0,
	//                     "ENTITY_NAME": "Smith",
	//                     "RECORD_SUMMARY": [
	//                         {
	//                             "DATA_SOURCE": "TEST",
	//                             "RECORD_COUNT": 1,
	//                             "FIRST_SEEN_DT": "2022-12-06 16:13:27.272",
	//                             "LAST_SEEN_DT": "2022-12-06 16:13:27.272"
	//                         }
	//                     ],
	//                     "LAST_SEEN_DT": "2022-12-06 16:13:27.272"
	//                 }
	//             ]
	//         }
	//     ]
	// }
}

func ExampleSzengine_WhySearch() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	attributes := `{"NAMES": [{"NAME_TYPE": "PRIMARY", "NAME_LAST": "Smith"}], "EMAIL_ADDRESS": "bsmith@work.com"}`
	searchProfile := "SEARCH"
	flags := senzing.SzNoFlags

	entityID := getEntityID(ctx, szEngine, truthset.CustomerRecords["1001"])

	result, err := szEngine.WhySearch(ctx, attributes, entityID, searchProfile, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(
		jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "WHY_RESULTS": [
	//         {
	//             "ENTITY_ID": 100001,
	//             "MATCH_INFO": {
	//                 "WHY_KEY": "+PNAME+EMAIL",
	//                 "WHY_ERRULE_CODE": "SF1",
	//                 "MATCH_LEVEL_CODE": "POSSIBLY_RELATED"
	//             }
	//         }
	//     ],
	//     "ENTITIES": [
	//         {
	//             "RESOLVED_ENTITY": {
	//                 "ENTITY_ID": 100001
	//             }
	//         }
	//     ]
	// }
}

// ----------------------------------------------------------------------------
// Interface methods - Destructive, so order needs to be maintained
// ----------------------------------------------------------------------------

func ExampleSzengine_ReevaluateEntity() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	entityID := getEntityIDForRecord(ctx, szEngine, "CUSTOMERS", "1001")

	flags := senzing.SzWithoutInfo

	result, err := szEngine.ReevaluateEntity(ctx, entityID, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(result)
	// Output:
}

func ExampleSzengine_ReevaluateEntity_withInfo() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	entityID := getEntityIDForRecord(ctx, szEngine, "CUSTOMERS", "1001")
	flags := senzing.SzWithInfo

	result, err := szEngine.ReevaluateEntity(ctx, entityID, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "AFFECTED_ENTITIES": [
	//         {
	//             "ENTITY_ID": 100001
	//         }
	//     ]
	// }
}

func ExampleSzengine_ReevaluateRecord() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := senzing.SzWithoutInfo

	result, err := szEngine.ReevaluateRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(result)
	// Output:
}

func ExampleSzengine_ReevaluateRecord_withInfo() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	dataSourceCode := "CUSTOMERS"
	recordID := "1001"
	flags := senzing.SzWithInfo

	result, err := szEngine.ReevaluateRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "DATA_SOURCE": "CUSTOMERS",
	//     "RECORD_ID": "1001",
	//     "AFFECTED_ENTITIES": [
	//         {
	//             "ENTITY_ID": 100001
	//         }
	//     ]
	// }
}

func ExampleSzengine_DeleteRecord() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	flags := senzing.SzWithoutInfo

	result, err := szEngine.DeleteRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(result)
	// Output:
}

func ExampleSzengine_DeleteRecord_withInfo() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szAbstractFactory := createSzAbstractFactory(ctx)

	defer func() { handleError(szAbstractFactory.Close(ctx)) }()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		handleError(err)
		return
	}

	defer func() { handleError(szEngine.Destroy(ctx)) }()

	dataSourceCode := "CUSTOMERS"
	recordID := "1003"
	flags := senzing.SzWithInfo

	result, err := szEngine.DeleteRecord(ctx, dataSourceCode, recordID, flags)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Println(jsonutil.PrettyPrint(result, jsonIndentation))
	// Output:
	// {
	//     "DATA_SOURCE": "CUSTOMERS",
	//     "RECORD_ID": "1003",
	//     "AFFECTED_ENTITIES": []
	// }
}

// ----------------------------------------------------------------------------
// Logging and observing
// ----------------------------------------------------------------------------

func ExampleSzengine_SetLogLevel() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)

	err := szEngine.SetLogLevel(ctx, logging.LevelInfoName)
	if err != nil {
		handleError(err)
		return
	}
	// Output:
}

func ExampleSzengine_SetObserverOrigin() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szEngine.SetObserverOrigin(ctx, origin)
	// Output:
}

func ExampleSzengine_GetObserverOrigin() {
	// For more information, visit
	// https://github.com/senzing-garage/sz-sdk-go-core/blob/main/szengine/szengine_examples_test.go
	ctx := context.TODO()
	szEngine := getSzEngine(ctx)
	origin := "Machine: nn; Task: UnitTest"
	szEngine.SetObserverOrigin(ctx, origin)
	result := szEngine.GetObserverOrigin(ctx)

	fmt.Println(result)
	// Output: Machine: nn; Task: UnitTest
}
