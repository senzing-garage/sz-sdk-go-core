#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include "libg2configmgr.h"

typedef void *(*resize_buffer_type)(void *, size_t);

struct G2ConfigMgr_addConfig_result
{
    long long configID;
    int returnCode;
};

struct G2ConfigMgr_getConfig_result
{
    char *config;
    int returnCode;
};
struct G2ConfigMgr_getConfigList_result
{
    char *configList;
    int returnCode;
};

struct G2ConfigMgr_getDefaultConfigID_result
{
    long long configID;
    int returnCode;
};

struct G2ConfigMgr_addConfig_result G2ConfigMgr_addConfig_helper(const char *configStr, const char *configComments);
struct G2ConfigMgr_getConfig_result G2ConfigMgr_getConfig_helper(const long long configID);
struct G2ConfigMgr_getConfigList_result G2ConfigMgr_getConfigList_helper();
struct G2ConfigMgr_getDefaultConfigID_result G2ConfigMgr_getDefaultConfigID_helper();
