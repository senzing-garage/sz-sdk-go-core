#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include "libg2config.h"
#include "g2config.h"

void *G2config_resizeStringBuffer(void *ptr, size_t size)
{
    // allocate new buffer
    return realloc(ptr, size);
}

struct G2Config_addDataSource_result G2Config_addDataSource_helper(uintptr_t configHandle, const char *inputJson)
{
    size_t charBufferSize = 0;
    char *charBuffer = NULL;
    char **charBufferPtr = &charBuffer;
    resize_buffer_type resizeFuncPointer = &G2config_resizeStringBuffer;
    int returnCode = G2Config_addDataSource((void *)configHandle, inputJson, charBufferPtr, &charBufferSize, resizeFuncPointer);
    struct G2Config_addDataSource_result result;
    result.response = *charBufferPtr;
    result.returnCode = returnCode;
    return result;
}

int G2config_close_helper(uintptr_t configHandle)
{
    int returnCode = G2Config_close((void *)configHandle);
    return returnCode;
}

struct G2Config_create_result G2config_create_helper()
{
    ConfigHandle configHandle;
    int returnCode = G2Config_create(&configHandle);
    struct G2Config_create_result result;
    result.response = configHandle;
    result.returnCode = returnCode;
    return result;
}

int G2Config_deleteDataSource_helper(uintptr_t configHandle, const char *inputJson)
{
    int returnCode = G2Config_deleteDataSource((void *)configHandle, inputJson);
    return returnCode;
}

struct G2Config_listDataSources_result G2Config_listDataSources_helper(uintptr_t configHandle)
{
    size_t charBufferSize = 0;
    char *charBuffer = NULL;
    char **charBufferPtr = &charBuffer;
    resize_buffer_type resizeFuncPointer = &G2config_resizeStringBuffer;
    int returnCode = G2Config_listDataSources((void *)configHandle, charBufferPtr, &charBufferSize, resizeFuncPointer);
    struct G2Config_listDataSources_result result;
    result.response = *charBufferPtr;
    result.returnCode = returnCode;
    return result;
}

struct G2Config_load_result G2Config_load_helper(const char *inputJson)
{
    ConfigHandle configHandle;
    int returnCode = G2Config_load(inputJson, &configHandle);
    struct G2Config_load_result result;
    result.response = configHandle;
    result.returnCode = returnCode;
    return result;
}

struct G2Config_save_result G2Config_save_helper(uintptr_t configHandle)
{
    size_t charBufferSize = 0;
    char *charBuffer = NULL;
    char **charBufferPtr = &charBuffer;
    resize_buffer_type resizeFuncPointer = &G2config_resizeStringBuffer;
    int returnCode = G2Config_save((void *)configHandle, charBufferPtr, &charBufferSize, resizeFuncPointer);
    struct G2Config_save_result result;
    result.response = *charBufferPtr;
    result.returnCode = returnCode;
    return result;
}
