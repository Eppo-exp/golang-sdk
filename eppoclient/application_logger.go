package eppoclient

import "regexp"

// maskSensitiveInfo replaces sensitive information (like apiKey or sdkKey)
// in the error message with 'XXXXXX' to prevent exposure of these keys in
// logs or error messages.
func maskSensitiveInfo(errMsg string) string {
	re := regexp.MustCompile(`(apiKey|sdkKey)=[^&]*`)
	return re.ReplaceAllString(errMsg, "$1=XXXXXX")
}
