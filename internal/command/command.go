// Package command contains the command flags and implementations
package command

import (
	"errors"
	"flag"
)

type Command struct {
	shouldKeepSameDomain bool
	shouldFindEmails     bool
	yesToAll             bool
	initialHost          string
}

func (c *Command) KeepSameDomain() bool { return c.shouldKeepSameDomain }
func (c *Command) FindEmails() bool     { return c.shouldFindEmails }
func (c *Command) YesToAll() bool       { return c.yesToAll }
func (c *Command) InitialHost() string  { return c.initialHost }

func (c *Command) Parse() error {
	yesToAllFlag := flag.Bool("y", false, "-y -> say yes to all questions (search until is every site done)")
	domainFlag := flag.Bool("keepDomain", false, "--keepDomain -> keep the search on same domain")
	emailsFlag := flag.Bool("emails", false, "--emails -> search also for emails")
	initialHostFlag := flag.String("h", "", "-h [initial host]")
	flag.Parse()

	if *yesToAllFlag {
		c.yesToAll = *yesToAllFlag
	}

	if *domainFlag {
		c.shouldKeepSameDomain = *domainFlag
	}

	if *emailsFlag {
		c.shouldFindEmails = *emailsFlag
	}

	if *initialHostFlag != "" {
		c.initialHost = *initialHostFlag
	}

	if c.initialHost == "" {
		return errors.New("the initial host is required, inform it with -h")
	}

	return nil
}

func NewCommand() *Command { return &Command{} }
