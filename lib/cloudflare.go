package lib

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
)

var (
	client *cloudflare.API
)

func Query(query string) ([]cloudflare.D1Result, error) {
	queryResults, err := client.QueryD1Database(context.TODO(),
		cloudflare.AccountIdentifier(getEnv("CF_ACCOUNT_ID")),
		cloudflare.QueryD1DatabaseParams{
			DatabaseID: getEnv("D1_ID"),
			SQL:        query,
		},
	)

	if err != nil {
		return nil, err
	}

	return queryResults, nil
}

func SetupCloudflareClient() {
	apiKey := getEnv("CF_API_KEY")
	accountIdentifier := getEnv("CF_ACCOUNT_ID")

	var err error
	client, err = cloudflare.NewWithAPIToken(apiKey)

	if err != nil {
		panic(err)
	}

	println("Cloudflare client initialized")
	database, err := client.GetD1Database(
		context.Background(),
		cloudflare.AccountIdentifier(accountIdentifier),
		getEnv("D1_ID"),
	)
	if err != nil {
		panic(err)
	}

	print(fmt.Sprintf("Database: %+v \n", database))

	queryResults, err := Query(`
		CREATE TABLE IF NOT EXISTS activity (
		    			id STRING PRIMARY KEY,
		    			snapshot_time TIMESTAMP,
		    			mouse_activity OBJECT,
		    			keyboard_presses INTEGER,
		                window_activity OBJECT
				)
	`)

	if err != nil {
		panic(err)
	}

	println(fmt.Sprintf("Query results: %+v \n", queryResults))
}
