# Debug a local Galasa test within Microsoft vscode


## Prepare to launch the java debugger
Make sure you have the ["Debugger for Java" plugin](https://github.com/microsoft/vscode-java-debug) installed, from Microsoft.

Add the following text to your ${workspace}/.vscode/
```
{
    "version": "0.2.0",
    "configurations": [
    {
        "type": "java",
        "name": "Debug (Attach to Galasa test on 2970)",
        "projectName": "banking",
        "request": "attach",
        "hostName": "localhost",
        "port": 2970
    }
    ]
}
```

This can be done manually, or you can use the "create a launch.json" file which is available when you open the "Run and Debug" side-bar view, as in this diagram:
[view of Run and Debug side-panel](create-launch-json.png)
