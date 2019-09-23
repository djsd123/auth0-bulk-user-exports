package auth

import "os"

var BEARERTOKEN = os.Getenv("BEARER_TOKEN")
