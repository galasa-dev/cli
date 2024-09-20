## galasactl auth tokens get

Get a list of authentication tokens

### Synopsis

Get a list of tokens used for authentication with the Galasa API server

```
galasactl auth tokens get [flags]
```

### Options

```
  -h, --help          Displays the options for the 'auth tokens get' command.
      --user string   Optional. Retrieves a list of access tokens for the user with the given username.
```

### Options inherited from parent commands

```
  -b, --bootstrap string    Bootstrap URL. Should start with 'http://' or 'file://'. If it starts with neither, it is assumed to be a fully-qualified path. If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties
      --galasahome string   Path to a folder where Galasa will read and write files and configuration settings. The default is '${HOME}/.galasa'. This overrides the GALASA_HOME environment variable which may be set instead.
  -l, --log string          File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl auth tokens](galasactl_auth_tokens.md)	 - Queries tokens in an ecosystem

