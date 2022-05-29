package indicator

import "strings"

type Mention struct {
	TopicName string `json:"topicName" db:"topicName"`
	Mention   string `json:"mention" db:"mention"`
}

func (m Mention) Equal(mention Mention) bool {
	return strings.EqualFold(m.TopicName, mention.TopicName) && strings.EqualFold(m.Mention, mention.Mention)
}

func (tmc Mention) Map() map[string]string {
	outMap := make(map[string]string)
	outMap["topic_name"] = tmc.TopicName
	outMap["mention"] = tmc.Mention
	return outMap
}

func NewMention(topicName string, mention []byte) (Mention, error) {
	newMention := Mention{
		TopicName: topicName,
		Mention:   string(mention),
	}
	return newMention, nil
}

func MustNewMention(topicName string, mention []byte) Mention {
	newMention, err := NewMention(topicName, mention)
	if err != nil {
		panic(err)
	}
	return newMention
}
