# Galasa CLI

The Galasa command line interface (Galasa CLI) is used to interact with the Galasa ecosystem or local development environment.


## Environment variables

### GALASA_BOOTSTRAP
Most commands will need a reference to the Galasa bootstrap file or url.  This can be provided with the `--bootstrap` flag or the `GALASA_BOOTSTRAP` environment variable. If both are provided, the explicit flag takes precedence.

### GALASA_HOME
`galasactl` commands assume that `${HOME}/.galasa` is a folder which is writeable, and contains property files, and storage for test results which are run using local JVMs.

This value can be overridden using the `GALASA_HOME` environment variable.

The `--galasahome` command-line flag can override the `GALASA_HOME` environment variable on a call-by-call basis.

Note: If you change this to a non-existent and/or non-initialised folder path, then 
you will have to create and re-initialise the folder using the `galasactl local init` command again. That command will respect the `GALASA_HOME` variable and will create the folder and initialise it were it not to exist.


## Syntax
The syntax is documented in generated documentation [here](docs/generated/galasactl.md)


## Setting up your environment
Run this command to set up your environment for future calls to the galasactl tooling.
```
galasactl local init
```

If you wish to build content which depends on the very latest/bleeding-edge of galasa code, 
add the `--development` flag. This will set things up to include the maven repositories 
in any new {HOME}/.m2/settings.xml file, if that doesn't already exist.

This command respects the GALASA_HOME environment variable, and will set up a number of files
there, or in {HOME}/.galasa otherwise.


## Creating an example project

`galasactl` can be used to create near-empty test projects to lay-down an initial structure 
prior to fleshing out with more tests. This can provide a boost to productivity when 
starting a new OSGi bundle containing Galasa tests.

### Examples

Getting help:-

```
galasactl --help
galasactl project --help
galasactl project create --help
```

Create a folder tree which can be built into an OSGi bundle (without an OBR):
```
galasactl project create --package dev.galasa.example.banking
```

Create a folder tree which can be built into an OSGi bundle (with an OBR):
```
galasactl project create --package dev.galasa.example.banking --obr
```

Create a folder tree which has two bundles, each aiming to test different features of an application
(while also viewing the tooling log on the console)
```
galasactl project create --package dev.galasa.example.banking --features payee,account --obr --log -
```


### Building the example project

Maven and Gradle are both build tools, which read metadata from files which guide how the code within a module should be built. Maven and Gradle use different formats for these build files.

By default, the `galasactl project create` generates a project which includes a Maven build mechanism. The `--maven` flag being present also explicitly tells the tool to generate Maven build artifacts (`pom.xml` files).

To create a project which includes a Gradle build mechanism, add the `--gradle` flag. This tells the tool to add generated artifacts which direct a Gradle build.

You can use the `--maven` and `--gradle` flags together to produce a project which contains both Maven and Gradle build infrastructure files.

Create a folder tree which can be built using either Maven or Gradle:
```
galasactl project create --package dev.galasa.example.banking --features payee,account --obr --gradle --maven
```

To build a project with Maven artifacts, use `mvn clean install`

To build a project with Gradle artifacts, use `gradle build publishToMavenLocal`

If you wish the generated code to depend upon the very latest/bleeding-edge of galasa code, then add the `--development` flag. This will add extra settings to the `settings.gradle` file of the parent and the `build.gradle` file of any of the child test projects.

## runs prepare

The purpose of `runs prepare` is to build a portfolio of tests, possibly from multiple test streams.  This portfolio can then be used in the `runs submit` command.

### Examples

Getting help:-

```
galasactl --help
galasactl runs --help
galasactl runs prepare --help
```

Finding out the version of the galasactl tool:-
```
galasactl --version
```

Will yield output such as this:
```
galasactl version 0.20.0-alpha-2305ba574524af5cd0fba59de18411582f470de5
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
          --override zos.default.lpar=MYLPAR
          --override zos.default.cluster=MYPLEXCLUSTER
```

Building a portfolio over mulitple selections and overrides:-

```
galasactl runs prepare
          --portfolio test.yaml
          --stream inttests
          --package test.package.one
          --override zos.default.lpar=MYLPAR
          --override zos.default.cluster=MYPLEXCLUSTER

galasactl runs prepare
          --portfolio test.yaml
          --append
          --stream inttests
          --package test.package.two
          --override zos.default.lpar=MYLPARB
          --override zos.default.cluster=MYPLEXCLUSTER
```

## runs submit

The purpose of `runs submit` is to submit and monitor tests in the Galasa ecosystem.  Tests can be input from a portfolio or using the same commands as the `runs prepare` command, but not both.

Note: The `--log -` directs logging information to the stderr console.
Omit this option if you do not want to see logging, or specify `--log myFileName.txt` if you wish 
to capture log information in a file.

### Examples

Getting help:-

```
galasactl --help
galasactl runs --help
galasactl runs runs --help
```

Running tests from a portfolio (and see the log on the console):-

```
galasactl runs submit --log -
          --portfolio test.yaml
          --poll 5
          --progress 1
          --throttle 5
```

Submitting tests without a portfolio (and see the log on the console) :-

```
galasactl runs submit --log -
          --class test.package.one/Test1
          --class test.package.one/Test2
```

Providing overrides for all tests during this run, note, overrides in the portfolio will take precedence over these
(and see the log on the console) :-

