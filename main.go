package main

import (
	"fmt"

	"github.com/rafael-gsr/go-web-crawler/internal/command"
	"github.com/rafael-gsr/go-web-crawler/internal/crawler"
	"github.com/rafael-gsr/go-web-crawler/utils"
)

func craw() {
	command := command.NewCommand()
	command.Parse()

	optionsHeader := ""
	if command.KeepSameDomain() {
		optionsHeader = fmt.Sprintf("%s - keeping same domain -", optionsHeader)
	}

	if command.FindEmails() {
		optionsHeader = fmt.Sprintf("%s - finding emails -", optionsHeader)
	}

	fmt.Println(optionsHeader)

	crawler := crawler.NewCrawler(command.InitialHost(), command.KeepSameDomain(), command.FindEmails())

crawlling:
	for crawler.Links().Len() > 0 {
		crawler.SetupNext()

		fmt.Println("Current links:")
		crawler.PrintLinks()

		fmt.Printf("\nNext to crawl %s \n", crawler.Current())

		if command.YesToAll() {
			crawler.Crawl()
			continue
		}

		userAction, err := utils.DoUserWantsToContinue("Continue?")
		if err != nil {
			fmt.Printf("Error while parsing user answer: %s\n", err)
			return
		}

		switch userAction {
		case "quit":
			break crawlling

		case "no":
			crawler.Skip()
			continue

		case "yes":
			fmt.Println("Start crawling ...")
			crawler.Crawl()
			continue
		}
	}

	fmt.Println("\n--==//==--")
	fmt.Println("\nSkipped: ")
	crawler.PrintSkipped()

	fmt.Println("\nVisited: ")
	crawler.PrintVisited()

	if command.FindEmails() {
		fmt.Println("\nEmails: ")
		crawler.PrintEmails()
	}
}

func main() { craw() }
