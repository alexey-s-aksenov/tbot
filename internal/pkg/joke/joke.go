package joke

// package joke

import (
	"bytes"
	"errors"
	"hash/fnv"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"

	xhtml "golang.org/x/net/html"
)

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

type duppsStore struct {
	dupps   [20]uint32
	current int
}

var fnvHash = fnv.New32a()

var store = duppsStore{}

func getHash(s string) uint32 {
	fnvHash.Write([]byte(s))
	defer fnvHash.Reset()
	return fnvHash.Sum32()
}

func checkAndStoreDupps(s string) bool {
	hash := getHash(s)
	if hashInDupps(hash) {
		return true
	}
	index := store.current
	store.dupps[index] = hash
	index = index + 1
	if index > 19 {
		index = 0
	}
	store.current = index
	return false
}

func hashInDupps(h uint32) bool {
	for _, v := range store.dupps {
		if v == h {
			return true
		}
	}
	return false
}

func (liveGetWebRequest) FetchBytes(site string) ([]byte, error) {
	resp, err := http.Get(site)
	if err != nil {
		return nil, errors.New("joke.go: Error in GetJokeBash func: " + err.Error())
	}

	defer resp.Body.Close()

	// re := regexp.MustCompile(`\r?\n`)
	// reT := regexp.MustCompile(`\t`)

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("joke.go: Error reading from reader: " + err.Error())
	}

	return buf, nil
}

// GetJokeBash получает html по ссылке и возвращает цитату с наивысшим рейтингом
func (j *bashorg) GetJoke() (string, error) {
	request := liveGetWebRequest{}
	return innerGetJoke(request)
}

func innerGetJoke(request getWebRequest) (string, error) {
	site := "https://bash.im/random"

	buf, err := request.FetchBytes(site)
	if err != nil {
		return "", errors.New("joke.go: Failed to get data by URL: " + err.Error())
	}
	root, err := xhtml.Parse(bytes.NewReader(buf))

	joke, err := getQuotes(root)
	if err != nil {
		return "", errors.New("joke.go: Error while parsing html: " + err.Error())
	}
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
				re := regexp.MustCompile(`\r?\n`)
				reT := regexp.MustCompile(`\t`)
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					if c.Data == "br" {
						quote = quote + "\n"
					} else if c.Data == "div" {

					} else {
						text := re.ReplaceAllString(c.Data, "")
						text = reT.ReplaceAllString(text, "")
						quote = quote + text
					}
				}
				q.text = xhtml.UnescapeString(quote)
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
	if checkAndStoreDupps(qq[0].text) && len(qq) > 1 {
		return qq[1].text, nil
	}
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
