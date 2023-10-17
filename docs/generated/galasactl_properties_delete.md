## galasactl properties delete

Delete a property in a namespace.

### Synopsis

Delete a property and its value in a namespace

```
galasactl properties delete [flags]
```

### Options

```
  -h, --help   Help command for  properties delete
```

### Options inherited from parent commands

```
  -b, --bootstrap string    Bootstrap URL. Should start with 'http://' or 'file://'. If it starts with neither, it is assumed to be a fully-qualified path. If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties
      --galasahome string   Path to a folder where Galasa will read and write files and configuration settings. The default is '${HOME}/.galasa'. This overrides the GALASA_HOME environment variable which may be set instead.
  -l, --log string          File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
  -n, --name string         Name of a property in the namespace. It has no default value.
  -s, --namespace string    Namespace. A mandatory flag that describes a container for a collection of properties. It has no default value.
```

### SEE ALSO

* [galasactl properties](galasactl_properties.md)	 - Manages properties in an ecosystem

