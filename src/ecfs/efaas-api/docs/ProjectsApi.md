# \ProjectsApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetProject**](ProjectsApi.md#GetProject) | **Get** /projects/{project} | Get project resource
[**ListProjects**](ProjectsApi.md#ListProjects) | **Get** /projects | List projects
[**ProjectAddUsers**](ProjectsApi.md#ProjectAddUsers) | **Post** /projects/{project}/addUsers | Add users to project
[**ProjectRemoveUsers**](ProjectsApi.md#ProjectRemoveUsers) | **Post** /projects/{project}/removeUsers | Remove users from project
[**RegisterProject**](ProjectsApi.md#RegisterProject) | **Post** /projects | Register project for use with the service


# **GetProject**
> Projects GetProject($project)

Get project resource


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **project** | **string**| The project numeric id | 

### Return type

[**Projects**](Projects.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListProjects**
> []Projects ListProjects()

List projects


### Parameters
This endpoint does not need any parameter.

### Return type

[**[]Projects**](Projects.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ProjectAddUsers**
> Projects ProjectAddUsers($project, $payload)

Add users to project


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **project** | **string**| The project numeric id | 
 **payload** | [**Users**](Users.md)|  | 

### Return type

[**Projects**](Projects.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ProjectRemoveUsers**
> Projects ProjectRemoveUsers($project, $payload)

Remove users from project


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **project** | **string**| The project numeric id | 
 **payload** | [**Users**](Users.md)|  | 

### Return type

[**Projects**](Projects.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **RegisterProject**
> Projects RegisterProject($payload)

Register project for use with the service


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **payload** | [**Projects**](Projects.md)|  | 

### Return type

[**Projects**](Projects.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

