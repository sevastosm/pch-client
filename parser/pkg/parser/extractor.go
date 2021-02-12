package parser

import (
	"errors"
	"fmt"
	"github.com/sermojohn/postgres-client/pkg/domain"
	"regexp"
	"strconv"
)

var (
	nonceCookieRegExp    = regexp.MustCompile("pch_nonce([a-z0-9A-Z]*)\\=")
	totalNeighborsRegExp = regexp.MustCompile("Total number of neighbors (\\d+)\\n?")
	ribEntriesRegExp     = regexp.MustCompile("RIB entries (\\d+),?")
	peersPegExp          = regexp.MustCompile("Peers (\\d+),?")
	localASNumberPegExp  = regexp.MustCompile("local AS number (\\d+)\\n?")
)

func ExtractNonceFromCookie(cookie string) (Nonce, error) {
	if cookie == "" {
		return "", errors.New("no nonce can be extract from empty cookie")
	}

	res, err := extractIXPIntStringItem(cookie, nonceCookieRegExp)
	if err != nil {
		return "", errors.New("failed to fetch nonce value from cookie")
	}
	return Nonce(res[1]), nil
}

func ExtractIXPSummary(summaryStr string) (*domain.BGPSummary, error) {
	asn, err := extractIXPIntItem(summaryStr, localASNumberPegExp)
	if err != nil {
		return nil, err
	}

	rib, err := extractIXPIntItem(summaryStr, ribEntriesRegExp)
	if err != nil {
		return nil, err
	}

	peers, err := extractIXPIntItem(summaryStr, peersPegExp)
	if err != nil {
		return nil, err
	}

	neigh, err := extractIXPIntItem(summaryStr, totalNeighborsRegExp)
	if err != nil {
		return nil, err
	}

	return &domain.BGPSummary{
		LocalASNumber:          asn,
		RIBEntries:             rib,
		NumberOfPeers:          peers,
		TotalNumberOfNeighbors: neigh,
	}, nil
}

func extractIXPIntItem(summaryStr string, pat *regexp.Regexp) (int, error) {
	res := pat.FindStringSubmatch(summaryStr)
	if len(res) != 2 {
		return 0, fmt.Errorf("could not parse section with %v\ninput %s", pat.String(), summaryStr)
	}
	item, err := strconv.Atoi(res[1])
	if err != nil {
		return 0, err
	}
	return item, nil
}

func extractIXPIntStringItem(summaryStr string, pat *regexp.Regexp) (string, error) {
	res := pat.FindStringSubmatch(summaryStr)
	if len(res) != 2 {
		return "", fmt.Errorf("could not parse section with %v\ninput %s", pat.String(), summaryStr)
	}
	return res[1], nil
}
