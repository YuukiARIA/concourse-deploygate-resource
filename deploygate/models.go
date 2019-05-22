package deploygate

type Results struct {
	Name        string `json:"name"`
	PackageName string `json:"package_name"`
	OSName      string `json:"os_name"`
	Path        string `json:"path"`
	Revision    int    `json:"revision"`
	VersionCode int    `json:"version_code"`
	VersionName string `json:"version_name"`
	Message     string `json:"message"`
}

type ErrorResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Because string `json:"because"`
}

type SuccessResponse struct {
	Error   bool    `json:"error"`
	Results Results `json:"results"`
}

type Response struct {
	SuccessResponse *SuccessResponse
	ErrorResponse   *ErrorResponse
}
