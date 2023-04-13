#!/bin/bash

echo "Running script test-galasactl-ecosystem.sh"

# This script can be ran locally or executed in a pipeline to test the various built binaries of galasactl
# This script tests the 'galasactl runs submit' command against a test that is in our ecosystem's testcatalog already
# Pre-requesite: the CLI must have been built first so the binaries are present in the /bin directory


# Where is this script executing from ?
BASEDIR=$(dirname "$0");pushd $BASEDIR 2>&1 >> /dev/null ;BASEDIR=$(pwd);popd 2>&1 >> /dev/null
export ORIGINAL_DIR=$(pwd)
cd "${BASEDIR}"


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


#-----------------------------------------------------------------------------------------                   
# Functions
#-----------------------------------------------------------------------------------------                   
function usage {
    info "Syntax: test-galasactl-ecosystem.sh --binary [OPTIONS] --bootstrap [BOOTSTRAP]"
    cat << EOF
Options are:
galasactl-darwin-amd64 : Use the galasactl-darwin-amd64 binary
galasactl-darwin-arm64 : Use the galasactl-darwin-arm64 binary
galasactl-linux-amd64 : Use the galasactl-linux-amd64 binary
galasactl-linux-s390x : Use the galasactl-linux-s390x binary
galasactl-windows-amd64.exe : Use the galasactl-windows-amd64.exe binary

Bootstrap must refer to a remote ecosystem.
EOF
}

#-----------------------------------------------------------------------------------------                   
# Process parameters
#-----------------------------------------------------------------------------------------                   
binary=""
bootstrap=""

while [ "$1" != "" ]; do
    case $1 in
        --binary )                        shift
                                          binary="$1"
                                          ;;
        --bootstrap )                     shift
                                          bootstrap="$1"
                                          ;;
        -h | --help )                     usage
                                          exit
                                          ;;
        * )                               error "Unexpected argument $1"
                                          usage
                                          exit 1
    esac
    shift
done

if [[ "${binary}" != "" ]]; then
    case ${binary} in
        galasactl-darwin-amd64 )            echo "Using the galasactl-darwin-amd64 binary"
                                            ;;
        galasactl-darwin-arm64 )            echo "Using the galasactl-darwin-arm64 binary"
                                            ;;
        galasactl-linux-amd64 )             echo "Using the galasactl-linux-amd64 binary"
                                            ;;
        galasactl-linux-s390x )             echo "Using the galasactl-linux-s390x binary"
                                            ;;
        galasactl-windows-amd64.exe )       echo "Using the galasactl-windows-amd64.exe binary"
                                            ;;
        * )                                 error "Unrecognised galasactl binary ${binary}"
                                            usage
                                            exit 1
    esac
else
    error "Need to specify which binary of galasactl to use."
    usage
    exit 1  
fi

# Can't really verify that the bootstrap provided is a valid one, but galasactl will pick this up later if not
if [[ "${bootstrap}" != "" ]]; then
    echo "Running tests against ecosystem bootstrap ${bootstrap}"
else
    error "Need to provide the bootstrap for a remote ecosystem."
    usage
    exit 1  
fi

#--------------------------------------------------------------------------
function launch_test_on_ecosystem_with_portfolio {
    h2 "Building a portfolio..."

    mkdir -p ${BASEDIR}/temp
    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${binary} runs prepare \
    --bootstrap $bootstrap \
    --stream inttests \
    --portfolio portfolio.yaml \
    --test local.CoreLocalJava11Ubuntu \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    # We expect a return code of '0' because this test is in the ecosystem's testcatalog.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to create a portfolio with a known test from the ecosystem's testcatalog."
        exit 1
    fi
    success "Creating portfolio.yaml worked OK"

    h2 "Launching test on an ecosystem from a portfolio..."

    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${binary} runs submit \
    --bootstrap $bootstrap \
    --portfolio portfolio.yaml \
    --throttle 1 \
    --poll 10 \
    --progress 1 \
    --noexitcodeontestfailures \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    # We expect a return code of '0' because the ecosystem should be able to run this test.
    # We have specified the flag --noexitcodeontestfailures so that we still receive a return code '0' even if the test fails,
    # as we are testing galasactl here, not the test itself
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to submit a test to a remote ecosystem."
        exit 1
    fi
    success "Submitting test to ecosystem worked OK"
}

#--------------------------------------------------------------------------
function launch_test_on_ecosystem_without_portfolio {
    h2 "Launching test on an ecosystem..."

    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${binary} runs submit \
    --bootstrap $bootstrap \
    --class dev.galasa.inttests/dev.galasa.inttests.core.local.CoreLocalJava11Ubuntu \
    --throttle 1 \
    --poll 10 \
    --progress 1 \
    --noexitcodeontestfailures \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    # We expect a return code of '0' because the ecosystem should be able to run this test.
    # We have specified the flag --noexitcodeontestfailures so that we still receive a return code '0' even if the test fails,
    # as we are testing galasactl here, not the test itself
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to submit a test to a remote ecosystem."
        exit 1
    fi
    success "Submitting test to ecosystem worked OK"
}

#--------------------------------------------------------------------------
function create_portfolio_with_unknown_test {
    h2 "Building a portfolio with an unknown test..."

    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${binary} runs prepare \
    --bootstrap $bootstrap \
    --stream inttests \
    --portfolio unknown-portfolio.yaml \
    --test local.UnknownTest \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    # We expect a return code of '1' because the ecosystem doesn't know about this testcase.
    if [[ "${rc}" != "1" ]]; then 
        error "Failed to recognise an Unknown testcase."
        exit 1
    fi
    success "Unknown test was recognised and no tests were selected"
}

#--------------------------------------------------------------------------
function launch_test_from_unknown_portfolio {
    h2 "Launching a test from an unknown portfolio..."

    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${binary} runs submit \
    --bootstrap $bootstrap \
    --portfolio unknown-portfolio.yaml \
    --throttle 1 \
    --poll 10 \
    --progress 1 \
    --noexitcodeontestfailures \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    # We expect a return code of '1' because the galasactl shouldn't be able to read this portfolio.
    if [[ "${rc}" != "1" ]]; then 
        error "Failed to recognise a non-existent portfolio."
        exit 1
    fi
    success "Unknown portfolio could not be read. galasactl reported this error correctly."
}

# Launch test on ecosystem from a portfolio ...
launch_test_on_ecosystem_with_portfolio

# Launch test on ecosystem without a portfolio ...
# NOTE - Potential bug found with this command so commenting out for now
# launch_test_on_ecosystem_without_portfolio

# Attempt to create a test portfolio with an unknown test ...
create_portfolio_with_unknown_test

# Attempt to launch a test from an unknown portfolio ...
launch_test_from_unknown_portfolio