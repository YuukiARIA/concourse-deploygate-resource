package models

type Source struct {
	ApiKey string `json:"api_key"`
	Owner  string `json:"owner"`
}

type Version struct {
	Platform    string `json:"platform"`
	PackageName string `json:"package_name"`
	Revision    string `json:"revision"`
}

type PutParams struct {
	File             string `json:"file"`
	Message          string `json:"message"`
	MessageFile      string `json:"message_file"`
	ReleaseNote      string `json:"release_note"`
	DistributionKey  string `json:"distribution_key"`
	DistributionName string `json:"distribution_name"`
	DisableNotify    bool   `json:"disable_notify"`
	Visibility       string `json:"visibility"`
}

type PutRequest struct {
	Source Source    `json:"source"`
	Params PutParams `json:"params"`
}

type MetadataEntry struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type PutResponse struct {
	Version  Version         `json:"version"`
	Metadata []MetadataEntry `json:"metadata"`
}
