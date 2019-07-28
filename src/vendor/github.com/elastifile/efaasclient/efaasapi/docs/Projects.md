# Projects

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | Project numeric id, which is automatically assigned when you create Google cloud project. | [default to null]
**Name** | **string** | Project ID, which is a unique identifier for the project. | [default to null]
**DisplayName** | **string** | Project display name. | [default to null]
**AllowedUsers** | [**[]AllowedUser**](AllowedUser.md) | List of users allowed to access resources on the specified project. | [optional] [default to null]
**AlphaEnabled** | **bool** | Alpha features enabled on this project | [optional] [default to null]
**Status** | **string** | The status of the project, which can be one of the following: PENDING_APPROVAL, ENABLED, or DISABLED. | [optional] [default to null]
**CreationTimestamp** | [**time.Time**](time.Time.md) | [Output Only] Creation timestamp in RFC3339 text format. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


