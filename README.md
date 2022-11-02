# Galasa cli

The Galasa cli is used to interact with the Galasa ecosystem or local development environment.

Most commands will need a reference to the Galasa bootstrap file or url.  This can be provided with the `--bootstrap` flag or the `GALASA_BOOTSTRAP`environment variable.

## runs prepare

The purpose of `runs prepare` is to build a portfolio of tests, possibly from multiple test streams.  This portfolio can then be used in the `runs submit` command.

### Examples

Getting help:-

```
galasactl --help
galasactl runs --help
galasactl runs prepare --help
```

Selecting tests from a test steam:-

```
galasactl runs prepare
          --portfolio test.yaml
          --stream inttests
          --package test.package.one
          --package test.package.two
```

Selecting tests using regex:-

```
galasactl runs prepare
          --portfolio test.yaml
          --stream inttests
          --package '.*age.*'
          --regex
```

Selecting tests without a test stream:-

```
galasactl runs prepare
          --portfolio test.yaml
          --class test.package.one/Test1
          --class test.package.one/Test2
```

Providing test specific overrides:-

```
galasactl runs prepare
          --portfolio test.yaml
          --stream inttests
          --package test.package.one
          --override zos.default.lpar=MV2C
          --override zos.default.cluster=PLEX2
```

Building a portfolio over mulitple selections and overrides:-

```
galasactl runs prepare
          --portfolio test.yaml
          --stream inttests
          --package test.package.one
          --override zos.default.lpar=MV2C
          --override zos.default.cluster=PLEX2

galasactl runs prepare
          --portfolio test.yaml
          --append
          --stream inttests
          --package test.package.two
          --override zos.default.lpar=MV2D
          --override zos.default.cluster=PLEX2
```

## runs submit

The purpose of `runs submit` is to submit and monitor tests in the Galasa ecosystem.  Tests can be input from a portfolio or using the same commands as the `runs prepare` command, but not both.

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

## How to build locally
- Clone the `framework` repository so it is a sibling project of this repository
  - This repo genrates a golang client for the openapi REST interface offered by the framework.
  - The openapi.yaml file is kept in the `framework` repository.
- Run the `build-locally.sh` script.
- Binary executable programs appear in the `bin` folder


