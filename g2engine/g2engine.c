#include <stdlib.h>
#include <stdio.h>
#include "libg2.h"
#include "g2engine.h"

void *G2_resizeStringBuffer(void *ptr, size_t size)
{
    // allocate new buffer
    return realloc(ptr, size);
}

struct G2_addRecordWithInfo_result G2_addRecordWithInfo_helper(const char *dataSourceCode, const char *recordID, const char *jsonData, const char *loadID, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_addRecordWithInfo(dataSourceCode, recordID, jsonData, loadID, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_addRecordWithInfo_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_addRecordWithInfoWithReturnedRecordID_result G2_addRecordWithInfoWithReturnedRecordID_helper(const char *dataSourceCode, const char *jsonData, const char *loadID, const long long flags)
{
    size_t charBufferSize = 1;
    size_t recordIDBufSize = 256;
    char recordIDBuf[recordIDBufSize];
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_addRecordWithInfoWithReturnedRecordID(dataSourceCode, jsonData, loadID, flags, recordIDBuf, recordIDBufSize, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_addRecordWithInfoWithReturnedRecordID_result result;
    result.recordID = recordIDBuf;
    result.withInfo = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_checkRecord_result G2_checkRecord_helper(const char *record, const char *recordQueryList)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_checkRecord(record, recordQueryList, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_checkRecord_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_deleteRecordWithInfo_result G2_deleteRecordWithInfo_helper(const char *dataSourceCode, const char *recordID, const char *loadID, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_deleteRecordWithInfo(dataSourceCode, recordID, loadID, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_deleteRecordWithInfo_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

int G2_closeExport_helper(uintptr_t responseHandle)
{
    int returnCode = G2_closeExport((void *)responseHandle);
    return returnCode;
}

struct G2_exportConfigAndConfigID_result G2_exportConfigAndConfigID_helper()
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    long long configID;
    int returnCode = G2_exportConfigAndConfigID(&charBuffer, &charBufferSize, resizeFuncPointer, &configID);
    struct G2_exportConfigAndConfigID_result result;
    result.configID = configID;
    result.config = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_exportConfig_result G2_exportConfig_helper()
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_exportConfig(&charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_exportConfig_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_exportCSVEntityReport_result G2_exportCSVEntityReport_helper(const char *csvColumnList, const long long flags)
{
    ExportHandle exportHandle;
    int returnCode = G2_exportCSVEntityReport(csvColumnList, flags, &exportHandle);
    struct G2_exportCSVEntityReport_result result;
    result.exportHandle = exportHandle;
    result.returnCode = returnCode;
    return result;
}

struct G2_exportJSONEntityReport_result G2_exportJSONEntityReport_helper(const long long flags)
{
    ExportHandle exportHandle;
    int returnCode = G2_exportJSONEntityReport(flags, &exportHandle);
    struct G2_exportJSONEntityReport_result result;
    result.exportHandle = exportHandle;
    result.returnCode = returnCode;
    return result;
}

struct G2_findInterestingEntitiesByEntityID_result G2_findInterestingEntitiesByEntityID_helper(long long entityID, long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findInterestingEntitiesByEntityID(entityID, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findInterestingEntitiesByEntityID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findInterestingEntitiesByRecordID_result G2_findInterestingEntitiesByRecordID_helper(const char *dataSourceCode, const char *recordID, long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findInterestingEntitiesByRecordID(dataSourceCode, recordID, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findInterestingEntitiesByRecordID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findNetworkByEntityID_result G2_findNetworkByEntityID_helper(const char *entityList, const int maxDegree, const int buildOutDegree, const int maxEntities)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findNetworkByEntityID(entityList, maxDegree, buildOutDegree, maxEntities, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findNetworkByEntityID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findNetworkByEntityID_V2_result G2_findNetworkByEntityID_V2_helper(const char *entityList, const int maxDegree, const int buildOutDegree, const int maxEntities, long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findNetworkByEntityID_V2(entityList, maxDegree, buildOutDegree, maxEntities, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findNetworkByEntityID_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findNetworkByRecordID_result G2_findNetworkByRecordID_helper(const char *recordList, const int maxDegree, const int buildOutDegree, const int maxEntities)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findNetworkByRecordID(recordList, maxDegree, buildOutDegree, maxEntities, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findNetworkByRecordID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findNetworkByRecordID_V2_result G2_findNetworkByRecordID_V2_helper(const char *recordList, const int maxDegree, const int buildOutDegree, const int maxEntities, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findNetworkByRecordID_V2(recordList, maxDegree, buildOutDegree, maxEntities, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findNetworkByRecordID_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findPathByEntityID_result G2_findPathByEntityID_helper(const long long entityID1, const long long entityID2, const int maxDegree)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findPathByEntityID(entityID1, entityID2, maxDegree, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findPathByEntityID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findPathByEntityID_V2_result G2_findPathByEntityID_V2_helper(const long long entityID1, const long long entityID2, const int maxDegree, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findPathByEntityID_V2(entityID1, entityID2, maxDegree, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findPathByEntityID_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findPathByRecordID_result G2_findPathByRecordID_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2, const int maxDegree)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findPathByRecordID(dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findPathByRecordID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findPathByRecordID_V2_result G2_findPathByRecordID_V2_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2, const int maxDegree, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findPathByRecordID_V2(dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findPathByRecordID_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findPathExcludingByEntityID_result G2_findPathExcludingByEntityID_helper(const long long entityID1, const long long entityID2, const int maxDegree, const char *excludedEntities)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findPathExcludingByEntityID(entityID1, entityID2, maxDegree, excludedEntities, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findPathExcludingByEntityID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findPathExcludingByEntityID_V2_result G2_findPathExcludingByEntityID_V2_helper(const long long entityID1, const long long entityID2, const int maxDegree, const char *excludedEntities, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findPathExcludingByEntityID_V2(entityID1, entityID2, maxDegree, excludedEntities, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findPathExcludingByEntityID_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findPathExcludingByRecordID_result G2_findPathExcludingByRecordID_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2, const int maxDegree, const char *excludedRecords)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findPathExcludingByRecordID(dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedRecords, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findPathExcludingByRecordID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findPathExcludingByRecordID_V2_result G2_findPathExcludingByRecordID_V2_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2, const int maxDegree, const char *excludedRecords, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findPathExcludingByRecordID_V2(dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedRecords, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findPathExcludingByRecordID_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findPathIncludingSourceByEntityID_result G2_findPathIncludingSourceByEntityID_helper(const long long entityID1, const long long entityID2, const int maxDegree, const char *excludedEntities, const char *requiredDsrcs)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findPathIncludingSourceByEntityID(entityID1, entityID2, maxDegree, excludedEntities, requiredDsrcs, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findPathIncludingSourceByEntityID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findPathIncludingSourceByEntityID_V2_result G2_findPathIncludingSourceByEntityID_V2_helper(const long long entityID1, const long long entityID2, const int maxDegree, const char *excludedEntities, const char *requiredDsrcs, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findPathIncludingSourceByEntityID_V2(entityID1, entityID2, maxDegree, excludedEntities, requiredDsrcs, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findPathIncludingSourceByEntityID_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findPathIncludingSourceByRecordID_result G2_findPathIncludingSourceByRecordID_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2, const int maxDegree, const char *excludedRecords, const char *requiredDsrcs)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findPathIncludingSourceByRecordID(dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedRecords, requiredDsrcs, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findPathIncludingSourceByRecordID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_findPathIncludingSourceByRecordID_V2_result G2_findPathIncludingSourceByRecordID_V2_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2, const int maxDegree, const char *excludedRecords, const char *requiredDsrcs, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_findPathIncludingSourceByRecordID_V2(dataSourceCode1, recordID1, dataSourceCode2, recordID2, maxDegree, excludedRecords, requiredDsrcs, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_findPathIncludingSourceByRecordID_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_fetchNext_result G2_fetchNext_helper(uintptr_t exportHandle)
{
    size_t charBufferSize = 65535;
    char *charBuffer = (char *)malloc(charBufferSize);
    int returnCode = G2_fetchNext((void *)exportHandle, charBuffer, charBufferSize);
    struct G2_fetchNext_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_getActiveConfigID_result G2_getActiveConfigID_helper()
{
    long long configID;
    int returnCode = G2_getActiveConfigID(&configID);
    struct G2_getActiveConfigID_result result;
    result.configID = configID;
    result.returnCode = returnCode;
    return result;
}

struct G2_getEntityByEntityID_result G2_getEntityByEntityID_helper(const long long entityID)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_getEntityByEntityID(entityID, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_getEntityByEntityID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_getEntityByEntityID_V2_result G2_getEntityByEntityID_V2_helper(const long long entityID, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_getEntityByEntityID_V2(entityID, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_getEntityByEntityID_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_getEntityByRecordID_result G2_getEntityByRecordID_helper(const char *dataSourceCode, const char *recordID)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_getEntityByRecordID(dataSourceCode, recordID, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_getEntityByRecordID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_getEntityByRecordID_V2_result G2_getEntityByRecordID_V2_helper(const char *dataSourceCode, const char *recordID, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_getEntityByRecordID_V2(dataSourceCode, recordID, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_getEntityByRecordID_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_getRecord_result G2_getRecord_helper(const char *dataSourceCode, const char *recordID)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_getRecord(dataSourceCode, recordID, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_getRecord_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_getRecord_V2_result G2_getRecord_V2_helper(const char *dataSourceCode, const char *recordID, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_getRecord_V2(dataSourceCode, recordID, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_getRecord_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_getRedoRecord_result G2_getRedoRecord_helper()
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_getRedoRecord(&charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_getRedoRecord_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_getRepositoryLastModifiedTime_result G2_getRepositoryLastModifiedTime_helper()
{
    long long repositoryLastModifiedTime;
    int returnCode = G2_getActiveConfigID(&repositoryLastModifiedTime);
    struct G2_getRepositoryLastModifiedTime_result result;
    result.time = repositoryLastModifiedTime;
    result.returnCode = returnCode;
    return result;
}

struct G2_getVirtualEntityByRecordID_result G2_getVirtualEntityByRecordID_helper(const char *recordList)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_getVirtualEntityByRecordID(recordList, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_getVirtualEntityByRecordID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_getVirtualEntityByRecordID_V2_result G2_getVirtualEntityByRecordID_V2_helper(const char *recordList, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_getVirtualEntityByRecordID_V2(recordList, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_getVirtualEntityByRecordID_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_howEntityByEntityID_result G2_howEntityByEntityID_helper(const long long entityID)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_howEntityByEntityID(entityID, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_howEntityByEntityID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_howEntityByEntityID_V2_result G2_howEntityByEntityID_V2_helper(const long long entityID, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_howEntityByEntityID_V2(entityID, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_howEntityByEntityID_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_processRedoRecord_result G2_processRedoRecord_helper()
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_processRedoRecord(&charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_processRedoRecord_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_processRedoRecordWithInfo_result G2_processRedoRecordWithInfo_helper(const long long flags)
{
    size_t withInfoBufferSize = 1;
    size_t responseBufferSize = 1;
    char *responseBuffer = (char *)malloc(responseBufferSize);
    char *withInfoBuffer = (char *)malloc(withInfoBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_processRedoRecordWithInfo(flags, &responseBuffer, &responseBufferSize, &withInfoBuffer, &withInfoBufferSize, resizeFuncPointer);
    struct G2_processRedoRecordWithInfo_result result;
    result.response = responseBuffer;
    result.withInfo = withInfoBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_processWithInfo_result G2_processWithInfo_helper(const char *record, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_processWithInfo(record, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_processWithInfo_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_processWithResponse_result G2_processWithResponse_helper(const char *record)
{
    size_t charBufferSize = 65535;
    char *charBuffer = (char *)malloc(charBufferSize);
    int returnCode = G2_processWithResponse(record, charBuffer, charBufferSize);
    struct G2_processWithResponse_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_processWithResponseResize_result G2_processWithResponseResize_helper(const char *record)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_processWithResponseResize(record, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_processWithResponseResize_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_reevaluateEntityWithInfo_result G2_reevaluateEntityWithInfo_helper(const long long entityID, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_reevaluateEntityWithInfo(entityID, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_reevaluateEntityWithInfo_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_reevaluateRecordWithInfo_result G2_reevaluateRecordWithInfo_helper(const char *dataSourceCode, const char *recordID, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_reevaluateRecordWithInfo(dataSourceCode, recordID, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_reevaluateRecordWithInfo_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_replaceRecordWithInfo_result G2_replaceRecordWithInfo_helper(const char *dataSourceCode, const char *recordID, const char *jsonData, const char *loadID, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_replaceRecordWithInfo(dataSourceCode, recordID, jsonData, loadID, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_replaceRecordWithInfo_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_searchByAttributes_result G2_searchByAttributes_helper(const char *jsonData)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_searchByAttributes(jsonData, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_searchByAttributes_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_searchByAttributes_V2_result G2_searchByAttributes_V2_helper(const char *jsonData, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_searchByAttributes_V2(jsonData, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_searchByAttributes_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_stats_result G2_stats_helper()
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_stats(&charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_stats_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_whyEntities_result G2_whyEntities_helper(const long long entityID1, const long long entityID2)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_whyEntities(entityID1, entityID2, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_whyEntities_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_whyEntities_V2_result G2_whyEntities_V2_helper(const long long entityID1, const long long entityID2, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_whyEntities_V2(entityID1, entityID2, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_whyEntities_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_whyEntityByEntityID_result G2_whyEntityByEntityID_helper(const long long entityID1)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_whyEntityByEntityID(entityID1, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_whyEntityByEntityID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_whyEntityByEntityID_V2_result G2_whyEntityByEntityID_V2_helper(const long long entityID1, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_whyEntityByEntityID_V2(entityID1, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_whyEntityByEntityID_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_whyEntityByRecordID_result G2_whyEntityByRecordID_helper(const char *dataSourceCode, const char *recordID)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_whyEntityByRecordID(dataSourceCode, recordID, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_whyEntityByRecordID_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_whyEntityByRecordID_V2_result G2_whyEntityByRecordID_V2_helper(const char *dataSourceCode, const char *recordID, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_whyEntityByRecordID_V2(dataSourceCode, recordID, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_whyEntityByRecordID_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_whyRecords_result G2_whyRecords_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_whyRecords(dataSourceCode1, recordID1, dataSourceCode2, recordID2, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_whyRecords_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2_whyRecords_V2_result G2_whyRecords_V2_helper(const char *dataSourceCode1, const char *recordID1, const char *dataSourceCode2, const char *recordID2, const long long flags)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2_resizeStringBuffer;
    int returnCode = G2_whyRecords_V2(dataSourceCode1, recordID1, dataSourceCode2, recordID2, flags, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2_whyRecords_V2_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}
