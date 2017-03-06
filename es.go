package main

/*
  ElasticSearch package
*/

import (
	"context"
	"fmt"
	"log"

	elastic "gopkg.in/olivere/elastic.v5"
)

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

// EnsureIndexExists creates the main index if it does not exist.
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
	}

	return nil
}

// Put upserts the SQS message into ElasticSearch's own index.
func (es *ESClient) Put(fromSQS SQSMessage) error {
	docID := esID(fromSQS.CompanyID, fromSQS.EmployeeID, fromSQS.AbsMonth)
	log.Printf("Indexing payroll with id '%s'\n", docID)

	response, err := es.client.Index().
		Index(es.index).
		Type("payroll").
		Id(docID).
		BodyJson(fromSQS.Content).
		Do(es.ctx)

	if err != nil {
		return err
	}

	if !response.Created {
		log.Printf("[WARN] payroll with id '%s' has not been created\n", docID)
	}

	return nil
}

// Flush well, flushes changes to disk.
func (es *ESClient) Flush() error {
	_, err := es.client.Flush().Index(es.index).Do(es.ctx)
	return err
}

// esID builds the elastic search ID based on provided args.
func esID(company, employee string, month int) string {
	return fmt.Sprintf("%s%s%d", company, employee, month)
}
