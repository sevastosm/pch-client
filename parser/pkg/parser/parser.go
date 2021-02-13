package parser

import (
	"encoding/json"
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

// PCHParser fetch data from PCH extracts data and triggers their processing
type PCHParser interface {
	// FetchSummaries fetch summaries and trigger processing
	FetchSummaries(processor func(response *FetchResponse)) error
}

func New(cl *http.Client, config config.ClientConfig) PCHParser {
	return &pchParser{hc: cl, config: config}
}

type pchParser struct {
	hc     *http.Client
	config config.ClientConfig
}

func (ip *pchParser) FetchSummaries(proc func(response *FetchResponse)) error {
	initResp, err := ip.initRequest()
	if err != nil {
		return err
	}

	nonce := initResp.Nonce
	for _, server := range filterServers(initResp.Servers, ip.config) {
		fetchResp, err := ip.fetchServerData(nonce, server)
		if err != nil {
			return err
		}
		// skip failure to fetch data
		if fetchResp == nil {
			continue
		}

		proc(fetchResp)
		applyDelay(ip.config.ParserRateLimitDelayMillis)
		nonce = fetchResp.Nonce
	}
	return nil
}

func (p *pchParser) initRequest() (*InitResponse, error) {
	res, err := http.Get(PCHInitURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	var output InitResponse
	output.Nonce, err = ExtractNonceFromCookie(res.Header.Get("Set-Cookie"))
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
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
					IXP:      strings.TrimSpace(optionContent[0]),
					City:     strings.TrimSpace(optionContent[1]),
					Country:  strings.TrimSpace(optionContent[2]),
					Protocol: p.config.IPVersion,
				},
				ItemID: optionID,
			}
			output.Servers = append(output.Servers, option)
		}
	})

	log.Println("[parser] initiated")
	return &output, nil
}

func (p *pchParser) fetchServerData(nonce Nonce, server IXPServerOption) (*FetchResponse, error) {
	params := url.Values{}
	params.Add("router", strconv.Itoa(server.ItemID))
	params.Add("pch_nonce", string(nonce))
	params.Add("args", "")

	if p.config.IPVersion == config.IPv6 {
		params.Add("query", "v6_summary")
	} else {
		params.Add("query", "summary")
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
	res, err := p.hc.Do(req)
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
		log.Println("[parser] expected single result - skipping")
		return nil, nil
	}

	summary, err := ExtractIXPSummary(results[0].Result)
	if err != nil {
		log.Printf("[parser] failed to extract summary, error: %v\n", err)
		return nil, nil
	}

	log.Printf("[parser] received summary for IXP: %v, IP: %s\n", server, p.config.IPVersion)
	return &FetchResponse{
		Server:  server,
		Summary: *summary,
		Nonce:   nonce,
	}, nil
}

func applyDelay(delayMillis int64) {
	if delayMillis > 0 {
		log.Printf("[parser] delay %dms", delayMillis)
		time.Sleep(time.Duration(delayMillis) * time.Millisecond)
	}
}

func filterServers(servers IXPServerOptions, clientConfig config.ClientConfig) []IXPServerOption {
	if ixp := clientConfig.IXP; ixp != "" {
		servers = servers.filterBy(func(opt *IXPServerOption) bool {
			return strings.ToLower(opt.IXP) == strings.ToLower(ixp)
		})
	}

	if city := clientConfig.City; city != "" {
		servers = servers.filterBy(func(opt *IXPServerOption) bool {
			return strings.ToLower(opt.City) == strings.ToLower(city)
		})
	}

	if country := clientConfig.Country; country != "" {
		servers = servers.filterBy(func(opt *IXPServerOption) bool {
			return strings.ToLower(opt.Country) == strings.ToLower(country)
		})
	}

	if limit := clientConfig.ServerLimit; limit > 0 {
		if limit < len(servers) {
			servers = servers[:clientConfig.ServerLimit]
		}
	}

	log.Printf("[parser] filtered IXP servers to size %d\n", len(servers))
	return servers
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
