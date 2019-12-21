package pget

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

// Options struct for parse command line arguments
type Options struct {
	Help      bool   `short:"h" long:"help"`
	Version   bool   `short:"v" long:"version"`
	Procs     int    `short:"p" long:"procs"`
	Output    string `short:"o" long:"output"`
	TargetDir string `short:"d" long:"target-dir"`
	Timeout   int    `short:"t" long:"timeout"`
	UserAgent string `short:"u" long:"user-agent"`
	Referer   string `short:"r" long:"referer"`
	Update    bool   `long:"check-update"`
	Trace     bool   `long:"trace"`
}

func (opts *Options) parse(argv []string) ([]string, error) {
	/*p := flags.NewParser(opts, flags.PrintErrors)
	args, err := p.ParseArgs(argv)

	if err != nil {
		os.Stderr.Write(opts.usage())
		return nil, errors.New("invalid command line options: " + err.Error())
	}

	return args, nil*/
}

func (opts Options) usage() []byte {
	buf := bytes.Buffer{}

	fmt.Fprintf(&buf, msg+
		`Usage: pget [options] URL
  Options:
  -h,  --help                   print usage and exit
  -v,  --version                display the version of pget and exit
  -p,  --procs <num>            split ratio to download file
  -o,  --output <filename>      output file to <filename>
  -d,  --target-dir <path>    	path to the directory to save the downloaded file, filename will be taken from url
  -t,  --timeout <seconds>      timeout of checking request in seconds
  -u,  --user-agent <agent>     identify as <agent>
  -r,  --referer <referer>      identify as <referer>
  --check-update                check if there is update available
  --trace                       display detail error messages
`)
	return buf.Bytes()
}
