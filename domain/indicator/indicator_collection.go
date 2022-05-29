package indicator

import "fmt"

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
