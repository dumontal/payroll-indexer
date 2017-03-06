package main

import (
	"context"
	"fmt"
	"log"

	elastic "gopkg.in/olivere/elastic.v5"
)

/*
  ElasticSearch package
*/

// ESClient is an ElasticSearch client.
// Starting with elastic.v5, you must pass a context to execute each service
type ESClient struct {
	client *elastic.Client
	url    string
	index  string
	ctx    context.Context
}

// NewESClient provides a new ElasticSearch client.
func NewESClient(url, index string) *ESClient {
	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	info, code, err := client.Ping(url).Do(ctx)
	if err != nil {
		panic(err)
	}

	log.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	return &ESClient{
		client: client,
		url:    url,
		index:  index,
		ctx:    ctx,
	}
}

func (es *ESClient) EnsureIndexExists() error {
	exists, err := es.client.IndexExists(es.index).Do(es.ctx)
	if err != nil {
		return err
	}

	if !exists {
		// Create a new index.
		createIndex, err := es.client.CreateIndex(es.index).Do(es.ctx)
		if err != nil {
			return err
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
		fmt.Println(createIndex)
	}

	return nil
}
