package ui

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/briandowns/spinner"
	ghterm "github.com/cli/go-gh/pkg/term"
	"github.com/mattn/go-colorable"
)

// ColorMode represents coloring mode for the prompt.
type ColorMode string

const (
	// ColorAuto indicates the coloring is turned on/off according to the terminal settings.
	ColorAuto ColorMode = "auto"
	// ColorAlways indicates the output is always colored.
	ColorAlways ColorMode = "always"
	// ColorNever indicates the output is never colored.
	ColorNever ColorMode = "never"
)

// Color represents the coloring mode.
var Color string

// NewPrompter returns a new Prompter.
func NewPrompter() Prompter {
	var enableColor bool
	switch ColorMode(Color) {
	case ColorAuto:
		enableColor = ghterm.FromEnv().IsColorEnabled()
	case ColorNever:
		enableColor = false
	case ColorAlways:
		enableColor = true
	}
	if !enableColor {
		core.DisableColor = true
	}

	prompter := &surveyPrompter{
		in:  os.Stdin,
		out: os.Stdout,
		err: colorable.NewColorable(os.Stderr),

		warnIcon: "!",
		errIcon:  "X",
		doneIcon: "✓",

		spinnerCharID:   14,
		spinnerDuration: 50 * time.Millisecond,

		cf: &ColorFormatter{enabled: enableColor},
	}

	// for windows, translate ANSI escape sequence.
	if colorableStdout := colorable.NewColorable(os.Stdout); colorableStdout != os.Stdout {
		prompter.out = &fdWriter{
			fd:     os.Stdout.Fd(),
			Writer: colorableStdout,
		}
	}
	return prompter
}

// Prompter is an interface that provides functions for the prompt.
type Prompter interface {
	Ask(q survey.Prompt, response interface{}, opts ...survey.AskOpt) error

	Printf(format string, a ...any)
	Println(a ...any)
	Warn(msg string)
	Error(msg string)
	Done(msg string)

	StartSpinner(msg string, newline bool) *spinner.Spinner
	ColorFormatter() *ColorFormatter
}

var _ Prompter = &surveyPrompter{}

type surveyPrompter struct {
	in  terminal.FileReader
	out terminal.FileWriter
	err io.Writer

	warnIcon string
	errIcon  string
	doneIcon string

	spinnerCharID   int
	spinnerDuration time.Duration
	spinner         *spinner.Spinner

	cf *ColorFormatter
}

// Ask gets a user input depends on survey.Prompt settings.
func (p *surveyPrompter) Ask(q survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	if p.spinner != nil {
		// if spinner exists, stop it to prevent interference with prompts.
		p.spinner.Stop() // Stop() can be called multiple times.
	}
	opts = append(opts, survey.WithStdio(p.in, p.out, p.err))
	return survey.AskOne(q, response, opts...)
}

// Printf formats string and prints it.
func (p *surveyPrompter) Printf(format string, a ...any) {
	_, _ = fmt.Fprintf(p.out, format, a...)
}

// Println prints string with a newline.
func (p *surveyPrompter) Println(a ...any) {
	_, _ = fmt.Fprintln(p.out, a...)
}

// Done prints done message.
func (p *surveyPrompter) Done(msg string) {
	_, _ = fmt.Fprintf(p.out, "%s %s\n", p.cf.Green(p.doneIcon), msg)
}

// Warn prints warning message.
func (p *surveyPrompter) Warn(msg string) {
	_, _ = fmt.Fprintln(p.err, p.cf.Yellow(fmt.Sprintf("%s %s", p.warnIcon, msg)))
}

// Error prints error message.
func (p *surveyPrompter) Error(msg string) {
	_, _ = fmt.Fprintln(p.err, p.cf.Red(fmt.Sprintf("%s %s", p.errIcon, msg)))
}

// StartSpinner creates and start a spinner. If newline is true, adding newline to the message.
func (p *surveyPrompter) StartSpinner(msg string, newline bool) *spinner.Spinner {
	suffix := "\n"
	if newline {
		suffix = "\n\n"
	}
	s := spinner.New(spinner.CharSets[p.spinnerCharID], p.spinnerDuration)
	s.Prefix = fmt.Sprintf("%s... ", msg)
	s.FinalMSG = fmt.Sprintf("%s... %s%s", msg, p.ColorFormatter().Green("done"), suffix)
	s.Start()

	p.spinner = s
	return s
}

// ColorFormatter returns the ColorFormatter.
func (p *surveyPrompter) ColorFormatter() *ColorFormatter {
	return p.cf
}

var _ terminal.FileWriter = &fdWriter{}

type fdWriter struct {
	io.Writer
	fd uintptr
}

func (w *fdWriter) Fd() uintptr {
	return w.fd
}
