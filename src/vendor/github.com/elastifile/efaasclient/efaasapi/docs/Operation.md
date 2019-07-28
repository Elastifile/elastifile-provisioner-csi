# Operation

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | [Output Only] The unique identifier for the resource. This identifier is defined by the server. | [optional] [default to null]
**Name** | **string** | [Output Only] Name of the resource. | [optional] [default to null]
**Description** | **string** | [Output Only] A textual description of the resource. | [optional] [default to null]
**ClientOperationId** | **string** | [Output Only] The value of requestId if you provided it in the request. Not present otherwise. | [optional] [default to null]
**OperationType** | **string** | [Output Only] The type of operation, such as insert, update, or delete, and so on. | [optional] [default to null]
**TargetLink** | **string** | [Output Only] The URL of the resource that the operation modifies. | [optional] [default to null]
**TargetId** | **string** | [Output Only] The unique target ID | [optional] [default to null]
**Status** | **string** | [Output Only] The status of the operation, which can be one of the following: PENDING, RUNNING, or DONE. | [optional] [default to null]
**StatusMessage** | **string** | [Output Only] An optional textual description of the current status of the operation. | [optional] [default to null]
**User** | **string** | [Output Only] User who requested the operation, for example: user@example.com | [optional] [default to null]
**Progress** | **int32** | [Output Only] An optional progress indicator that ranges from 0 to 100. There is no requirement that this be linear or support any granularity of operations. This should not be used to guess when the operation will be complete. This number should monotonically increase as the operation progresses. | [optional] [default to null]
**InsertTime** | **string** | [Output Only] The time that this operation was requested. This value is in RFC3339 text format. | [optional] [default to null]
**StartTime** | **string** | [Output Only] The time that this operation was started by the server. This value is in RFC3339 text format. | [optional] [default to null]
**EndTime** | **string** | [Output Only] The time that this operation was completed. This value is in RFC3339 text format. | [optional] [default to null]
**Error_** | [***ModelError**](error.md) | [Output Only] If errors are generated during processing of the operation, this field will be populated. | [optional] [default to null]
**Warnings** | [**[]Warnings**](warnings.md) | [Output Only] If warning messages are generated during processing of the operation, this field will be populated. | [optional] [default to null]
**HttpErrorStatusCode** | **int32** | [Output Only] This field contains the HTTP error status code that was returned. For example, a 404 means the resource was not found. | [optional] [default to null]
**HttpErrorMessage** | **string** | [Output Only] This field contains the HTTP error message that was returned, such as NOT FOUND. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


