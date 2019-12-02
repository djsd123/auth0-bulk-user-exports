# Auth0 Bulk User Exports

##### Query Auth0's API and export data in csv format to either local file system or an S3 bucket

### Prerequisites

- OIDC Client ID
- OIDC Client Secret 
- Go version [1.13](https://blog.golang.org/go1.13)

You'll need to have an auth0 [application](https://auth0.com/docs/applications) in order to retrieve these values

**Note**

This app requests temporary credentials every time it executes.  You'll need to ensure the resulting bearer tokens have the following scopes

- `create:client_grants`
- `read:client_grants`
- `read:connections`
- `read:users`

### Usage

```bash
go run main.go
```
##### Or

```bash
GOOS=linux go build -o auth0-bulk-user-export
```

```bash
./auth0-bulk-user-exports
```

##### Or

```bash
docker image build -t auth0-bulk-user-exports .
```

```bash
docker container run -it --rm -v /tmp:/tmp --env CLIENT_ID=41Tmoz00wBN1... --env CLIENT_SECRET=StdNYnUaPuv9iMEfKsiLEZ0GUTe... --name auth0-bulk-user-exports auth0-bulk-user-exports
```

#### Write data to local file: `~/Dowloads/userdata.csv`

```bash
export CLIENT_ID="41Tmoz00wBN1..."
export CLIENT_SECRET="StdNYnUaPuv9iMEfKsiLEZ0GUTe...."
export FILE_PATH="~/Dowloads/userdata.csv"
```

#### Write to S3

```bash
export CLIENT_ID="41Tmoz00wBN1..."
export CLIENT_SECRET="StdNYnUaPuv9iMEfKsiLEZ0GUTe...."
export ENV=aws
export BUCKET=my-bucket
```

__Using [SSM](https://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-parameter-store.html)__

```bash
export SSM_PATH=/my/ssm/params/path/secrets/
export ENV=aws
export BUCKET=my-bucket

```

[ssm]: https://docs.aws.amazon.com/systems-manager/latest/APIReference/API_GetParametersByPath.html

### Configuration

| Env Variable  | Default  | Description                                |
|---------------|----------|--------------------------------------------|
| `CLIENT_ID` | (**Required**) | The Client ID of the auth0 application used for this app |
| `CLIENT_SECRET` | (**Required**) | The Client Secret of the auth0 application used for this app |
| `SSM_PATH`    | | [SSM Parameter Path][ssm] to retrieve `CLIENT_ID` and `CLIENT_SECRET` values. **Required If** `CLIENT_ID` and `CLIENT_SECRET` environment variables are unset | 
| `API_URL`     | (**Required**) | Auth0 management API endpoint i.e. https://`DOMAIN`.eu.auth0.com |
| `CONNECTION_NAME` | `github` | Config param for auth0.  Which connection to target when querying the API (https://auth0.com/docs/identityproviders) |
| `ENV` | (**Write Data Locally**) | **Do not set** to write locally or set to `aws` to write data to `S3` |
| `ROLE_ARN` | | The IAM role arn this app uses to write to `S3` and fetch `ssm` parameters |
| `FILE_PATH` | `/tmp/userdata.csv` | File path when writing locally. Only works when `ENV` is **not** set |
| `BUCKET` | `auth0-bulk-user-exports` | The `S3` bucket to write to when `ENV=aws` is set.  The resulting key will be suffixed with the date i.e `userdata-22-09-2019` |
| `JOB_CONFIG_FILE` | [job_config.json](./job_config.json) | Define what [user Attributes](https://auth0.com/docs/users/references/user-profile-structure#user-profile-attributes) to export |

### Test

To run tests, `cd` to the root of this project

```bash
go test ./...
```
