#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include "libg2config.h"

typedef void *(*resize_buffer_type)(void *, size_t);

struct G2Config_addDataSource_result
{
    char *response;
    int returnCode;
};

struct G2Config_create_result
{
    void *response;
    int returnCode;
};

struct G2Config_listDataSources_result
{
    char *response;
    int returnCode;
};

struct G2Config_save_result
{
    char *response;
    int returnCode;
};

struct G2Config_addDataSource_result G2Config_addDataSource_helper(uintptr_t configHandle, const char *inputJson);
int G2config_close_helper(uintptr_t configHandle);
struct G2Config_create_result G2config_create_helper();
int G2Config_deleteDataSource_helper(uintptr_t configHandle, const char *inputJson);
struct G2Config_listDataSources_result G2Config_listDataSources_helper(uintptr_t configHandle);
int G2Config_load_helper(uintptr_t configHandle, const char *inputJson);
struct G2Config_save_result G2Config_save_helper(uintptr_t configHandle);
