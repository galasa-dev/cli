## galasactl auth login

Log in to a Galasa ecosystem using an existing access token

### Synopsis

Log in to a Galasa ecosystem using an existing access token stored in the 'galasactl.properties' file in your GALASA_HOME directory. If you do not have an access token, request one through your ecosystem's web user interface and follow the instructions on the web user interface to populate the 'galasactl.properties' file.

```
galasactl auth login [flags]
```

### Options

```
  -b, --bootstrap string   Bootstrap URL. Should start with 'http://' or 'file://'. If it starts with neither, it is assumed to be a fully-qualified path. If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties
  -h, --help               Displays the options for the 'auth login' command.
```

### Options inherited from parent commands

```
      --galasahome string   Path to a folder where Galasa will read and write files and configuration settings. The default is '${HOME}/.galasa'. This overrides the GALASA_HOME environment variable which may be set instead.
  -l, --log string          File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl auth](galasactl_auth.md)	 - Manages authentication with a Galasa ecosystem

