/*
 * Elastifile FaaS API
 *
 * Elastifile Filesystem as a Service API
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package efaasapi

type CapacityUnits struct {
	// Increment steps of capacity units, currently a size of a node in bytes.
	UnitSize int32 `json:"unitSize"`
	// Minimum capacity units supported by this class (min # of nodes)
	Min int32 `json:"min,omitempty"`
	// Maximum capacity units supported by this class (max # of nodes)
	Max int32 `json:"max,omitempty"`
}
