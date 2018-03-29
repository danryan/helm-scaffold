package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/danryan/helm-scaffold/scaffold"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/strvals"
)

type runVars struct {
	verbose     bool
	force       bool
	dryRun      bool
	templateDir string
	prefix      string
	delims      []string
	chartPath   string
	values      []string
}

var cmd = &cobra.Command{
	Use:   "scaffold TYPE [CHART] [flags]",
	Short: fmt.Sprintf("create templates from other templates"),
	RunE:  run,
}

var r = &runVars{}

func main() {
	f := cmd.Flags()
	f.BoolVarP(&r.verbose, "verbose", "v", false, "also render templates to STDOUT")
	f.BoolVarP(&r.force, "force", "f", false, "force overwriting templates, even if they already exist")
	f.BoolVarP(&r.dryRun, "dry-run", "r", false, "only run through the process; do not write any files")
	f.StringVar(&r.templateDir, "templates", path.Join(os.Getenv("HELM_PLUGIN_DIR"), "templates"), "directory to look for templates")
	f.StringVarP(&r.prefix, "prefix", "p", "", "prefix for the generated template filename (`helm scaffold configmap chart -p foo -> foo-configmap.yaml`")
	f.StringSliceVarP(&r.delims, "delims", "d", []string{"<%", "%>"}, "default left and right template delimiters")
	f.StringArrayVar(&r.values, "set", []string{}, "set values on the command line (can specify multiple times or separate values with commas: key1=val1,key2=val2)")
	// cmd.MarkFlagRequired("name")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("Error: template type is required")
	}
	kind := args[0]

	if len(args) < 2 {
		// assume chart is in current directory
		args = append(args, ".")
	}

	cp, err := filepath.Abs(args[1])
	if err != nil {
		return errors.New("Error: chart directory does not exist")
	}
	r.chartPath = cp

	c, err := chartutil.Load(r.chartPath)
	if err != nil {
		return err
	}

	vv, err := vals(r.values)
	if err != nil {
		return err
	}
	config := &chart.Config{Raw: string(vv), Values: map[string]*chart.Value{}}

	vals, err := chartutil.ToRenderValues(c, config, chartutil.ReleaseOptions{})
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(r.templateDir)
	if err != nil {
		return err
	}

	templates := make(map[string]scaffold.Template, len(files))

	for _, file := range files {
		// skip templates that begin with `__`.
		if strings.HasPrefix(file.Name(), "__") {
			continue
		}

		contents, err := ioutil.ReadFile(path.Join(r.templateDir, file.Name()))
		if err != nil {
			return err
		}
		file := strings.Split(file.Name(), ".")
		kind := file[0]
		ext := file[1]
		tpl := scaffold.Template{
			Kind:    kind,
			Ext:     ext,
			Content: string(contents),
		}
		if r.prefix != "" {
			tpl.Prefix = r.prefix
		}
		templates[kind] = tpl
	}
	options := &scaffold.EngineOptions{
		Verbose:    r.verbose,
		DryRun:     r.dryRun,
		ChartPath:  r.chartPath,
		LeftDelim:  r.delims[0],
		RightDelim: r.delims[1],
		Prefix:     r.prefix,
	}

	scaffold := scaffold.New(templates[kind], vals, options)
	if err := scaffold.Render(); err != nil {
		return fmt.Errorf("Error: rendering template: %s", err)
	}

	return nil
}

func vals(values []string) ([]byte, error) {
	base := map[string]interface{}{}

	for _, value := range values {
		if err := strvals.ParseInto(value, base); err != nil {
			return []byte{}, fmt.Errorf("failed parsing --set data: %s", err)
		}
	}

	return yaml.Marshal(base)
}
