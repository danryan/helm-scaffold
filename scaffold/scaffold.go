package scaffold

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/spf13/afero"
	jww "github.com/spf13/jwalterweatherman"
	"k8s.io/helm/pkg/chartutil"
)

type Engine struct {
	Template Template
	Values   chartutil.Values
	ChartFs  afero.Fs
	Options  *EngineOptions
}

type EngineOptions struct {
	LeftDelim  string
	RightDelim string
	ChartPath  string
	Verbose    bool
	DryRun     bool
	Prefix     string
}

type Template struct {
	Prefix  string
	Kind    string
	Ext     string
	Content string
	// Values   chartutil.Values
}

func New(tpl Template, vals chartutil.Values, opts *EngineOptions) *Engine {
	engine := &Engine{
		Template: tpl,
		Values:   vals,
		Options:  opts,
	}

	engine.ChartFs = afero.NewBasePathFs(afero.NewOsFs(), engine.Options.ChartPath)

	return engine
}

// Render creates, parses and renders the template
func (e *Engine) Render() error {
	tpl := template.New("gotpl").Funcs(sprig.TxtFuncMap())
	tpl.Delims(e.Options.LeftDelim, e.Options.RightDelim)

	t, err := tpl.Parse(e.Template.Content)
	if err != nil {
		return err
	}
	// spew.Dump(e.Template)
	if e.Options.Verbose {
		e.verboseRender(t)
	}
	if !e.Options.DryRun {
		e.fileRender(t)
	}

	return nil
}

func (t *Template) filename() string {
	if t.Prefix != "" {
		return fmt.Sprintf("%s-%s.%s", t.Prefix, t.Kind, t.Ext)
	}
	return fmt.Sprintf("%s.%s", t.Kind, t.Ext)
}

func (e *Engine) verboseRender(tpl *template.Template) error {
	err := tpl.Execute(os.Stdout, e.Values)
	if err != nil {
		return err
	}
	return nil
}

func (e *Engine) fileRender(tpl *template.Template) error {
	// f, err := e.ChartFs.Open(path.Join("templates", e.Template.filename()))
	path := path.Join(e.Options.ChartPath, "templates", e.Template.filename())
	jww.INFO.Printf("Writing template to %s", path)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Sync()

	w := bufio.NewWriter(f)
	defer w.Flush()

	if err := tpl.Execute(w, e.Values); err != nil {
		return fmt.Errorf("Error: error executing template: %s", err)
	}

	return nil
}
