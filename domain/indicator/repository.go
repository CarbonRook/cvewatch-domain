package indicator

import (
	"context"
	"errors"
)

var (
	ErrIndicatorNotFound      = errors.New("the indicator was not found")
	ErrIndicatorAlreadyExists = errors.New("the indicator already exists")
)

type IndicatorRepository interface {
	Add(ctx context.Context, indicator Indicator) error
	GetById(ctx context.Context, id string) (*Indicator, error)
	GetByLink(ctx context.Context, link string) (*Indicator, error)
	//GetByTriggerName(ctx context.Context, triggerName string) (*IndicatorCollection, error)
	//GetByMatch(ctx context.Context, match string) (*IndicatorCollection, error)
	//GetBySource(ctx context.Context, source string) (*IndicatorCollection, error)
	//GetBySourceId(ctx context.Context, sourceId string) (Indicator, error)
	//GetLatest(ctx context.Context) (Indicator, error)
	//GetBetween(ctx context.Context, start time.Time, end time.Time) (*IndicatorCollection, error)
}
