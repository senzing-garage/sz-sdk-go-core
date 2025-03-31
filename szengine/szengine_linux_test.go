//go:build linux

package szengine_test

var expectedExportCsvEntityReport = []string{
	`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
	`3,0,"","","CUSTOMERS","1001"`,
	`3,0,"RESOLVED","+NAME+DOB+PHONE","CUSTOMERS","1002"`,
	`3,0,"RESOLVED","+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
}

var expectedExportCsvEntityReportIterator = []string{
	`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
	`13,0,"","","CUSTOMERS","1001"`,
	`13,0,"RESOLVED","+NAME+DOB+PHONE","CUSTOMERS","1002"`,
	`13,0,"RESOLVED","+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
}

var expectedExportCsvEntityReportIteratorNilCsvColumnList = []string{
	`RESOLVED_ENTITY_ID,RELATED_ENTITY_ID,MATCH_LEVEL_CODE,MATCH_KEY,DATA_SOURCE,RECORD_ID`,
	`19,0,"","","CUSTOMERS","1001"`,
	`19,0,"RESOLVED","+NAME+DOB+PHONE","CUSTOMERS","1002"`,
	`19,0,"RESOLVED","+NAME+DOB+EMAIL","CUSTOMERS","1003"`,
}
