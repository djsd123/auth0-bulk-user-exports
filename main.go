package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
	"time"

	session2 "github.com/aws/aws-sdk-go/aws/session"
	env "github.com/djsd123/auth0-bulk-user-exports/default_env_value"
	connection "github.com/djsd123/auth0-bulk-user-exports/get_connection"
	"github.com/djsd123/auth0-bulk-user-exports/job"
)

var (
	apiUrl         = env.SetDefaultEnvValue("API_URL", "https://alpha-analytics-moj.eu.auth0.com/api/v2")
	connectionName = env.SetDefaultEnvValue("CONNECTION_NAME", "github")
	platform       = os.Getenv("ENV")
	session        = session2.Must(session2.NewSession())
	creds          = stscreds.NewCredentials(session, env.SetDefaultEnvValue("ROLE_ARN", "arn:aws:iam::593291632749:role/restricted-admin-data"))
	s3svc          = s3.New(session, &aws.Config{Credentials: creds, Region: aws.String("eu-west-1"), CredentialsChainVerboseErrors: aws.Bool(true)})
)

func main() {

	connectionId, err := connection.GetConnection(fmt.Sprintf("%s/connections", apiUrl), connectionName)
	if err != nil {
		log.Fatalf("Error while retrieving connection id: %s", err)
	}

	jobConfig := fmt.Sprintf(`{
		"connection_id": "%s",
		"format": "csv", 
		"fields": [
			{"name": "email"}, { "name": "nickname", "export_as": "username"}
		]
	}`, *connectionId)

	bulkUserExportJob, err := job.CreateJob(fmt.Sprintf("%s/jobs/users-exports", apiUrl), jobConfig)
	if err != nil {
		log.Fatalf("Error occurred while creating job: \n%s", err)
	}

	resultLocation, err := job.WaitForJobCompletion(fmt.Sprintf("%s/jobs", apiUrl), *bulkUserExportJob)
	if err != nil {
		log.Fatal(err)
	}

	exportData, err := job.GetUserExport(*resultLocation)
	if err != nil {
		log.Fatalf("Error trying to download userdata: \n%s", err)
	}

	bucket := env.SetDefaultEnvValue("BUCKET", "auth0-userdata")
	key := fmt.Sprintf("userdata-%s.csv", time.Now().Format("02-01-2006"))

	if platform == "aws" {

		result, err := job.UploadUserExportToS3(s3svc, exportData, bucket, key)
		if err != nil {
			log.Print(fmt.Errorf("Error while uploading userdata to S3: \n%s", err))
		} else {
			log.Printf("FILE=%s successfully uploaded to BUCKET=%s, VERSIONID=%p", key, bucket, result.VersionId)
		}
	} else {
		ok, err := job.WriteLocalFile(exportData, env.SetDefaultEnvValue("FILE_PATH", "/tmp/userdata.csv"))
		if err != nil {
			log.Fatal(err)
		} else {
			log.Print(*ok)
		}
	}
}
