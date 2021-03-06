/*
 * Vendor API V1
 *
 * Apps documentation
 *
 * API version: 1.0.0
 * Contact: info@replicated.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type BodyPromoteRelease struct {
	Channels     []string `json:"channels"`
	Label        string   `json:"label"`
	ReleaseNotes string   `json:"release_notes"`
	Required     bool     `json:"required"`
}
