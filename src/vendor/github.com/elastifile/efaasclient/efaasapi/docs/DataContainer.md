# DataContainer

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | [Output Only] The unique identifier for the resource. This identifier is defined by the server. | [optional] [default to null]
**Name** | **string** | Filesystem name | [default to null]
**Description** | **string** | Filesystem description | [default to null]
**QuotaType** | **string** | Supported values are: auto and fixed. Use auto if you have one filesystem, the size of the filesystem will be the same as the instance size. Use fixed if you have more than one filesystem, and set the filesystem size through filesystemQuota. | [default to null]
**HardQuota** | **int64** | Set the size of a filesystem if filesystemQuotaType is set to fixed. If it is set to auto, this value is ignored and quota is the instance total size. | [optional] [default to 0]
**Utilization** | [***Utilization**](utilization.md) | instance utilization metrics | [optional] [default to null]
**Exports** | [**[]Export**](export.md) |  | [optional] [default to null]
**Snapshots** | [***SnapshotSchedule**](snapshot_schedule.md) | Snapshot object | [optional] [default to null]
**Accessors** | [***Accessors**](accessors.md) | Defines the access rights to the File System. This is a list of access rights configured by the client for the file system. | [optional] [default to null]
**CreationTimestamp** | **string** | [Output Only] Creation timestamp in RFC3339 text format. | [optional] [default to null]
**UpdateTimestamp** | **string** | [Output Only] Update timestamp in RFC3339 text format. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


