package domain

import "testing"

func TestUnmarshalIndicatorFromDatabase(t *testing.T) {
	factory, err := NewIndicatorFactory(IndicatorFactoryConfig{})
	if err != nil {
		t.Errorf("Failed to initialise IndicatorFactory")
	}
	indicator, err := factory.UnmarshalIndicatorFromDatabase(
		"1",
		"Test post",
		1484,
		"2021-10-02 21:32:59.100",
		"2021-10-02 21:33:05.100",
		"https://reddit.com/r/netsec/testing",
		"Reddit",
		"qwfy433",
		[]string{
			"https://reddit.com/r/netsec/reference",
		},
	)

	if err != nil {
		t.Errorf("Failed to unmarshal indicator from DB")
	}

	t.Logf("Indicator:\n%s", indicator.String())
}
