package scrapper

import (
	"github.com/andybalholm/cascadia"
	"github.com/RulzUrLibrary/api/utils"
	"golang.org/x/net/html"
	"net/url"
	"strconv"
)

const SENSCRITIQUE_NAME = "Sens Critique"
const SENSCRITIQUE_URL = "https://www.senscritique.com"
const SENSCRITIQUE_MAX_NOTATION = 10

var (
	matchSCNotation = cascadia.MustCompile("a.erra-global")
	matchSCDetail   = cascadia.MustCompile("a.elco-anchor")
)

func getLink(node *html.Node) string {
	if node == nil {
		return ""
	}
	link, ok := getAttrs(node)["href"]
	if ok {
		link = SENSCRITIQUE_URL + link
	}
	return link
}

func (s *Scrapper) SensCritique(title string) (notation utils.Notation, err error) {
	var u *url.URL
	var q url.Values

	u, err = url.Parse(SENSCRITIQUE_URL)
	if err != nil {
		return
	}
	u, err = u.Parse("/recherche")
	if err != nil {
		return
	}
	q = u.Query()
	q.Set("query", title)
	u.RawQuery = q.Encode()

	err = s.Parse(u.String(), func(doc *html.Node) error {
		link := matchSCDetail.MatchFirst(doc)
		res, _ := strconv.ParseFloat(getText(matchSCNotation.MatchFirst(doc)), 32)
		if err != nil {
			return nil // we don't care if we got nothing
		}
		notation = utils.Notation{SENSCRITIQUE_NAME, float32(res), SENSCRITIQUE_MAX_NOTATION, getLink(link)}
		return nil
	})
	return
}
