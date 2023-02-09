#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>
#include "libg2.h"

typedef void *(*resize_buffer_type)(void *, size_t);

struct G2_addRecordWithInfo_result
{
    char *response;
    int returnCode;
};

struct G2_addRecordWithInfoWithReturnedRecordID_result
{
    char *recordID;
    char *withInfo;
    int returnCode;
};

struct G2_checkRecord_result
{
    char *response;
    int returnCode;
};

struct G2_deleteRecordWithInfo_result
{
    char *response;
    int returnCode;
};

struct G2_exportConfigAndConfigID_result
{
    long long configID;
    char *config;
    int returnCode;
};

struct G2_exportConfig_result
{
    char *response;
    int returnCode;
};

struct G2_exportCSVEntityReport_result
{
    void *exportHandle;
    int returnCode;
};
struct G2_exportJSONEntityReport_result
{
    void *exportHandle;
    int returnCode;
};
struct G2_fetchNext_result
{
    char *response;
    int returnCode;
};

struct G2_findInterestingEntitiesByEntityID_result
{
    char *response;
    int returnCode;
};

struct G2_findInterestingEntitiesByRecordID_result
{
    char *response;
    int returnCode;
};

struct G2_findNetworkByEntityID_result
{
    char *response;
    int returnCode;
};

struct G2_findNetworkByEntityID_V2_result
{
    char *response;
    int returnCode;
};

struct G2_findNetworkByRecordID_result
{
    char *response;
    int returnCode;
};

struct G2_findNetworkByRecordID_V2_result
{
    char *response;
    int returnCode;
};

struct G2_findPathByEntityID_result
{
    char *response;
    int returnCode;
};

struct G2_findPathByEntityID_V2_result
{
    char *response;
    int returnCode;
};

struct G2_findPathByRecordID_result
{
    char *response;
    int returnCode;
};

struct G2_findPathByRecordID_V2_result
{
    char *response;
    int returnCode;
};

struct G2_findPathExcludingByEntityID_result
{
    char *response;
    int returnCode;
};

struct G2_findPathExcludingByEntityID_V2_result
{
    char *response;
    int returnCode;
};

struct G2_findPathExcludingByRecordID_result
{
    char *response;
    int returnCode;
};

struct G2_findPathExcludingByRecordID_V2_result
{
    char *response;
    int returnCode;
};

struct G2_findPathIncludingSourceByEntityID_result
{
    char *response;
    int returnCode;
};

struct G2_findPathIncludingSourceByEntityID_V2_result
{
    char *response;
    int returnCode;
};

struct G2_findPathIncludingSourceByRecordID_result
{
    char *response;
    int returnCode;
};

struct G2_findPathIncludingSourceByRecordID_V2_result
{
    char *response;
    int returnCode;
};

struct G2_getActiveConfigID_result
{
    long long configID;
    int returnCode;
};

struct G2_getEntityByEntityID_result
{
    char *response;
    int returnCode;
};

struct G2_getEntityByEntityID_V2_result
{
    char *response;
    int returnCode;
};

struct G2_getEntityByRecordID_result
{
    char *response;
    int returnCode;
};

struct G2_getEntityByRecordID_V2_result
{
    char *response;
    int returnCode;
};

struct G2_getRecord_result
{
    char *response;
    int returnCode;
};

struct G2_getRecord_V2_result
{
    char *response;
    int returnCode;
};

struct G2_getRedoRecord_result
{
    char *response;
    int returnCode;
};

struct G2_getRepositoryLastModifiedTime_result
{
    long long time;
    int returnCode;
};

struct G2_getVirtualEntityByRecordID_result
{
    char *response;
    int returnCode;
};

struct G2_getVirtualEntityByRecordID_V2_result
{
    char *response;
    int returnCode;
};

struct G2_howEntityByEntityID_result
{
    char *response;
    int returnCode;
};

struct G2_howEntityByEntityID_V2_result
{
    char *response;
    int returnCode;
};

struct G2_processRedoRecord_result
{
    char *response;
    int returnCode;
};

struct G2_processRedoRecordWithInfo_result
{
    char *response;
    char *withInfo;
    int returnCode;
};

struct G2_processWithInfo_result
{
    char *response;
    int returnCode;
};

struct G2_processWithResponse_result
{
    char *response;
    int returnCode;
};

struct G2_processWithResponseResize_result
{
    char *response;
    int returnCode;
};

struct G2_reevaluateEntityWithInfo_result
{
    char *response;
    int returnCode;
};

struct G2_reevaluateRecordWithInfo_result
{
    char *response;
    int returnCode;
};

