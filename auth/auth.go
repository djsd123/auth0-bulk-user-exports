package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type AccessToken struct {
	AccessToken string `json:"access_token"` // Valid for 6hrs
}

func GetToken(apiUrl, clientId, clientSecret string) (tkn *AccessToken, err error) {

	audience := fmt.Sprintf("%s/api/v2/", apiUrl)

	url := fmt.Sprintf("%s/oauth/token", apiUrl)

	payload := strings.NewReader(fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s&audience=%s", clientId, clientSecret, audience))

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &tkn)
	if err != nil {
		return nil, err
	}

	return tkn, nil
}
