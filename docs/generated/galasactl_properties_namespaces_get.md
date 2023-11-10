## galasactl properties namespaces get

Get a list of namespaces.

### Synopsis

Get a list of namespaces within the CPS

```
galasactl properties namespaces get [flags]
```

### Options

```
      --format string   output format for the data returned. Supported formats are: 'raw', 'summary', 'yaml'. (default "summary")
  -h, --help            Displays the options for the namespaces get command.
```

### Options inherited from parent commands

```
  -b, --bootstrap string    Bootstrap URL. Should start with 'http://' or 'file://'. If it starts with neither, it is assumed to be a fully-qualified path. If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties
      --galasahome string   Path to a folder where Galasa will read and write files and configuration settings. The default is '${HOME}/.galasa'. This overrides the GALASA_HOME environment variable which may be set instead.
  -l, --log string          File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl properties namespaces](galasactl_properties_namespaces.md)	 - Queries namespaces in an ecosystem

