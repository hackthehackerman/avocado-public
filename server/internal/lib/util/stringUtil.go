package util

import (
	"bytes"
	"encoding/json"
)

func SafeSubtring(s string, maxLength int) (ss string) {
	return s[0:Min(len(s), maxLength)]
}

func PrettyString(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}
