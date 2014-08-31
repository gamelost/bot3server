package wikiquote

import (
	// i would like to go on the record to state that go.net/html sucks balls.
	"code.google.com/p/go.net/html"
	iniconf "code.google.com/p/goconf/conf"
	"errors"
	"fmt"
	"github.com/gamelost/bot3server/server"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var shortcuts = map[string]string{
	"archer": "Archer_(TV_series)",
	"b5":     "Babylon_5",
	"nge":    "Neon_Genesis_Evangelion",
	"tng":    "Star_Trek:_The_Next_Generation",
	"vb":     "The_Venture_Bros.",
}

// TODO eventually we want to use mongodb, but for now ...
var quotes = make(map[string][][]string, 0)

type WikiQuoteService struct {
	server.BotHandlerService
}

func stripHTML(node *html.Node) string {
	// ugly, but elements can be nested in quotes --
	// italicized links to wikipedia, for example.
	if node.FirstChild != nil {
		innards := node
		for innards.FirstChild != nil {
			innards = innards.FirstChild
		}
		return innards.Data
	}
	return node.Data
}

func extractStatement(node *html.Node) string {
	var what string
	for child := node; child != nil; child = child.NextSibling {
		if child.FirstChild != nil {
			// translate <b>...</b> to *...*
			// and <i>...</i> to /.../
			var token string
			switch child.Data {
			case "b":
				token = "*"
			case "i":
				token = "/"
			default:
				token = ""
			}
			stripped := stripHTML(child)
			what += token + stripped + token
		} else {
			what += child.Data
		}
	}
	return what
}

func extractWikiQuoteSection(node *html.Node) []string {
	// start of quote.
	lines := make([]string, 0)

	for quote := node.FirstChild; quote != nil; quote = quote.NextSibling {
		if quote.Type == html.ElementNode && quote.Data == "dd" {
			line := quote.FirstChild
			// extract name.
			name := stripHTML(line)

			// extract statement.
			var statement string
			if line.NextSibling != nil {
				statement = extractStatement(line.NextSibling)
			}

			// post-process.
			result := name + statement
			result = strings.Replace(result, "…", "...", -1) // personal pet peeve
			lines = append(lines, result)
		}
	}
	// end of quote.
	return lines
}

func extractWikiQuotes(node *html.Node, accum [][]string) [][]string {
	// TODO here we assume that the quote section begins with
	// <dl>. This is not always true.
	if node.Type == html.ElementNode && node.Data == "dl" {
		section := extractWikiQuoteSection(node)
		accum = append(accum, section)
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		accum = extractWikiQuotes(child, accum)
	}
	return accum
}

func parseWikiQuotePage(input string) (quotes [][]string, err error) {
	response, err := http.Get(input)
	if err != nil {
		errors.New(fmt.Sprintf("Could not retrieve: %s\n", err))
		return
	}

	defer response.Body.Close()

	doc, err := html.Parse(response.Body)
	if err != nil {
		errors.New(fmt.Sprintf("Could not parse: %s\n", err))
		return
	}

	var accum = make([][]string, 0)
	results := extractWikiQuotes(doc, accum)
	return results, nil
}

func randomlyQuote(page string, regex string) []string {

	var wq_index string

	// do we have a shortcut?
	if _, ok := shortcuts[page]; ok {
		wq_index = shortcuts[page]
	} else {
		wq_index = page
	}

	// do we already have this wikiquote page parsed?
	if _, ok := quotes[wq_index]; !ok {
		full_url := "http://en.wikiquote.org/wiki/" + wq_index
		result, err := parseWikiQuotePage(full_url)
		if err != nil {
			error_message := fmt.Sprintf("Error: %s\n", err)
			return []string{error_message}
		}
		quotes[wq_index] = result
	}

	// get the quotes we're looking for.
	ng := rand.New(rand.NewSource(time.Now().UnixNano()))
	wq_page := quotes[wq_index]

	// if we have to match a regex, it's going to be a bit slow.
	// using mongodb should speed things up somewhat.
	if regex != "" {
		regex = "(?i)" + regex // make case-insensitive
		matches := make([][]string, 0)
		for _, blockquote := range wq_page {
			for _, quote := range blockquote {
				hit, _ := regexp.MatchString(regex, quote)
				if hit {
					matches = append(matches, blockquote)
					break
				}
			}
		}
		if len(matches) == 0 {
			return []string{"No matches found."}
		}
		// randomly return a quote
		index := ng.Intn(len(matches))
		return matches[index]
	} else {
		// randomly return a quote
		index := ng.Intn(len(wq_page))
		return wq_page[index]
	}
}

func parseInput(what string) (string, string) {

	input := strings.TrimPrefix(what, "!wq ")
	args := strings.SplitAfter(input, " ")

	var wq_key string
	var regexp string

	switch len(args) {
	case 0:
		// randomly select a wiki page.
		for wq_key = range quotes {
			// stupid. go doesn't have a better way of doing this?
			break
		}
	case 1:
		wq_key = strings.TrimSpace(args[0])
	default:
		wq_key = strings.TrimSpace(args[0])
		regexp = strings.TrimSpace(args[1])
	}

	return wq_key, regexp
}

func (svc *WikiQuoteService) NewService(config *iniconf.ConfigFile, publishToIRCChan chan *server.BotResponse) server.BotHandler {
	newSvc := &WikiQuoteService{}
	newSvc.Config = config
	newSvc.PublishToIRCChan = publishToIRCChan
	return newSvc
}

func (svc *WikiQuoteService) DispatchRequest(botRequest *server.BotRequest) {
	botResponse := svc.CreateBotResponse(botRequest)

	what := botRequest.Text()
	page, regex := parseInput(what)
	quote := randomlyQuote(page, regex)

	botResponse.SetMultipleLineResponse(quote)
	svc.PublishBotResponse(botResponse)
}
