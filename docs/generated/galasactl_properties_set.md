## galasactl properties set

Set the details of properties in a namespace.

### Synopsis

Set the details of a property in a namespace. If the property does not exist, a new property is created, otherwise the value for that property will be updated.

```
galasactl properties set [flags]
```

### Options

```
  -h, --help           help for set
      --value string   the value of the property you want to create
```

### Options inherited from parent commands

```
  -b, --bootstrap string    Bootstrap URL. Should start with 'http://' or 'file://'. If it starts with neither, it is assumed to be a fully-qualified path. If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties
      --galasahome string   Path to a folder where Galasa will read and write files and configuration settings. The default is '${HOME}/.galasa'. This overrides the GALASA_HOME environment variable which may be set instead.
  -l, --log string          File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
  -n, --name string         Name of a property in the namespace. It has no default value.
  -s, --namespace string    Namespace. A container for a collection of properties. It has no default value.
```

### SEE ALSO

* [galasactl properties](galasactl_properties.md)	 - Manages properties in an ecosystem

