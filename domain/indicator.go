package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Indicator struct {
	ID           string
	Title        string
	Score        int
	CreatedDate  time.Time
	AccessedDate time.Time
	Link         string
	Source       string
	SourceId     string
	References   []string
}

func (indicator *Indicator) String() string {
	return fmt.Sprintf("[%s] %s %s (%d)", indicator.CreatedDate.UTC().Format(time.RFC3339), indicator.ID, indicator.Title, indicator.Score)
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

func (icollection *IndicatorCollection) Last() Indicator {
	return icollection.Indicators[len(icollection.Indicators)-1]
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
	source string,
	sourceId string,
	references []string) Indicator {
	return Indicator{
		ID:           uuid.New().String(),
		Title:        title,
		Score:        score,
		CreatedDate:  createdDate,
		AccessedDate: accessedDate,
		Link:         link,
		Source:       source,
		SourceId:     sourceId,
		References:   references,
	}
}

func (f IndicatorFactory) UnmarshalIndicatorFromDatabase(id string, title string, score int, created string, accessed string, link string, source string, sourceId string, references []string) (*Indicator, error) {
	createdDate, err := time.Parse("2006-01-02 15:04:05.000", created)

	if err != nil {
		return nil, err
	}

	accessedDate, err := time.Parse("2006-01-02 15:04:05.000", accessed)

	if err != nil {
		return nil, err
	}

	return &Indicator{
		ID:           id,
		Title:        title,
		Score:        score,
		CreatedDate:  createdDate,
		AccessedDate: accessedDate,
		Link:         link,
		Source:       source,
		SourceId:     sourceId,
		References:   references,
	}, nil
}

func (f IndicatorFactory) NewIndicatorCollection() IndicatorCollection {
	return IndicatorCollection{
		Indicators: []Indicator{},
	}
}

type IndicatorRepository interface {
	Store(ctx context.Context, indicator Indicator) error
	GetById(ctx context.Context, id string) (Indicator, error)
	GetByCve(ctx context.Context, cveId string) (*IndicatorCollection, error)
	GetBySource(ctx context.Context, source string) (*IndicatorCollection, error)
	GetBySourceId(ctx context.Context, sourceId string) (Indicator, error)
	GetLatest(ctx context.Context) (Indicator, error)
	GetBetween(ctx context.Context, start time.Time, end time.Time) (*IndicatorCollection, error)
}

type IndicatorService interface {
}
