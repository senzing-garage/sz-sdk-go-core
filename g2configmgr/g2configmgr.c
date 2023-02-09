#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include "libg2configmgr.h"
#include "g2configmgr.h"

void *G2ConfigMgr_resizeStringBuffer(void *ptr, size_t size)
{
    // allocate new buffer
    return realloc(ptr, size);
}

struct G2ConfigMgr_addConfig_result G2ConfigMgr_addConfig_helper(const char *configStr, const char *configComments)
{
    long long configID;
    int returnCode = G2ConfigMgr_addConfig(configStr, configComments, &configID);
    struct G2ConfigMgr_addConfig_result result;
    result.configID = configID;
    result.returnCode = returnCode;
    return result;
}

struct G2ConfigMgr_getConfig_result G2ConfigMgr_getConfig_helper(const long long configID)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2ConfigMgr_resizeStringBuffer;
    int returnCode = G2ConfigMgr_getConfig(configID, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2ConfigMgr_getConfig_result result;
    result.config = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2ConfigMgr_getConfigList_result G2ConfigMgr_getConfigList_helper()
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2ConfigMgr_resizeStringBuffer;
    int returnCode = G2ConfigMgr_getConfigList(&charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2ConfigMgr_getConfigList_result result;
    result.configList = charBuffer;
    result.returnCode = returnCode;
    return result;
}

struct G2ConfigMgr_getDefaultConfigID_result G2ConfigMgr_getDefaultConfigID_helper()
{
    long long configID;
    int returnCode = G2ConfigMgr_getDefaultConfigID(&configID);
    struct G2ConfigMgr_getDefaultConfigID_result result;
    result.configID = configID;
    result.returnCode = returnCode;
    return result;
}
