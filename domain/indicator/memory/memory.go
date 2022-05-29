package memory

import (
	"context"
	"sync"

	"github.com/carbonrook/cvewatch-domain/domain/indicator"
)

type MemoryRepository struct {
	Factory    indicator.IndicatorFactory
	Collection indicator.IndicatorCollection
	lock       *sync.RWMutex
}

func (mr MemoryRepository) Add(ctx context.Context, indicator indicator.Indicator) error {
	mr.lock.Lock()
	defer mr.lock.Unlock()
	mr.Collection.Append(indicator)
	return nil
}

func (mr MemoryRepository) GetById(ctx context.Context, id string) (*indicator.Indicator, error) {
	mr.lock.RLock()
	defer mr.lock.RUnlock()
	for _, indicatorMatch := range mr.Collection.Indicators {
		if indicatorMatch.Id == id {
			return &indicatorMatch, nil
		}
	}
	return nil, indicator.ErrIndicatorNotFound
}

func (mr MemoryRepository) GetByLink(ctx context.Context, link string) (*indicator.Indicator, error) {
	mr.lock.RLock()
	defer mr.lock.RUnlock()
	for _, indicatorMatch := range mr.Collection.Indicators {
		if indicatorMatch.Link == link {
			return &indicatorMatch, nil
		}
	}
	return nil, indicator.ErrIndicatorNotFound
}

func NewMemoryIndicatorRepository() (indicator.IndicatorRepository, error) {
	factory, err := indicator.NewIndicatorFactory("")
	if err != nil {
		return MemoryRepository{}, err
	}
	collection := factory.MustNewIndicatorCollection()
	return MemoryRepository{
		Factory:    factory,
		Collection: collection,
		lock:       &sync.RWMutex{},
	}, nil
}

func MustNewMemoryIndicatorRepository() indicator.IndicatorRepository {
	repo, err := NewMemoryIndicatorRepository()
	if err != nil {
		panic(err)
	}
	return repo
}
