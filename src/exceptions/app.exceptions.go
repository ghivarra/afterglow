package exceptions

import "time"

type AppException struct {
	Code      int    `json:"code"`
	Status    string `json:"status"`
	Message   string `json:"message"`
	ErrorData any    `json:"errorData,omitempty"`
	Reason    error  `json:"reason"`
	Timestamp string `json:"timestamp"`
}

func (e *AppException) Error() string {
	if e.Reason != nil {
		return e.Reason.Error()
	}

	return e.Message
}

func NewAppException(code int, message string, reason error, errorData any) *AppException {
	return &AppException{
		Code:      code,
		Message:   message,
		Reason:    reason,
		ErrorData: errorData,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}