struct G2_replaceRecordWithInfo_result
{
    char *response;
    int returnCode;
};

struct G2_searchByAttributes_result
{
    char *response;
    int returnCode;
};

struct G2_searchByAttributes_V2_result
{
    char *response;
    int returnCode;
};

struct G2_stats_result
{
    char *response;
    int returnCode;
};

struct G2_whyEntities_result
{
    char *response;
    int returnCode;
};

struct G2_whyEntities_V2_result
{
    char *response;
    int returnCode;
};

struct G2_whyEntityByEntityID_result
{
    char *response;
    int returnCode;
};

struct G2_whyEntityByEntityID_V2_result
{
    char *response;
    int returnCode;
};

struct G2_whyEntityByRecordID_result
{
    char *response;
    int returnCode;
};

struct G2_whyEntityByRecordID_V2_result
{
    char *response;
    int returnCode;
};

struct G2_whyRecords_result
{
    char *response;
    int returnCode;
};

struct G2_whyRecords_V2_result
{
    char *response;
    int returnCode;
};

void *G2_resizeStringBuffer(void *ptr, size_t size);
struct G2_addRecordWithInfo_result G2_addRecordWithInfo_helper(const char *dataSourceCode, const char *recordID, const char *jsonData, const char *loadID, const long long flags);
struct G2_addRecordWithInfoWithReturnedRecordID_result G2_addRecordWithInfoWithReturnedRecordID_helper(const char *dataSourceCode, const char *jsonData, const char *loadID, const long long flags);
struct G2_checkRecord_result G2_checkRecord_helper(const char *record, const char *recordQueryList);
int G2_closeExport_helper(uintptr_t responseHandle);
struct G2_deleteRecordWithInfo_result G2_deleteRecordWithInfo_helper(const char *dataSourceCode, const char *recordID, const char *loadID, const long long flags);
struct G2_exportConfigAndConfigID_result G2_exportConfigAndConfigID_helper();
struct G2_exportConfig_result G2_exportConfig_helper();
struct G2_exportCSVEntityReport_result G2_exportCSVEntityReport_helper(const char *csvColumnList, const long long flags);
struct G2_exportJSONEntityReport_result G2_exportJSONEntityReport_helper(const long long flags);
struct G2_findInterestingEntitiesByEntityID_result G2_findInterestingEntitiesByEntityID_helper(long long entityID, long long flags);
struct G2_findInterestingEntitiesByRecordID_result G2_findInterestingEntitiesByRecordID_helper(const char *dataSourceCode, const char *recordID, long long flags);
struct G2_findNetworkByEntityID_result G2_findNetworkByEntityID_helper(const char *entityList, const int maxDegree, const int buildOutDegree, const int maxEntities);
struct G2_findNetworkByEntityID_V2_result G2_findNetworkByEntityID_V2_helper(const char *entityList, const int maxDegree, const int buildOutDegree, const int maxEntities, long long flags);
struct G2_findNetworkByRecordID_result G2_findNetworkByRecordID_helper(const char *recordList, const int maxDegree, const int buildOutDegree, const int maxEntities);
struct G2_findNetworkByRecordID_V2_result G2_findNetworkByRecordID_V2_helper(const char *recordList, const int maxDegree, const int buildOutDegree, const int maxEntities, const long long flags);
struct G2_findPathByEntityID_result G2_findPathByEntityID_helper(const long long entityID1, const long long entityID2, const int maxDegree);
struct G2_findPathByEntityID_V2_result G2_findPathByEntityID_V2_helper(const long long entityID1, const long long entityID2, const int maxDegree, const long long flags);
struct G2_findPathByRecordID_result G2_findPathByRecordID_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2, const int maxDegree);
struct G2_findPathByRecordID_V2_result G2_findPathByRecordID_V2_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2, const int maxDegree, const long long flags);
struct G2_findPathExcludingByEntityID_result G2_findPathExcludingByEntityID_helper(const long long entityID1, const long long entityID2, const int maxDegree, const char *excludedEntities);
struct G2_findPathExcludingByEntityID_V2_result G2_findPathExcludingByEntityID_V2_helper(const long long entityID1, const long long entityID2, const int maxDegree, const char *excludedEntities, const long long flags);
struct G2_findPathExcludingByRecordID_result G2_findPathExcludingByRecordID_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2, const int maxDegree, const char *excludedRecords);
struct G2_findPathExcludingByRecordID_V2_result G2_findPathExcludingByRecordID_V2_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2, const int maxDegree, const char *excludedRecords, const long long flags);
struct G2_findPathIncludingSourceByEntityID_result G2_findPathIncludingSourceByEntityID_helper(const long long entityID1, const long long entityID2, const int maxDegree, const char *excludedEntities, const char *requiredDsrcs);
struct G2_findPathIncludingSourceByEntityID_V2_result G2_findPathIncludingSourceByEntityID_V2_helper(const long long entityID1, const long long entityID2, const int maxDegree, const char *excludedEntities, const char *requiredDsrcs, const long long flags);
struct G2_findPathIncludingSourceByRecordID_result G2_findPathIncludingSourceByRecordID_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2, const int maxDegree, const char *excludedRecords, const char *requiredDsrcs);
struct G2_findPathIncludingSourceByRecordID_V2_result G2_findPathIncludingSourceByRecordID_V2_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2, const int maxDegree, const char *excludedRecords, const char *requiredDsrcs, const long long flags);
struct G2_fetchNext_result G2_fetchNext_helper(uintptr_t exportHandle);
struct G2_getActiveConfigID_result G2_getActiveConfigID_helper();
struct G2_getEntityByEntityID_result G2_getEntityByEntityID_helper(const long long entityID);
struct G2_getEntityByEntityID_V2_result G2_getEntityByEntityID_V2_helper(const long long entityID, const long long flags);
struct G2_getEntityByRecordID_result G2_getEntityByRecordID_helper(const char *dataSourceCode, const char *recordID);
struct G2_getEntityByRecordID_V2_result G2_getEntityByRecordID_V2_helper(const char *dataSourceCode, const char *recordID, const long long flags);
struct G2_getRecord_result G2_getRecord_helper(const char *dataSourceCode, const char *recordID);
struct G2_getRecord_V2_result G2_getRecord_V2_helper(const char *dataSourceCode, const char *recordID, const long long flags);
struct G2_getRedoRecord_result G2_getRedoRecord_helper();
struct G2_getRepositoryLastModifiedTime_result G2_getRepositoryLastModifiedTime_helper();
struct G2_getVirtualEntityByRecordID_result G2_getVirtualEntityByRecordID_helper(const char *recordList);
struct G2_getVirtualEntityByRecordID_V2_result G2_getVirtualEntityByRecordID_V2_helper(const char *recordList, const long long flags);
struct G2_howEntityByEntityID_result G2_howEntityByEntityID_helper(const long long entityID);
struct G2_howEntityByEntityID_V2_result G2_howEntityByEntityID_V2_helper(const long long entityID, const long long flags);
struct G2_processRedoRecord_result G2_processRedoRecord_helper();
struct G2_processRedoRecordWithInfo_result G2_processRedoRecordWithInfo_helper(const long long flags);
struct G2_processWithInfo_result G2_processWithInfo_helper(const char *record, const long long flags);
struct G2_processWithResponse_result G2_processWithResponse_helper(const char *record);
struct G2_processWithResponseResize_result G2_processWithResponseResize_helper(const char *record);
struct G2_reevaluateEntityWithInfo_result G2_reevaluateEntityWithInfo_helper(const long long entityID, const long long flags);
struct G2_reevaluateRecordWithInfo_result G2_reevaluateRecordWithInfo_helper(const char *dataSourceCode, const char *recordID, const long long flags);
struct G2_replaceRecordWithInfo_result G2_replaceRecordWithInfo_helper(const char *dataSourceCode, const char *recordID, const char *jsonData, const char *loadID, const long long flags);
struct G2_searchByAttributes_result G2_searchByAttributes_helper(const char *jsonData);
struct G2_searchByAttributes_V2_result G2_searchByAttributes_V2_helper(const char *jsonData, const long long flags);
struct G2_stats_result G2_stats_helper();
struct G2_whyEntities_result G2_whyEntities_helper(const long long entityID1, const long long entityID2);
struct G2_whyEntities_V2_result G2_whyEntities_V2_helper(const long long entityID1, const long long entityID2, const long long flags);
struct G2_whyEntityByEntityID_result G2_whyEntityByEntityID_helper(const long long entityID1);
struct G2_whyEntityByEntityID_V2_result G2_whyEntityByEntityID_V2_helper(const long long entityID1, const long long flags);
struct G2_whyEntityByRecordID_result G2_whyEntityByRecordID_helper(const char *dataSourceCode, const char *recordID);
struct G2_whyEntityByRecordID_V2_result G2_whyEntityByRecordID_V2_helper(const char *dataSourceCode, const char *recordID, const long long flags);
struct G2_whyRecords_result G2_whyRecords_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2);
struct G2_whyRecords_V2_result G2_whyRecords_V2_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2, const long long flags);
