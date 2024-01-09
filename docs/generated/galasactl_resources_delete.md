## galasactl resources delete

Delete Galasa Ecosystem resources.

### Synopsis

Delete Galasa Ecosystem resources in a file.

```
galasactl resources delete [flags]
```

### Options

```
  -h, --help   Displays the options for the 'resources delete' command.
```

### Options inherited from parent commands

```
  -b, --bootstrap string    Bootstrap URL. Should start with 'http://' or 'file://'. If it starts with neither, it is assumed to be a fully-qualified path. If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties
  -f, --file string         The file containing yaml definitions of resources to be applied manipulated by this command. This can be a fully-qualified path or path relative to the current directory.Example: my_resources.yaml
      --galasahome string   Path to a folder where Galasa will read and write files and configuration settings. The default is '${HOME}/.galasa'. This overrides the GALASA_HOME environment variable which may be set instead.
  -l, --log string          File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl resources](galasactl_resources.md)	 - Manages resources in an ecosystem

