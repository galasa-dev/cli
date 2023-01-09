## galasactl project create

Creates a new Galasa project

### Synopsis

Creates a new Galasa test project with optional OBR project and build process files

```
galasactl project create [flags]
```

### Options

```
      --force            Force-overwrite files which already exist.
  -h, --help             help for create
      --obr              An OSGi Object Bundle Resource (OBR) project is needed.
      --package string   Java package name for tests we create. Forms part of the project name, maven/gradle group/artifact ID, and OSGi bundle name. For example: com.myco.myproduct.myapp
```

### Options inherited from parent commands

```
  -l, --log string   File to which log information will be sent. Any folder referred to must exist. An existing file will be over-written.
```

### SEE ALSO

* [galasactl project](galasactl_project.md)	 - Manipulate local project source code

