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


## Reference Material

### Syntax
Full syntax, with descriptions of every parameter is available [here](./docs/generated/galasactl.md)

### Errors
See [here](./docs/generated/errors-list.md) for a list of error messages the `galasactl` tool can output.

### Known limitations
- Go programs can sometimes struggle to resolve DNS names, especially when a working over a virtual private network (VPN). In such situations, you may notice that a bootstrap file cannot be found with `galasactl` but can be found by a desktop browser, or curl command.
In such situations we recommend that you manually add the host detail to the `/etc/hosts` file, 
to avoid DNS being involved in the resolution mechanism.

## Building locally
To build the cli tools locally, use the `./build-locally.sh --help` script for instructions.

## Built artifacts
Download built artifacts:

Browse the following web site and download whichever built binary files you wish:

- Bleeding edge/Unstable : https://development.galasa.dev/main/binary/cli/

## Docker images containing the command-line tools
The build process builds some docker images with the command-line tools installed.
This could be useful when wishing to embed a usage of the command-line within a build process which can use a docker image.

- Bleeding edge/Unstable : `docker pull harbor.galasa.dev/galasadev/galasa-cli-amd64:main`

### How to use the docker image
The docker image has the `galasactl` tool on the path of the docker image when it starts up.
So, invoke the `galasactl` without installing on your local machine, using the docker image like this:
```
docker run harbor.galasa.dev/galasadev/galasa-cli-amd64:main galasactl --version
```



