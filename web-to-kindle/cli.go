package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/darkhelmet/env"
	"github.com/darkhelmet/tinderizer"
	"github.com/darkhelmet/tinderizer/job"
)

var (
	url                 = flag.String("url", "", "The URL to download")
	mercuryToken        = flag.String("mercury", env.StringDefault("MERCURY_TOKEN", ""), "The Mercury API token to use")
	postmarkToken       = flag.String("postmark", env.StringDefault("POSTMARK_TOKEN", ""), "The Postmark API token to use")
	kindlegenBinaryPath = flag.String("kindlegen", env.StringDefault("KINDLEGEN", mustLookPath(fmt.Sprintf("kindlegen-%s", runtime.GOOS))), "The path to the Kindlegen binary")
	from                = flag.String("from", "kindle@darkhelmetlive.com", "The FROM email address configured in Postmark")
	to                  = flag.String("to", "", "The email address to send to")
)

func mustLookPath(path string) string {
	path, err := exec.LookPath(path)
	if err != nil {
		return ""
	}
	return path
}

func init() {
	flag.Parse()
}

func main() {
	logger := log.New(os.Stdout, "[tinderizer] ", log.LstdFlags|log.Lmicroseconds)
	app := tinderizer.New(*mercuryToken, *postmarkToken, *from, *kindlegenBinaryPath, logger)
	app.RunOne(false)
	jerb, _ := job.New(*to, *url)
	app.Queue(*jerb)
	app.Shutdown()
}
