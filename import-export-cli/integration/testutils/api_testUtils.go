/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package testutils

import (
	"github.com/stretchr/testify/assert"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/apim"
	"github.com/wso2/product-apim-tooling/import-export-cli/integration/base"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func AddAPI(t *testing.T, client *apim.Client, username string, password string) *apim.API {
	client.Login(username, password)
	api := client.GenerateSampleAPIData(username)
	doClean := true
	id := client.AddAPI(t, api, username, password, doClean)
	api = client.GetAPI(id)
	return api
}

func AddAPIWithoutCleaning(t *testing.T, client *apim.Client, username string, password string) *apim.API {
	client.Login(username, password)
	api := client.GenerateSampleAPIData(username)
	doClean := false
	id := client.AddAPI(t, api, username, password, doClean)
	api = client.GetAPI(id)
	return api
}

func AddAPIToTwoEnvs(t *testing.T, client1 *apim.Client, client2 *apim.Client, username string, password string) (*apim.API, *apim.API) {
	client1.Login(username, password)
	api := client1.GenerateSampleAPIData(username)
	doClean := true
	id1 := client1.AddAPI(t, api, username, password, doClean)
	api1 := client1.GetAPI(id1)

	client2.Login(username, password)
	id2 := client2.AddAPI(t, api, username, password, doClean)
	api2 := client2.GetAPI(id2)

	return api1, api2
}

func AddAPIFromOpenAPIDefinition(t *testing.T, client *apim.Client, username string, password string) *apim.API {
	path := "testdata/petstore.yaml"
	client.Login(username, password)
	additionalProperties := client.GenerateAdditionalProperties(username)
	id := client.AddAPIFromOpenAPIDefinition(t, path, additionalProperties, username, password)
	api := client.GetAPI(id)
	return api
}

func AddAPIFromOpenAPIDefinitionToTwoEnvs(t *testing.T, client1 *apim.Client, client2 *apim.Client, username string, password string) (*apim.API, *apim.API) {
	path := "testdata/petstore.yaml"
	client1.Login(username, password)
	additionalProperties := client1.GenerateAdditionalProperties(username)
	id1 := client1.AddAPIFromOpenAPIDefinition(t, path, additionalProperties, username, password)
	api1 := client1.GetAPI(id1)

	client2.Login(username, password)
	id2 := client2.AddAPIFromOpenAPIDefinition(t, path, additionalProperties, username, password)
	api2 := client2.GetAPI(id2)

	return api1, api2
}

func AddAPIProductFromJSON(t *testing.T, client *apim.Client, username string, password string, apisList map[string]*apim.API) *apim.APIProduct {
	client.Login(username, password)
	path := "testdata/SampleAPIProduct.json"
	doClean := true
	id := client.AddAPIProductFromJSON(t, path, username, password, apisList, doClean)
	apiProduct := client.GetAPIProduct(id)
	return apiProduct
}

func getAPI(t *testing.T, client *apim.Client, name string, username string, password string) *apim.API {
	client.Login(username, password)
	apiInfo := client.GetAPIByName(name)
	return client.GetAPI(apiInfo.ID)
}

func getAPIs(client *apim.Client, username string, password string) *apim.APIList {
	client.Login(username, password)
	return client.GetAPIs()
}

func deleteAPI(t *testing.T, client *apim.Client, apiID string, username string, password string) {
	time.Sleep(2000 * time.Millisecond)
	client.Login(username, password)
	client.DeleteAPI(apiID)
}

func deleteAPIByCtl(t *testing.T, args *ApiImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "delete", "api", "-n", args.Api.Name, "-v", args.Api.Version, "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

func PublishAPI(client *apim.Client, username string, password string, apiID string) {
	time.Sleep(2000 * time.Millisecond)
	client.Login(username, password)
	client.PublishAPI(apiID)
}

func UnsubscribeAPI(client *apim.Client, username string, password string, apiID string) {
	client.Login(username, password)
	client.DeleteSubscriptions(apiID)
}

func getResourceURL(apim *apim.Client, api *apim.API) string {
	port := 8280 + apim.GetPortOffset()
	return "http://" + apim.GetHost() + ":" + strconv.Itoa(port) + api.Context + "/" + api.Version + "/menu"
}

func getEnvAPIExportPath(envName string) string {
	return filepath.Join(utils.DefaultExportDirPath, utils.ExportedApisDirName, envName)
}

func exportAPI(t *testing.T, name string, version string, provider string, env string) (string, error) {
	var output string
	var err error

	if provider == "" {
		output, err = base.Execute(t, "export-api", "-n", name, "-v", version, "-e", env, "-k", "--verbose")
	} else {
		output, err = base.Execute(t, "export-api", "-n", name, "-v", version, "-r", provider, "-e", env, "-k", "--verbose")
	}

	t.Cleanup(func() {
		base.RemoveAPIArchive(t, getEnvAPIExportPath(env), name, version)
	})

	return output, err
}

func importAPI(t *testing.T, sourceEnv string, api *apim.API, client *apim.Client) (string, error) {
	fileName := base.GetAPIArchiveFilePath(t, sourceEnv, api.Name, api.Version)
	output, err := base.Execute(t, "import-api", "-f", fileName, "-e", client.EnvName, "-k", "--verbose", "--preserve-provider=false")

	t.Cleanup(func() {
		client.DeleteAPIByName(api.Name)
	})

	return output, err
}

func importAPIPreserveProvider(t *testing.T, sourceEnv string, api *apim.API, client *apim.Client) (string, error) {
	fileName := base.GetAPIArchiveFilePath(t, sourceEnv, api.Name, api.Version)
	output, err := base.Execute(t, "import-api", "-f", fileName, "-e", client.EnvName, "-k", "--verbose")

	t.Cleanup(func() {
		client.DeleteAPIByName(api.Name)
	})

	return output, err
}

func importAPIPreserveProviderFailure(t *testing.T, sourceEnv string, api *apim.API, client *apim.Client) (string, error) {
	fileName := base.GetAPIArchiveFilePath(t, sourceEnv, api.Name, api.Version)
	output, err := base.Execute(t, "import-api", "-f", fileName, "-e", client.EnvName, "-k", "--verbose")
	return output, err
}

func listAPIs(t *testing.T, args *ApiImportExportTestArgs) (string, error) {
	output, err := base.Execute(t, "list", "apis", "-e", args.SrcAPIM.EnvName, "-k", "--verbose")
	return output, err
}

func ValidateAPIExportFailure(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl env
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Attempt exporting api from env
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportAPI(t, args.Api.Name, args.Api.Version, args.ApiProvider.Username, args.SrcAPIM.GetEnvName())

	// Validate that export failed
	assert.False(t, base.IsAPIArchiveExists(t, getEnvAPIExportPath(args.SrcAPIM.GetEnvName()),
		args.Api.Name, args.Api.Version))
}

func ValidateAPIExportImport(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportAPI(t, args.Api.Name, args.Api.Version, args.Api.Provider, args.SrcAPIM.GetEnvName())

	assert.True(t, base.IsAPIArchiveExists(t, getEnvAPIExportPath(args.SrcAPIM.GetEnvName()),
		args.Api.Name, args.Api.Version))

	// Import api to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importAPIPreserveProvider(t, args.SrcAPIM.GetEnvName(), args.Api, args.DestAPIM)

	// Give time for newly imported API to get indexed, or else getAPI by name will fail
	time.Sleep(1 * time.Second)

	// Get App from env 2
	importedAPI := getAPI(t, args.DestAPIM, args.Api.Name, args.ApiProvider.Username, args.ApiProvider.Password)

	// Validate env 1 and env 2 API is equal
	validateAPIsEqual(t, args.Api, importedAPI)
}

func ValidateAPIExport(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())
	base.SetupEnv(t, args.DestAPIM.GetEnvName(), args.DestAPIM.GetApimURL(), args.DestAPIM.GetTokenURL())

	// Export api from env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	exportAPI(t, args.Api.Name, args.Api.Version, args.ApiProvider.Username, args.SrcAPIM.GetEnvName())

	assert.True(t, base.IsAPIArchiveExists(t, getEnvAPIExportPath(args.SrcAPIM.GetEnvName()),
		args.Api.Name, args.Api.Version))
}

func ValidateAPIImport(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Import api to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	importAPI(t, args.SrcAPIM.GetEnvName(), args.Api, args.DestAPIM)

	// Give time for newly imported API to get indexed, or else getAPI by name will fail
	time.Sleep(1 * time.Second)

	// Get App from env 2
	importedAPI := getAPI(t, args.DestAPIM, args.Api.Name, args.ApiProvider.Username, args.ApiProvider.Password)

	// Validate env 1 and env 2 API is equal
	validateAPIsEqualCrossTenant(t, args.Api, importedAPI)
}

func ValidateAPIImportFailure(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Import api to env 2
	base.Login(t, args.DestAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	// importAPIPreserveProviderFailure is used to eleminate cleaning the API after importing
	result, err := importAPIPreserveProviderFailure(t, args.SrcAPIM.GetEnvName(), args.Api, args.DestAPIM)

	assert.NotNil(t, err, "Expected error was not returned")
	assert.Contains(t, base.GetValueOfUniformResponse(result), "Exit status 1")
}

func validateAPIsEqual(t *testing.T, api1 *apim.API, api2 *apim.API) {
	t.Helper()

	api1Copy := apim.CopyAPI(api1)
	api2Copy := apim.CopyAPI(api2)

	same := "override_with_same_value"
	// Since the APIs are from too different envs, their respective ID will defer.
	// Therefore this will be overriden to the same value to ensure that the equality check will pass.
	api1Copy.ID = same
	api2Copy.ID = same

	api1Copy.CreatedTime = same
	api2Copy.CreatedTime = same

	api1Copy.LastUpdatedTime = same
	api2Copy.LastUpdatedTime = same

	// Sort member collections to make equality chack possible
	apim.SortAPIMembers(&api1Copy)
	apim.SortAPIMembers(&api2Copy)

	assert.Equal(t, api1Copy, api2Copy, "API obejcts are not equal")
}

func ValidateAPIsList(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// List APIs of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	output, _ := listAPIs(t, args)

	apisList := args.SrcAPIM.GetAPIs()

	validateListAPIsEqual(t, output, apisList)
}

func validateListAPIsEqual(t *testing.T, apisListFromCtl string, apisList *apim.APIList) {

	for _, api := range apisList.List {
		// If the output string contains the same API ID, then decrement the count
		if strings.Contains(apisListFromCtl, api.ID) {
			apisList.Count = apisList.Count - 1
		}
	}

	// Count == 0 means that all the APIs from apisList were in apisListFromCtl
	assert.Equal(t, apisList.Count, 0, "API lists are not equal")
}

func validateAPIsEqualCrossTenant(t *testing.T, api1 *apim.API, api2 *apim.API) {
	t.Helper()

	api1Copy := apim.CopyAPI(api1)
	api2Copy := apim.CopyAPI(api2)

	same := "override_with_same_value"
	// Since the APIs are from too different envs, their respective ID will defer.
	// Therefore this will be overridden to the same value to ensure that the equality check will pass.
	api1Copy.ID = same
	api2Copy.ID = same

	api1Copy.CreatedTime = same
	api2Copy.CreatedTime = same

	api1Copy.LastUpdatedTime = same
	api2Copy.LastUpdatedTime = same

	// The contexts and providers will differ since this is a cross tenant import
	// Therefore this will be overridden to the same value to ensure that the equality check will pass.
	api1Copy.Context = same
	api2Copy.Context = same

	api1Copy.Provider = same
	api2Copy.Provider = same

	// Sort member collections to make equality check possible
	apim.SortAPIMembers(&api1Copy)
	apim.SortAPIMembers(&api2Copy)

	assert.Equal(t, api1Copy, api2Copy, "API obejcts are not equal")
}

func ValidateAPIDelete(t *testing.T, args *ApiImportExportTestArgs) {
	t.Helper()

	// Setup apictl envs
	base.SetupEnv(t, args.SrcAPIM.GetEnvName(), args.SrcAPIM.GetApimURL(), args.SrcAPIM.GetTokenURL())

	// Delete an API of env 1
	base.Login(t, args.SrcAPIM.GetEnvName(), args.CtlUser.Username, args.CtlUser.Password)

	time.Sleep(1 * time.Second)
	apisListBeforeDelete := args.SrcAPIM.GetAPIs()

	deleteAPIByCtl(t, args)

	apisListAfterDelete := args.SrcAPIM.GetAPIs()
	time.Sleep(1 * time.Second)

	// Validate whether the expected number of API count is there
	assert.Equal(t, apisListBeforeDelete.Count, apisListAfterDelete.Count+1, "Expected number of APIs not deleted")

	// Validate that the delete is a success
	validateAPIIsDeleted(t, args.Api, apisListAfterDelete)
}

func validateAPIIsDeleted(t *testing.T, api *apim.API, apisListAfterDelete *apim.APIList) {
	for _, existingAPI := range apisListAfterDelete.List {
		assert.NotEqual(t, existingAPI.ID, api.ID, "API delete is not successful")
	}
}

func ImportApiFromProject(t *testing.T, projectName string, envName string) (string, error) {
	projectPath, _ := filepath.Abs(projectName)
	return base.Execute(t, "import-api", "-f", projectPath, "-e", envName, "-k")
}

func ImportApiFromProjectWithUpdate(t *testing.T, projectName string, envName string) (string, error) {
	projectPath, _ := filepath.Abs(projectName)
	return base.Execute(t, "import-api", "-f", projectPath, "-e", envName, "-k", "--update")
}

func ExportApisWithOneCommand(t *testing.T, args *InitTestArgs) (string, error) {
	output, error := base.Execute(t, "export-apis", "-e", args.SrcAPIM.GetEnvName(), "-k")
	return output, error
}
