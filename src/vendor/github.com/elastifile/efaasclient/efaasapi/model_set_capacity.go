/*
 * Elastifile FaaS API
 *
 * Elastifile Filesystem as a Service API
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package efaasapi

type SetCapacity struct {
	// The number of storage capacity units provisioned
	ProvisionedCapacityUnits float32 `json:"provisionedCapacityUnits"`
	// The unit used for capacity, possible values are: Steps, Bytes.  Default value is Steps.
	CapacityUnitType string `json:"capacityUnitType,omitempty"`
}
