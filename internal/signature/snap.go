package signature

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/url"
)

// SignSNAP generates an HMAC-SHA512 signature for a BCA SNAP API request.
//
// The request body is minified before its SHA-256 hash is calculated.
// The resulting signature is Base64-encoded.
//
// The signature is generated from the following string:
//
//	HTTPMethod:RelativeURL:AccessToken:SHA256(RequestBody):Timestamp
//
// The client secret is used as the HMAC key.
func SignSNAP(method, relativeURL, accessToken, requestBody, timestamp, clientSecret string) (string, error) {
	minifiedBody, err := minifyJSON(requestBody)
	if err != nil {
		return "", err
	}

	bodyHash := sha256.Sum256([]byte(minifiedBody))
	bodyHashHex := hex.EncodeToString(bodyHash[:])

	canonicalURL, err := canonicalizeURL(relativeURL)
	if err != nil {
		return "", err
	}

	stringToSign := method + ":" + canonicalURL + ":" + accessToken + ":" + bodyHashHex + ":" + timestamp
	mac := hmac.New(sha512.New, []byte(clientSecret))
	mac.Write([]byte(stringToSign))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func minifyJSON(body string) (string, error) {
	if body == "" {
		return "", nil
	}

	var compact bytes.Buffer

	if err := json.Compact(&compact, []byte(body)); err != nil {
		return "", err
	}

	return compact.String(), nil
}
func canonicalizeURL(relativeURL string) (string, error) {
	u, err := url.Parse(relativeURL)
	if err != nil {
		return "", err
	}

	query := u.Query()

	u.RawQuery = query.Encode()

	return u.RequestURI(), nil
}
