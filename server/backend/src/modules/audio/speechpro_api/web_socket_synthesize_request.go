/*
 * TTS documentation
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 1.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package speechpro_api

type WebSocketSynthesizeRequest struct {

	// Text for synthesize to speech
	Text *WebSocketTextParam `json:"text"`

	// Name of name
	VoiceName string `json:"voice_name"`

	// Format of response audio
	Audio string `json:"audio"`
}
