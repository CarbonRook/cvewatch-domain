package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Indicator struct {
	Id           string                   `json:"id" db:"id"`
	Title        string                   `json:"title" db:"title"`
	Score        int                      `json:"score" db:"score"`
	CreatedDate  time.Time                `json:"createdDate" db:"createdDate"`
	AccessedDate time.Time                `json:"accessedDate" db:"accessedDate"`
	Link         string                   `json:"link" db:"link"`
	Source       string                   `json:"source" db:"source"`
	SourceId     string                   `json:"sourceId" db:"sourceId"`
	References   []string                 `json:"references" db:"references"`
	TriggeredOn  []TriggerMatchCollection `json:"triggeredOn,omitempty" db:"triggeredOn"`
	Tags         []string                 `json:"tags,omitempty" db:"tags"`
}

func (indicator *Indicator) String() string {
	return fmt.Sprintf("[%s] %s %s (%d)", indicator.CreatedDate.UTC().Format(time.RFC3339), indicator.Id, indicator.Title, indicator.Score)
}

func (indicator *Indicator) AddTriggerMatchCollection(trigger TriggerMatchCollection) {
	indicator.TriggeredOn = append(indicator.TriggeredOn, trigger)
}

func (indicator *Indicator) AddTag(tag string) {
	indicator.Tags = append(indicator.Tags, tag)
}

func (indicator *Indicator) AddReference(reference string) {
	indicator.References = append(indicator.References, reference)
}

func (indicator *Indicator) Map() map[string]interface{} {

	triggeredOn := make(map[string][]string)
	for _, trigger := range indicator.TriggeredOn {
		triggeredOn[trigger.TriggerName] = trigger.Matches
	}

	outMap := make(map[string]interface{})
	outMap["id"] = indicator.Id
	outMap["title"] = indicator.Title
	outMap["score"] = indicator.Score
	outMap["created_date"] = indicator.CreatedDate.Format(time.RFC3339)
	outMap["accessed_date"] = indicator.AccessedDate.Format(time.RFC3339)
	outMap["link"] = indicator.Link
	outMap["source"] = indicator.Source
	outMap["source_id"] = indicator.SourceId
	outMap["references"] = indicator.References
	outMap["triggered_on"] = triggeredOn
	outMap["tags"] = indicator.Tags
	return outMap
}

type IndicatorCollection struct {
	Indicators []Indicator
}

func (icollection *IndicatorCollection) CumulativeScore() (cumulativeScore int) {
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

type IndicatorFactory struct {
	factoryConfig IndicatorFactoryConfig
}

func (f IndicatorFactory) Config() IndicatorFactoryConfig {
	return f.factoryConfig
}

func (f IndicatorFactory) IsZero() bool {
	return f == IndicatorFactory{}
}

type IndicatorFactoryConfig struct {
	Source string
}

func NewIndicatorFactory(ifc IndicatorFactoryConfig) (*IndicatorFactory, error) {
	return &IndicatorFactory{factoryConfig: ifc}, nil
}

func (f IndicatorFactory) NewIndicator(
	title string,
	score int,
	createdDate time.Time,
	accessedDate time.Time,
	link string,
	sourceId string,
	references []string) Indicator {
	return Indicator{
		Id:           uuid.New().String(),
		Title:        title,
		Score:        score,
		CreatedDate:  createdDate,
		AccessedDate: accessedDate,
		Link:         link,
		Source:       f.factoryConfig.Source,
		SourceId:     sourceId,
		References:   references,
		TriggeredOn:  []TriggerMatchCollection{},
		Tags:         []string{},
	}
}

func (f IndicatorFactory) UnmarshallFromMap(indicatorMap map[string]interface{}) (*Indicator, error) {
	parsedIndicator := Indicator{}
	parsedIndicator.Id = indicatorMap["id"].(string)
	parsedIndicator.Title = indicatorMap["title"].(string)
	parsedIndicator.Score = indicatorMap["score"].(int)
	parsedIndicator.Link = indicatorMap["link"].(string)
	parsedIndicator.Source = indicatorMap["source"].(string)
	parsedIndicator.SourceId = indicatorMap["source_id"].(string)
	parsedIndicator.References = indicatorMap["references"].([]string)
	parsedIndicator.TriggeredOn = indicatorMap["triggered_on"].([]TriggerMatchCollection)
	parsedIndicator.Tags = indicatorMap["tags"].([]string)

	createdDate, err := time.Parse("2006-01-02 15:04:05.000", indicatorMap["created_date"].(string))
	if err != nil {
		return nil, err
	}
	parsedIndicator.CreatedDate = createdDate

	accessedDate, err := time.Parse("2006-01-02 15:04:05.000", indicatorMap["accessed_date"].(string))
	if err != nil {
		return nil, err
	}
	parsedIndicator.AccessedDate = accessedDate

	return &parsedIndicator, nil
}

func (f IndicatorFactory) NewIndicatorCollection() (IndicatorCollection, error) {
	return IndicatorCollection{
		Indicators: []Indicator{},
	}, nil
}

func (f IndicatorFactory) MustNewIndicatorCollection() IndicatorCollection {
	collection, err := f.NewIndicatorCollection()
	if err != nil {
		panic(err)
	}
	return collection
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
