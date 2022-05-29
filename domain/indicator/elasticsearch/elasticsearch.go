package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/carbonrook/cvewatch-domain/domain/indicator"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
)

type ElasticsearchIndicatorRepository struct {
	Client           *elasticsearch.Client
	IndexName        string
	IndicatorFactory indicator.IndicatorFactory
}

func (eir ElasticsearchIndicatorRepository) Add(ctx context.Context, indicator indicator.Indicator) error {
	body, err := json.Marshal(indicator)
	if err != nil {
		return err
	}
	request := esapi.IndexRequest{
		Index:      eir.IndexName,
		DocumentID: indicator.Id,
		Body:       strings.NewReader(string(body)),
		Refresh:    "true",
	}

	res, err := request.Do(ctx, eir.Client)
	if err != nil {
		return fmt.Errorf("failed to get response from cluster: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to index document %s: %s", indicator.Id, err)
	}
	return nil
}

type elasticSearchResponse struct {
	Took     int  `json:"took"`
	TimedOut bool `json:"timed_out"`
	Shards   struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Skipped    int `json:"skipped"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string              `json:"_index"`
			Type   string              `json:"_type"`
			Id     string              `json:"_id"`
			Score  float64             `json:"_score"`
			Source indicator.Indicator `json:"_source"`
		}
	} `json:"hits"`
}

func (eir ElasticsearchIndicatorRepository) searchWithQuery(ctx context.Context, query bytes.Buffer) ([]indicator.Indicator, error) {
	res, err := eir.Client.Search(
		eir.Client.Search.WithContext(ctx),
		eir.Client.Search.WithIndex(eir.IndexName),
		eir.Client.Search.WithBody(&query),
		eir.Client.Search.WithTrackTotalHits(true),
		eir.Client.Search.WithPretty(),
	)
	if err != nil {
		return []indicator.Indicator{}, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return []indicator.Indicator{}, fmt.Errorf("failed to retrieve document")
	}

	var r elasticSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return []indicator.Indicator{}, fmt.Errorf("error parsing the response: %s", err)
	}

	indicatorCollection := eir.IndicatorFactory.MustNewIndicatorCollection()
	for _, hit := range r.Hits.Hits {
		indicatorCollection.Append(hit.Source)
	}

	return indicatorCollection.Indicators, nil
}

type elasticGetResponse struct {
	Index          string              `json:"_index"`
	Id             string              `json:"_id"`
	Version        int                 `json:"_version"`
	SequenceNumber int                 `json:"_seq_no"`
	PrimaryTerm    int                 `json:"_primary_term"`
	Found          bool                `json:"found"`
	Source         indicator.Indicator `json:"_source"`
}

func (eir ElasticsearchIndicatorRepository) GetById(ctx context.Context, id string) (*indicator.Indicator, error) {

	res, err := eir.Client.Get(
		eir.IndexName,
		id,
		eir.Client.Get.WithContext(ctx),
		eir.Client.Get.WithPretty(),
		eir.Client.Get.WithSource(),
	)
	if err != nil {
		return &indicator.Indicator{}, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return &indicator.Indicator{}, fmt.Errorf("failed to retrieve document")
	}

	var r elasticGetResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return &indicator.Indicator{}, fmt.Errorf("error parsing the response: %s", err)
	}

	if !r.Found {
		return &indicator.Indicator{}, indicator.ErrIndicatorNotFound
	}

	retrievedIndicator := &r.Source

	return retrievedIndicator, nil
}

func (eir ElasticsearchIndicatorRepository) GetByLink(ctx context.Context, link string) (*indicator.Indicator, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"constant_score": map[string]interface{}{
				"filter": map[string]interface{}{
					"term": map[string]interface{}{
						"link": link,
					},
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	matches, err := eir.searchWithQuery(ctx, buf)
	if err != nil {
		return nil, err
	}

	if len(matches) > 0 {
		return &matches[len(matches)-1], nil
	}

	return nil, nil
}

func (eir ElasticsearchIndicatorRepository) GetByMention(ctx context.Context, mention indicator.Mention) ([]indicator.Indicator, error) {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{"term": map[string]interface{}{"mentions.topicName": mention.TopicName}},
					{"term": map[string]interface{}{"mentions.mention": mention.Mention}},
				},
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	matches, err := eir.searchWithQuery(ctx, buf)
	if err != nil {
		return nil, err
	}

	return matches, nil
}

type elasticsearchErrorResponse struct {
	Error struct {
		RootCause []struct {
			Type   string `json:"type"`
			Reason string `json:"reason"`
		} `json:"root_cause"`
	} `json:"error"`
	Status int `json:"status"`
}

func NewElasticsearchIndicatorRepository(config elasticsearch.Config, indexName string) (indicator.IndicatorRepository, error) {
	client, err := elasticsearch.NewClient(config)
	if err != nil {
		return ElasticsearchIndicatorRepository{}, err
	}

	res, err := client.API.Indices.Get([]string{indexName})
	if err != nil {
		return ElasticsearchIndicatorRepository{}, fmt.Errorf("failed to connect to cluster: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return ElasticsearchIndicatorRepository{}, fmt.Errorf("received error response from cluster: %s", res.String())
	}

	fmt.Println("successfully connected to cluster")
	return ElasticsearchIndicatorRepository{
		Client:           client,
		IndexName:        indexName,
		IndicatorFactory: indicator.IndicatorFactory{},
	}, nil
}

func MustNewElasticsearchIndicatorRepository(config elasticsearch.Config, indexName string) indicator.IndicatorRepository {
	repo, err := NewElasticsearchIndicatorRepository(config, indexName)
	if err != nil {
		panic(err)
	}
	return repo
}
