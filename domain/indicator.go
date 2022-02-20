package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Indicator struct {
	Id           string    `json:"id" db:"id"`
	Title        string    `json:"title" db:"title"`
	Body         string    `json:"body" db:"body"`
	Score        int64     `json:"score" db:"score"`
	CreatedDate  time.Time `json:"createdDate" db:"createdDate"`
	AccessedDate time.Time `json:"accessedDate" db:"accessedDate"`
	Link         string    `json:"link" db:"link"`
	Source       string    `json:"source" db:"source"`
	SourceId     string    `json:"sourceId" db:"sourceId"`
	References   []string  `json:"references" db:"references"`
	Mentions     []Mention `json:"mentions,omitempty" db:"mentions"`
	Tags         []string  `json:"tags,omitempty" db:"tags"`
}

func (indicator *Indicator) String() string {
	return fmt.Sprintf("[%s] %s %s (%d)", indicator.CreatedDate.UTC().Format(time.RFC3339), indicator.Id, indicator.Title, indicator.Score)
}

func (indicator *Indicator) AddMention(mention Mention) {
	for _, existingMention := range indicator.Mentions {
		if existingMention.Equal(mention) {
			return
		}
	}
	indicator.Mentions = append(indicator.Mentions, mention)
}

func (indicator *Indicator) AddTag(tag string) {
	for _, existingTag := range indicator.Tags {
		if tag == existingTag {
			return
		}
	}
	indicator.Tags = append(indicator.Tags, tag)
}

func (indicator *Indicator) AddReference(reference string) {
	for _, existingReference := range indicator.References {
		if reference == existingReference {
			return
		}
	}
	indicator.References = append(indicator.References, reference)
}

func (indicator *Indicator) Map() map[string]interface{} {

	mentions := make(map[string][]string)
	for _, mention := range indicator.Mentions {
		if _, ok := mentions[mention.TopicName]; !ok {
			mentions[mention.TopicName] = []string{}
		}
		mentions[mention.TopicName] = append(mentions[mention.TopicName], mention.Mention)
	}

	outMap := make(map[string]interface{})
	outMap["id"] = indicator.Id
	outMap["title"] = indicator.Title
	outMap["body"] = indicator.Body
	outMap["score"] = indicator.Score
	outMap["createdDate"] = indicator.CreatedDate.Format(time.RFC3339)
	outMap["accessedDate"] = indicator.AccessedDate.Format(time.RFC3339)
	outMap["link"] = indicator.Link
	outMap["source"] = indicator.Source
	outMap["sourceId"] = indicator.SourceId
	outMap["references"] = indicator.References
	outMap["mentions"] = mentions
	outMap["tags"] = indicator.Tags
	return outMap
}

type IndicatorCollection struct {
	Indicators []Indicator
}

func (icollection *IndicatorCollection) CumulativeScore() (cumulativeScore int64) {
	cumulativeScore = 0
	for _, indicator := range icollection.Indicators {
		cumulativeScore += indicator.Score
	}
	return cumulativeScore
}

func (icollection *IndicatorCollection) Length() int {
	return len(icollection.Indicators)
}

func (icollection *IndicatorCollection) AverageScore() float64 {
	return float64(icollection.CumulativeScore()) / float64(icollection.Length())
}

func (icollection *IndicatorCollection) Append(i Indicator) {
	icollection.Indicators = append(icollection.Indicators, i)
}

func (icollection *IndicatorCollection) Extend(inCollection *IndicatorCollection) {
	icollection.Indicators = append(icollection.Indicators, inCollection.Indicators...)
}

func (icollection *IndicatorCollection) IsEmpty() bool {
	return len(icollection.Indicators) == 0
}

func (icollection *IndicatorCollection) First() (Indicator, error) {
	if icollection.IsEmpty() {
		return Indicator{}, fmt.Errorf("indicator collection empty")
	}
	return icollection.Indicators[0], nil
}

func (icollection *IndicatorCollection) Last() (Indicator, error) {
	if icollection.IsEmpty() {
		return Indicator{}, fmt.Errorf("indicator collection empty")
	}
	return icollection.Indicators[len(icollection.Indicators)-1], nil
}

