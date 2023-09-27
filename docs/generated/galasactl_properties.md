## galasactl properties

Manages properties in an ecosystem

### Synopsis

Allows interaction with the CPS to Initiate, query and maintain properties in Galasa Ecosystem

### Options

```
  -b, --bootstrap string   Bootstrap URL. Should start with 'http://' or 'file://'. If it starts with neither, it is assumed to be a fully-qualified path. If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. Examples: http://galasa-cicsk8s.hursley.ibm.com/bootstrap , file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties
  -h, --help               help for properties
  -n, --name string        Name of a property in the namespace.It has no default value.
  -s, --namespace string   Namespace. A container for a collection of properties. It has no default value.
```

### Options inherited from parent commands

```
      --galasahome string   Path to a folder where Galasa will read and write files and configuration settings. The default is '${HOME}/.galasa'. This overrides the GALASA_HOME environment variable which may be set instead.
  -l, --log string          File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl](galasactl.md)	 - CLI for Galasa
* [galasactl properties get](galasactl_properties_get.md)	 - Get the details of properties in a namespace.
