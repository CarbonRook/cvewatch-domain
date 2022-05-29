package indicator

import (
	"time"

	"github.com/google/uuid"
)

type IndicatorFactory struct {
	Source string
}

func (f IndicatorFactory) IsZero() bool {
	return f == (IndicatorFactory{})
}

func (f IndicatorFactory) NewIndicator() (Indicator, error) {
	return Indicator{
		Id:         uuid.New().String(),
		Source:     f.Source,
		References: []string{},
		Mentions:   []Mention{},
		Tags:       []string{},
	}, nil
}

func (f IndicatorFactory) MustNewIndicator() Indicator {
	indicator, err := f.NewIndicator()
	if err != nil {
		panic(err)
	}
	return indicator
}

func (f IndicatorFactory) UnmarshallFromMap(indicatorMap map[string]interface{}) (*Indicator, error) {
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

func NewIndicatorFactory(source string) (IndicatorFactory, error) {
	return IndicatorFactory{Source: source}, nil
}

func MustNewIndicatorFactory(source string) IndicatorFactory {
	factory, err := NewIndicatorFactory(source)
	if err != nil {
		panic(err)
	}
	return factory
}
