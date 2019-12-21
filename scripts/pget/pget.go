package pget

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"os"
	"runtime"
	"strings"
)

const (
	version = "0.0.6"
)

// Pget structs
type Pget struct {
	Trace bool
	Utils
	TargetDir  string
	Procs      int
	URLs       []string
	TargetURLs []string
	args       []string
	timeout    int
	useragent  string
	referer    string
}

type ignore struct {
	err error
}

type cause interface {
	Cause() error
}

// New for pget package
func New() *Pget {
	return &Pget{
		Trace:   false,
		Utils:   &Data{},
		Procs:   runtime.NumCPU(), // default
		timeout: 10,
	}
}

// ErrTop get important message from wrapped error message
func (pget Pget) ErrTop(err error) error {
	for e := err; e != nil; {
		switch e.(type) {
		case ignore:
			return nil
		case cause:
			e = e.(cause).Cause()
		default:
			return e
		}
	}

	return nil
}

// Run execute methods in pget package
func (pget *Pget) Run() error {
	if err := pget.Ready(); err != nil {
		return pget.ErrTop(err)
	}

	if err := pget.Checking(); err != nil {
		return errors.New("failed to check header: " + err.Error())
	}

	if err := pget.Download(); err != nil {
		return err
	}

	if err := pget.Utils.BindwithFiles(pget.Procs); err != nil {
		return err
	}

	return nil
}

// Ready method define the variables required to Download.
func (pget *Pget) Ready() error {
	if procs := os.Getenv("GOMAXPROCS"); procs == "" {
		runtime.GOMAXPROCS(pget.Procs)
	}

	var opts Options
	if err := pget.parseOptions(&opts, os.Args[1:]); err != nil {
		return errors.New("failed to parse command line args: " + err.Error())
	}

	if opts.Trace {
		pget.Trace = opts.Trace
	}

	if opts.Procs > 2 {
		pget.Procs = opts.Procs
	}

	if opts.Timeout > 0 {
		pget.timeout = opts.Timeout
	}

	if err := pget.parseURLs(); err != nil {
		return errors.New("failed to parse of url: " + err.Error())
	}

	if opts.Output != "" {
		pget.Utils.SetFileName(opts.Output)
	}

	if opts.UserAgent != "" {
		pget.useragent = opts.UserAgent
	}

	if opts.Referer != "" {
		pget.referer = opts.Referer
	}

	if opts.TargetDir != "" {
		info, err := os.Stat(opts.TargetDir)
		if err != nil {
			if !os.IsNotExist(err) {
				return errors.New("target dir is invalid: " + err.Error())
			}

			if err := os.MkdirAll(opts.TargetDir, 0755); err != nil {
				return errors.New("failed to create diretory at " + opts.TargetDir + ", error: " + err.Error())
			}

		} else if !info.IsDir() {
			return errors.New("target dir is not a valid directory")
		}
		opts.TargetDir = strings.TrimSuffix(opts.TargetDir, "/")
	}
	pget.TargetDir = opts.TargetDir

	return nil
}

// Error for options: version, usage
func (i ignore) Error() string {
	return i.err.Error()
}

func (i ignore) Cause() error {
	return i.err
}

func (pget *Pget) parseOptions(opts *Options, argv []string) error {

	if len(argv) == 0 {
		os.Stdout.Write(opts.usage())
		return nil
	}

	o, err := opts.parse(argv)
	if err != nil {
		return errors.New("failed to parse command line options: " + err.Error())
	}

	if opts.Help {
		os.Stdout.Write(opts.usage())
		return nil
	}

	if opts.Version {
		return nil
	}

	pget.args = o

	return nil
}

func (pget *Pget) parseURLs() error {

	// find url in args
	for _, argv := range pget.args {
		if govalidator.IsURL(argv) {
			pget.URLs = append(pget.URLs, argv)
		}
	}

	if len(pget.URLs) < 1 {
		fmt.Fprintf(os.Stdout, "Please input url separate with space or newline\n")
		fmt.Fprintf(os.Stdout, "Start download at ^D\n")

		// scanning url from stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			scan := scanner.Text()
			urls := strings.Split(scan, " ")
			for _, url := range urls {
				if govalidator.IsURL(url) {
					pget.URLs = append(pget.URLs, url)
				}
			}
		}

		if err := scanner.Err(); err != nil {
			return errors.New("failed to parse url from stdin: " + err.Error())
		}

		if len(pget.URLs) < 1 {
			return errors.New("urls not found in the arguments passed")
		}
	}

	return nil
}
