package domain

import (
	"regexp"
)

type Trigger struct {
	Name          string        `json:"name"`
	Regex         string        `json:"regex"`
	CompiledRegex regexp.Regexp `json:"-"`
}

func NewTrigger(name string, regex string) (Trigger, error) {
	compiledRegex, err := regexp.Compile(regex)
	if err != nil {
		return Trigger{}, err
	}

	trigger := Trigger{
		Name:          name,
		Regex:         regex,
		CompiledRegex: *compiledRegex,
	}

	return trigger, nil
}

func MustNewTrigger(name string, regex string) Trigger {
	trigger, err := NewTrigger(name, regex)
	if err != nil {
		panic(err)
	}
	return trigger
}

type TriggerMatchCollection struct {
	TriggerName string   `json:"triggerName"`
	Matches     []string `json:"matches"`
}

func NewTriggerMatchCollection(triggerName string, matches [][]byte) (TriggerMatchCollection, error) {
	collection := TriggerMatchCollection{
		TriggerName: triggerName,
		Matches:     []string{},
	}
	for _, match := range matches {
		collection.Matches = append(collection.Matches, string(match))
	}
	return collection, nil
}

func MustNewTriggerMatchCollection(triggerName string, matches [][]byte) TriggerMatchCollection {
	collection, err := NewTriggerMatchCollection(triggerName, matches)
	if err != nil {
		panic(err)
	}
	return collection
}
