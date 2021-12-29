package domain

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"testing"
	"time"
)

var testIndicator Indicator
var testTrigger Trigger

func TestMain(m *testing.M) {
	var err error
	testTrigger, err = NewTrigger("cve", `(?i)cve-\d+-\d+`)
	if err != nil {
		log.Panicf("cannot create test trigger")
	}

	indicatorFactory, err := NewIndicatorFactory(IndicatorFactoryConfig{})
	if err != nil {
		log.Fatalf("failed to create indicator factory: %s", err)
	}

	triggerMatchCollection, err := NewTriggerMatchCollection(testTrigger.Name, [][]byte{[]byte("CVE-2021-44228")})
	if err != nil {
		log.Fatalf("failed to create triggerMatchCollection: %s", err)
	}

	createdDate, _ := time.Parse("2006-01-02 15:04:05.000", "2021-12-13T20:16:57Z")
	accessedDate, _ := time.Parse("2006-01-02 15:04:05.000", "2021-12-13T20:17:36.602452Z")
	testIndicator = indicatorFactory.NewIndicator(
		"Logpresso CVE-2021-44228-Scanner (Log4j Vulnerability)",
		1,
		createdDate,
		accessedDate,
		"https://reddit.com/r/sysadmin/comments/rfoz5d/logpresso_cve202144228scanner_log4j_vulnerability/",
		"Reddit",
		"t3_rfoz5d",
		[]string{
			"https://www.reddit.com/r/sysadmin/comments/rfoz5d/logpresso_cve202144228scanner_log4j_vulnerability/",
			"https://github.com/logpresso/CVE-2021-44228-Scanner/releases/download/v1.2.3/logpresso-log4j2-scan-1.2.3.jar",
			"https://github.com/logpresso/CVE-2021-44228-Scanner/releases/download/v1.2.3/logpresso-log4j2-scan-1.2.3.jar",
		},
		triggerMatchCollection,
		[]string{
			"cve-2021-44228",
			"scanner",
			"log4j",
		})

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestUnmarshallIndicatorFromJson(t *testing.T) {
	indicatorJson, err := json.Marshal(testIndicator)
	if err != nil {
		t.Errorf("failed to marshall indicator to json: %s", err)
	}

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

func TestUnmarshalIndicatorFromDatabase(t *testing.T) {
	factory, err := NewIndicatorFactory(IndicatorFactoryConfig{})
	if err != nil {
		t.Errorf("failed to initialise IndicatorFactory")
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
		[]string{
			"netsec", "cve-9999-1234",
		},
	)

	if err != nil {
		t.Errorf("failed to unmarshal indicator from DB")
	}

	t.Logf("indicator:\n%s", indicator.String())
}
