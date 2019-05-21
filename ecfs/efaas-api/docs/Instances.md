# Instances

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | [Output Only] The unique identifier for the resource. This identifier is defined by the server. | [optional] [default to null]
**Name** | **string** | The name of the resource, provided by the client when initially creating the resource. The resource name must be 1-63 characters long, and comply with RFC1035 | [default to null]
**Description** | **string** | [Output Only] A textual description of the resource. | [optional] [default to null]
**ServiceClass** | **string** | ServiceClass name | [default to null]
**ServiceClassId** | **string** | ServiceClass resource unique id | [optional] [default to null]
**ServiceClassDescription** | **string** | ServiceClass descriptive name | [optional] [default to null]
**ProvisionedCapacityUnits** | **float32** | The number of storage capacity units provisioned | [default to null]
**CapacityUnitType** | **string** | The unit used for capacity, possible values are: Steps, Bytes.  Default value is Steps. | [optional] [default to null]
**AllocatedCapacity** | **int64** | The allocated capacity in bytes | [optional] [default to null]
**Region** | **string** | Region name for this request, required if serviceClass.serviceProtection.protectionMode is set to &#39;multi&#39;. | [optional] [default to null]
**Zone** | **string** | Zone name for this request, required if serviceClass.serviceProtection.protectionMode is set to &#39;single&#39;. | [optional] [default to null]
**Network** | **string** | Name of your VPC network connected with service producer network. | [default to null]
**NetworkProject** | **string** | The host project id if using a shared VPC network. | [optional] [default to null]
**Filesystems** | [**[]DataContainer**](data_container.md) |  | [optional] [default to null]
**Status** | **string** | [Output Only] The status of the operation, which can be one of the following: PENDING, RUNNING, or DONE. | [optional] [default to null]
**StatusMessage** | **string** | [Output Only] An optional textual description of the current status of the operation. | [optional] [default to null]
**CreationTimestamp** | **string** | [Output Only] Creation timestamp in RFC3339 text format. | [optional] [default to null]
**UpdateTimestamp** | **string** | [Output Only] Update timestamp in RFC3339 text format. | [optional] [default to null]
**Utilization** | [**Utilization**](utilization.md) | instance utilization metrics | [optional] [default to null]
**ServiceHealth** | [**ServiceHealth**](serviceHealth.md) | instance health metrics | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


