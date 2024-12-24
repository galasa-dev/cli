## galasactl secrets

Manage secrets stored in the Galasa service's credentials store

### Synopsis

The parent command for operations to manipulate secrets in the Galasa service's credentials store

### Options

```
  -b, --bootstrap string                 Bootstrap URL. Should start with 'http://' or 'file://'. If it starts with neither, it is assumed to be a fully-qualified path. If missing, it defaults to use the 'bootstrap.properties' file in your GALASA_HOME. Example: http://example.com/bootstrap, file:///user/myuserid/.galasa/bootstrap.properties , file://C:/Users/myuserid/.galasa/bootstrap.properties
  -h, --help                             Displays the options for the 'secrets' command.
      --rate-limit-retries int           The maximum number of retries that should be made when requests to the Galasa Service fail due to rate limits being exceeded. Must be a whole number. Defaults to 3 retries (default 3)
      --rate-limit-retry-backoff float   The amount of time in seconds to wait before retrying a command if it failed due to rate limits being exceeded. Defaults to 1 second. (default 1)
```

### Options inherited from parent commands

```
      --galasahome string   Path to a folder where Galasa will read and write files and configuration settings. The default is '${HOME}/.galasa'. This overrides the GALASA_HOME environment variable which may be set instead.
  -l, --log string          File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl](galasactl.md)	 - CLI for Galasa
* [galasactl secrets delete](galasactl_secrets_delete.md)	 - Deletes a secret from the credentials store
* [galasactl secrets get](galasactl_secrets_get.md)	 - Get secrets from the credentials store
* [galasactl secrets set](galasactl_secrets_set.md)	 - Creates or updates a secret in the credentials store

