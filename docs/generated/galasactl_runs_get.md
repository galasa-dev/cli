## galasactl runs get

Get the details of a test runname which ran or is running.

### Synopsis

Get the details of a test runname which ran or is running, displaying the results to the caller

```
galasactl runs get [flags]
```

### Options

```
      --age string      the age of the test run(s) we want information about. Supported formats are: 'FROM' or 'FROM:TO', where FROM and TO are each ages, made up of an integer and a time-unit qualifier. Supported time-units are days('d'), weeks('w') and hours ('h'). If missing, the TO part is defaulted to '0h'. Examples: '--age 1d' , '--age 6h:1h' 
      --format string   output format for the data returned. Supported formats are: 'summary', 'details' or 'raw' (default "summary")
  -h, --help            help for get
      --name string     the name of the test run we want information about
      --result string   the result of the test run we want information about
```

### Options inherited from parent commands

```
  -b, --bootstrap string    Bootstrap URL
      --galasahome string   Path to a folder where Galasa will read and write files and configuration settings. The default is '${HOME}/.galasa'. This overrides the GALASA_HOME environment variable which may be set instead.
  -l, --log string          File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl runs](galasactl_runs.md)	 - Manage test runs in the ecosystem

