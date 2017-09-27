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
	matchSearchLink  = cascadia.MustCompile("#result_0 a")
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

func ParseIndex(_url string, isbn string) (string, error) {
	var form *html.Node

	err := Parse(_url, parse(&form, matchIndexForm))
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

func ParseSearch(_url string) (href string, err error) {

	err = Parse(_url, func(doc *html.Node) error {
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

func ParseInfo(_url string, book *Book) error {
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

	err := Parse(_url, parsingFn)
	if err != nil {
		return err
	}

	book.Title = getText(title)
	book.Title, book.Serie, book.Number = getTitle(book.Title)
	book.Price = getPrice(price)
	book.Description = getDescription(description)

	for _, author := range authors {
		book.Authors = append(book.Authors, &utils.Author{Name: getAuthor(getText(author))})
	}
	if book.Title == "" && book.Serie == "" {
		return utils.ErrParsingProduct
	}

	return book.DownloadAsset(getThumb(thumbnail))
}

func Amazon(isbn string) (book Book, err error) {
	var url string
	if url, err = ParseIndex(amazon, isbn); err != nil {
		return
	}
	if url, err = ParseSearch(url); err != nil {
		return
	}
	book = Book{&utils.Book{Isbn: isbn}}
	err = ParseInfo(url, &book)
	return
}
