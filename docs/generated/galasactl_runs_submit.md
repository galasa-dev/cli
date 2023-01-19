## galasactl runs submit

submit a list of tests to the ecosystem

### Synopsis

Submit a list of tests to the ecosystem, monitor them and wait for them to complete

```
galasactl runs submit [flags]
```

### Options

```
      --bundle strings             bundles of which tests will be selected from, bundles are selected if the name contains this string, or if --regex is specified then matches the regex
      --class strings              test class names, for building a portfolio when a stream/test catalog is not available
  -g, --group string               the group name to assign the test runs to, if not provided, a psuedo unique id will be generated
  -h, --help                       help for submit
      --noexitcodeontestfailures   set to true if you don't want an exit code to be returned from galasactl if a test fails
      --override strings           overrides to be sent with the tests (overrides in the portfolio will take precedence)
      --package strings            packages of which tests will be selected from, packages are selected if the name contains this string, or if --regex is specified then matches the regex
      --poll int                   Optional. The interval time in seconds between successive polls of the ecosystem for the status of the test runs. Defaults to 30 seconds. (default 30)
  -p, --portfolio string           portfolio containing the tests to run
      --progress int               in minutes, how often the cli will report the overall progress of the test runs, -1 or less will disable progress reports (default 5)
      --regex                      Test selection is performed by using regex
      --reportjson string          json file to record the final results in
      --reportjunit string         junit xml file to record the final results in
      --reportyaml string          yaml file to record the final results in
      --requestor string           the requestor id to be associated with the test runs (default "mcobbett")
      --requesttype string         the type of request, used to allocate a run name (default "CLI")
  -s, --stream string              test stream to extract the tests from
      --tag strings                tags of which tests will be selected from, tags are selected if the name contains this string, or if --regex is specified then matches the regex
      --test strings               test names which will be selected if the name contains this string, or if --regex is specified then matches the regex
      --throttle int               how many test runs can be submitted in parallel, 0 or less will disable throttling (default 3)
      --throttlefile string        a file where the current throttle is stored and monitored, used to dynamically change the throttle
      --trace                      Trace to be enabled on the test runs
```

### Options inherited from parent commands

```
  -b, --bootstrap string   Bootstrap URL
  -l, --log string         File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl runs](galasactl_runs.md)	 - Manage test runs in the ecosystem

