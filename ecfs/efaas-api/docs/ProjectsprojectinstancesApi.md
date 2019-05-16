# \ProjectsprojectinstancesApi

All URIs are relative to *https://localhost/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**CreateInstance**](ProjectsprojectinstancesApi.md#CreateInstance) | **Post** /projects/{project}/instances | 
[**DeleteInstanceItem**](ProjectsprojectinstancesApi.md#DeleteInstanceItem) | **Delete** /projects/{project}/instances/{name} | Deletes the specified Instance resource
[**EditInstance**](ProjectsprojectinstancesApi.md#EditInstance) | **Post** /projects/{project}/instances/{name}/editParallel | Edit instance
[**GetInstance**](ProjectsprojectinstancesApi.md#GetInstance) | **Get** /projects/{project}/instances/{name} | :param name:
[**GetInstanceConstraints**](ProjectsprojectinstancesApi.md#GetInstanceConstraints) | **Get** /projects/{project}/instances/{name}/getConstraints | Get resource capacity constraints
[**GetInstanceStatistics**](ProjectsprojectinstancesApi.md#GetInstanceStatistics) | **Get** /projects/{project}/instances/{name}/statistics | Get resource statistics
[**ListInstances**](ProjectsprojectinstancesApi.md#ListInstances) | **Get** /projects/{project}/instances | Return list of instances
[**PostInstanceSetCapacity**](ProjectsprojectinstancesApi.md#PostInstanceSetCapacity) | **Post** /projects/{project}/instances/{name}/setCapacity | Sets capacity to instance
[**PostInstanceSetScheduling**](ProjectsprojectinstancesApi.md#PostInstanceSetScheduling) | **Post** /projects/{project}/instances/{name}/setScheduling | Update instance snapshot scheduling
[**SetsAccessorsForTheSpecifiedInstanceToTheDataIncludedInTheRequest_**](ProjectsprojectinstancesApi.md#SetsAccessorsForTheSpecifiedInstanceToTheDataIncludedInTheRequest_) | **Post** /projects/{project}/instances/{name}/setAccessors | Sets accessors for the specified instance to the data included in the request


# **CreateInstance**
> Operation CreateInstance($project, $payload, $requestId)




### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **project** | **string**|  | 
 **payload** | [**Instances**](Instances.md)|  | 
 **requestId** | **string**| An optional request ID to identify requests. Specify a unique request ID so that if you must retry your request, the server will know to ignore the request if it has already been completed. For example, consider a situation where you make an initial request and the request times out. If you make the request again with the same request ID, the server can check if original operation with the same request ID was received, and if so, will ignore the second request. This prevents clients from accidentally creating duplicate commitments.  The request ID must be a valid UUID with the exception that zero UUID is not supported (00000000-0000-0000-0000-000000000000). | [optional] 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **DeleteInstanceItem**
> Operation DeleteInstanceItem($name, $project, $force, $requestId)

Deletes the specified Instance resource


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| Instance name | 
 **project** | **string**|  | 
 **force** | **string**| [Experimental] Force operation, even if resource is not in ready state, possible values: true/false, on/off. Default false. | [optional] 
 **requestId** | **string**| An optional request ID to identify requests. Specify a unique request ID so that if you must retry your request, the server will know to ignore the request if it has already been completed. For example, consider a situation where you make an initial request and the request times out. If you make the request again with the same request ID, the server can check if original operation with the same request ID was received, and if so, will ignore the second request. This prevents clients from accidentally creating duplicate commitments.  The request ID must be a valid UUID with the exception that zero UUID is not supported (00000000-0000-0000-0000-000000000000). | [optional] 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **EditInstance**
> Operation EditInstance($name, $project, $payload, $requestId)

Edit instance


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| Resource name | 
 **project** | **string**|  | 
 **payload** | [**EditParallel**](EditParallel.md)|  | 
 **requestId** | **string**| An optional request ID to identify requests. Specify a unique request ID so that if you must retry your request, the server will know to ignore the request if it has already been completed. For example, consider a situation where you make an initial request and the request times out. If you make the request again with the same request ID, the server can check if original operation with the same request ID was received, and if so, will ignore the second request. This prevents clients from accidentally creating duplicate commitments.  The request ID must be a valid UUID with the exception that zero UUID is not supported (00000000-0000-0000-0000-000000000000). | [optional] 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetInstance**
> Instances GetInstance($name, $project)

:param name:

:rtype:


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| Instance name | 
 **project** | **string**|  | 

### Return type

[**Instances**](Instances.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetInstanceConstraints**
> CapacityUnits GetInstanceConstraints($name, $project)

Get resource capacity constraints


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| Resource name | 
 **project** | **string**|  | 

### Return type

[**CapacityUnits**](capacityUnits.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetInstanceStatistics**
> Statistics GetInstanceStatistics($name, $project)

Get resource statistics


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| Resource name | 
 **project** | **string**|  | 

### Return type

[**Statistics**](statistics.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListInstances**
> []Instances ListInstances($project)

Return list of instances


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **project** | **string**|  | 

### Return type

[**[]Instances**](Instances.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostInstanceSetCapacity**
> Operation PostInstanceSetCapacity($name, $project, $payload)

Sets capacity to instance


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| Resource name | 
 **project** | **string**|  | 
 **payload** | [**SetCapacity**](SetCapacity.md)|  | 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostInstanceSetScheduling**
> Operation PostInstanceSetScheduling($name, $project, $payload)

Update instance snapshot scheduling


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| Resource name | 
 **project** | **string**|  | 
 **payload** | [**Snapshots**](Snapshots.md)|  | 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **SetsAccessorsForTheSpecifiedInstanceToTheDataIncludedInTheRequest_**
> Operation SetsAccessorsForTheSpecifiedInstanceToTheDataIncludedInTheRequest_($name, $project, $payload)

Sets accessors for the specified instance to the data included in the request


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string**| Resource name | 
 **project** | **string**|  | 
 **payload** | [**Accessors**](Accessors.md)|  | 

### Return type

[**Operation**](Operation.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

