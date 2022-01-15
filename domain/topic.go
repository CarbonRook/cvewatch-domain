package domain

import (
	"regexp"
)

type Topic interface {
	Name() string
	Mentioned(string) (bool, error)
	Mentions(string) ([]Mention, error)
}

type RegexTopic struct {
	name          string
	regex         string
	compiledRegex regexp.Regexp
}

func (rt RegexTopic) Name() string {
	return rt.name
}

func (rt RegexTopic) Mentions(post string) ([]Mention, error) {
	mentions := []Mention{}

	matches := rt.compiledRegex.FindAll([]byte(post), -1)
	for _, match := range matches {
		mention := MustNewMention(rt.name, match)
		mentions = append(mentions, mention)
	}

	return mentions, nil
}

func (rt RegexTopic) Mentioned(post string) (bool, error) {
	mentions, err := rt.Mentions(post)
	if err != nil {
		return false, err
	}
	return len(mentions) > 0, nil
}

func NewRegexTopic(name string, regex string) (Topic, error) {
	compiledRegex, err := regexp.Compile(regex)
	if err != nil {
		return RegexTopic{}, err
	}

	trigger := RegexTopic{
		name:          name,
		regex:         regex,
		compiledRegex: *compiledRegex,
	}

	return trigger, nil
}

func MustNewRegexTopic(name string, regex string) Topic {
	trigger, err := NewRegexTopic(name, regex)
	if err != nil {
		panic(err)
	}
	return trigger
}
