## galasactl runs download

Download the artifacts of a test run which ran.

### Synopsis

Download the artifacts of a test run which ran and store them in a directory within the current working directory

```
galasactl runs download [flags]
```

### Options

```
      --destination string   The folder we want to download test run artifacts into. Sub-folders will be created within this location (default ".")
      --force                force artifacts to be overwritten if they already exist
  -h, --help                 Displays the options for the 'runs download' command.
      --name string          the name of the test run we want information about
```

### Options inherited from parent commands

```
  -b, --bootstrap string    Bootstrap URL. Should start with 'http://' or 'file://'. If it starts with neither, it is assumed to be a fully-qualified path. If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties
      --galasahome string   Path to a folder where Galasa will read and write files and configuration settings. The default is '${HOME}/.galasa'. This overrides the GALASA_HOME environment variable which may be set instead.
  -l, --log string          File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl runs](galasactl_runs.md)	 - Manage test runs in the ecosystem

