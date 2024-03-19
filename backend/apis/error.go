package apis

import "dam/enums"

type ErrorResponse struct {
	Message string      `json:"message"`
	Code    enums.Error `json:"code"`
}
