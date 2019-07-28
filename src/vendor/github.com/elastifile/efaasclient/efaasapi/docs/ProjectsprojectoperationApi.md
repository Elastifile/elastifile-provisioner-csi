# \ProjectsprojectoperationApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetOperation**](ProjectsprojectoperationApi.md#GetOperation) | **Get** /projects/{project}/operation/{resourceId} | Return operation resource
[**ListOperations**](ProjectsprojectoperationApi.md#ListOperations) | **Get** /projects/{project}/operation | Return list of operations


# **GetOperation**
> Operation GetOperation(ctx, resourceId, project)
Return operation resource

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **resourceId** | **string**| Operation resource id | 
  **project** | **string**|  | 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListOperations**
> []Operation ListOperations(ctx, project)
Return list of operations

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **project** | **string**|  | 

### Return type

[**[]Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

