package recoveryservicesbackup

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
)

// BackupResourceVaultConfigsClient is the open API 2.0 Specs for Azure RecoveryServices Backup service
type BackupResourceVaultConfigsClient struct {
	ManagementClient
}

// NewBackupResourceVaultConfigsClient creates an instance of the BackupResourceVaultConfigsClient client.
func NewBackupResourceVaultConfigsClient(subscriptionID string) BackupResourceVaultConfigsClient {
	return NewBackupResourceVaultConfigsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewBackupResourceVaultConfigsClientWithBaseURI creates an instance of the BackupResourceVaultConfigsClient client.
func NewBackupResourceVaultConfigsClientWithBaseURI(baseURI string, subscriptionID string) BackupResourceVaultConfigsClient {
	return BackupResourceVaultConfigsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get fetches resource vault config.
//
// vaultName is the name of the recovery services vault. resourceGroupName is the name of the resource group where the
// recovery services vault is present.
func (client BackupResourceVaultConfigsClient) Get(vaultName string, resourceGroupName string) (result BackupResourceVaultConfigResource, err error) {
	req, err := client.GetPreparer(vaultName, resourceGroupName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "recoveryservicesbackup.BackupResourceVaultConfigsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "recoveryservicesbackup.BackupResourceVaultConfigsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "recoveryservicesbackup.BackupResourceVaultConfigsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client BackupResourceVaultConfigsClient) GetPreparer(vaultName string, resourceGroupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vaultName":         autorest.Encode("path", vaultName),
	}

	const APIVersion = "2016-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/Subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupconfig/vaultconfig", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client BackupResourceVaultConfigsClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client BackupResourceVaultConfigsClient) GetResponder(resp *http.Response) (result BackupResourceVaultConfigResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Update updates vault security config.
//
// vaultName is the name of the recovery services vault. resourceGroupName is the name of the resource group where the
// recovery services vault is present. parameters is resource config request
func (client BackupResourceVaultConfigsClient) Update(vaultName string, resourceGroupName string, parameters BackupResourceVaultConfigResource) (result BackupResourceVaultConfigResource, err error) {
	req, err := client.UpdatePreparer(vaultName, resourceGroupName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "recoveryservicesbackup.BackupResourceVaultConfigsClient", "Update", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "recoveryservicesbackup.BackupResourceVaultConfigsClient", "Update", resp, "Failure sending request")
		return
	}

	result, err = client.UpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "recoveryservicesbackup.BackupResourceVaultConfigsClient", "Update", resp, "Failure responding to request")
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client BackupResourceVaultConfigsClient) UpdatePreparer(vaultName string, resourceGroupName string, parameters BackupResourceVaultConfigResource) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vaultName":         autorest.Encode("path", vaultName),
	}

	const APIVersion = "2016-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/Subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupconfig/vaultconfig", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare(&http.Request{})
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client BackupResourceVaultConfigsClient) UpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req)
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client BackupResourceVaultConfigsClient) UpdateResponder(resp *http.Response) (result BackupResourceVaultConfigResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
