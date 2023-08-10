package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/littlehawk93/columba/tracking"
)

const (
	// userAgentWin10Edge Win10 MS Edge User Agent
	userAgentWin10Edge       string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246"
	cleanLocationRegexCutset string = `(^[^A-Za-z0-9]+)|([^A-Za-z0-9]+$)`
)

var DebugRequests bool = false

// GoqueryTrackingEventParser definition for a function that returns a tracking event from a GoQuery HTML selection
type GoqueryTrackingEventParser func(doc *goquery.Selection) (tracking.Event, error)

// ChromeDpTrackingEventParser definition for a function that returns a tracking event from a ChromeDP HTML node
type ChromeDpTrackingEventParser func(node *cdp.Node) ([]tracking.Event, error)

// ParseTrackingEventsURL helper function for parsing HTML from a given URL and parsing a set of events from the returned response
func ParseTrackingEventsURL(url, selector string, parser GoqueryTrackingEventParser) ([]tracking.Event, error) {

	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("[ParseTrackingEventsURL] error retrieving URL: %w", err)
	}

	defer res.Body.Close()

	return ParseTrackingEventsReader(res.Body, selector, parser)
}

func ParseTrackingEventsHeadlessChrome(url, selector string, options tracking.Options, provActions []chromedp.Action, parser ChromeDpTrackingEventParser) ([]tracking.Event, error) {

	tmpDir, err := ioutil.TempDir(os.TempDir(), "columba-chrome-headless")

	if err != nil {
		return nil, fmt.Errorf("[ParseTrackingEventsHeadlessChrome] unable to create temp directory: %w", err)
	}

	defer func() {
		time.Sleep(1 * time.Second)
		os.RemoveAll(tmpDir)
	}()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("window-size", "1280x800"),
		chromedp.UserDataDir(tmpDir),
		chromedp.UserAgent(options.UserAgent),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, time.Duration(options.TimeoutSeconds)*time.Second)
	defer cancel()

	if err := chromedp.Run(ctx); err != nil {
		return nil, fmt.Errorf("[ParseTrackingEventsHeadlessChrome] error connecting to headless chrome: %w", err)
	}

	nodes := make([]*cdp.Node, 0)

	actions := []chromedp.Action{
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.BySearch),
	}

	actions = append(actions, provActions...)

	actions = append(actions, chromedp.Nodes(selector, &nodes))

	if err := chromedp.Run(ctx, actions...); err != nil {
		return nil, fmt.Errorf("[ParseTrackingEventsHeadlessChrome] error executing chromedp actions: %w", err)
	}

	results := make([]tracking.Event, 0)

	for _, n := range nodes {
		events, err := parser(n)

		if err != nil {
			return nil, fmt.Errorf("[ParseTrackingEventsHeadlessChrome] error parsing HTML node: %w", err)
		}

		if len(events) > 0 {
			for _, eve := range events {
				if eve.EventText != "" || eve.Location != nil || eve.Timestamp != nil {
					results = append(results, eve)
				}
			}
		}
	}

	return results, nil
}

// ParseTrackingEventsReader helper function for parsing HTML from a given reader and parsing a set of events from the returned response
func ParseTrackingEventsReader(r io.Reader, selector string, parser GoqueryTrackingEventParser) ([]tracking.Event, error) {

	results := make([]tracking.Event, 0)

	document, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return results, err
	}

	document.Find(selector).Each(func(i int, s *goquery.Selection) {

		if err != nil {
			return
		}

		evt, err2 := parser(s)

		if err2 != nil {
			err = err2
		} else if evt.EventText != "" || evt.Location != nil || evt.Timestamp != nil {
			results = append(results, evt)
		}
	})

	return results, err
}

// ElemExistsWithData returns the first matching element or nil if no element found and returns true if the elem has any text in it
func ElemExistsWithData(s *goquery.Selection, selector string) (*goquery.Selection, bool) {

	elem := s.Find(selector).First()

	if elem == nil {
		return nil, false
	}

	return elem, len(strings.TrimSpace(elem.Text())) > 0
}

// GetUserAgent retreive a user agent to make Columba requests appear as desktop client apps to tracking websites
func GetUserAgent() string {
	return userAgentWin10Edge
}

// CleanLocationWord cleans a particular city or state name for a location to standardize naming format across all shipping providers
func CleanLocationWord(word string) string {
	return strings.ToUpper(regexp.MustCompile(cleanLocationRegexCutset).ReplaceAllString(word, ""))
}

// CreateJSONRequest creates a new HTTP request, writes JSON data to the request body and sets the User-Agent and Content-Type headers. Returns any errors encountered
func CreateJSONRequest[T any](method, url string, data T) (*http.Request, error) {

	b, err := json.Marshal(data)

	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(b)

	req, err := CreateRequest(method, url, r)

	if DebugRequests {
		log.Println(string(b))
	}

	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// CreateRequest creates a new HTTP request and sets the User-Agent header
func CreateRequest(method, url string, r io.Reader) (*http.Request, error) {

	req, err := http.NewRequest(method, url, r)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", GetUserAgent())

	if DebugRequests {
		log.Println("REQUEST")
		log.Printf("%s: %s\n\n\n", method, url)
	}

	return req, nil
}

// GetResponseJSON executes a HTTP request using a HTTP client and parses the JSON data from the response. Returns any errors encountered
func GetResponseJSON[T any](client *http.Client, req *http.Request, data T) error {

	bytes, err := GetResponseBytes(client, req)

	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &data)
}

// GetResponseBytes executes a HTTP request using a HTTP client and reads all bytes from the response. Returns any errors encountered
func GetResponseBytes(client *http.Client, req *http.Request) ([]byte, error) {

	if DebugRequests {
		log.Println("COOKIES")

		for _, c := range client.Jar.Cookies(req.URL) {
			log.Printf("%s - %s\n", c.Name, c.Value)
		}
		log.Print("\n\n\n")
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)

	if err == nil && DebugRequests {
		log.Println("RESPONSE")
		log.Println(string(b))
		log.Print("\n\n\n")
	}
	return b, err
}

// NewClientWithCookies
func NewClientWithCookies() (*http.Client, error) {

	jar, err := cookiejar.New(nil)

	if err != nil {
		return nil, err
	}
	return &http.Client{Jar: jar}, nil
}
