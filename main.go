package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/cep21/benchparse"
	"github.com/pkg/errors"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

type Application struct {
	fs         flag.FlagSet
	config     config
	parameters []string
	log        Logger
	osExit     func(int)
}

type Logger struct {
	verbosity int
	logger    *log.Logger
}

func (l *Logger) Log(verbosity int, msg string, fmtArgs ...interface{}) {
	if l.verbosity >= verbosity {
		l.logger.Printf(msg, fmtArgs...)
	}
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
	pt, err := toPlotType(c.plot)
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

func toFilterPairs(s string) []filterPair {
	parts := strings.Split(s, "/")
	ret := make([]filterPair, 0, len(parts))
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		kv := strings.SplitN(p, "=", 2)
		if len(kv) == 1 {
			ret = append(ret, filterPair{
				key: p,
			})
		} else {
			ret = append(ret, filterPair{
				key:   kv[0],
				value: kv[1],
			})
		}
	}
	return ret
}

func toPlotType(s string) (plotType, error) {
	if s == "" || s == "bar" {
		return plotTypeBar, nil
	}
	if s == "line" {
		return plotTypeLine, nil
	}
	return plotType(0), errors.New("unknown plot type " + s)
}

type plotType int

const (
	_ plotType = iota
	plotTypeBar
	plotTypeLine
)

type filterPair struct {
	key   string
	value string
}

type parsedConfig struct {
	title   string
	filters []filterPair
	group   []string
	plot    plotType
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
	log: Logger{
		logger: log.New(os.Stderr, "benchdraw", log.LstdFlags),
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
	run, err := a.readBenchmarks(pcfg)
	a.log.Log(3, "benchmarks: %s", run)
	if err != nil {
		return errors.Wrap(err, "unable to read benchmark data")
	}
	filteredResults := filterBenchmarks(run.Results, pcfg.filters, pcfg.y)
	a.log.Log(3, "filtered results: %s", filteredResults)
	uniqueKeys := uniqueValuesForKey(filteredResults, pcfg.x)
	a.log.Log(3, "uniqueKeys: %s", uniqueKeys)
	var groupSet stringSet
	for _, g := range pcfg.group {
		groupSet.add(g)
	}
	// Each group is a line in our graph
	a.log.Log(3, "groupSet: %v", groupSet)
	grouped := groupBenchmarks(filteredResults, groupSet, pcfg.x)
	a.log.Log(3, "grouped: %v", grouped)
	normalize(grouped)
	a.log.Log(3, "normalize: %v", grouped)

	plotLines := make([]plotLine, 0, len(grouped))
	for _, g := range grouped {
		// For this line in our graph, compute the X values
		allVals := valuesByX(g, pcfg.x, pcfg.y, uniqueKeys)
		pl := plotLine{
			name:   g.nominalLineName(allSingleKey(grouped), a.config.filter),
			values: allVals,
		}
		a.log.Log(3, "nominal=%v plot=%v", pl.name, pl)
		plotLines = append(plotLines, pl)
		a.log.Log(3, "plot line: %v", pl)
	}
	p, err := a.createPlot(pcfg, plotLines, uniqueKeys.order)
	if err != nil {
		return errors.Wrap(err, "unable to make plot")
	}
	if err := savePlot(pcfg, p, plotLines, uniqueKeys); err != nil {
		return errors.Wrap(err, "unable to save plot")
	}
	return nil
}

func allSingleKey(groups []*benchmarkGroup) bool {
	if len(groups) <= 1 {
		return true
	}
	if len(groups[0].values.order) > 1 {
		return false
	}
	expectedKey := groups[0].values.order[0]
	for i := 1; i < len(groups); i++ {
		if len(groups[i].values.order) > 1 {
			return false
		}
		if groups[0].values.order[0] != expectedKey {
			return false
		}
	}
	return true
}

func savePlot(pcfg *parsedConfig, p *plot.Plot, lines []plotLine, set stringSet) error {
	x := float64(30*(len(lines))*(len(set.items)) + 290)
	wt, err := p.WriterTo(vg.Points(x), vg.Points(x/2), pcfg.imageFormat)
	if err != nil {
		return errors.Wrap(err, "unable to make plot writer")
	}
	if _, err := wt.WriteTo(pcfg.output); err != nil {
		return errors.Wrap(err, "unable to write plotter to output")
	}
	return nil
}

type plotLine struct {
	name   string
	values [][]float64
}

func (a *Application) setupFlags() error {
	a.fs.StringVar(&a.config.plot, "plot", "bar", "Which picture type to plot.  Valid values [bar,box]")
	a.fs.StringVar(&a.config.filter, "filter", "", "Filter which benchmarks to graph.  See README for filter syntax")
	a.fs.StringVar(&a.config.title, "title", "", "A title for your graph.  If empty, will use filter")
	a.fs.StringVar(&a.config.group, "group", "", "Pick benchmarks tags to group together")
	a.fs.StringVar(&a.config.x, "x", "", "Pick unit for the X axis")
	a.fs.StringVar(&a.config.y, "y", "ns/op", "Pick unit for the Y axis")
	a.fs.StringVar(&a.config.input, "input", "-", "Input file to read from.  - means stdin")
	a.fs.StringVar(&a.config.output, "output", "-", "Output file to write to.  - means stdout")
	a.fs.StringVar(&a.config.format, "format", "svg", "Which image format to render.  Must be supported by gonum/plot.  You probably want the default.")
	a.fs.IntVar(&a.log.verbosity, "v", 0, "Higher the value, the more verbose the output.  Max value is 4")
	if err := a.fs.Parse(a.parameters); err != nil {
		return errors.Wrap(err, "unable to parse cli parameters")
	}
	return nil
}

func (a *Application) readBenchmarks(cfg *parsedConfig) (*benchparse.Run, error) {
	d := benchparse.Decoder{}
	run, err := d.Decode(cfg.input)
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode benchmark format")
	}
	return run, nil
}

