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
      --class strings      test class names, for building a portfolio when a stream/test catalog is not available
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
  -b, --bootstrap string   Bootstrap URL
  -l, --log string         File to which log information will be sent
```

### SEE ALSO

* [galasactl runs](galasactl_runs.md)	 - Manage test runs in the ecosystem

