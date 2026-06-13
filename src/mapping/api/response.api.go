package api

// interface for single content responses
type Response[DataType any, ErrorType any] struct {
	Status    string    `json:"status"` // HTTP STATUS TEXT => error / success
	Message   string    `json:"message"`
	Data      DataType  `json:"data"`
	Errors    ErrorType `json:"errors"`
	Timestamp *string   `json:"timestamp,omitempty"`
}

// interface for index-type data
type Index struct {
	Total    IndexTotal `json:"total"`
	Contents any        `json:"contents"`
}

// total data struct
type IndexTotal struct {
	FilteredData   int64 `json:"filteredData"`
	UnfilteredData int64 `json:"unfilteredData"`
}
