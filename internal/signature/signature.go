package signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Sign generates an HMAC-SHA256 signature for a BCA API request.
//
// The signature is generated using the API secret as the HMAC key and a
// colon-separated string containing the HTTP method, relative URL, access
// token, SHA-256 hash of the request body, and timestamp.
//
// See the BCA Developer API documentation for more information:
// https://developer.bca.co.id/Dokumentasi#signature
func Sign(method, relativeURL, accessToken, requestBody, timestamp, apiSecret string) string {
	stringToSign := buildStringToSign(method, relativeURL, accessToken, requestBody, timestamp)

	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(stringToSign))

	return hex.EncodeToString(mac.Sum(nil))
}

func buildStringToSign(method, relativeURL, accessToken, requestBody, timestamp string) string {
	bodyHash := sha256.Sum256([]byte(requestBody))
	bodyHashHex := hex.EncodeToString(bodyHash[:])

	return method + ":" + relativeURL + ":" + accessToken + ":" + bodyHashHex + ":" + timestamp
}
