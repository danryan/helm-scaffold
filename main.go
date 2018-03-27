package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

var (
	flagVerbose bool
	flagForce   bool
)

func main() {
	cmd := &cobra.Command{
		Use:   "scaffold [flags] TYPE [CHART]",
		Short: fmt.Sprintf("create templates from other templates"),
		RunE:  run,
	}

	f := cmd.Flags()
	f.BoolVarP(&flagVerbose, "verbose", "v", false, "render templates to STDOUT also")
	f.BoolVarP(&flagForce, "force", "f", false, "force overwriting templates, even if they already exist")
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("template type is required")
	}
	typ := args[0]

	if len(args) < 2 {
		// assume chart is in current directory
		args = append(args, ".")
	}

	c, err := chartutil.Load(args[1])
	if err != nil {
		return err
	}

	config := &chart.Config{Values: map[string]*chart.Value{}}

	vals, err := chartutil.ToRenderValues(c, config, chartutil.ReleaseOptions{})
	if err != nil {
		return err
	}

	bp := afero.NewBasePathFs(afero.NewOsFs(), "./templates")
	files, err := afero.ReadDir(bp, ".")
	if err != nil {
		return err
	}

	outputs := make(map[string]string, len(files))

	for _, file := range files {
		name := strings.Split(file.Name(), ".")

		contents, err := afero.ReadFile(bp, file.Name())
		if err != nil {
			return err
		}

		outputs[name[0]] = string(contents)

		// t.Execute(os.Stdout, vals)
	}
	tpl := template.New(typ)
	tpl.Delims("((", "))")
	tpl.Parse(outputs[typ])
	tpl.Execute(os.Stdout, vals)
	// tpl := template.New("template")
	// tpl.Delims("((", "))")
	// tpl.Parse(templateMap(t))
	// tpl.Execute(os.Stdout, vals)

	return nil
}
