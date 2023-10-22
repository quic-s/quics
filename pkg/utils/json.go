package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func UnmarshalRequestBody(r *http.Request, dstStruct any) error {
	buf := make([]byte, r.ContentLength)

	n, err := r.Body.Read(buf)
	if err != nil && err != io.EOF {
		return err
	}

	if n == 0 {
		return fmt.Errorf("empty body")
	}

	err = json.Unmarshal(buf, dstStruct)
	if err != nil {
		return err
	}

	return nil
}
