package exceptions

import "time"

type AppException struct {
	Code      int    `json:"code"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	Errors    any    `json:"errors,omitempty"`
	Reason    error  `json:"reason"`
	Timestamp string `json:"timestamp"`
}

func (e *AppException) Error() string {
	if e.Reason != nil {
		return e.Reason.Error()
	}

	return e.Message
}

func NewAppException(code int, message string, errors any) *AppException {
	return &AppException{
		Code:      code,
		Message:   message,
		Errors:    errors,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
