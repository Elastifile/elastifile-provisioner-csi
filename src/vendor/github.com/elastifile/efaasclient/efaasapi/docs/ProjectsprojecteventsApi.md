# \ProjectsprojecteventsApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**AckEvent**](ProjectsprojecteventsApi.md#AckEvent) | **Put** /projects/{project}/events/{id}/ack | 
[**CountEvents**](ProjectsprojecteventsApi.md#CountEvents) | **Get** /projects/{project}/events/count | 
[**GetEvent**](ProjectsprojecteventsApi.md#GetEvent) | **Get** /projects/{project}/events/{id} | 
[**ListEvents**](ProjectsprojecteventsApi.md#ListEvents) | **Get** /projects/{project}/events | 
[**UnAckEvent**](ProjectsprojecteventsApi.md#UnAckEvent) | **Put** /projects/{project}/events/{id}/unack | 


# **AckEvent**
> Events AckEvent(ctx, id, project)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| The event identifier | 
  **project** | **string**|  | 

### Return type

[**Events**](Events.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **CountEvents**
> EventsCount CountEvents(ctx, project)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **project** | **string**|  | 

### Return type

[**EventsCount**](EventsCount.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetEvent**
> Events GetEvent(ctx, id, project)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| The event identifier | 
  **project** | **string**|  | 

### Return type

[**Events**](Events.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **ListEvents**
> []Events ListEvents(ctx, project)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **project** | **string**|  | 

### Return type

[**[]Events**](Events.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **UnAckEvent**
> Events UnAckEvent(ctx, id, project)


### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **id** | **string**| The event identifier | 
  **project** | **string**|  | 

### Return type

[**Events**](Events.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [cloud-console-partner-local](../README.md#cloud-console-partner-local), [google-id](../README.md#google-id)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

