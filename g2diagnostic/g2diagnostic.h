#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include "libg2diagnostic.h"

// typedef void* EntityListBySizeHandle;
typedef void *(*resize_buffer_type)(void *, size_t);

struct G2Diagnostic_checkDBPerf_result
{
    char *response;
    int returnCode;
};

struct G2Diagnostic_findEntitiesByFeatureIDs_result
{
    char *response;
    int returnCode;
};

struct G2Diagnostic_getDataSourceCounts_result
{
    char *response;
    int returnCode;
};

struct G2Diagnostic_getDBInfo_result
{
    char *response;
    int returnCode;
};

struct G2Diagnostic_getEntityDetails_result
{
    char *response;
    int returnCode;
};

struct G2Diagnostic_getEntityListBySize_result
{
    void *response;
    int returnCode;
};

struct G2Diagnostic_getEntityResume_result
{
    char *response;
    int returnCode;
};

struct G2Diagnostic_getEntitySizeBreakdown_result
{
    char *response;
    int returnCode;
};

struct G2Diagnostic_getFeature_result
{
    char *response;
    int returnCode;
};

struct G2Diagnostic_getGenericFeatures_result
{
    char *response;
    int returnCode;
};

struct G2Diagnostic_getMappingStatistics_result
{
    char *response;
    int returnCode;
};

struct G2Diagnostic_getRelationshipDetails_result
{
    char *response;
    int returnCode;
};

struct G2Diagnostic_getResolutionStatistics_result
{
    char *response;
    int returnCode;
};

void *G2Diagnostic_resizeStringBuffer(void *ptr, size_t size);
struct G2Diagnostic_checkDBPerf_result G2Diagnostic_checkDBPerf_helper(int secondsToRun);
int G2Diagnostic_closeEntityListBySize_helper(uintptr_t entityListBySizeHandle);
int G2Diagnostic_fetchNextEntityBySize_helper(uintptr_t configHandle, char *responseBuf, const size_t bufSize);
struct G2Diagnostic_findEntitiesByFeatureIDs_result G2Diagnostic_findEntitiesByFeatureIDs_helper(const char *features);
struct G2Diagnostic_getDataSourceCounts_result G2Diagnostic_getDataSourceCounts_helper();
struct G2Diagnostic_getDBInfo_result G2Diagnostic_getDBInfo_helper();
struct G2Diagnostic_getEntityDetails_result G2Diagnostic_getEntityDetails_helper(const long long entityID, const int includeInternalFeatures);
struct G2Diagnostic_getEntityListBySize_result G2Diagnostic_getEntityListBySize_helper(const size_t entitySize);
struct G2Diagnostic_getEntityResume_result G2Diagnostic_getEntityResume_helper(const long long entityID);
struct G2Diagnostic_getEntitySizeBreakdown_result G2Diagnostic_getEntitySizeBreakdown_helper(const size_t minimumEntitySize, const int includeInternalFeatures);
struct G2Diagnostic_getFeature_result G2Diagnostic_getFeature_helper(const long long libFeatID);
struct G2Diagnostic_getGenericFeatures_result G2Diagnostic_getGenericFeatures_helper(const char *featureType, const size_t maximumEstimatedCount);
struct G2Diagnostic_getMappingStatistics_result G2Diagnostic_getMappingStatistics_helper(const int includeInternalFeatures);
struct G2Diagnostic_getRelationshipDetails_result G2Diagnostic_getRelationshipDetails_helper(const long long relationshipID, const int includeInternalFeatures);
struct G2Diagnostic_getResolutionStatistics_result G2Diagnostic_getResolutionStatistics_helper();
