package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	datastructure "github.com/rafael-gsr/go-web-crawler/internal/data_structure"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func PrintList(list []string) {
	for _, item := range list {
		fmt.Println(item)
	}
}

func main() {
	toCraw := datastructure.NewSet[string]()
	toCraw.Add("http://www.bancocn.com")

	visited := datastructure.NewSet[string]()

	for toCraw.Len() > 0 {
		crawledPageUrl := toCraw.GetFirst()
		toCraw.Remove(crawledPageUrl)

		if visited.Contains(crawledPageUrl) {
			continue
		}

		fmt.Printf("Crawling '%s' \n", crawledPageUrl)
		links, err := GetPageLinks(crawledPageUrl)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, link := range links {
			if !(visited.Contains(link) || toCraw.Contains(link)) {
				toCraw.Add(link)
				fmt.Println(link)
			}
		}

		visited.Add(crawledPageUrl)

		isUserChoosing := true

		for isUserChoosing {

			response, err := DoUserWantsToContinue()
			if err != nil {
				fmt.Println(err)
				return
			}

			switch response {

			case "quit":
				isUserChoosing = false
				fmt.Printf("\n\nVisited sites:\n")
				PrintList(visited.GetList())
				return

			case "yes":
				isUserChoosing = false
				continue

			case "no":
				toCraw.Remove(toCraw.GetFirst())
				PrintList(toCraw.GetList())

			default:
				continue

			}
		}
	}
}

func GetPageLinks(pageUrl string) (links []string, err error) {
	res, err := http.Get(pageUrl)
	if err != nil {
		return
	}

	parsedHTML, err := html.Parse(res.Body)
	if err != nil {
		return
	}

	for des := range parsedHTML.Descendants() {
		isLinkNode := des.DataAtom == atom.A

		if isLinkNode {
			for _, attrbt := range des.Attr {

				isLinkAttribute := attrbt.Key == "href"

				if isLinkAttribute {
					link := attrbt.Val
					isExternalLink := len(link) > 3 && link[:4] == "http"

					if !isExternalLink {
						completeInternalPath, err := url.JoinPath(pageUrl, link)
						if err != nil {
							continue
						}
						link = completeInternalPath
					}

					links = append(links, link)

				}
			}
		}
	}

	return
}

func DoUserWantsToContinue() (string, error) {
	YES := "yes"
	NO := "no"
	QUIT := "quit"

	valid := false
	answer := ""

	for !valid {
		fmt.Println("Craw next page? [Y/n/q]")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		switch strings.ToLower(strings.TrimSpace(response)) {
		case "y":
			answer = YES
			valid = true
			continue

		case "":
			answer = YES
			valid = true
			continue

		case "n":
			answer = NO
			valid = true
			continue

		case "q":
			answer = QUIT
			valid = true
			continue
		}

	}

	return answer, nil
}
