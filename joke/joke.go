package joke

// package joke

import (
	"bytes"
	"errors"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	xhtml "golang.org/x/net/html"

	"github.com/djimenez/iconv-go"
)

// func main() {
// 	fmt.Println(GetJokeBash())
// }

type quote struct {
	rating int64
	text   string
}

func (q quote) String() string {
	return fmt.Sprintf("rating: %v, text: %v", q.rating, q.text)
}

type bashorg struct {
}

// GetJokeBash получает html по ссылке и возвращает цитату с наивысшим рейтингом
func (j *bashorg) GetJoke() (string, error) {
	site := "https://bash.im/random"
	resp, err := http.Get(site)
	if err != nil {
		log.Printf("joke.go: Error in GetJokeBash func: %s", err)
		return "", err
	}

	defer resp.Body.Close()
	//body, err := ioutil.ReadAll(resp.Body)

	utfBody, err := iconv.NewReader(resp.Body, "windows-1251", "utf-8")
	if err != nil {
		log.Printf("joke.go: Error converting ccsid to UTF: %s", err)
		return "", nil
	}
	//buf := make([]byte, 100)
	re := regexp.MustCompile(`\r?\n`)
	reT := regexp.MustCompile(`\t`)

	buf, err := ioutil.ReadAll(utfBody)
	if err != nil {
		log.Printf("joke.go: Error reading from reader: %s", err)
		return "", nil
	}
	input := re.ReplaceAllString(string(buf), "")
	input = reT.ReplaceAllString(input, "")

	doc, err := xhtml.Parse(strings.NewReader(input))

	joke, err := getQuotes(doc)
	if err != nil {
		log.Printf("joke.go: Error while parsing html: %s", err)
		return "", nil
	}
	// fmt.Println(quotes)
	return joke, nil

}

func getQuotes(root *xhtml.Node) (string, error) {

	body, err := getBody(root)
	if err != nil {
		log.Printf("joke.go: Error in getQuotes func: %s", err)
		return "", err
	}

	body, err = getDivBody(body)
	if err != nil {
		log.Printf("joke.go: Error in getQuotes func: %s", err)
		return "", err
	}
	m, err := getQuotesList(body)

	return getTheBest(m), nil
}

//  get html <body>
func getBody(doc *xhtml.Node) (*xhtml.Node, error) {
	var b *xhtml.Node
	var f func(*xhtml.Node)
	f = func(n *xhtml.Node) {
		if n.Type == xhtml.ElementNode && n.Data == "body" {
			b = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if b != nil {
		return b, nil
	}
	return nil, errors.New("Missing <body> in the node tree")
}

// get <div id="body"
func getDivBody(doc *xhtml.Node) (*xhtml.Node, error) {
	var b *xhtml.Node
	var f func(*xhtml.Node)
	f = func(n *xhtml.Node) {
		if n.Type == xhtml.ElementNode && n.Data == "div" && checkAttr(n, "id", "body") {
			b = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if b != nil {
		return b, nil
	}
	return nil, errors.New("Missing <div id=body> in the node tree")
}

// create slice of quotes from <div class="quote">
func getQuotesList(doc *xhtml.Node) ([]*quote, error) {
	s := make([]*quote, 0)

	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == xhtml.ElementNode && c.Data == "div" && checkAttr(c, "class", "quote") {
			quot := new(quote)
			// walk through all tags inside each <div class="quote">
			for q := c.FirstChild; q != nil; q = q.NextSibling {
				// geting rating
				if q.Type == xhtml.ElementNode && q.Data == "div" && checkAttr(q, "class", "actions") {
					for qq := q.FirstChild; qq != nil; qq = qq.NextSibling {
						if qq.Type == xhtml.ElementNode && qq.Data == "span" && checkAttr(qq, "class", "rating-o") {
							//fmt.Println(qq.FirstChild.FirstChild.Data)
							quot.rating, _ = strconv.ParseInt(qq.FirstChild.FirstChild.Data, 10, 32)
						}
					}
				}
				// getting quote it self
				if q.Type == xhtml.ElementNode && q.Data == "div" && checkAttr(q, "class", "text") {
					//xtml fails to parse correctly
					dirtyText := nodeToString(q)
					end := strings.Index(dirtyText, "</div>")
					dirtyText = dirtyText[len("<div class=\"text\">"):end]
					dirtyText = html.UnescapeString(dirtyText)
					quot.text = strings.Replace(dirtyText, "<br/>", "\n", -1)
				}
			}
			s = append(s, quot)
			if quot.rating <= 0 || quot.text == "" {
				fmt.Println("!!!!!!!!!!!!!", nodeToString(c))
			}
		}
	}

	return s, nil
}

func checkAttr(n *xhtml.Node, att string, val string) bool {
	for _, attr := range n.Attr {
		if attr.Key == att && attr.Val == val {
			return true
		}
	}

	return false
}

func nodeToString(doc *xhtml.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	xhtml.Render(w, doc)
	return buf.String()
}

func getTheBest(q []*quote) string {
	var result string
	var found int64
	for _, qq := range q {
		if qq.rating > found {
			result = qq.text
			found = qq.rating
		}
	}
	return result
}
