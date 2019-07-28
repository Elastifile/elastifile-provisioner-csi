# \ProjectsprojectinstancesApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AddFilesystem**](ProjectsprojectinstancesApi.md#AddFilesystem) | **Post** /projects/{project}/instances/{name}/filesystem | Add filesystem
[**CreateInstance**](ProjectsprojectinstancesApi.md#CreateInstance) | **Post** /projects/{project}/instances | Create an instances
[**DeleteFilesystem**](ProjectsprojectinstancesApi.md#DeleteFilesystem) | **Delete** /projects/{project}/instances/{name}/filesystem/{filesystem_id} | Delete filesystem
[**DeleteInstanceItem**](ProjectsprojectinstancesApi.md#DeleteInstanceItem) | **Delete** /projects/{project}/instances/{name} | Deletes the specified Instance resource
[**FilesystemManualCreateSnapshot**](ProjectsprojectinstancesApi.md#FilesystemManualCreateSnapshot) | **Post** /projects/{project}/instances/{name}/filesystem/{filesystem_id}/snapshots | Create manual snapshot
[**GetInstance**](ProjectsprojectinstancesApi.md#GetInstance) | **Get** /projects/{project}/instances/{name} | Get an instance
[**GetInstanceConstraints**](ProjectsprojectinstancesApi.md#GetInstanceConstraints) | **Get** /projects/{project}/instances/{name}/getConstraints | Get resource capacity constraints
[**GetInstanceStatistics**](ProjectsprojectinstancesApi.md#GetInstanceStatistics) | **Get** /projects/{project}/instances/{name}/statistics | Get resource statistics
[**ListInstances**](ProjectsprojectinstancesApi.md#ListInstances) | **Get** /projects/{project}/instances | Return list of instances
[**PostInstanceSetCapacity**](ProjectsprojectinstancesApi.md#PostInstanceSetCapacity) | **Post** /projects/{project}/instances/{name}/setCapacity | Sets capacity to instance
[**SetAccessorsToFilesystem**](ProjectsprojectinstancesApi.md#SetAccessorsToFilesystem) | **Post** /projects/{project}/instances/{name}/filesystem/{filesystem_id}/setAccessors | Filesystem set accessors
[**SetFilesystemDescription**](ProjectsprojectinstancesApi.md#SetFilesystemDescription) | **Post** /projects/{project}/instances/{name}/filesystem/{filesystem_id}/setDescription | Update filesystem description
[**SetFilesystemSnapshotScheduling**](ProjectsprojectinstancesApi.md#SetFilesystemSnapshotScheduling) | **Post** /projects/{project}/instances/{name}/filesystem/{filesystem_id}/setScheduling | Filesystem set snapshot scheduling
[**UpdateFilesystemQuota**](ProjectsprojectinstancesApi.md#UpdateFilesystemQuota) | **Post** /projects/{project}/instances/{name}/filesystem/{filesystem_id}/setQuota | Filesystem quota update


# **AddFilesystem**
> Operation AddFilesystem(ctx, name, project, payload)
Add filesystem

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Resource name | 
  **project** | **string**|  | 
  **payload** | [**DataContainerAdd**](DataContainerAdd.md)|  | 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CreateInstance**
> Operation CreateInstance(ctx, project, payload, optional)
Create an instances

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **project** | **string**|  | 
  **payload** | [**Instances**](Instances.md)|  | 
 **optional** | ***CreateInstanceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a CreateInstanceOpts struct

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

# **DeleteFilesystem**
> Operation DeleteFilesystem(ctx, name, filesystemId, project)
Delete filesystem

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Resource name | 
  **filesystemId** | **string**| Filesystem id | 
  **project** | **string**|  | 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteInstanceItem**
> Operation DeleteInstanceItem(ctx, name, project, optional)
Deletes the specified Instance resource

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Instance name | 
  **project** | **string**|  | 
 **optional** | ***DeleteInstanceItemOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a DeleteInstanceItemOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


 **force** | **optional.String**| [Experimental] Force operation, even if resource is not in ready state, possible values: true/false, on/off. Default false. | 
 **requestId** | **optional.String**| An optional request ID to identify requests. Specify a unique request ID so that if you must retry your request, the server will know to ignore the request if it has already been completed. For example, consider a situation where you make an initial request and the request times out. If you make the request again with the same request ID, the server can check if original operation with the same request ID was received, and if so, will ignore the second request. This prevents clients from accidentally creating duplicate commitments.  The request ID must be a valid UUID with the exception that zero UUID is not supported (00000000-0000-0000-0000-000000000000). | 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **FilesystemManualCreateSnapshot**
> Operation FilesystemManualCreateSnapshot(ctx, name, filesystemId, project, payload, optional)
Create manual snapshot

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Resource name | 
  **filesystemId** | **string**| Filesystem id | 
  **project** | **string**|  | 
  **payload** | [**Snapshot**](Snapshot.md)|  | 
 **optional** | ***FilesystemManualCreateSnapshotOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a FilesystemManualCreateSnapshotOpts struct

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

# **GetInstance**
> Instances GetInstance(ctx, name, project)
Get an instance

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Instance name | 
  **project** | **string**|  | 

### Return type

[**Instances**](Instances.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetInstanceConstraints**
> CapacityUnits GetInstanceConstraints(ctx, name, project)
Get resource capacity constraints

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Resource name | 
  **project** | **string**|  | 

### Return type

[**CapacityUnits**](capacityUnits.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetInstanceStatistics**
> Statistics GetInstanceStatistics(ctx, name, project)
Get resource statistics

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Resource name | 
  **project** | **string**|  | 

### Return type

[**Statistics**](statistics.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListInstances**
> []Instances ListInstances(ctx, project)
Return list of instances

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **project** | **string**|  | 

### Return type

[**[]Instances**](Instances.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostInstanceSetCapacity**
> Operation PostInstanceSetCapacity(ctx, name, project, payload)
Sets capacity to instance

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Resource name | 
  **project** | **string**|  | 
  **payload** | [**SetCapacity**](SetCapacity.md)|  | 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SetAccessorsToFilesystem**
> Operation SetAccessorsToFilesystem(ctx, name, filesystemId, project, payload)
Filesystem set accessors

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Resource name | 
  **filesystemId** | **string**| Filesystem id | 
  **project** | **string**|  | 
  **payload** | [**Accessors**](Accessors.md)|  | 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SetFilesystemDescription**
> Operation SetFilesystemDescription(ctx, name, filesystemId, project, payload)
Update filesystem description

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Resource name | 
  **filesystemId** | **string**| Filesystem id | 
  **project** | **string**|  | 
  **payload** | [**UpdateDesciption**](UpdateDesciption.md)|  | 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SetFilesystemSnapshotScheduling**
> Operation SetFilesystemSnapshotScheduling(ctx, name, filesystemId, project, payload)
Filesystem set snapshot scheduling

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Resource name | 
  **filesystemId** | **string**| Filesystem id | 
  **project** | **string**|  | 
  **payload** | [**SnapshotSchedule**](SnapshotSchedule.md)|  | 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UpdateFilesystemQuota**
> Operation UpdateFilesystemQuota(ctx, name, filesystemId, project, payload)
Filesystem quota update

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **name** | **string**| Resource name | 
  **filesystemId** | **string**| Filesystem id | 
  **project** | **string**|  | 
  **payload** | [**UpdateQuota**](UpdateQuota.md)|  | 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

