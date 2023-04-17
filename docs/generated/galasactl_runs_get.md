## galasactl runs get

Get the details of a test runname which ran or is running.

### Synopsis

Get the details of a test runname which ran or is running, displaying the results to the caller

```
galasactl runs get [flags]
```

### Options

```
  -h, --help             help for get
      --output string    output format for the data returned. Supported formats are: summary (default "summary")
      --runname string   the name of the test run we want information about
```

### Options inherited from parent commands

```
  -b, --bootstrap string    Bootstrap URL
      --galasahome string   Path to a folder where Galasa will read and write files and configuration settings. The default is '${HOME}/.galasa'. This overrides the GALASA_HOME environment variable which may be set instead.
  -l, --log string          File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl runs](galasactl_runs.md)	 - Manage test runs in the ecosystem

