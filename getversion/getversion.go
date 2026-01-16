package getversion

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/senzing-garage/sz-sdk-go-core/szproduct"
)

type SenzingVersionResponse struct {
	Version string `json:"VERSION"`
}

/*
The GetSenzingVersion function returns an integer in the form MMmmPP where:
- MM is the Major
- mm is the Minor
- PP is the Patch
For instance 40103 is semantic version 4.1.3.

Limitations: neither Major, Minor, nor Patch can be greater than 99.

Input
  - ctx: a context.

Output
  - An integer in the base-ten form MMmmPP.
*/
func GetSenzingVersion(ctx context.Context) int {

	var result int

	szProduct := &szproduct.Szproduct{}
	versionJSON, err := szProduct.GetVersion(ctx)

	senzingVersionResponse := SenzingVersionResponse{}
	err = json.Unmarshal([]byte(versionJSON), &senzingVersionResponse)
	if err != nil {
		panic(err)
	}

	// Parse semantic version string (Major.Minor.Patch)
	parts := strings.Split(senzingVersionResponse.Version, ".")
	if len(parts) >= 3 {
		major, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(err)
		}
		minor, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}
		patch, err := strconv.Atoi(parts[2])
		if err != nil {
			panic(err)
		}
		result = (major * 10000) + (minor * 100) + patch
	}

	return result
}
