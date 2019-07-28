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
> Operation CreateShare(ctx, project, resourceId, payload, optional)
Create share for the specified snapshot resource

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **project** | **string**|  | 
  **resourceId** | **string**|  | 
  **payload** | [**SnapshotShareCreate**](SnapshotShareCreate.md)|  | 
 **optional** | ***CreateShareOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a CreateShareOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **requestId** | **optional.String**| An optional request ID to identify requests. Specify a unique request ID so that if you must retry your request, the server will know to ignore the request if it has already been completed. For example, consider a situation where you make an initial request and the request times out. If you make the request again with the same request ID, the server can check if original operation with the same request ID was received, and if so, will ignore the second request. This prevents clients from accidentally creating duplicate commitments.  The request ID must be a valid UUID with the exception that zero UUID is not supported (00000000-0000-0000-0000-000000000000). | 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteShare**
> Operation DeleteShare(ctx, project, resourceId, optional)
Delete share for the specified snapshot resource

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **project** | **string**|  | 
  **resourceId** | **string**|  | 
 **optional** | ***DeleteShareOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeleteShareOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **requestId** | **optional.String**| An optional request ID to identify requests. Specify a unique request ID so that if you must retry your request, the server will know to ignore the request if it has already been completed. For example, consider a situation where you make an initial request and the request times out. If you make the request again with the same request ID, the server can check if original operation with the same request ID was received, and if so, will ignore the second request. This prevents clients from accidentally creating duplicate commitments.  The request ID must be a valid UUID with the exception that zero UUID is not supported (00000000-0000-0000-0000-000000000000). | 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteSnapshot**
> Operation DeleteSnapshot(ctx, project, resourceId, optional)
Deletes the specified snapshot

Deleting a snapshot removes its data permanently and is irreversible

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **project** | **string**|  | 
  **resourceId** | **string**|  | 
 **optional** | ***DeleteSnapshotOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeleteSnapshotOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **requestId** | **optional.String**| An optional request ID to identify requests. Specify a unique request ID so that if you must retry your request, the server will know to ignore the request if it has already been completed. For example, consider a situation where you make an initial request and the request times out. If you make the request again with the same request ID, the server can check if original operation with the same request ID was received, and if so, will ignore the second request. This prevents clients from accidentally creating duplicate commitments.  The request ID must be a valid UUID with the exception that zero UUID is not supported (00000000-0000-0000-0000-000000000000). | 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetSnapshot**
> Snapshots GetSnapshot(ctx, project, resourceId)
Returns a specified snapshot

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **project** | **string**|  | 
  **resourceId** | **string**|  | 

### Return type

[**Snapshots**](Snapshots.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListInstanceSnapshots**
> []Snapshots ListInstanceSnapshots(ctx, project, instance, optional)
Return list of instance snapshots

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **project** | **string**|  | 
  **instance** | **string**|  | 
 **optional** | ***ListInstanceSnapshotsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ListInstanceSnapshotsOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **requestId** | **optional.String**| An optional request ID to identify requests. Specify a unique request ID so that if you must retry your request, the server will know to ignore the request if it has already been completed. For example, consider a situation where you make an initial request and the request times out. If you make the request again with the same request ID, the server can check if original operation with the same request ID was received, and if so, will ignore the second request. This prevents clients from accidentally creating duplicate commitments.  The request ID must be a valid UUID with the exception that zero UUID is not supported (00000000-0000-0000-0000-000000000000). | 

### Return type

[**[]Snapshots**](Snapshots.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListSnapshots**
> []Snapshots ListSnapshots(ctx, project, optional)
Return list of instances snapshots for the specified project

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **project** | **string**|  | 
 **optional** | ***ListSnapshotsOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ListSnapshotsOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **filesystem** | **optional.String**| Filesystem id | 
 **instance** | **optional.String**| The name of the resource, provided by the client when initially creating the resource. The resource name must be 1-63 characters long, and comply with RFC1035 | 

### Return type

[**[]Snapshots**](Snapshots.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

