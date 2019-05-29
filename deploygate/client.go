package deploygate

type Client struct {
	ApiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{ApiKey: apiKey}
}
