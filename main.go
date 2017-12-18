package main

import (
	"errors"
	"fmt"
	"html/template"
	"os"

	"github.com/spf13/cobra"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

func main() {
	cmd := &cobra.Command{
		Use:   "scaffold [flags] TYPE CHART",
		Short: fmt.Sprintf("create templates from other templates"),
		RunE:  run,
	}

	// f := cmd.Flags()

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("chart is required")
	}
	c, err := chartutil.Load(args[0])
	if err != nil {
		return err
	}

	if len(args) < 2 {
		return errors.New("template type is required")
	}
	t := args[1]

	config := &chart.Config{Values: map[string]*chart.Value{}}

	vals, err := chartutil.ToRenderValues(c, config, chartutil.ReleaseOptions{})
	if err != nil {
		return err
	}

	tpl := template.New("template")
	tpl.Delims("%%", "%%")
	tpl.Parse(templateMap(t))
	tpl.Execute(os.Stdout, vals)

	return nil
}
