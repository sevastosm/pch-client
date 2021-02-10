package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sermojohn/postgres-client/pkg/config"
	"github.com/sermojohn/postgres-client/pkg/domain"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	PCHInitURL  = "https://www.pch.net/tools/looking_glass/"
	PCHQueryURL = "https://www.pch.net/tools/looking_glass_query"
)

type IXPParser interface {
	// InitIXPServers config data for requesting actual IXP data
	InitIXPServers() (*InitResponse, error)
	// FetchIXPData fetches and extracts IXP entry data for specified options
	FetchIXPData(nonce Nonce, option IXPServerOption) (*FetchResponse, error)
	// ForEachSummary fetch and process summaries
	ForEachSummary(amount int, delayMillis int64, processor func(response *FetchResponse) error) error
}

func New(cl *http.Client, config config.ParserConfig) IXPParser {
	return &ixpParser{hc: cl, config: config}
}

type ixpParser struct {
	hc     *http.Client
	config config.ParserConfig
}

func (is *ixpParser) InitIXPServers() (*InitResponse, error) {
	res, err := http.Get(PCHInitURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	var output InitResponse
	output.Nonce, err = ExtractNonceFromCookie(res.Header.Get("Set-Cookie"))
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the IXP options
	doc.Find("#looing_glass_form .router_sort_city > option").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		optionContent := strings.Split(s.Text(), ",")
		s.Attr("value")

		if optionValue, found := s.Attr("value"); found {
			optionID, _ := strconv.Atoi(optionValue)
			option := IXPServerOption{
				IXPServer: domain.IXPServer{
					IXP:     optionContent[0],
					City:    optionContent[1],
					Country: optionContent[2],
				},
				ItemID: optionID,
			}
			output.Servers = append(output.Servers, option)
		}
	})

	log.Println("[parser] initiated")
	return &output, nil
}

func (ip *ixpParser) FetchIXPData(nonce Nonce, server IXPServerOption) (*FetchResponse, error) {
	params := url.Values{}
	params.Add("router", strconv.Itoa(server.ItemID))
	params.Add("pch_nonce", string(nonce))
	params.Add("args", "")
	switch ip.config.IPVersion{
	case "v6":
		params.Add("query", "v6_summary")
	case "v4":
		params.Add("query", "summary")
	default:
		return nil, errors.New("invalid IP version")
	}

	reqURL, err := url.Parse(PCHQueryURL)
	if err != nil {
		return nil, err
	}
	reqURL.RawQuery = params.Encode()

	req := &http.Request{
		URL:    reqURL,
		Header: map[string][]string{},
	}
	req.AddCookie(&http.Cookie{
		Name:  "pch_nonce" + string(nonce),
		Value: string(nonce),
	})
	res, err := ip.hc.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := downloadFileAndRead(res)
	if err != nil {
		return nil, err
	}

	var results []domain.QueryResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	if len(results) != 1 {
		return nil, errors.New("expected single result")
	}

	summary, err := ExtractIXPSummary(results[0].Result)
	if err != nil {
		log.Printf("failed to extract summary, error: %v\n", err)
		return nil, nil
	}

	log.Printf("[parser] received summary for IXP: %v, IP: %s", server, ip.config.IPVersion)
	return &FetchResponse{
		Server:  server,
		Summary: *summary,
		Nonce:   nonce,
	}, nil
}

func (ip *ixpParser) ForEachSummary(amount int, delayMillis int64, proc func(response *FetchResponse) error) error {
	initResp, err := ip.InitIXPServers()
	if err != nil {
		return err
	}

	nonce := initResp.Nonce
	for _, server := range selectServers(amount, initResp.Servers) {
		fetchResp, err := ip.FetchIXPData(nonce, server)
		if err != nil {
			return err
		}

		err = proc(fetchResp)
		if err != nil {
			return err
		}
		applyDelay(delayMillis)
		nonce = fetchResp.Nonce
	}
	return nil
}

func applyDelay(delayMillis int64) {
	if delayMillis > 0 {
		log.Printf("[parser] delay %dms", delayMillis)
		time.Sleep(time.Duration(delayMillis) * time.Millisecond)
	}
}

func selectServers(amount int, servers []IXPServerOption) []IXPServerOption {
	if amount == 0 {
		return servers
	}

	if amount > len(servers) {
		return servers
	}

	return servers[:amount]
}

func downloadFileAndRead(response *http.Response) ([]byte, error) {
	f, err := ioutil.TempFile(os.TempDir(), "ixp-")
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = f.Close()
		_ = os.Remove(f.Name())
	}()

	_, err = io.Copy(f, response.Body)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(f.Name())

}
