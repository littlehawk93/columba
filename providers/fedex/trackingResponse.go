package fedex

import (
	"fmt"
	"strings"
	"time"

	"github.com/littlehawk93/columba/providers/utils"
	"github.com/littlehawk93/columba/tracking"
)

const (
	scanEventDateTimeFormat string = "2006-01-02 15:04:05 -0700"
)

type responsePackage struct {
	ScanEventList []scanEvent `json:"scanEventList"`
}

type responseOutput struct {
	Packages []responsePackage `json:"packages"`
}

type trackingResponse struct {
	TransactionID string         `json:"transactionId"`
	Output        responseOutput `json:"output"`
}

type scanEvent struct {
	DateStr       string `json:"date"`
	TimeStr       string `json:"time"`
	TimeOffsetStr string `json:"gmtOffset"`
	Status        string `json:"status"`
	Location      string `json:"scanLocation"`
	Details       string `json:"scanDetails"`
}

func (me scanEvent) GetTimestamp() (*time.Time, error) {

	if me.DateStr == "" && me.TimeStr == "" && me.TimeOffsetStr == "" {
		return nil, nil
	}

	timestampStr := strings.TrimSpace(fmt.Sprintf("%s %s %s", strings.TrimSpace(me.DateStr), strings.TrimSpace(me.TimeStr), strings.TrimSpace(strings.ReplaceAll(me.TimeOffsetStr, ":", ""))))

	t, err := time.Parse(scanEventDateTimeFormat[:len(timestampStr)], timestampStr)
	return &t, err
}

func (me scanEvent) GetLocation() *tracking.Location {

	words := strings.Fields(me.Location)

	if len(words) < 2 {
		return nil
	}

	for i := range words {
		words[i] = utils.CleanLocationWord(words[i])
	}

	stateIndex := -1

	for i, word := range words {
		if tracking.IsStateAbbreviation(word) && stateIndex == -1 {
			stateIndex = i
			break
		}
	}

	if stateIndex == -1 {
		return nil
	}

	cityStr := ""

	if stateIndex > 0 {
		cityStr = strings.Join(words[:stateIndex], " ")
	}

	return &tracking.Location{
		State: words[stateIndex],
		City:  cityStr,
	}
}
