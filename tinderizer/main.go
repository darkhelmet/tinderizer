package main

import (
	"flag"
	"log"
	"os"
	"os/exec"

	"github.com/darkhelmet/env"
	"github.com/darkhelmet/tinderizer"
	J "github.com/darkhelmet/tinderizer/job"
)

var (
	logger       = log.New(os.Stdout, "[tinderizer] ", env.IntDefault("LOG_FLAGS", log.LstdFlags|log.Lmicroseconds))
	url          string
	email        string
	kindlegen    string
	mercuryToken string
	pmToken      string
	from         string
)

func init() {
	flag.StringVar(&url, "url", "", "the url to send")
	flag.StringVar(&email, "email", "", "the email to send to")
	path, _ := exec.LookPath("kindlegen")
	flag.StringVar(&kindlegen, "kindlegen", path, "the path to the kindlegen binary")

	flag.StringVar(&mercuryToken, "mercury", env.StringDefault("MERCURY_TOKEN", ""), "the mercury token")
	flag.StringVar(&pmToken, "postmark", env.StringDefault("POSTMARK_TOKEN", ""), "the postmark token")
	flag.StringVar(&from, "from", env.StringDefault("FROM", ""), "the from address")

	flag.Parse()
}

func check(args ...string) {
	for _, str := range args {
		if str == "" {
			flag.PrintDefaults()
			os.Exit(1)
		}
	}
}

func main() {
	check(url, email, mercuryToken, pmToken, from, kindlegen)
	app := tinderizer.New(mercuryToken, pmToken, from, kindlegen, logger)
	app.Run(1)

	job, err := J.New(email, url)
	if err != nil {
		log.Fatalf("failed building job: %s", err)
	}

	app.Queue(*job)
	app.Shutdown()
}
