package internal

import (
	"io"

	"github.com/pkg/errors"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

// Plotter knows how to draw a picture to a writer
type Plotter struct {
}

// ToPlotType converts a string name to a known plot type
func ToPlotType(s string) (PlotType, error) {
	if s == "" || s == "bar" {
		return PlotTypeBar, nil
	}
	if s == "line" {
		return PlotTypeLine, nil
	}
	return PlotType(0), errors.New("unknown plot type " + s)
}

type PlotType int

const (
	_ PlotType = iota
	// PlotTypeBar is a bar graph
	PlotTypeBar
	// PlotTypeLine is a line graph
	PlotTypeLine
)

// Plot will write to out this plot.
func (l *Plotter) Plot(log Logger, out io.Writer, imgFormat string, pt PlotType, title string, x string, y string, lines []PlotLine, uniqueKeys OrderedStringSet) error {
	p, err := l.createPlot(log, pt, title, x, y, lines, uniqueKeys.Order)
	if err != nil {
		return errors.Wrap(err, "unable to make plot")
	}
	if err := l.savePlot(out, p, imgFormat, lines, uniqueKeys); err != nil {
		return errors.Wrap(err, "unable to save plot")
	}
	return nil
}

// PlotLine is a line to plot.  It has a name (used in the legend) and values for each x index.  It assumes integer
// indexes.
type PlotLine struct {
	Name   string
	Values [][]float64
}

func (l *Plotter) savePlot(out io.Writer, p *plot.Plot, imageFormat string, lines []PlotLine, set OrderedStringSet) error {
	x := float64(30*(len(lines))*(len(set.Items)) + 290)
	wt, err := p.WriterTo(vg.Points(x), vg.Points(x/2), imageFormat)
	if err != nil {
		return errors.Wrap(err, "unable to make plot writer")
	}
	if _, err := wt.WriteTo(out); err != nil {
		return errors.Wrap(err, "unable to write plotter to output")
	}
	return nil
}

func (l *Plotter) createPlot(log Logger, pt PlotType, title string, x string, y string, lines []PlotLine, nominalX []string) (*plot.Plot, error) {
	p, err := plot.New()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create initial plot")
	}
	p.Title.Text = title
	p.Y.Label.Text = y
	p.X.Label.Text = x
	log.Log(2, "nominal x: %v", nominalX)
	p.NominalX(nominalX...)
	p.Legend.Top = true
	for i, line := range lines {
		pl, err := l.makePlotter(log, pt, lines, line, i)
		if err != nil {
			return nil, errors.Wrap(err, "unable to make plotter")
		}
		p.Add(pl)
		if asT, ok := pl.(plot.Thumbnailer); ok {
			p.Legend.Add(line.Name, asT)
		}
	}
	return p, nil
}

func (l *Plotter) addBar(log Logger, line PlotLine, offset int, numLines int) (*plotter.BarChart, error) {
	w := vg.Points(30)
	log.Log(2, "adding line %s", line.Name)
	groupValues := aggregatePlotterValues(line.Values, meanAggregation)
	log.Log(2, "Values: %v", groupValues)
	bar, err := plotter.NewBarChart(plotter.YValues{XYer: groupValues}, w)
	if err != nil {
		return nil, errors.Wrap(err, "unable to make bar chart")
	}
	bar.LineStyle.Width = 0
	bar.Offset = w * vg.Points(float64(numLines/-2+offset))
	bar.Color = plotutil.Color(offset)
	return bar, nil
}

func (l *Plotter) addLine(log Logger, line PlotLine, offset int) (*plotter.Line, error) {
	log.Log(2, "adding line %s", line.Name)
	groupValues := aggregatePlotterValues(line.Values, meanAggregation)
	log.Log(2, "Values: %v", groupValues)
	pline, err := plotter.NewLine(groupValues)
	if err != nil {
		return nil, errors.Wrap(err, "unable to make bar chart")
	}
	pline.LineStyle.Width = 1
	pline.Color = plotutil.Color(offset)
	return pline, nil
}

func (l *Plotter) makePlotter(log Logger, pt PlotType, lines []PlotLine, line PlotLine, index int) (plot.Plotter, error) {
	if pt == PlotTypeBar {
		return l.addBar(log, line, index, len(lines))
	}
	return l.addLine(log, line, index)
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
