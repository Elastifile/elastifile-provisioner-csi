# Snapshots

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | [Output Only] The unique identifier for the snapshot. This identifier is defined by the server. | [optional] [default to null]
**Name** | **string** | The name of the resource, provided by the client when initially creating the resource. The resource name must be 1-63 characters long, and comply with RFC1035 | [default to null]
**Retention** | **float32** | Snapshot retention policy. The number of days to hold the snapshot till automatic deletion. Default 0, meaning no retention policy defined. | [default to null]
**InstanceId** | **string** | [Output Only] The filesystem instance id that this snapshot was taken for. | [optional] [default to null]
**InstanceName** | **string** | [Output Only] The filesystem instance name that this snapshot was taken for. | [optional] [default to null]
**FilesystemId** | **string** | [Output Only] The filesystem id that this snapshot was taken for. | [optional] [default to null]
**FilesystemName** | **string** | [Output Only] The filesystem name that this snapshot was taken for. | [optional] [default to null]
**Region** | **string** | Snapshot region location. | [optional] [default to null]
**Size** | **int32** | [Output Only] Snapshot size in bytes. | [optional] [default to null]
**Schedule** | **string** | Snapshot scheduling Daily, Weekly, Monthly or Manual. | [optional] [default to null]
**Share** | [**Share**](share.md) | [Output Only] If exists, this object includes the snapshot share parameters. | [optional] [default to null]
**CreationTimestamp** | **string** | [Output Only] Creation timestamp in RFC3339 text format. | [optional] [default to null]
**DeletionTime** | **string** |  | [optional] [default to null]
**Status** | **string** | [Output Only] The status of the snapshot. A snapshot can be used to mount a previous copy of the filesystem, only after the snapshot has been successfully created and the status is set to READY. Possible values are PENDING, READY. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


