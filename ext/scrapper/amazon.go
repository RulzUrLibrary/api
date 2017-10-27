package scrapper

import (
	"github.com/andybalholm/cascadia"
	"github.com/golang/glog"
	"github.com/rulzurlibrary/api/utils"
	"golang.org/x/net/html"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

const amazon = "http://amazon.fr"

var (
	matchIndexForm   = cascadia.MustCompile("form")
	matchSearchLink  = cascadia.MustCompile("#result_0 a.s-access-detail-page")
	matchTitle       = cascadia.MustCompile("#productTitle")
	matchPrice       = cascadia.MustCompile(".swatchElement.selected .a-color-price")
	matchDescription = cascadia.MustCompile("#bookDescription_feature_div noscript")
	matchThumb       = cascadia.MustCompile("#imgBlkFront")
	matchAuthors     = []cascadia.Selector{
		cascadia.MustCompile(".contributorNameID"),
		cascadia.MustCompile(".author a"),
	}
	regexpTitle = regexp.MustCompile(
		`(.*?)\s*(?:-|,)*\s*(?i:Vol\.|volume|tome|t)\s*(\d+)\s*(?:-|,|\:)*\s*(.*[a-zA-Z].*)?`,
	)
	// small alias
	ParseTitle = getTitle
)

func parse(node **html.Node, sel cascadia.Selector) ParseFn {
	return func(doc *html.Node) error {
		if *node = sel.MatchFirst(doc); *node == nil {
			return utils.ErrNoProduct
		} else {
			return nil
		}
	}
}

func getAttrs(node *html.Node) map[string]string {
	attrs := make(map[string]string)

	for _, attr := range node.Attr {
		attrs[attr.Key] = attr.Val
	}
	return attrs
}

func getText(node *html.Node) (text string) {
	if node == nil {
		return
	}
	var gettext func(node *html.Node)
	gettext = func(node *html.Node) {
		if node.Type == html.TextNode {
			text += node.Data
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			gettext(c)
		}
	}
	gettext(node)
	return strings.TrimSpace(text)
}

func getTitle(title string) (t string, s string, n int) {
	res := regexpTitle.FindStringSubmatch(title)
	if len(res) == 0 {
		t = title
	} else {
		s = res[1]
		n, _ = strconv.Atoi(res[2])
		t = res[3]
	}
	return
}

func getAuthor(author string) string {
	firstname, surname := []string{}, []string{}

	for _, w := range strings.Split(author, " ") {
		if w == strings.ToUpper(w) {
			surname = append(surname, w)
		} else {
			firstname = append(firstname, w)
		}
	}

	return strings.Join(append(firstname, surname...), " ")
}

func getThumb(node *html.Node) string {
	return strings.TrimSpace(getAttrs(node)["src"])
}

func getPrice(node *html.Node) float32 {
	splitted := strings.Split(strings.TrimSpace(getText(node)), " ")
	price := strings.Replace(splitted[len(splitted)-1], ",", ".", -1)

	value, err := strconv.ParseFloat(price, 32)
	if err != nil {
		glog.Errorf("%s", err)
	}

	return float32(value)
}

func getDescription(node *html.Node) string {
	doc, err := html.Parse(strings.NewReader(getText(node)))
	if err != nil {
		glog.Errorf("%s", err)
	}
	return getText(cascadia.MustCompile("div").MatchFirst(doc))

}

func getAuthors(authors []*html.Node) *utils.Authors {
	authorList := utils.Authors{}
	for _, author := range authors {
		authorList = append(authorList, utils.Author{Name: getAuthor(getText(author))})
	}
	return &authorList
}

func (s *Scrapper) AmazonParseIndex(_url string, isbn string) (string, error) {
	var form *html.Node

	err := s.Parse(_url, parse(&form, matchIndexForm))
	if err != nil {
		return "", err
	}
	base, err := url.Parse(_url) // parse url for adding some elements
	if err != nil {
		return "", err
	}
	rel, err := url.Parse(getAttrs(form)["action"])
	if err != nil {
		return "", err
	}

	base = base.ResolveReference(rel)
	query := base.Query()
	query.Set("field-keywords", isbn)
	base.RawQuery = query.Encode()

	return base.String(), nil
}

func (s *Scrapper) AmazonParseSearch(_url string) (href string, err error) {

	err = s.Parse(_url, func(doc *html.Node) error {
		for _, link := range matchSearchLink.MatchAll(doc) {
			href = getAttrs(link)["href"]
			if strings.Contains(href, "ebook") {
				continue
			}
			return nil
		}
		return utils.ErrNoProduct
	})
	return
}

func (s *Scrapper) AmazonParseInfo(_url string, isbn string) (book utils.Book, err error) {
	var title *html.Node
	var price *html.Node
	var description *html.Node
	var thumbnail *html.Node
	var authors []*html.Node

	var parsingFn = func(doc *html.Node) error {
		title = matchTitle.MatchFirst(doc)
		price = matchPrice.MatchFirst(doc)
		description = matchDescription.MatchFirst(doc)
		thumbnail = matchThumb.MatchFirst(doc)
		for _, matcher := range matchAuthors {
			authors = matcher.MatchAll(doc)
			if len(authors) != 0 {
				break
			}
		}
		if title == nil || thumbnail == nil {
			return utils.ErrParsingProduct
		}
		return nil
	}

	err = s.Parse(_url, parsingFn)
	if err != nil {
		return
	}

	book.Isbn = isbn
	book.Title, book.Serie, book.Number = getTitle(getText(title))
	book.Price = getPrice(price)
	book.Description = getDescription(description)
	book.Authors = getAuthors(authors)

	if book.Title == "" && book.Serie == "" {
		return book, utils.ErrParsingProduct
	}

	return book, s.DownloadThumb(getThumb(thumbnail), isbn)
}

func (s *Scrapper) Amazon(isbn string) (_ utils.Book, err error) {
	var url string
	if url, err = s.AmazonParseIndex(amazon, isbn); err != nil {
		return
	}
	if url, err = s.AmazonParseSearch(url); err != nil {
		return
	}
	return s.AmazonParseInfo(url, isbn)
}
