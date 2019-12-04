package connection

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const accessToken = "12345abcde"

func TestGetConnection(t *testing.T) {

	const expectedConnId = "2"
	mockData := []byte(fmt.Sprintf(`[
			{"id": "1", "name": "google"},
			{"id": "%s", "name": "github"}
		]`, expectedConnId))

	auth0Mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write(mockData)
		if err != nil {
			t.Fatalf("error occurred during write: \n%s", err)
		}
	}))
	defer auth0Mock.Close()

	conn, err := GetConnection(auth0Mock.URL, accessToken, "github")
	assert.Nil(t, err)
	assert.Equal(t, expectedConnId, *conn)
}

func TestGetConnection2(t *testing.T) {

	auth0Mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer auth0Mock.Close()

	_, err := GetConnection(auth0Mock.URL, accessToken, "github")
	assert.Errorf(t, err, "expected %d got %d: \n%s", http.StatusOK, http.StatusServiceUnavailable, err)
	assert.Error(t, err, err)
}
