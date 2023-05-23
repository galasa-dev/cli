## galasactl runs download

Download the artifacts of a test run which ran.

### Synopsis

Download the artifacts of a test run which ran and store them in a directory within the current working directory

```
galasactl runs download [flags]
```

### Options

```
      --force         force artifacts to be overwritten if they already exist
  -h, --help          help for download
      --name string   the name of the test run we want information about
```

### Options inherited from parent commands

```
  -b, --bootstrap string    Bootstrap URL
      --galasahome string   Path to a folder where Galasa will read and write files and configuration settings. The default is '${HOME}/.galasa'. This overrides the GALASA_HOME environment variable which may be set instead.
  -l, --log string          File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl runs](galasactl_runs.md)	 - Manage test runs in the ecosystem

