#!/usr/bin/env bash


# Where is this script executing from ?
BASEDIR=$(dirname "$0");pushd $BASEDIR 2>&1 >> /dev/null ;BASEDIR=$(pwd);popd 2>&1 >> /dev/null
# echo "Running from directory ${BASEDIR}"
export ORIGINAL_DIR=$(pwd)
# cd "${BASEDIR}"


#--------------------------------------------------------------------------
#
# Set Colors
#
#--------------------------------------------------------------------------
bold=$(tput bold)
underline=$(tput sgr 0 1)
reset=$(tput sgr0)

red=$(tput setaf 1)
green=$(tput setaf 76)
white=$(tput setaf 7)
tan=$(tput setaf 202)
blue=$(tput setaf 25)

#--------------------------------------------------------------------------
#
# Headers and Logging
#
#--------------------------------------------------------------------------
underline() { printf "${underline}${bold}%s${reset}\n" "$@"
}
h1() { printf "\n${underline}${bold}${blue}%s${reset}\n" "$@"
}
h2() { printf "\n${underline}${bold}${white}%s${reset}\n" "$@"
}
debug() { printf "${white}%s${reset}\n" "$@"
}
info() { printf "${white}➜ %s${reset}\n" "$@"
}
success() { printf "${green}✔ %s${reset}\n" "$@"
}
error() { printf "${red}✖ %s${reset}\n" "$@"
}
warn() { printf "${tan}➜ %s${reset}\n" "$@"
}
bold() { printf "${bold}%s${reset}\n" "$@"
}
note() { printf "\n${underline}${bold}${blue}Note:${reset} ${blue}%s${reset}\n" "$@"
}


#--------------------------------------------------------------------------
# 
# Main script logic
#
#--------------------------------------------------------------------------
# Create a temporary folder which is never checked in.
h2 "Making sure the tools folder is present."
mkdir -p tools
rc=$? ; if [[ "${rc}" != "0" ]]; then error "Failed to ensure the tools folder is present. rc=${rc}" ; exit 1 ; fi
success "OK"

#--------------------------------------------------------------------------
# Download the open api generator tool if we've not got it already.
export OPENAPI_GENERATOR_CLI_VERSION="6.2.0"
export OPENAPI_GENERATOR_CLI_JAR=${BASEDIR}/tools/openapi-generator-cli.jar
if [[ ! -e ${OPENAPI_GENERATOR_CLI_JAR} ]]; then
    info "The openapi generator tool is not available, so download it."
    export OPENAPI_GENERATOR_CLI_SITE="https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli"
    wget ${OPENAPI_GENERATOR_CLI_SITE}/${OPENAPI_GENERATOR_CLI_VERSION}/openapi-generator-cli-${OPENAPI_GENERATOR_CLI_VERSION}.jar \
    -O ${OPENAPI_GENERATOR_CLI_JAR}
    rc=$? ; if [[ "${rc}" != "0" ]]; then error "Failed to download the open api generator tool. rc=${rc}" ; exit 1 ; fi
    success "Downloaded OK"
fi

#--------------------------------------------------------------------------
# Invoke the generator
h2 "Generate the openapi client go code..."
./genapi.sh 2>&1 > tools/generate-log.txt
rc=$? ; if [[ "${rc}" != "0" ]]; then cat tools/generate-log.txt ; error "Failed to generate the code from the yaml file. rc=${rc}" ; exit 1 ; fi
rm -f tools/generate-log.txt
success "Code generation OK"

#--------------------------------------------------------------------------
# Invoke the generator again with different parameters
h2 "Generate the openapi client go code... part II"
./generate.sh 2>&1 > tools/generate-log.txt
rc=$? ; if [[ "${rc}" != "0" ]]; then cat tools/generate-log.txt ; error "Failed to generate II the code from the yaml file. rc=${rc}" ; exit 1 ; fi
rm -f tools/generate-log.txt
success "Code generation part II - OK"

#--------------------------------------------------------------------------
# Build the executables
h2 "Cleaning the binaries out..."
make clean
rc=$? ; if [[ "${rc}" != "0" ]]; then error "Failed to build binary executable galasactl programs. rc=${rc}" ; exit 1 ; fi
success "Binaries cleaned up - OK"

h2 "Building new binaries..."
make all
rc=$? ; if [[ "${rc}" != "0" ]]; then error "Failed to build binary executable galasactl programs. rc=${rc}" ; exit 1 ; fi
success "New binaries built - OK"

#--------------------------------------------------------------------------
h2 "Use the results.."
info "Binary executable programs are found in the 'bin' folder."
ls bin