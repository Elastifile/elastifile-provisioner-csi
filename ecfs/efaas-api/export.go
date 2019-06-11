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

type Export struct {

	// [Output Only] The unique identifier for the resource. This identifier is defined by the server.
	Id string `json:"id,omitempty"`

	// Filesystem name
	Name string `json:"name,omitempty"`

	// Export path
	Path string `json:"path,omitempty"`

	// The NFS service mount point to use for accessing the filesystem.
	NfsMountPoint string `json:"nfsMountPoint,omitempty"`
}