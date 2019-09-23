package get_connection

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/djsd123/auth0-bulk-user-exports/auth"
)

var (
	BEARER_TOKEN = auth.BEARERTOKEN
)

type Connection struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GetConnection(url, connectionName string) (*string, error) {

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", BEARER_TOKEN))
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		err := fmt.Errorf("expected %d got %d: \n%s", http.StatusOK, res.StatusCode, err)
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	conns, err := handleData(body)
	if err != nil {
		return nil, err
	}

	return findConnectionByName(conns, connectionName), nil
}

func findConnectionByName(ConnectionSlice []Connection, ConnectionName string) (connectionId *string) {

	for _, c := range ConnectionSlice {
		if c.Name == ConnectionName {
			connectionId = &c.Id
			break
		}
	}
	return
}

func handleData(body []byte) ([]Connection, error) {

	var connections = []Connection{}
	err := json.Unmarshal(body, &connections)

	return connections, err
}
