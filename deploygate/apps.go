package deploygate

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetApps(organization string) (*GetAppsResponse, error) {
	url := fmt.Sprintf("https://deploygate.com/api/organizations/%s/apps", organization)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.ApiKey)

	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resObj := new(GetAppsResponse)
	if err := json.NewDecoder(res.Body).Decode(resObj); err != nil {
		return nil, err
	}
	return resObj, nil
}

func (c *Client) GetApp(organization, platform, packageName string) (*GetAppResponse, error) {
	url := fmt.Sprintf("https://deploygate.com/api/organizations/%s/platforms/%s/apps/%s", organization, platform, packageName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", c.ApiKey)

	httpClient := http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resObj := new(GetAppResponse)
	if err := json.NewDecoder(res.Body).Decode(resObj); err != nil {
		return nil, err
	}
	return resObj, nil
}
