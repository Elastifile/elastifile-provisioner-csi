# \AuthApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**Login**](AuthApi.md#Login) | **Post** /auth/login | Authenticates using JWT from Google OAuth


# **Login**
> Login Login($payload)

Authenticates using JWT from Google OAuth

Validates the Google Oauth token.


### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **payload** | [**Login**](Login.md)|  | 

### Return type

[**Login**](login.md)

### Authorization

[cloud-console-partner](../README.md#cloud-console-partner), [cloud-console-partner-autopush](../README.md#cloud-console-partner-autopush), [google-id](../README.md#google-id), [cloud-console-partner-local](../README.md#cloud-console-partner-local)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

