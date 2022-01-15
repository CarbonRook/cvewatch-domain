package domain

type Mention struct {
	TopicName string `json:"topicName" db:"topicName"`
	Mention   string `json:"mention" db:"mention"`
}

func (tmc Mention) Map() map[string]string {
	outMap := make(map[string]string)
	outMap["topicName"] = tmc.TopicName
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