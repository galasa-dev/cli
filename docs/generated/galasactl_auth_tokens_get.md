## galasactl auth tokens get

Get a list of authentication tokens

### Synopsis

Get a list of tokens used for authenticating with the Galasa API server

```
galasactl auth tokens get [flags]
```

### Options

```
      --format string   output format for the data returned. Supported formats are: 'summary'. (default "summary")
  -h, --help            Displays the options for the 'auth tokens get' command.
```

### Options inherited from parent commands

```
      --galasahome string   Path to a folder where Galasa will read and write files and configuration settings. The default is '${HOME}/.galasa'. This overrides the GALASA_HOME environment variable which may be set instead.
  -l, --log string          File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl auth tokens](galasactl_auth_tokens.md)	 - Queries tokens in an ecosystem

