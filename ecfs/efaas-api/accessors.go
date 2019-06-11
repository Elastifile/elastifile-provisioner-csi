/* 
 * Elastifile FaaS API
 *
 * Elastifile Filesystem as a Service API
 *
 * OpenAPI spec version: 2.0
 * 
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 */

package EfaasApi

type Accessors struct {

	Items []AccessorItems `json:"items,omitempty"`

	// Specifies a fingerprint for this request, which is essentially a hash of the accessors contents and used for optimistic locking. The fingerprint is initially generated by server and changes after every request to modify or update accessors. You must always provide an up-to-date fingerprint hash in order to update or change accessors.  To see the latest fingerprint, make get() request to the instance.   A base64-encoded string.
	Fingerprint string `json:"fingerprint,omitempty"`
}