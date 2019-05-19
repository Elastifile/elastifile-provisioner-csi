# \ProjectsprojectoperationApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetOperation**](ProjectsprojectoperationApi.md#GetOperation) | **Get** /projects/{project}/operation/{resourceId} | Return operation resource
[**ListOperations**](ProjectsprojectoperationApi.md#ListOperations) | **Get** /projects/{project}/operation | Return list of operations


# **GetOperation**
> []Operation GetOperation($resourceId, $project)

Return operation resource


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **resourceId** | **string**| Operation resource id | 
 **project** | **string**|  | 

### Return type

[**[]Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListOperations**
> []Operation ListOperations($project)

Return list of operations


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **project** | **string**|  | 

### Return type

[**[]Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

