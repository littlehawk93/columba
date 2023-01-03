package ups

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/littlehawk93/columba/providers/utils"
	"github.com/littlehawk93/columba/tracking"
)

const (
	shipmentProgressDetailTimestampFormat string = "01/02/2006 3:04 PM"
)

type trackingResponse struct {
	Details []trackDetails `json:"trackDetails"`
}

type trackDetails struct {
	ProgressDetails []shipmentProgressDetail `json:"shipmentProgressActivities"`
}

type shipmentProgressDetail struct {
	DateStr  string `json:"date"`
	TimeStr  string `json:"time"`
	Location string `json:"location"`
	Activity string `json:"activityScan"`
}

func (me shipmentProgressDetail) GetLocation() *tracking.Location {

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

func (me shipmentProgressDetail) GetTimestamp() (*time.Time, error) {

	if me.DateStr == "" && me.TimeStr == "" {
		return nil, nil
	}

	timestampString := fmt.Sprintf("%s %s", me.DateStr, me.TimeStr)
	timestampString = strings.TrimSpace(regexp.MustCompile(`[^0-9AMP:\/\s]`).ReplaceAllString(timestampString, ""))

	t, err := time.Parse(shipmentProgressDetailTimestampFormat, timestampString)
	return &t, err
}