```
galasactl runs submit --log -
          --portfolio test.yaml
          --override zos.default.lpar=MYLPAR
          --override zos.default.cluster=MYPLEXCLUSTERA
```

## runs submit local

This command sequence causes the specified tests to be executed within the local JVM server.

Note: Runnning test running locally does not benefit from the features of running within a Galasa Ecosystem,
such as cleaning up resources when things fail and arbitrating contention for limited resources 
between competing tests. It should only be used during test development to verify that the test is 
behaving correctly.

### Example : Run a single test in the local JVM.
```
galasactl runs submit local --log -
          --obr mvn:dev.galasa.example.banking/dev.galasa.example.banking.obr/0.0.1-SNAPSHOT/obr
          --class dev.galasa.example.banking.account/dev.galasa.example.banking.account.TestAccount
```

- The --log - parameter indicates that debugging information should be sent to the console.
- The --obr indicates where the tool can find an OBR which refers to the bundle where all the tests are housed.
- The --class parameter tells the tool which test class to run. The string is in the format of `<osgi-bundle-id>/<fully-qualified-java-class>`. All the test methods within the class will be run. You can use multiple such flags to test multiple classes.
- The `JAVA_HOME` environment variable should be set to refer to the JVM to use in which the test will be launched.

- The `--throttle 1` option would mean all your tests run sequentially. A higher throttle value means that local tests run in parallel.

To configure a JVM with special options, such as `-Xms20m` and other JVM options, you can set the optional parameter `framework.jvm.local.launch.options` in your bootstrap properties to hold a space-separated list of extra options which will be used when the JVM running your test in a local JVM is launched.

### Debugging a single test which runs in the local JVM
The `galasactl runs submit local` command has an option `--debug` which causes the test to be launched in 'debug mode'.
The test will attempt to connect with a JDB java debugger based on some configuration parameters.

The 'port' used to connect the testcase to the java debugger needs to be configured.
- The default value is 2970
- The above value can be overridden by adding the optional property `galasactl.jvm.local.launch.debug.port` into the `bootstrap.properties` file.
  - For example: `galasactl.jvm.local.launch.debug.port=2971`
  - This parameter is ignored if the `--debug` argument isn't supplied to the `galasactl runs submit local` command.
- The above value can be overridden by using the optional argument `--debugPort` to the `galasactl runs submit local` command.

The port value itself must be an unsigned integer.

Your IDE would typically need to be configured to connect to the same port the testcase is using.

The 'mode' used to control the connection from the local JVM to the debugger can be `listen` or `attach`.
- `listen` configured into the galasactl configuration means that the JVM launching opens a port and pauses to listen for traffic on that port, 
waiting for a JDB debugger to connect to that port.
- `attach` configured into the galasactl configuration means that the JVM launching attempts to attach to the debug port, which has previously been opened
by the JDB debugger.

Configure the debug mode like this:
- The default value is `listen`.
- The above default value can be overridden by adding the optional property `galasactl.jvm.local.launch.debug.mode` into the `bootstrap.properties` file.
  - For example: `galasactl.jvm.local.launch.debug.mode=attach`
- The above value can be overridden by using the optional argument `--debugMode` to the `galasactl runs submit local` command.

Your IDE would typically need to be configured with the opposite type of connection mode in order to attach the JDB debugger to the running Galasa test.
For example: If your `galasactl` is configured to `listen`, then start the test first, and configure your IDE to attach to the same port.
If your `galasactl` is configured to `attach`, then start your JDB debugger first, so it is there waiting for the testcase to attach to the debug port when 
the testcase is launched.

To configure the different IDEs to connect to a local testcase in `--debug` mode, follow these instructions:
- For Microsoft vscode see [here](./docs/vscode/debug_in_vscode.md)
- For IntelliJ see [here](./docs/intellij/debug_in_intellij.md)
- For Eclipse see [here](./docs/eclipse/debug_in_eclipse.md)


## runs get
This command retrieves information about a historic run on an ecosystem.
Several formats are supported including: 'summary', 'details', 'raw' 
```
galasactl runs get --name C1234 --format details
```
For a complete list of supported formatters try running the command with a known to be bad formatter name. For example:
```
galasactl runs get --name C1234 --format badFormatterName
```
For a complete list of supported parameters see [here](./docs/generated/galasactl_runs_get.md).

## runs download

This command downloads all artifacts for a test run that are stored in an ecosystem's RAS.
The artifacts are stored in a directory within the working directory where the command is executed. The name of the created directory corresponds to the run name that was provided (e.g. `./C123/...`).

### Examples

All artifacts for a run named "C1234" can be downloaded to a directory named "C1234" in the current working directory using the following command:

```
galasactl runs download --name C1234
```

If a run directory named "C1234" already exists, artifacts stored within the directory can be overwritten with the `--force` flag as follows:

```
galasactl runs download --name C1234 --force
```

A complete list of supported parameters for the `runs download` command is available [here](./docs/generated/galasactl_runs_download.md).

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

Built artifacts include:

- galasactl-darwin-amd64 
- galasactl-darwin-arm64
- galasactl-linux-amd64 
- galasactl-linux-s390x 
- galasactl-windows-amd64.exe

Browse the following web site and download whichever built binary files you wish:

- Latest (and previous) stable releases: https://github.com/galasa-dev/cli/releases
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