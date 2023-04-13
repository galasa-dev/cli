## galasactl runs prepare

prepares a list of tests

### Synopsis

Prepares a list of tests from a test catalog providing specific overrides if required

```
galasactl runs prepare [flags]
```

### Options

```
      --append             Append tests to existing portfolio
      --bundle strings     bundles of which tests will be selected from, bundles are selected if the name contains this string, or if --regex is specified then matches the regex
      --class strings      test class names, for building a portfolio when a stream/test catalog is not available. The format of each entry is osgi-bundle-name/java-class-name . Java class names are fully qualified. No .class suffix is needed.
  -h, --help               help for prepare
      --override strings   overrides to be sent with the tests (overrides in the portfolio will take precedence)
      --package strings    packages of which tests will be selected from, packages are selected if the name contains this string, or if --regex is specified then matches the regex
  -p, --portfolio string   portfolio to add tests to
      --regex              Test selection is performed by using regex
  -s, --stream string      test stream to extract the tests from
      --tag strings        tags of which tests will be selected from, tags are selected if the name contains this string, or if --regex is specified then matches the regex
      --test strings       test names which will be selected if the name contains this string, or if --regex is specified then matches the regex
```

### Options inherited from parent commands

```
  -b, --bootstrap string    Bootstrap URL
      --galasahome string   Path to a folder where Galasa will read and write files and configuration settings. The default is '${HOME}/.galasa'. This overrides the GALASA_HOME environment variable which may be set instead.
  -l, --log string          File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl runs](galasactl_runs.md)	 - Manage test runs in the ecosystem

