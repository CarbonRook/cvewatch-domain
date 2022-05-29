package indicator

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

var testIndicator Indicator
var testTopic Topic

func TestMain(m *testing.M) {
	var err error
	testTopic, err = NewRegexTopic("cve", `(?i)cve-\d+-\d+`)
	if err != nil {
		log.Panicf("cannot create test trigger")
	}

	indicatorFactory, err := NewIndicatorFactory("reddit")
	if err != nil {
		log.Fatalf("failed to create indicator factory: %s", err)
	}

	mention, err := NewMention(testTopic.Name(), []byte("CVE-2021-44228"))
	mentions := []Mention{mention}
	if err != nil {
		log.Fatalf("failed to create triggerMatchCollection: %s", err)
	}

	createdDate, err := time.Parse(time.RFC3339, "2021-12-13T20:16:57Z")
	if err != nil {
		log.Fatalf("failed to parse createdDate: %s", err.Error())
	}
	accessedDate, err := time.Parse(time.RFC3339, "2021-12-13T20:17:36.602452Z")
	if err != nil {
		log.Fatalf("failed to parse accessedDate: %s", err.Error())
	}
	testIndicator = indicatorFactory.MustNewIndicator()
	testIndicator.Title = "Logpresso CVE-2021-44228-Scanner (Log4j Vulnerability)"
	testIndicator.Score = 1
	testIndicator.CreatedDate = createdDate
	testIndicator.AccessedDate = accessedDate
	testIndicator.Link = "https://reddit.com/r/sysadmin/comments/rfoz5d/logpresso_cve202144228scanner_log4j_vulnerability/"
	testIndicator.SourceId = "t3_rfoz5d"
	testIndicator.References = []string{
		"https://www.reddit.com/r/sysadmin/comments/rfoz5d/logpresso_cve202144228scanner_log4j_vulnerability/",
		"https://github.com/logpresso/CVE-2021-44228-Scanner/releases/download/v1.2.3/logpresso-log4j2-scan-1.2.3.jar",
		"https://github.com/logpresso/CVE-2021-44228-Scanner/releases/download/v1.2.3/logpresso-log4j2-scan-1.2.3.jar",
	}
	testIndicator.Mentions = mentions
	testIndicator.Tags = []string{"cve-2021-44228", "scanner", "log4j"}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestUnmarshallIndicatorFromJson(t *testing.T) {
	indicatorJson, err := json.Marshal(testIndicator)
	if err != nil {
		t.Errorf("failed to marshall indicator to json: %s", err)
	}
	log.Println(string(indicatorJson))

	var indicator Indicator
	err = json.Unmarshal([]byte(indicatorJson), &indicator)
	if err != nil {
		t.Errorf("failed to unmarshall indicator json: %s", err)
	}

	if !reflect.DeepEqual(indicator, testIndicator) {
		t.Errorf("indicators are not equal after marhsalling")
	}

	log.Printf("successfully unmarshalled")
}
