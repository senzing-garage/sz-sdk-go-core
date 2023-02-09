#include <stdlib.h>
#include <stdio.h>
#include "libg2diagnostic.h"
#include "g2diagnostic.h"

void *G2Diagnostic_resizeStringBuffer(void *ptr, size_t size)
{
    // allocate new buffer
    return realloc(ptr, size);
}

struct G2Diagnostic_checkDBPerf_result G2Diagnostic_checkDBPerf_helper(int secondsToRun)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(1);
    resize_buffer_type resizeFuncPointer = &G2Diagnostic_resizeStringBuffer;
    int returnCode = G2Diagnostic_checkDBPerf(secondsToRun, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2Diagnostic_checkDBPerf_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

int G2Diagnostic_closeEntityListBySize_helper(uintptr_t entityListBySizeHandle)
{
    int returnCode = G2Diagnostic_closeEntityListBySize((void *)entityListBySizeHandle);
    return returnCode;
}

int G2Diagnostic_fetchNextEntityBySize_helper(uintptr_t entityListBySizeHandle, char *responseBuf, const size_t bufSize)
{
    int returnCode = G2Diagnostic_fetchNextEntityBySize((void *)entityListBySizeHandle, responseBuf, bufSize);
    return returnCode;
}

struct G2Diagnostic_findEntitiesByFeatureIDs_result G2Diagnostic_findEntitiesByFeatureIDs_helper(const char *features)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(1);
    resize_buffer_type resizeFuncPointer = &G2Diagnostic_resizeStringBuffer;
    int returnCode = G2Diagnostic_findEntitiesByFeatureIDs(features, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2Diagnostic_findEntitiesByFeatureIDs_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2Diagnostic_getDataSourceCounts_result G2Diagnostic_getDataSourceCounts_helper()
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(1);
    resize_buffer_type resizeFuncPointer = &G2Diagnostic_resizeStringBuffer;
    int returnCode = G2Diagnostic_getDataSourceCounts(&charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2Diagnostic_getDataSourceCounts_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2Diagnostic_getDBInfo_result G2Diagnostic_getDBInfo_helper()
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(1);
    resize_buffer_type resizeFuncPointer = &G2Diagnostic_resizeStringBuffer;
    int returnCode = G2Diagnostic_getDBInfo(&charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2Diagnostic_getDBInfo_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2Diagnostic_getEntityDetails_result G2Diagnostic_getEntityDetails_helper(const long long entityID, const int includeInternalFeatures)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(1);
    resize_buffer_type resizeFuncPointer = &G2Diagnostic_resizeStringBuffer;
    int returnCode = G2Diagnostic_getEntityDetails(entityID, includeInternalFeatures, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2Diagnostic_getEntityDetails_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2Diagnostic_getEntityListBySize_result G2Diagnostic_getEntityListBySize_helper(const size_t entitySize)
{
    EntityListBySizeHandle handle;
    int returnCode = G2Diagnostic_getEntityListBySize(entitySize, &handle);
    struct G2Diagnostic_getEntityListBySize_result result;
    result.response = handle;
    result.returnCode = returnCode;
    return result;
}

struct G2Diagnostic_getEntityResume_result G2Diagnostic_getEntityResume_helper(const long long entityID)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(1);
    resize_buffer_type resizeFuncPointer = &G2Diagnostic_resizeStringBuffer;
    int returnCode = G2Diagnostic_getEntityResume(entityID, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2Diagnostic_getEntityResume_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2Diagnostic_getEntitySizeBreakdown_result G2Diagnostic_getEntitySizeBreakdown_helper(const size_t minimumEntitySize, const int includeInternalFeatures)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(1);
    resize_buffer_type resizeFuncPointer = &G2Diagnostic_resizeStringBuffer;
    int returnCode = G2Diagnostic_getEntitySizeBreakdown(minimumEntitySize, includeInternalFeatures, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2Diagnostic_getEntitySizeBreakdown_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2Diagnostic_getFeature_result G2Diagnostic_getFeature_helper(const long long libFeatID)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(1);
    resize_buffer_type resizeFuncPointer = &G2Diagnostic_resizeStringBuffer;
    int returnCode = G2Diagnostic_getFeature(libFeatID, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2Diagnostic_getFeature_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2Diagnostic_getGenericFeatures_result G2Diagnostic_getGenericFeatures_helper(const char *featureType, const size_t maximumEstimatedCount)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(1);
    resize_buffer_type resizeFuncPointer = &G2Diagnostic_resizeStringBuffer;
    int returnCode = G2Diagnostic_getGenericFeatures(featureType, maximumEstimatedCount, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2Diagnostic_getGenericFeatures_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2Diagnostic_getMappingStatistics_result G2Diagnostic_getMappingStatistics_helper(const int includeInternalFeatures)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(1);
    resize_buffer_type resizeFuncPointer = &G2Diagnostic_resizeStringBuffer;
    int returnCode = G2Diagnostic_getMappingStatistics(includeInternalFeatures, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2Diagnostic_getMappingStatistics_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2Diagnostic_getRelationshipDetails_result G2Diagnostic_getRelationshipDetails_helper(const long long relationshipID, const int includeInternalFeatures)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(1);
    resize_buffer_type resizeFuncPointer = &G2Diagnostic_resizeStringBuffer;
    int returnCode = G2Diagnostic_getRelationshipDetails(relationshipID, includeInternalFeatures, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2Diagnostic_getRelationshipDetails_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2Diagnostic_getResolutionStatistics_result G2Diagnostic_getResolutionStatistics_helper()
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(1);
    resize_buffer_type resizeFuncPointer = &G2Diagnostic_resizeStringBuffer;
    int returnCode = G2Diagnostic_getResolutionStatistics(&charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2Diagnostic_getResolutionStatistics_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}
