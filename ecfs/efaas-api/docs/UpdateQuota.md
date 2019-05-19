# UpdateQuota

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**QuotaType** | **string** | Supported values are: auto and fixed. Use auto if you have one filesystem, the size of the filesystem will be the same as the instance size. Use fixed if you have more than one filesystem, and set the filesystem size through filesystemQuota. | [default to null]
**HardQuota** | **int64** | Set the size of a filesystem if filesystemQuotaType is set to fixed. If it is set to auto, this value is ignored and quota is the instance total size. | [default to 0]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


