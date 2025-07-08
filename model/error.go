package model

import "fmt"

type DError struct {
	Code    int16  `json:"code"`
	Message string `json:"message"`
}

func (e *DError) Error() string {
	return fmt.Sprintf("{\"code\": %d, \"message\": %s}", e.Code, e.Message)
}
