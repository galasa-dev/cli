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
      --local                      set to true if you don't want an exit code to be returned from galasactl if a test fails
      --noexitcodeontestfailures   set to true if you don't want an exit code to be returned from galasactl if a test fails
      --override strings           overrides to be sent with the tests (overrides in the portfolio will take precedence). Each override is of the form 'name=value'. Multiple instances of this flag can be used. For example --override=prop1=val1 --override=prop2=val2
      --package strings            packages of which tests will be selected from, packages are selected if the name contains this string, or if --regex is specified then matches the regex
      --poll int                   Optional. The interval time in seconds between successive polls of the ecosystem for the status of the test runs. Defaults to 30 seconds. If less than 1, then default value is used. (default 30)
  -p, --portfolio string           portfolio containing the tests to run
      --progress int               in minutes, how often the cli will report the overall progress of the test runs, -1 or less will disable progress reports. Defaults to 5 minutes. If less than 1, then default value is used. (default 5)
      --regex                      Test selection is performed by using regex
      --reportjson string          json file to record the final results in
      --reportjunit string         junit xml file to record the final results in
      --reportyaml string          yaml file to record the final results in
      --requestor string           the requestor id to be associated with the test runs. Defaults to the current user id (default "mcobbett")
      --requesttype string         the type of request, used to allocate a run name. Defaults to CLI. (default "CLI")
  -s, --stream string              test stream to extract the tests from
      --tag strings                tags of which tests will be selected from, tags are selected if the name contains this string, or if --regex is specified then matches the regex
      --test strings               test names which will be selected if the name contains this string, or if --regex is specified then matches the regex
      --throttle int               how many test runs can be submitted in parallel, 0 or less will disable throttling. Default is 3 (default 3)
      --throttlefile string        a file where the current throttle is stored. Periodically the throttle value is read from the file used. Someone with edit access to the file can change it which dynamically takes effect. Long-running large portfolios can be throttled back to nothing (paused) using this mechanism (if throttle is set to 0). And they can be resumed (un-paused) if the value is set back. This facility can allow the tests to not show a failure when the system under test is taken out of service for maintainence.
```

### Options inherited from parent commands

```
  -b, --bootstrap string   Bootstrap URL
  -l, --log string         File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl runs](galasactl_runs.md)	 - Manage test runs in the ecosystem

