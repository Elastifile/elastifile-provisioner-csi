# \ProjectsprojectsnapshotsApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateShare**](ProjectsprojectsnapshotsApi.md#CreateShare) | **Post** /projects/{project}/snapshots/{resourceId}/share | Create share for the specified snapshot resource
[**DeleteShare**](ProjectsprojectsnapshotsApi.md#DeleteShare) | **Delete** /projects/{project}/snapshots/{resourceId}/share | Delete share for the specified snapshot resource
[**DeleteSnapshot**](ProjectsprojectsnapshotsApi.md#DeleteSnapshot) | **Delete** /projects/{project}/snapshots/{resourceId} | Deletes the specified snapshot
[**GetSnapshot**](ProjectsprojectsnapshotsApi.md#GetSnapshot) | **Get** /projects/{project}/snapshots/{resourceId} | Returns a specified snapshot
[**ListInstanceSnapshots**](ProjectsprojectsnapshotsApi.md#ListInstanceSnapshots) | **Get** /projects/{project}/snapshots/instances/{instance} | Return list of instance snapshots
[**ListSnapshots**](ProjectsprojectsnapshotsApi.md#ListSnapshots) | **Get** /projects/{project}/snapshots | Return list of instances snapshots for the specified project


# **CreateShare**
> Operation CreateShare($project, $resourceId, $payload, $requestId)

Create share for the specified snapshot resource


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **project** | **string**|  | 
 **resourceId** | **string**|  | 
 **payload** | [**SnapshotShareCreate**](SnapshotShareCreate.md)|  | 
 **requestId** | **string**| An optional request ID to identify requests. Specify a unique request ID so that if you must retry your request, the server will know to ignore the request if it has already been completed. For example, consider a situation where you make an initial request and the request times out. If you make the request again with the same request ID, the server can check if original operation with the same request ID was received, and if so, will ignore the second request. This prevents clients from accidentally creating duplicate commitments.  The request ID must be a valid UUID with the exception that zero UUID is not supported (00000000-0000-0000-0000-000000000000). | [optional] 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteShare**
> Operation DeleteShare($project, $resourceId, $requestId)

Delete share for the specified snapshot resource


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **project** | **string**|  | 
 **resourceId** | **string**|  | 
 **requestId** | **string**| An optional request ID to identify requests. Specify a unique request ID so that if you must retry your request, the server will know to ignore the request if it has already been completed. For example, consider a situation where you make an initial request and the request times out. If you make the request again with the same request ID, the server can check if original operation with the same request ID was received, and if so, will ignore the second request. This prevents clients from accidentally creating duplicate commitments.  The request ID must be a valid UUID with the exception that zero UUID is not supported (00000000-0000-0000-0000-000000000000). | [optional] 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteSnapshot**
> Operation DeleteSnapshot($project, $resourceId, $requestId)

Deletes the specified snapshot

Deleting a snapshot removes its data permanently and is irreversible


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **project** | **string**|  | 
 **resourceId** | **string**|  | 
 **requestId** | **string**| An optional request ID to identify requests. Specify a unique request ID so that if you must retry your request, the server will know to ignore the request if it has already been completed. For example, consider a situation where you make an initial request and the request times out. If you make the request again with the same request ID, the server can check if original operation with the same request ID was received, and if so, will ignore the second request. This prevents clients from accidentally creating duplicate commitments.  The request ID must be a valid UUID with the exception that zero UUID is not supported (00000000-0000-0000-0000-000000000000). | [optional] 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetSnapshot**
> Snapshots GetSnapshot($project, $resourceId)

Returns a specified snapshot


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **project** | **string**|  | 
 **resourceId** | **string**|  | 

### Return type

[**Snapshots**](Snapshots.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListInstanceSnapshots**
> []Snapshots ListInstanceSnapshots($project, $instance, $requestId)

Return list of instance snapshots


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **project** | **string**|  | 
 **instance** | **string**|  | 
 **requestId** | **string**| An optional request ID to identify requests. Specify a unique request ID so that if you must retry your request, the server will know to ignore the request if it has already been completed. For example, consider a situation where you make an initial request and the request times out. If you make the request again with the same request ID, the server can check if original operation with the same request ID was received, and if so, will ignore the second request. This prevents clients from accidentally creating duplicate commitments.  The request ID must be a valid UUID with the exception that zero UUID is not supported (00000000-0000-0000-0000-000000000000). | [optional] 

### Return type

[**[]Snapshots**](Snapshots.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListSnapshots**
> []Snapshots ListSnapshots($project, $filesystem, $instance)

Return list of instances snapshots for the specified project


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **project** | **string**|  | 
 **filesystem** | **string**| Filesystem id | [optional] 
 **instance** | **string**| The name of the resource, provided by the client when initially creating the resource. The resource name must be 1-63 characters long, and comply with RFC1035 | [optional] 

### Return type

[**[]Snapshots**](Snapshots.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

