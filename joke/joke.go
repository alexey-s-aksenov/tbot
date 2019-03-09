package joke

// package joke

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"

	xhtml "golang.org/x/net/html"
)

// func main() {
// 	fmt.Println(GetJokeBash())
// }

type quote struct {
	rating int64
	text   string
}

type bashorg struct {
}

type getWebRequest interface {
	FetchBytes(url string) ([]byte, error)
}

type liveGetWebRequest struct {
}

func (liveGetWebRequest) FetchBytes(url string) ([]byte, error) {
	resp, err := http.Get(site)
	if err != nil {
		log.Printf("joke.go: Error in GetJokeBash func: %s", err)
		return nil, err
	}

	defer resp.Body.Close()

	// re := regexp.MustCompile(`\r?\n`)
	// reT := regexp.MustCompile(`\t`)

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("joke.go: Error reading from reader: %s", err)
		return nil, err
	}

	return buf, nil
}

// GetJokeBash получает html по ссылке и возвращает цитату с наивысшим рейтингом
func (*bashorg) GetJoke() (string, error) {
	request := liveGetWebRequest{}
	return innerGetJoke(request)
}

func innerGetJoke(request getWebRequest) (string, error) {
	site := "https://bash.im/random"

	buf, err := request.FetchBytes(site)
	if err != nil {
		return "", err
	}
	root, err := xhtml.Parse(bytes.NewReader(buf))

	joke, err := getQuotes(root)
	if err != nil {
		log.Printf("joke.go: Error while parsing html: %s", err)
		return "", nil
	}
	// fmt.Println(quotes)
	return joke, nil
}

func getQuotesSlice(doc *xhtml.Node) []*xhtml.Node {
	var result = make([]*xhtml.Node, 0)
	var f func(n *xhtml.Node)
	f = func(n *xhtml.Node) {
		if n.Type == xhtml.ElementNode && n.Data == "article" && checkAttr(n, "class", "quote") {
			result = append(result, n)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return result
}

func extractData(nn []*xhtml.Node) []*quote {
	var result = make([]*quote, 0)
	for _, n := range nn {
		var q = new(quote)
		var f func(n *xhtml.Node)
		f = func(n *xhtml.Node) {
			if n.Type == xhtml.ElementNode &&
				n.Data == "div" &&
				checkAttr(n, "class", "quote__body") {
				var quote string
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					quote = quote + c.Data + "\n"
				}
			}
			if n.Type == xhtml.ElementNode && n.Data == "div" && checkAttr(n, "class", "quote__total") {
				q.rating, _ = strconv.ParseInt(n.FirstChild.Data, 10, 32)
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
		f(n)
		if q != nil {
			result = append(result, q)
		}
	}
	return result
}

func getQuotes(root *xhtml.Node) (string, error) {

	qq := extractData(getQuotesSlice(root))

	if len(qq) == 0 {
		return "", errors.New("joke.go: No quotes found")
	}

	sort.SliceStable(qq, func(i, j int) bool {
		return qq[i].rating > qq[j].rating
	})

	return qq[0].text, nil
}

func checkAttr(n *xhtml.Node, att string, val string) bool {
	for _, attr := range n.Attr {
		if attr.Key == att && attr.Val == val {
			return true
		}
	}
	return false
}

// func nodeToString(doc *xhtml.Node) string {
// 	var buf bytes.Buffer
// 	w := io.Writer(&buf)
// 	xhtml.Render(w, doc)
// 	return buf.String()
// }
