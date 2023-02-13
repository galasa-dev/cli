## galasactl runs submit local

submit a list of tests to be run on a local java virtual machine (JVM)

### Synopsis

Submit a list of tests to a local JVM, monitor them and wait for them to complete

```
galasactl runs submit local [flags]
```

### Options

```
      --bundle strings         bundles of which tests will be selected from, bundles are selected if the name contains this string, or if --regex is specified then matches the regex
      --class strings          test class names, for building a portfolio when a stream/test catalog is not available. The format of each entry is osgi-bundle-name/java-class-name . Java class names are fully qualified.
      --galasaVersion string   the version of galasa you want to use to run your tests. Defaults to 0.26.0 (default "0.26.0")
  -h, --help                   help for local
      --obr strings            The maven coordinates of the obr bundle(s) which refer to your test bundles. The format of this parameter is 'mvn:${TEST_OBR_GROUP_ID}/${TEST_OBR_ARTIFACT_ID}/${TEST_OBR_VERSION}/obr' Multiple instances of this flag can be used to describe multiple obr bundles.
      --package strings        packages of which tests will be selected from, packages are selected if the name contains this string, or if --regex is specified then matches the regex
      --regex                  Test selection is performed by using regex
      --remoteMaven string     the url of the remote maven where galasa bundles can be loaded from. Defaults to maven central https://repo.maven.apache.org/maven2 (default "https://repo.maven.apache.org/maven2")
  -s, --stream string          test stream to extract the tests from
      --tag strings            tags of which tests will be selected from, tags are selected if the name contains this string, or if --regex is specified then matches the regex
      --test strings           test names which will be selected if the name contains this string, or if --regex is specified then matches the regex
```

### Options inherited from parent commands

```
  -b, --bootstrap string   Bootstrap URL
  -l, --log string         File to which log information will be sent. Any folder referred to must exist. An existing file will be overwritten. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl runs submit](galasactl_runs_submit.md)	 - submit a list of tests to the ecosystem

