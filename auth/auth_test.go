package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetToken(t *testing.T) {

	const expectedAccessToken = "12345abcde"
	const expectedToExpireInSecs int32 = 21600
	const clientId = "c1l2i3n4t5ID"
	const clientSecret = "876d684bbcd5b3bec8f30c6d95f86ad3a68b"
	mockData := []byte(fmt.Sprintf(`{"access_token": "%s",
	"expires_in": %d}`, expectedAccessToken, expectedToExpireInSecs))

	auth0mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write(mockData)
		if err != nil {
			t.Fatalf("error occurred during write: \n%s", err)
		}
	}))
	defer auth0mock.Close()

	accessTkn, err := GetToken(auth0mock.URL, clientId, clientSecret)
	assert.Nil(t, err)
	if assert.NotNil(t, accessTkn) {
		assert.Equal(t, expectedAccessToken, accessTkn.AccessToken)
	}
}
