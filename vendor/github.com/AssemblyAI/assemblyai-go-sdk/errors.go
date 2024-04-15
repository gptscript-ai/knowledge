package assemblyai

// APIError represents an error returned by the AssemblyAI API.
type APIError struct {
	Status  int    `json:"-"`
	Message string `json:"error"`
}

// Error returns the API error message.
func (e APIError) Error() string {
	return e.Message
}
