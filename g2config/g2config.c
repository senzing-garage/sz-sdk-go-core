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
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2config_resizeStringBuffer;
    int returnCode = G2Config_addDataSource((void *)configHandle, inputJson, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2Config_addDataSource_result result;
    result.response = charBuffer;
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
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2config_resizeStringBuffer;
    int returnCode = G2Config_listDataSources((void *)configHandle, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2Config_listDataSources_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

int G2Config_load_helper(uintptr_t configHandle, const char *inputJson)
{
    int returnCode = G2Config_load(inputJson, (void *)configHandle);
    return returnCode;
}

struct G2Config_save_result G2Config_save_helper(uintptr_t configHandle)
{
    size_t charBufferSize = 1;
    char *charBuffer = (char *)malloc(charBufferSize);
    resize_buffer_type resizeFuncPointer = &G2config_resizeStringBuffer;
    int returnCode = G2Config_save((void *)configHandle, &charBuffer, &charBufferSize, resizeFuncPointer);
    struct G2Config_save_result result;
    result.response = charBuffer;
    result.returnCode = returnCode;
    return result;
}

// == DEBUG ===================================================================

int G2config_close_helper_debug(uintptr_t configHandle)
{
    printf(">>>> Close >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n");
    printf(" configHandle: %lu\n", configHandle);
    printf("&configHandle: %p\n", &configHandle);
    fflush(stdout);
    printf("<<<< Close <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n");
    int returnCode = G2Config_close((void *)configHandle);
    return returnCode;
}

void *G2config_create_helper_debug()
{
    ConfigHandle configHandle;
    printf(">>>> Create >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>\n");
    fflush(stdout);
    int returnCode = G2Config_create(&configHandle);
    printf("Return  code: %i\n", returnCode);
    printf("configHandle: %p\n", configHandle);
    printf("<<<< Create <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n");
    fflush(stdout);
    return configHandle;
}