func filterBenchmarks(in []benchparse.BenchmarkResult, filters []filterPair, unit string) []benchparse.BenchmarkResult {
	ret := make([]benchparse.BenchmarkResult, 0, len(in))
	for _, b := range in {
		// Benchmark must have a valid unit
		if _, exists := b.ValueByUnit(unit); !exists {
			continue
		}
		keys := b.AllKeyValuePairs()
		okToAdd := true
		// each filter must pass
		for _, f := range filters {
			val, exists := keys.Contents[f.key]
			if !exists || (f.value != "" && f.value != val) {
				okToAdd = false
				break
			}
		}
		if okToAdd {
			ret = append(ret, b)
		}
	}
	return ret
}

type benchmarkGroup struct {
	values  hashableMap
	results []benchparse.BenchmarkResult
}

func (b *benchmarkGroup) String() string {
	return fmt.Sprintf("vals=%v len_results=%d", b.values, len(b.results))
}

func (b *benchmarkGroup) nominalLineName(singleKey bool, filterName string) string {
	if singleKey && len(b.values.order) > 0 {
		return b.values.values[b.values.order[0]]
	}
	ret := make([]string, 0, len(b.values.order))
	for _, c := range b.values.order {
		ret = append(ret, c+"="+b.values.values[c])
	}
	if len(ret) == 0 {
		return ""
	}
	return "[" + strings.Join(ret, ",") + "]"
}

func makeKeys(r benchparse.BenchmarkResult) hashableMap {
	nameKeys := r.AllKeyValuePairs()
	var ret hashableMap
	for _, k := range nameKeys.Order {
		ret.insert(k, nameKeys.Contents[k])
	}
	return ret
}

func uniqueValuesForKey(in []benchparse.BenchmarkResult, key string) stringSet {
	var ret stringSet
	for _, b := range in {
		keys := makeKeys(b)
		if keyValue, exists := keys.values[key]; exists {
			ret.add(keyValue)
		}
	}
	return ret
}

// each returned benchmarkGroup will aggregate results by unique groups key/value pairs
func groupBenchmarks(in []benchparse.BenchmarkResult, groups stringSet, unit string) []*benchmarkGroup {
	ret := make([]*benchmarkGroup, 0, len(in))
	setMap := make(map[string]*benchmarkGroup)
	for _, b := range in {
		keysMap := makeKeys(b)
		var hm hashableMap
		if len(groups.order) == 0 {
			// Group by everything except unit
			for _, k := range keysMap.order {
				if k != unit {
					hm.insert(k, keysMap.values[k])
				}
			}
		} else {
			for _, ck := range groups.order {
				if configValue, exists := keysMap.values[ck]; exists {
					hm.insert(ck, configValue)
				}
			}
		}
		mapHash := hm.Hash()
		if existing, exists := setMap[mapHash]; exists {
			existing.results = append(existing.results, b)
		} else {
			bg := &benchmarkGroup{
				values:  hm,
				results: []benchparse.BenchmarkResult{b},
			}
			setMap[mapHash] = bg
			ret = append(ret, bg)
		}
	}
	return ret
}

// Normalize modifies in to remove key/value pairs that exist in every group
func normalize(in []*benchmarkGroup) {
	if len(in) == 0 {
		return
	}
	keysToRemove := make([]string, 0, len(in[0].values.values))
	for k, v := range in[0].values.values {
		canRemoveValue := true
	checkRestLoop:
		for i := 1; i < len(in); i++ {
			if !in[i].values.contains(k, v) {
				canRemoveValue = false
				break checkRestLoop
			}
		}
		if canRemoveValue {
			keysToRemove = append(keysToRemove, k)
		}
	}
	for _, k := range keysToRemove {
		for _, i := range in {
			i.values.remove(k)
		}
	}
}

