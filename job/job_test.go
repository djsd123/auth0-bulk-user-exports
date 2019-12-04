package job

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/stretchr/testify/assert"
	"log"
)

const accessToken = "12345abcde"

func TestCreateJob(t *testing.T) {

	const expectedJobId = "job_xyz"
	JobResult := `{
             "id": "%s"
	}`

	mockJobResponse := fmt.Sprintf(JobResult, expectedJobId)
	connectionId := "con_xyz"
	configTemplate, err := ioutil.ReadFile("../job_config.json")
	if err != nil {
		log.Fatal(err)
	}
	jobConfig := fmt.Sprintf(string(configTemplate), connectionId)

	auth0Mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(mockJobResponse))
		if err != nil {
			t.Fatalf("error occurred during write: \n%s", err)
		}
	}))
	defer auth0Mock.Close()

	job, err := CreateJob(auth0Mock.URL, accessToken, jobConfig)
	assert.Nil(t, err, err)
	if assert.NotNil(t, job) {
		assert.Equal(t, expectedJobId, *job)
		assert.Containsf(t, mockJobResponse, *job, "Error: expected response should contain %s", job)
	}
}

func TestWaitForJobCompletion(t *testing.T) {

	const expectedLocation = "https://auth0Mock/data/job_xyz.csv"

	statusCompleted := fmt.Sprintf(`{
		"status": "completed",
        "location": "%s"
    }`, expectedLocation)

	auth0Mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte(statusCompleted))
		if err != nil {
			t.Fatalf("error occurred during write: \n%s", err)
		}
	}))
	defer auth0Mock.Close()

	job := new(Job)
	location, err := WaitForJobCompletion(auth0Mock.URL, accessToken, job.Id)
	assert.Nil(t, err, err)
	if assert.NotNil(t, location) {
		assert.Equal(t, expectedLocation, *location)
	}
}

func TestGetUserExport(t *testing.T) {

	const expectedResponse = "job_xyz.csv"

	// Create mock gzip stream
	var buffer bytes.Buffer
	gz := gzip.NewWriter(&buffer)
	_, err := gz.Write([]byte(expectedResponse))
	err = gz.Close()

	auth0Mock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		_, err := w.Write(buffer.Bytes())
		if err != nil {
			t.Fatalf("error occurred during write: \n%s", err)
		}
	}))
	defer auth0Mock.Close()

	resp, err := GetUserExport(auth0Mock.URL)
	assert.Nil(t, err, err)
	defer resp.Close()

	// ioutil.ReadAll wraps bytes.Buffer's ReadFrom method, copying io.Reader contents to byte slice
	body, err := ioutil.ReadAll(resp)
	assert.Nil(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, expectedResponse, string(body))
	}
}

type mockS3Client struct {
	s3iface.S3API
}

func (m *mockS3Client) PutObject(in *s3.PutObjectInput) (out *s3.PutObjectOutput, err error) {
	return &s3.PutObjectOutput{
		VersionId: aws.String("test-version-id"),
	}, nil
}

func TestUploadUserExportToS3(t *testing.T) {

	const expectedVersionID = "test-version-id"
	mockSvc := &mockS3Client{}
	bucket := "MyS3Bucket"
	key := "userdata.csv"

	// Create mock gzip stream
	var buffer bytes.Buffer
	gz := gzip.NewWriter(&buffer)
	_, err := gz.Write([]byte("userdata.csv"))
	err = gz.Close()

	// Read gzipped data
	readBuf, err := ioutil.ReadAll(&buffer)
	reader := ioutil.NopCloser(bytes.NewReader(readBuf))

	result, err := UploadUserExportToS3(mockSvc, reader, bucket, key)
	assert.Nil(t, err, err)
	if assert.NotNil(t, result) {
		assert.Equal(t, expectedVersionID, *result.VersionId)
	}
}

func TestWriteLocalFile(t *testing.T) {

	const expectedLogMessage = "Data written to /tmp/userdata.csv"
	// Create mock gzip stream
	var buffer bytes.Buffer
	gz := gzip.NewWriter(&buffer)
	_, err := gz.Write([]byte("userdata.csv"))
	err = gz.Close()

	// Read gzipped data
	readBuf, err := ioutil.ReadAll(&buffer)
	reader := ioutil.NopCloser(bytes.NewReader(readBuf))

	result, err := WriteLocalFile(reader, "/tmp/userdata.csv")
	assert.Nil(t, err, err)
	if assert.NotNil(t, result) {
		assert.Equal(t, expectedLogMessage, *result)
	}
}
