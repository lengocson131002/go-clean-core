package es

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func getElasticClient() ElasticSearchClient {
	c, err := NewElasticSearchClient(
		WithEsAddresses([]string{"http://localhost:9200"}),
	)
	if err != nil {
		panic(err)
	}

	return c
}

func TestCreateIndex(t *testing.T) {
	var (
		c = getElasticClient()
	)
	err := c.CreateIndex(context.TODO(), "test.index.1")
	if err != nil {
		t.Error(err)
	}

	err = c.DeleteIndex(context.TODO(), []string{"test.index.1"})
	if err != nil {
		t.Error(err)
	}
}

type Data struct {
	Id   string
	Name string
}

func TestCreateDocument(t *testing.T) {
	var (
		c     = getElasticClient()
		index = "test.index.3"
		ctx   = context.Background()
	)

	data := Data{
		Id:   uuid.New().String(),
		Name: "test",
	}

	err := c.CreateIndex(ctx, index)
	if err != nil {
		t.Error(err)
	}

	err = c.Index(ctx, index, data, WithDocumentId(data.Id))
	if err != nil {
		t.Error(err)
	}

	res, err := c.Get(ctx, index, data.Id)
	if err != nil {
		t.Error(err)
	}

	if res == nil {
		t.Error("Empty get response")
	}

	data.Name = "updated"
	err = c.Update(ctx, index, data.Id, data)
	if err != nil {
		t.Error(err)
	}

	err = c.Delete(ctx, index, data.Id)
	if err != nil {
		t.Error(err)
	}

	err = c.DeleteIndex(context.TODO(), []string{"test.index.3"})
	if err != nil {
		t.Error(err)
	}

}
