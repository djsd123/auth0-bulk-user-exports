package connection

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Connection struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GetConnection(url, authToken, connectionName string) (*string, error) {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", authToken))
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected %d got %d: \n%s", http.StatusOK, res.StatusCode, res.Body)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	connections := []Connection{}
	err = json.Unmarshal(body, &connections)
	if err != nil {
		return nil, err
	}

	return findConnectionByName(connections, connectionName), nil
}

func findConnectionByName(ConnectionSlice []Connection, ConnectionName string) (connectionId *string) {

	for _, c := range ConnectionSlice {
		if c.Name == ConnectionName {
			return &c.Id
		}
	}
	return
}
