package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type SearchDocumentOptions struct {
	query map[string]interface{}
	sort  []string
}

type SearchDocumentOption func(*SearchDocumentOptions)

func WithSearchQuery(query map[string]interface{}) SearchDocumentOption {
	return func(so *SearchDocumentOptions) {
		so.query = query
	}
}

func WithSearchSort(sort []string) SearchDocumentOption {
	return func(so *SearchDocumentOptions) {
		so.sort = sort
	}
}

type IndexDocumentOptions struct {
	timeout time.Duration
	docId   string
}

type IndexDocumentOption func(*IndexDocumentOptions)

func WithIndexTimeout(timeout time.Duration) IndexDocumentOption {
	return func(io *IndexDocumentOptions) {
		io.timeout = timeout
	}
}

func WithDocumentId(id string) IndexDocumentOption {
	return func(io *IndexDocumentOptions) {
		io.docId = id
	}
}

type UpdateDocumentOptions struct{}

type UpdateDocumentOption func(*UpdateDocumentOptions)

type GetDocumentOptions struct{}

type GetDocumentOption func(*GetDocumentOptions)

type DeleteDocumentOptions struct{}

type DeleteDocumentOption func(*DeleteDocumentOptions)

type CreateIndexOptions struct{}

type CreateIndexOption func(*CreateIndexOptions)

type DeleteIndexOptions struct{}

type DeleteIndexOption func(*DeleteIndexOptions)

type ElasticSearchClient interface {
	Search(ctx context.Context, index string, opts ...SearchDocumentOption) ([]map[string]interface{}, error)
	Index(ctx context.Context, index string, data interface{}, opts ...IndexDocumentOption) error
	Get(ctx context.Context, index string, id string, opts ...GetDocumentOption) (map[string]interface{}, error)
	Update(ctx context.Context, index string, id string, data interface{}, opts ...UpdateDocumentOption) error
	Delete(ctx context.Context, index string, id string, opts ...DeleteDocumentOption) error
	CreateIndex(ctx context.Context, index string, opts ...CreateIndexOption) error
	DeleteIndex(ctx context.Context, indexes []string, opts ...DeleteIndexOption) error
}

type elasticSearchClient struct {
	client *elasticsearch.Client
}

func (e *elasticSearchClient) Search(ctx context.Context, index string, opts ...SearchDocumentOption) ([]map[string]interface{}, error) {
	var options SearchDocumentOptions
	for _, opt := range opts {
		opt(&options)
	}

	var (
		client = e.client
		r      map[string]interface{}
	)

	var requests []func(*esapi.SearchRequest)
	requests = append(requests, client.Search.WithIndex(index))
	requests = append(requests, client.Search.WithContext(ctx))
	requests = append(requests, client.Search.WithTrackTotalHits(true))
	requests = append(requests, client.Search.WithPretty())

	if len(options.query) != 0 {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(options.query); err != nil {
			return nil, err
		}
		requests = append(requests, client.Search.WithBody(&buf))
	}

	if len(options.sort) != 0 {
		requests = append(requests, client.Search.WithSort(options.sort...))
	}

	res, err := e.client.Search(requests...)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, e.parseErrorResponse(res)
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %w", err)
	}

	var output = make([]map[string]interface{}, 0)
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		_ = hit.(map[string]interface{})["_id"].(string)
		data := hit.(map[string]interface{})["_source"].(map[string]interface{})
		output = append(output, data)
	}
	return output, nil
}

func (e *elasticSearchClient) Index(ctx context.Context, index string, data interface{}, opts ...IndexDocumentOption) error {
	var (
		options IndexDocumentOptions
		client  = e.client
	)

	for _, opt := range opts {
		opt(&options)
	}

	dataByte, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("elasticsearch error: %w", err)
	}

	var requests []func(*esapi.IndexRequest)
	requests = append(requests, client.Index.WithContext(ctx))
	if options.timeout > 0 {
		requests = append(requests, client.Index.WithTimeout(options.timeout))
	}
	if len(options.docId) != 0 {
		requests = append(requests, client.Index.WithDocumentID(options.docId))
	}

	res, err := client.Index(index, bytes.NewReader(dataByte), requests...)
	if err != nil {
		return err
	}

	if res.IsError() {
		return e.parseErrorResponse(res)
	}

	return nil
}

func (e *elasticSearchClient) Get(ctx context.Context, index string, id string, opts ...GetDocumentOption) (map[string]interface{}, error) {
	var options GetDocumentOptions
	for _, opt := range opts {
		opt(&options)
	}

	var r map[string]interface{}

	res, err := e.client.Get(index, id)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, e.parseErrorResponse(res)
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %w", err)
	}

	return r["_source"].(map[string]interface{}), nil
}

func (e *elasticSearchClient) Update(ctx context.Context, index string, id string, data interface{}, opts ...UpdateDocumentOption) error {
	var options UpdateDocumentOptions
	for _, opt := range opts {
		opt(&options)
	}

	dataByte, err := json.Marshal(map[string]interface{}{
		"doc": data,
	})
	if err != nil {
		return fmt.Errorf("elasticsearch error: %w", err)
	}

	res, err := e.client.Update(index, id, bytes.NewReader(dataByte))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return e.parseErrorResponse(res)
	}

	return nil
}

func (e *elasticSearchClient) Delete(ctx context.Context, index string, id string, opts ...DeleteDocumentOption) error {
	var options DeleteDocumentOptions
	for _, opt := range opts {
		opt(&options)
	}

	res, err := e.client.Delete(index, id)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return e.parseErrorResponse(res)
	}

	return nil
}

func (e *elasticSearchClient) CreateIndex(ctx context.Context, index string, opts ...CreateIndexOption) error {
	var options CreateIndexOptions
	for _, opt := range opts {
		opt(&options)
	}

	res, err := e.client.Indices.Create(index)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return e.parseErrorResponse(res)
	}

	return nil
}

func (e *elasticSearchClient) DeleteIndex(ctx context.Context, indexes []string, opts ...DeleteIndexOption) error {
	var options DeleteIndexOptions
	for _, opt := range opts {
		opt(&options)
	}

	res, err := e.client.Indices.Delete(indexes)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.IsError() {
		return e.parseErrorResponse(res)
	}

	return nil
}

type EsClientOpions struct {
	addresses []string
	username  string
	password  string
}

type EsClientOption func(*EsClientOpions)

func WithEsAddresses(addreses []string) EsClientOption {
	return func(eco *EsClientOpions) {
		eco.addresses = addreses
	}
}

func WithEsUsername(username string) EsClientOption {
	return func(eco *EsClientOpions) {
		eco.username = username
	}
}

func WithEsPassword(password string) EsClientOption {
	return func(eco *EsClientOpions) {
		eco.password = password
	}
}

func NewElasticSearchClient(opts ...EsClientOption) (ElasticSearchClient, error) {
	var options EsClientOpions
	for _, opt := range opts {
		opt(&options)
	}

	cfg := elasticsearch.Config{}

	if options.addresses != nil {
		cfg.Addresses = options.addresses
	}

	if len(options.username) != 0 {
		cfg.Username = options.username
	}

	if len(options.password) != 0 {
		cfg.Password = options.password
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &elasticSearchClient{
		client: client,
	}, nil
}

func (c *elasticSearchClient) parseErrorResponse(res *esapi.Response) error {
	var e map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
		return err
	}

	esErr := ElasticSearchError{
		Status: res.Status(),
	}

	if e["error"] != nil {
		esErr.Type = e["error"].(map[string]interface{})["type"].(string)
		esErr.Reason = e["error"].(map[string]interface{})["reason"].(string)
	}

	return esErr
}
