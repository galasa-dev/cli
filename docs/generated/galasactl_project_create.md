## galasactl project create

Creates a new Galasa project

### Synopsis

Creates a new Galasa test project with optional OBR project and build process files

```
galasactl project create [flags]
```

### Options

```
      --features string   A comma-separated list of features you are testing. Defaults to "test". These must be able to form parts of a java package name. For example: "payee,account" (default "main")
      --force             Force-overwrite files which already exist.
  -h, --help              help for create
      --obr               An OSGi Object Bundle Resource (OBR) project is needed.
      --package string    Java package name for tests we create. Forms part of the project name, maven/gradle group/artifact ID, and OSGi bundle name. It may reflect the name of your organisation or company, the department, function or application under test. For example: dev.galasa.banking.example
```

### Options inherited from parent commands

```
  -l, --log string   File to which log information will be sent. Any folder referred to must exist. An existing file will be over-written. Specify "-" to log to stderr. Defaults to not logging.
```

### SEE ALSO

* [galasactl project](galasactl_project.md)	 - Manipulate local project source code

