# Helm Scaffold

This is a Helm plugin to help chart developers write templates faster by generating scaffolds from templates. Yes, templates for your templates.

## Usage

Generate templates from templates

```
$ helm scaffold [flags] TYPE CHART]
```

### Flags

```
  -d, --delims stringSlice   default left and right template delimiters (default [<%,%>])
  -r, --dry-run              only run through the process; do not write any files
  -f, --force                force overwriting templates, even if they already exist
  -h, --help                 help for scaffold
      --set string      set values on the command line (can specify multiple times or separate values with commas: key1=val1,key2=val2)
      --templates path     directory to look for templates (default "/Users/dan/.helm/plugins/helm-scaffold/templates")
  -v, --verbose              also render templates to STDOUT
```

### Install

```
$ helm plugin install https://github.com/danryan/helm-scaffold
```

This will fetch the latest binary release of `helm scaffold` and install it.