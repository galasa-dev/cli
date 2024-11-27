## galasactl secrets delete

Deletes a secret from the credentials store

### Synopsis

Deletes a secret from the credentials store

```
galasactl secrets delete [flags]
```

### Options

```
  -h, --help          Displays the options for the 'secrets delete' command.
      --name string   A mandatory flag that identifies the secret to be created or manipulated.
```

### Options inherited from parent commands

```
  -b, --bootstrap string    Bootstrap URL. Should start with 'http://' or 'file://'. If it starts with neither, it is assumed to be a fully-qualified path. If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties
      --galasahome string   Path to a folder where Galasa will read and write files and configuration settings. The default is '${HOME}/.galasa'. This overrides the GALASA_HOME environment variable which may be set instead.
  -l, --log string          File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl secrets](galasactl_secrets.md)	 - Manage secrets stored in the Galasa service's credentials store
