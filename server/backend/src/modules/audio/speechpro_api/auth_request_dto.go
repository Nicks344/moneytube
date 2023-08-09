/*
 * SessionService documentation
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package speechpro_api

// User information provided on authentication
type AuthRequestDto struct {

	// User name
	Username string `json:"username"`

	// User domain
	DomainId int64 `json:"domain_id"`

	// User password - planed text
	Password string `json:"password"`
}