type IndicatorFactory interface {
	Config() IndicatorFactoryConfig
	SetConfig(IndicatorFactoryConfig) error
	NewIndicator() (Indicator, error)
	MustNewIndicator() Indicator
	UnmarshallFromMap(indicatorMap map[string]interface{}) (*Indicator, error)
	NewIndicatorCollection() (IndicatorCollection, error)
	MustNewIndicatorCollection() IndicatorCollection
}

type IndicatorFactoryConfig struct {
	Source string
}

type DefaultIndicatorFactory struct {
	factoryConfig IndicatorFactoryConfig
}

func (f DefaultIndicatorFactory) Config() IndicatorFactoryConfig {
	return f.factoryConfig
}

func (f *DefaultIndicatorFactory) SetConfig(config IndicatorFactoryConfig) error {
	f.factoryConfig = config
	return nil
}

func (f DefaultIndicatorFactory) IsZero() bool {
	return f == (DefaultIndicatorFactory{})
}

func (f DefaultIndicatorFactory) NewIndicator() (Indicator, error) {
	return Indicator{
		Id:         uuid.New().String(),
		Source:     f.factoryConfig.Source,
		References: []string{},
		Mentions:   []Mention{},
		Tags:       []string{},
	}, nil
}

func (f DefaultIndicatorFactory) MustNewIndicator() Indicator {
	indicator, err := f.NewIndicator()
	if err != nil {
		panic(err)
	}
	return indicator
}

func (f DefaultIndicatorFactory) UnmarshallFromMap(indicatorMap map[string]interface{}) (*Indicator, error) {
	parsedIndicator := Indicator{}
	parsedIndicator.Id = indicatorMap["id"].(string)
	parsedIndicator.Title = indicatorMap["title"].(string)
	parsedIndicator.Body = indicatorMap["body"].(string)
	scoreFloat := indicatorMap["score"].(float64)
	parsedIndicator.Score = int64(scoreFloat)
	parsedIndicator.Link = indicatorMap["link"].(string)
	parsedIndicator.Source = indicatorMap["source"].(string)
	parsedIndicator.SourceId = indicatorMap["sourceId"].(string)
	parsedIndicator.References = indicatorMap["references"].([]string)
	parsedIndicator.Mentions = indicatorMap["mentions"].([]Mention)
	parsedIndicator.Tags = indicatorMap["tags"].([]string)

	createdDate, err := time.Parse("2006-01-02 15:04:05.000", indicatorMap["createdDate"].(string))
	if err != nil {
		return nil, err
	}
	parsedIndicator.CreatedDate = createdDate

	accessedDate, err := time.Parse("2006-01-02 15:04:05.000", indicatorMap["accessedDate"].(string))
	if err != nil {
		return nil, err
	}
	parsedIndicator.AccessedDate = accessedDate

	return &parsedIndicator, nil
}

func (f DefaultIndicatorFactory) NewIndicatorCollection() (IndicatorCollection, error) {
	return IndicatorCollection{
		Indicators: []Indicator{},
	}, nil
}

func (f DefaultIndicatorFactory) MustNewIndicatorCollection() IndicatorCollection {
	collection, err := f.NewIndicatorCollection()
	if err != nil {
		panic(err)
	}
	return collection
}

func NewIndicatorFactory(ifc IndicatorFactoryConfig) (IndicatorFactory, error) {
	return &DefaultIndicatorFactory{factoryConfig: ifc}, nil
}

type IndicatorRepository interface {
	Store(ctx context.Context, indicator Indicator) error
	GetById(ctx context.Context, id string) (Indicator, error)
	GetByTriggerName(ctx context.Context, triggerName string) (*IndicatorCollection, error)
	GetByMatch(ctx context.Context, match string) (*IndicatorCollection, error)
	GetBySource(ctx context.Context, source string) (*IndicatorCollection, error)
	GetBySourceId(ctx context.Context, sourceId string) (Indicator, error)
	GetLatest(ctx context.Context) (Indicator, error)
	GetBetween(ctx context.Context, start time.Time, end time.Time) (*IndicatorCollection, error)
}
