package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/danryan/helm-scaffold/scaffold"
	"github.com/spf13/cobra"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

var (
	flagVerbose     bool
	flagForce       bool
	flagDryRun      bool
	flagTemplateDir string
	flagName        string
	flagDelims      []string
)

var cmd = &cobra.Command{
	Use:   "scaffold [flags] TYPE [CHART]",
	Short: fmt.Sprintf("create templates from other templates"),
	RunE:  run,
}

func main() {
	f := cmd.Flags()
	f.BoolVarP(&flagVerbose, "verbose", "v", false, "also render templates to STDOUT")
	f.BoolVarP(&flagForce, "force", "f", false, "force overwriting templates, even if they already exist")
	f.BoolVarP(&flagDryRun, "dry-run", "r", false, "only run through the process; do not write any files")
	f.StringVar(&flagTemplateDir, "templates", path.Join(os.Getenv("HELM_PLUGIN_DIR"), "templates"), "directory to look for templates")
	f.StringVarP(&flagName, "name", "n", "", "name of the generated template")
	f.StringSliceVarP(&flagDelims, "delims", "d", []string{"<%", "%>"}, "default left and right template delimiters")
	cmd.MarkFlagRequired("name")

	// f.BoolVarP(&flagVerbose, "verbose", "v", false, "render templates to STDOUT also")
	if err := cmd.Execute(); err != nil {
		spew.Dump(err)
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

	chartDir, err := filepath.Abs(args[1])
	if err != nil {
		return errors.New("Error: chart directory does not exist")
	}

	c, err := chartutil.Load(chartDir)
	if err != nil {
		return err
	}

	config := &chart.Config{Values: map[string]*chart.Value{}}

	vals, err := chartutil.ToRenderValues(c, config, chartutil.ReleaseOptions{})
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(flagTemplateDir)
	if err != nil {
		return err
	}

	templates := make(map[string]scaffold.Template, len(files))

	for _, file := range files {
		// skip templates that begin with `__`.
		if strings.HasPrefix(file.Name(), "__") {
			continue
		}

		contents, err := ioutil.ReadFile(path.Join(flagTemplateDir, file.Name()))
		if err != nil {
			return err
		}
		file := strings.Split(file.Name(), ".")
		kind := file[0]
		ext := file[1]
		tpl := scaffold.Template{
			Name:    flagName,
			Kind:    kind,
			Ext:     ext,
			Content: string(contents),
		}
		templates[kind] = tpl
	}
	options := &scaffold.EngineOptions{
		Verbose:    flagVerbose,
		DryRun:     flagDryRun,
		ChartPath:  chartDir,
		LeftDelim:  flagDelims[0],
		RightDelim: flagDelims[1],
	}
	spew.Dump(options)

	scaffold := scaffold.New(templates[kind], vals, options)

	if err := scaffold.Render(); err != nil {
		return fmt.Errorf("Error: rendering template: %s", err)
	}
	return nil
}
