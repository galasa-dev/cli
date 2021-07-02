java -jar ~/openapi/openapi-generator-cli.jar generate -i ../framework/openapi.yaml -g go -o pkg/galasaapi --additional-properties=packageName=galasaapi --additional-properties=isGoSubmodule=true
