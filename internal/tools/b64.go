package tools

import "encoding/base64"

func B64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func B64Decode(s string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
