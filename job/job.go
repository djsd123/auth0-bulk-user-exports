package job

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/djsd123/auth0-bulk-user-exports/auth"
)

var (
	BEARER_TOKEN = auth.BEARERTOKEN
)

type Job struct {
	Id string `json:"id"`
}

type Status struct {
	Status   string `json:"status"`
	Location string `json:"location"` // Location is the unauthenticated temporary (60 secs) url used to download data
}

func CreateJob(url, jobConfig string) (*Job, error) {

	JobResponse := Job{}

	payload := strings.NewReader(jobConfig)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", BEARER_TOKEN))
	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("expected %d, got %d", http.StatusOK, res.StatusCode)
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	JobResponse, err = handleData(body)
	if err != nil {
		return nil, err
	}

	log.Printf("Created job %s", JobResponse.Id)

	return &JobResponse, nil
}

func handleData(body []byte) (Job, error) {
	var job = Job{}
	err := json.Unmarshal(body, &job)

	return job, err
}

func WaitForJobCompletion(url string, job Job) (resultURL *string, err error) {

	jobStatusUrl := fmt.Sprintf("%s/%s", url, job.Id)

	var status Status
	var body []byte

	for {
		body, err = requestAndHandleData("GET", jobStatusUrl, nil)
		if err != nil {
			return nil, err
		}

		err := json.Unmarshal(body, &status)
		if err != nil {
			return nil, err
		}

		log.Printf("Job: %s, Status: %s", job.Id, status.Status)
		if status.Status == "completed" {

			return &status.Location, nil
		}
		time.Sleep(5 * time.Second)
	}
}

func GetUserExport(resultURL string) (io.ReadCloser, error) {

	client := http.Client{Timeout: 1 * time.Minute}

	req, err := http.NewRequest("GET", resultURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept-Encoding", "gzip")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("expected %d, got %d", http.StatusOK, resp.StatusCode)
		return nil, err
	}

	return resp.Body, nil
}

func UploadUserExportToS3(awsClient s3iface.S3API, readCloser io.ReadCloser, bucket, key string) (*s3.PutObjectOutput, error) {

	reader, err := gzip.NewReader(readCloser)
	if err != nil {
		return nil, fmt.Errorf("error during Gunzip operation while uploading %s:\n%s", key, err)
	}

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	uploadInput := &s3.PutObjectInput{
		ACL:    aws.String("authenticated-read"),
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   aws.ReadSeekCloser(bytes.NewReader(body)),
	}

	defer readCloser.Close()

	result, err := awsClient.PutObject(uploadInput)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func requestAndHandleData(method string, url string, payload io.Reader) ([]byte, error) {

	req, _ := http.NewRequest(method, url, payload)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", BEARER_TOKEN))
	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("expected 200 got %d: %s", res.StatusCode, body)
		return nil, err
	}

	return body, err
}

func WriteLocalFile(data io.ReadCloser, filePath string) (*string, error) {

	reader, err := gzip.NewReader(data)
	if err != nil {
		return nil, fmt.Errorf("error during Gunzip operation while writing %s:\n%s", filePath, err)
	}

	defer reader.Close()

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(filePath, body, 0644)
	if err != nil {
		return nil, fmt.Errorf("Error writing file: %s \n%s", filePath, err)
	}

	success := fmt.Sprintf("Data written to %s", filePath)

	return &success, nil
}
