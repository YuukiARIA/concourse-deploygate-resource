package deploygate

type ResponseBase struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Because string `json:"because,omitempty"`
}

type UploadResponse struct {
	ResponseBase
	Results Results `json:"results"`
}

type GetAppsResponse struct {
	ResponseBase
	Apps []App `json:"applications"`
}

type GetAppResponse struct {
	ResponseBase
	App App `json:"application"`
}

type Results struct {
	Name        string `json:"name"`
	PackageName string `json:"package_name"`
	OSName      string `json:"os_name"`
	Path        string `json:"path"`
	Revision    int    `json:"revision"`
	VersionCode string `json:"version_code"`
	VersionName string `json:"version_name"`
	Message     string `json:"message"`
}

type App struct {
	Name            string    `json:"name"`
	PackageName     string    `json:"package_name"`
	Labels          AppLabels `json:"labels"`
	OSName          string    `json:"os_name"`
	CurrentRevision int       `json:"current_revision"`
	URL             string    `json:"url"`
	IconURL         string    `json:"icon_url"`
	Owner           Owner     `json:"owner"`
}

type AppLabels map[string]string

type Owner struct {
	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	URL         string     `json:"url"`
	Enterprise  Enterprise `json:"enterprise"`
}

type Enterprise struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	URL         string `json:"url"`
	IconURL     string `json:"icon_url"`
}
