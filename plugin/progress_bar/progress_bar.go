package progress_bar

import (
	"fmt"
	"strings"
)

type ProgressBar struct {
	config config
}

type config struct {
	sizeProgressBar  int
	sizeUser         int
	countIteration   int
	progressBarStart string
	progressBarEnd   string
	arrowBody        string
	arrowHead        string
	increment        int
	arrowFinal       string
	percentage       int
	currentState     int
	title            string
	debug            bool
}

func NewProgressBar(size int) *ProgressBar {
	return &ProgressBar{
		config: config{
			sizeProgressBar:  40,
			sizeUser:         size,
			countIteration:   0,
			progressBarStart: "[",
			progressBarEnd:   "]",
			arrowBody:        "=",
			arrowHead:        ">",
			increment:        0,
			arrowFinal:       "",
			percentage:       0,
			currentState:     0,
			title:            "",
			debug:            false,
		},
	}
}

func (p *ProgressBar) PrintOptions() {
	fmt.Println(
		"\nsizeProgressBar", p.config.sizeProgressBar,
		"\nsizeUser", p.config.sizeUser,
		"\ncountIteration", p.config.countIteration,
		"\nprogressBarStart", p.config.progressBarStart,
		"\nprogressBarEnd", p.config.progressBarEnd,
		"\narrowBody", p.config.arrowBody,
		"\narrowHead", p.config.arrowHead,
		"\nincrement", p.config.increment,
		"\narrowFinal", p.config.arrowFinal,
		"\npercentage", p.config.percentage,
		"\ncurrentState", p.config.currentState,
		"\ntitle", p.config.title,
		"\ndebug", p.config.debug,
	)
}

func (p *ProgressBar) OptionSizeProgressBar(size int) *ProgressBar {
	p.config.sizeProgressBar = size
	return p
}

func (p *ProgressBar) OptioProgressBarStart(start string) *ProgressBar {
	p.config.progressBarStart = start
	return p
}

func (p *ProgressBar) OptioProgressBarEnd(end string) *ProgressBar {
	p.config.progressBarEnd = end
	return p
}

func (p *ProgressBar) OptioArrowBody(arrow string) *ProgressBar {
	p.config.arrowBody = arrow
	return p
}

func (p *ProgressBar) OptioArrowHead(arrow string) *ProgressBar {
	p.config.arrowHead = arrow
	return p
}

func (p *ProgressBar) OptionTitle(title string) *ProgressBar {
	p.config.title = title
	return p
}

func (p *ProgressBar) OptionDebug(debug bool) *ProgressBar {
	p.config.debug = debug
	return p
}

func (p *ProgressBar) initValues() {
	if p.config.sizeProgressBar < p.config.sizeUser {
		p.config.sizeProgressBar = p.config.sizeUser
	}

	if p.config.increment == 0 {
		if p.config.sizeProgressBar != 0 && p.config.sizeUser != 0{
			p.config.increment = p.config.sizeProgressBar / p.config.sizeUser
		}
	}
}

func (p *ProgressBar) render() {
	if p.config.debug {
		p.debug()
	}

	fmt.Printf("\r\a%s", fmt.Sprintf("%s [%d/%d] %d %% %s %s %s", p.config.title, p.config.countIteration, p.config.sizeUser, p.config.percentage, p.config.progressBarStart, p.config.arrowFinal, p.config.progressBarEnd))
}

func (p *ProgressBar) debug() {
	fmt.Println("\nsizeProgressBar", p.config.sizeProgressBar,
		"\nsizeUser", p.config.sizeUser,
		"\nincrement", p.config.increment,
		"\nsizeBlank", p.config.sizeUser*p.config.increment,
		"\ncountIteration", p.config.countIteration,
		"\ncurrentState", p.config.currentState,
		"\npercentage", p.config.percentage,
	)
}

func (p *ProgressBar) PrintBlankProgressBar() *ProgressBar {
	p.initValues()
	p.config.arrowFinal = ""

	p.config.arrowFinal += strings.Repeat(" ", p.config.sizeUser*p.config.increment)

	p.render()
	return p
}

func (p *ProgressBar) Add() {
	p.initValues()
	if p.config.percentage != 100 {
		p.incrementProgressBar()
		p.render()
	}
}

func (p *ProgressBar) incrementProgressBar() {
	p.config.arrowFinal = ""
	p.config.currentState += p.config.increment

	p.config.arrowFinal += strings.Repeat(p.config.arrowBody, p.config.currentState-1) + p.config.arrowHead + strings.Repeat(" ", (p.config.sizeUser*p.config.increment)-p.config.currentState)

	if p.config.countIteration < p.config.sizeUser {
		p.config.countIteration++
	}

	p.config.percentage = p.config.countIteration * (100 / p.config.sizeUser)

	if len(p.config.arrowFinal) == p.config.currentState {
		p.config.percentage = 100
		p.config.currentState = p.config.sizeProgressBar
		p.config.progressBarEnd += "\n"
	}
}
