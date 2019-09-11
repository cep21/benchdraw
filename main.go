package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/cep21/benchdraw/internal"

	"github.com/pkg/errors"
)

type Application struct {
	benchreader internal.BenchmarkReader
	filter      internal.Filter
	grouper     internal.Grouper
	plotter     internal.Plotter
	fs          flag.FlagSet
	config      config
	parameters  []string
	log         internal.Logger
	osExit      func(int)
}

type config struct {
	filter string
	title  string
	group  string
	plot   string
	x      string
	y      string
	input  string
	output string
	format string
}

func filterEmpty(s []string) []string {
	ret := make([]string, 0, len(s))
	for _, i := range s {
		if len(i) > 0 {
			ret = append(ret, i)
		}
	}
	return ret
}

func (c config) parse() (*parsedConfig, error) {
	ret := parsedConfig{
		title:       c.title,
		filters:     toFilterPairs(c.filter),
		group:       filterEmpty(strings.Split(c.group, "/")),
		imageFormat: c.format,
		y:           c.y,
		x:           c.x,
	}
	if ret.title == "" {
		ret.title = c.filter
	}
	pt, err := internal.ToPlotType(c.plot)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to understand plot type %s", c.plot)
	}
	ret.plot = pt
	if c.input == "-" || c.input == "" {
		ret.input = os.Stdin
	} else {
		f, err := os.Open(c.input)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to open file for reading %s", c.input)
		}
		ret.input = f
		ret.onClose = append(ret.onClose, f.Close)
	}
	if c.output == "-" || c.output == "" {
		ret.output = os.Stdout
	} else {
		f, err := os.Create(c.output)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to open file for writing %s", c.output)
		}
		ret.output = f
		ret.onClose = append(ret.onClose, f.Close)
	}
	return &ret, nil
}

func toFilterPairs(s string) []internal.FilterPair {
	parts := strings.Split(s, "/")
	ret := make([]internal.FilterPair, 0, len(parts))
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 1 {
			ret = append(ret, internal.FilterPair{
				Key: p,
			})
		} else {
			ret = append(ret, internal.FilterPair{
				Key:   kv[0],
				Value: kv[1],
			})
		}
	}
	return ret
}

type parsedConfig struct {
	title   string
	filters []internal.FilterPair
	group   []string
	plot    internal.PlotType
	x       string
	y       string
	input   io.Reader
	output  io.Writer

	onClose     []func() error
	imageFormat string
}

func (p *parsedConfig) String() string {
	return fmt.Sprintf("[title=%s x=%s y=%s]", p.title, p.x, p.y)
}

func (p *parsedConfig) Close() error {
	var ret error
	for _, c := range p.onClose {
		if err := c(); err != nil {
			ret = err
		}
	}
	return ret
}

var mainInstance = &Application{
	parameters: os.Args[1:],
	log: internal.Logger{
		Logger: log.New(os.Stderr, "benchdraw", log.LstdFlags),
	},
	osExit: os.Exit,
}

func (a *Application) main() {
	if err := a.run(); err != nil {
		a.log.Log(0, "unable to run application: %s", err.Error())
		a.osExit(1)
		return
	}
}

func (a *Application) run() error {
	if err := a.setupFlags(); err != nil {
		return errors.Wrap(err, "unable to setup flags")
	}
	a.log.Log(1, "application startup")
	pcfg, err := a.config.parse()
	a.log.Log(1, "finished config parsing")
	if err != nil {
		return errors.Wrap(err, "unable to parse config")
	}
	a.log.Log(2, "parsed config: %s", pcfg)
	defer func() {
		if err := pcfg.Close(); err != nil {
			a.log.Log(1, "unable to shutdown config: %s", err)
		}
	}()
	run, err := a.benchreader.ReadBenchmarks(pcfg.input)
	a.log.Log(3, "benchmarks: %s", run)
	if err != nil {
		return errors.Wrap(err, "unable to read benchmark data")
	}
	filteredResults := a.filter.FilterBenchmarks(run.Results, pcfg.filters, pcfg.y)
	a.log.Log(3, "filtered Results: %s", filteredResults)
	uniqueKeys := filteredResults.UniqueValuesForKey(pcfg.x)
	a.log.Log(3, "uniqueKeys: %s", uniqueKeys)
	var groupSet internal.StringSet
	for _, g := range pcfg.group {
		groupSet.Add(g)
	}
	// Each group is a line in our graph
	a.log.Log(3, "groupSet: %v", groupSet)
	grouped := a.grouper.GroupBenchmarks(filteredResults, groupSet, pcfg.x)
	a.log.Log(3, "grouped: %v", grouped)
	grouped.Normalize()
	a.log.Log(3, "normalize: %v", grouped)

	plotLines := make([]internal.PlotLine, 0, len(grouped))
	for _, g := range grouped {
		// For this line in our graph, compute the X Values
		allVals := g.ValuesByX(pcfg.x, pcfg.y, uniqueKeys)
		pl := internal.PlotLine{
			Name:   g.NominalLineName(grouped.AllSingleKey()),
			Values: allVals,
		}
		a.log.Log(3, "nominal=%v plot=%v", pl.Name, pl)
		plotLines = append(plotLines, pl)
		a.log.Log(3, "plot line: %v", pl)
	}
	return a.plotter.Plot(a.log, pcfg.output, pcfg.imageFormat, pcfg.plot, pcfg.title, pcfg.x, pcfg.y, plotLines, uniqueKeys)
}

func (a *Application) setupFlags() error {
	a.fs.StringVar(&a.config.plot, "plot", "bar", "Which picture type to plot.  Valid Values [bar,box]")
	a.fs.StringVar(&a.config.filter, "filter", "", "Filter which benchmarks to graph.  See README for filter syntax")
	a.fs.StringVar(&a.config.title, "title", "", "A title for your graph.  If empty, will use filter")
	a.fs.StringVar(&a.config.group, "group", "", "Pick benchmarks tags to group together")
	a.fs.StringVar(&a.config.x, "x", "", "Pick unit for the X axis")
	a.fs.StringVar(&a.config.y, "y", "ns/op", "Pick unit for the Y axis")
	a.fs.StringVar(&a.config.input, "input", "-", "Input file to read from.  - means stdin")
	a.fs.StringVar(&a.config.output, "output", "-", "Output file to write to.  - means stdout")
	a.fs.StringVar(&a.config.format, "format", "svg", "Which image format to render.  Must be supported by gonum/plot.  You probably want the default.")
	a.fs.IntVar(&a.log.Verbosity, "v", 0, "Higher the Value, the more verbose the output.  Max Value is 4")
	if err := a.fs.Parse(a.parameters); err != nil {
		return errors.Wrap(err, "unable to parse cli parameters")
	}
	return nil
}

func main() {
	mainInstance.main()
}
