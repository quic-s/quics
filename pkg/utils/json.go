package utils

import (
	"encoding/json"
	"fmt"
)

func UnmarshalRequestBody(body []byte, dstStruct any) error {
	if len(body) == 0 {
		return fmt.Errorf("empty body")
	}

	err := json.Unmarshal(body, dstStruct)
	if err != nil {
		return err
	}

	return nil
}
