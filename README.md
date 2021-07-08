# Galasa cli

The Galasa cli is used to interact with the Galasa ecosystem or local development environment.

Most commands will need a reference to the Galasa bootstrap file or url.  This can be provided with the `--bootstrap` flag or the `GALASA_BOOTSTRAP`environment variable.

## runs assemble

The purpose of `runs assemble` is to build a portfolio of tests, possibly from multiple test streams.  This portfolio can then be used in the `runs submit` command.

### Examples

Getting help:-

```
galasactl --help
galasactl runs --help
galasactl runs assemble --help
```

Selecting tests from a test steam:-

```
galasactl runs assemble
          --portfolio test.yaml
          --stream inttests
          --package test.package.one
          --package test.package.two
```

Selecting tests using regex:-

```
galasactl runs assemble
          --portfolio test.yaml
          --stream inttests
          --package '.*age.*'
          --regex
```

Selecting tests without a test stream:-

```
galasactl runs assemble
          --portfolio test.yaml
          --class test.package.one/Test1
          --class test.package.one/Test2
```

Providing test specific overrides:-

```
galasactl runs assemble
          --portfolio test.yaml
          --stream inttests
          --package test.package.one
          --override zos.default.lpar=MV2C
          --override zos.default.cluster=PLEX2
```

Building a portfolio over mulitple selections and overrides:-

```
galasactl runs assemble
          --portfolio test.yaml
          --stream inttests
          --package test.package.one
          --override zos.default.lpar=MV2C
          --override zos.default.cluster=PLEX2

galasactl runs assemble
          --portfolio test.yaml
          --append
          --stream inttests
          --package test.package.two
          --override zos.default.lpar=MV2D
          --override zos.default.cluster=PLEX2
```

## runs submit

The purpose of `runs submit` is to submit and monitor tests in the Galasa ecosystem.  Tests can be input from a portfolio or using the same commands as the `runs assemble` command, but not both.

### Examples

Getting help:-

```
galasactl --help
galasactl runs --help
galasactl runs runs --help
```

Running tests from a portfolio:-

```
galasactl runs submit
          --portfolio test.yaml
          --poll 5
          --progress 1
          --throttle 5
```

Submitting tests without a portfolio:-

```
galasactl runs submit
          --class test.package.one/Test1
          --class test.package.one/Test2
```

Providing overrides for all tests during this run, note, overrides in the portfolio will take precedence over these:-

```
galasactl runs submit
          --portfolio test.yaml
          --override zos.default.lpar=MV2C
          --override zos.default.cluster=PLEX2
```


Things left to do:-
* test report yaml
* junit reports