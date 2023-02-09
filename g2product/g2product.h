#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include "libg2product.h"

typedef void *(*resize_buffer_type)(void *, size_t);

struct G2Product_validateLicenseFile_result
{
    char *response;
    int returnCode;
};

struct G2Product_validateLicenseStringBase64_result
{
    char *response;
    int returnCode;
};

struct G2Product_validateLicenseFile_result G2Product_validateLicenseFile_helper(const char *licenseFilePath);
struct G2Product_validateLicenseStringBase64_result G2Product_validateLicenseStringBase64_helper(const char *licenseString);
