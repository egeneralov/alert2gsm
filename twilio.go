package main

import (
	"encoding/json"
	"encoding/xml"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Twilio struct {
	AccountSID string
	AuthToken  string
	From       string

	queueLock sync.RWMutex
}

func (t *Twilio) QueueCall(to, xmlURL string) (TwilioResponse, error) {
	log.Debugf("Twilio.QueueCall(%v, %v)", to, xmlURL)
	var (
		answer TwilioResponse
		err    error
		req    *http.Request
		resp   *http.Response

		params = url.Values{}
	)

	params.Add("Url", xmlURL)
	params.Add("From", t.From)
	params.Add("To", to)
	req, err = http.NewRequest(
		"POST",
		urlPost,
		strings.NewReader(
			params.Encode(),
		),
	)
	if err != nil {
		return answer, err
	}
	req.SetBasicAuth(t.AccountSID, t.AuthToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	// lock call queue, only one call scheduling in one time
	t.queueLock.RLock()

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.queueLock.RUnlock()
		log.Error(err)
		return answer, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.queueLock.RUnlock()
		log.Error(err)
		return answer, err
	}

	err = json.Unmarshal(body, &answer)
	if err != nil {
		t.queueLock.RUnlock()
		log.Error(err)
		return answer, err
	}

	t.queueLock.RUnlock()
	return answer, nil
}

func (t *Twilio) GenerateXML(text []string, voice string) ([]byte, error) {
	log.Debugf("Twilio.GenerateXML(%+v, %v)", text, voice)
	var (
		payload = FormatRequest{}
	)

	for _, say := range text {
		payload.Say = append(
			payload.Say,
			FormatSay{
				Text:  say,
				Voice: voice,
			},
		)
		// payload.Pause = append(
		// 	payload.Pause,
		// 	FormatPause{
		// 		Length: "1",
		// 	},
		// )
	}

	return xml.MarshalIndent(payload, "", "  ")
}

type TwilioResponse struct {
	AccountSid      string      `json:"account_sid"`
	Annotation      interface{} `json:"annotation"`
	AnsweredBy      interface{} `json:"answered_by"`
	APIVersion      string      `json:"api_version"`
	CallerName      interface{} `json:"caller_name"`
	DateCreated     string      `json:"date_created"`
	DateUpdated     string      `json:"date_updated"`
	Direction       string      `json:"direction"`
	Duration        string      `json:"duration"`
	EndTime         string      `json:"end_time"`
	ForwardedFrom   string      `json:"forwarded_from"`
	From            string      `json:"from"`
	FromFormatted   string      `json:"from_formatted"`
	GroupSid        interface{} `json:"group_sid"`
	ParentCallSid   interface{} `json:"parent_call_sid"`
	PhoneNumberSid  string      `json:"phone_number_sid"`
	Price           string      `json:"price"`
	PriceUnit       string      `json:"price_unit"`
	Sid             string      `json:"sid"`
	StartTime       string      `json:"start_time"`
	Status          string      `json:"status"`
	SubresourceUris struct {
		Notifications     string `json:"notifications"`
		Recordings        string `json:"recordings"`
		Feedback          string `json:"feedback"`
		FeedbackSummaries string `json:"feedback_summaries"`
		Payments          string `json:"payments"`
		Events            string `json:"events"`
	} `json:"subresource_uris"`
	To          string      `json:"to"`
	ToFormatted string      `json:"to_formatted"`
	TrunkSid    interface{} `json:"trunk_sid"`
	URI         string      `json:"uri"`
	QueueTime   string      `json:"queue_time"`
}

type FormatSay struct {
	XMLName xml.Name `xml:"Say"`
	Text    string   `xml:",chardata"`
	Voice   string   `xml:"voice,attr,omitempty"`
}

type FormatRequest struct {
	XMLName xml.Name `xml:"Response"`
	Say     []FormatSay
	Text    string `xml:",chardata"`
	Play    string `xml:"Play,omitempty"`
	// Pause   []FormatPause `xml:"Pause"`
}

type FormatPause struct {
	XMLName xml.Name `xml:"Pause"`
	Text    string   `xml:",chardata"`
	Length  string   `xml:"length,attr"`
}