func meanAggregation(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

func valuesByX(in *benchmarkGroup, xDim string, unit string, allValues stringSet) [][]float64 {
	ret := make([][]float64, 0, len(allValues.order))
	for _, v := range allValues.order {
		allVals := make([]float64, 0, len(in.results))
		for _, b := range in.results {
			benchmarkKeys := makeKeys(b)
			if benchmarkKeys.values[xDim] != v {
				continue
			}
			if val, exists := b.ValueByUnit(unit); exists {
				allVals = append(allVals, val)
			}
		}
		ret = append(ret, allVals)
	}
	return ret
}

func (a *Application) addBar(line plotLine, offset int, numLines int) (*plotter.BarChart, error) {
	w := vg.Points(30)
	a.log.Log(2, "adding line %s", line.name)
	groupValues := aggregatePlotterValues(line.values, meanAggregation)
	a.log.Log(2, "values: %v", groupValues)
	bar, err := plotter.NewBarChart(plotter.YValues{XYer: groupValues}, w)
	if err != nil {
		return nil, errors.Wrap(err, "unable to make bar chart")
	}
	bar.LineStyle.Width = 0
	bar.Offset = w * vg.Points(float64(numLines/-2+offset))
	bar.Color = plotutil.Color(offset)
	return bar, nil
}

func (a *Application) addLine(line plotLine, offset int) (*plotter.Line, error) {
	a.log.Log(2, "adding line %s", line.name)
	groupValues := aggregatePlotterValues(line.values, meanAggregation)
	a.log.Log(2, "values: %v", groupValues)
	pline, err := plotter.NewLine(groupValues)
	if err != nil {
		return nil, errors.Wrap(err, "unable to make bar chart")
	}
	pline.LineStyle.Width = 1
	pline.Color = plotutil.Color(offset)
	return pline, nil
}

func (a *Application) makePlotter(cfg *parsedConfig, lines []plotLine, line plotLine, index int) (plot.Plotter, error) {
	if cfg.plot == plotTypeBar {
		return a.addBar(line, index, len(lines))
	} else {
		return a.addLine(line, index)
	}
}

func (a *Application) createPlot(cfg *parsedConfig, lines []plotLine, nominalX []string) (*plot.Plot, error) {
	p, err := plot.New()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create initial plot")
	}
	p.Title.Text = cfg.title
	p.Y.Label.Text = cfg.y
	p.X.Label.Text = cfg.x
	a.log.Log(2, "nominal x: %v", nominalX)
	p.NominalX(nominalX...)
	p.Legend.Top = true
	for i, line := range lines {
		pl, err := a.makePlotter(cfg, lines, line, i)
		if err != nil {
			return nil, errors.Wrap(err, "unable to make plotter")
		}
		p.Add(pl)
		if asT, ok := pl.(plot.Thumbnailer); ok {
			p.Legend.Add(line.name, asT)
		}
	}
	return p, nil
}

func aggregatePlotterValues(f [][]float64, aggregation func([]float64) float64) plotter.XYer {
	var ret plotter.XYs
	for i, x := range f {
		ret = append(ret, plotter.XY{
			X: float64(i),
			Y: aggregation(x),
		})
	}
	return ret
}

func main() {
	mainInstance.main()
}

type stringSet struct {
	items map[string]struct{}
	order []string
}

func (s *stringSet) contains(k string) bool {
	_, exists := s.items[k]
	return exists
}

func (s *stringSet) add(k string) {
	if s.contains(k) {
		return
	}
	if s.items == nil {
		s.items = make(map[string]struct{})
	}
	s.items[k] = struct{}{}
	s.order = append(s.order, k)
}

type hashableMap struct {
	values map[string]string
	order  []string
}

func mustWrite(_ int, err error) {
	if err != nil {
		panic(err)
	}
}

func mustNotError(err error) {
	if err != nil {
		panic(err)
	}
}

func (h *hashableMap) String() string {
	return fmt.Sprintf("%v", h.values)
}

func (h *hashableMap) contains(k string, v string) bool {
	current, exists := h.values[k]
	return exists && current == v
}

func (h *hashableMap) insert(k string, v string) {
	if _, exists := h.values[k]; exists {
		h.remove(k)
	}
	if h.values == nil {
		h.values = make(map[string]string)
	}
	h.values[k] = v
	h.order = append(h.order, k)
}

func (h *hashableMap) remove(k string) {
	if h.values == nil {
		return
	}
	delete(h.values, k)
	for i, o := range h.order {
		if o == k {
			h.order = append(h.order[:i], h.order[i+1:]...)
			return
		}
	}
}

func (h *hashableMap) Hash() string {
	type kv struct {
		k string
		v string
	}
	toSort := make([]kv, 0, len(h.values))
	for k, v := range h.values {
		toSort = append(toSort, kv{k: k, v: v})
	}
	sort.Slice(toSort, func(i, j int) bool {
		return toSort[i].k < toSort[j].k
	})
	var uid strings.Builder
	for _, s := range toSort {
		if uid.Len() != 0 {
			mustWrite(uid.WriteString(string([]byte{0, 0})))
		}
		mustWrite(uid.WriteString(s.k))
		mustNotError(uid.WriteByte(0))
		mustWrite(uid.WriteString(s.v))
	}
	return uid.String()
}
