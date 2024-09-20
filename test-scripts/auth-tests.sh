#!/bin/bash

#
# Copyright contributors to the Galasa project
#
# SPDX-License-Identifier: EPL-2.0
#
echo "Running script runs-tests.sh"
# This script can be ran locally or executed in a pipeline to test the various built binaries of galasactl
# This script tests the 'galasactl runs submit' command against a test that is in our ecosystem's testcatalog already
# Pre-requesite: the CLI must have been built first so the binaries are present in the /bin directory

if [[ "$CALLED_BY_MAIN" == "" ]]; then
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
    underline() { printf "${underline}${bold}%s${reset}\n" "$@" ;}
    h1() { printf "\n${underline}${bold}${blue}%s${reset}\n" "$@" ;}
    h2() { printf "\n${underline}${bold}${white}%s${reset}\n" "$@" ;}
    debug() { printf "${white}%s${reset}\n" "$@" ;}
    info() { printf "${white}➜ %s${reset}\n" "$@" ;}
    success() { printf "${green}✔ %s${reset}\n" "$@" ;}
    error() { printf "${red}✖ %s${reset}\n" "$@" ;}
    warn() { printf "${tan}➜ %s${reset}\n" "$@" ;}
    bold() { printf "${bold}%s${reset}\n" "$@" ;}
    note() { printf "\n${underline}${bold}${blue}Note:${reset} ${blue}%s${reset}\n" "$@" ;}

    #-----------------------------------------------------------------------------------------
    # Process parameters
    #-----------------------------------------------------------------------------------------
    bootstrap=""

    while [ "$1" != "" ]; do
        case $1 in
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

    # Can't really verify that the bootstrap provided is a valid one, but galasactl will pick this up later if not
    if [[ "${bootstrap}" == "" ]]; then
        export bootstrap="https://prod1-galasa-dev.cicsk8s.hursley.ibm.com/api/bootstrap"
        info "No bootstrap supplied. Defaulting the --bootstrap to be ${bootstrap}"
    fi

    info "Running tests against ecosystem bootstrap ${bootstrap}"

    #-----------------------------------------------------------------------------------------
    # Constants
    #-----------------------------------------------------------------------------------------
    export GALASA_TEST_NAME_SHORT="local.CoreLocalJava11Ubuntu"
    export GALASA_TEST_NAME_LONG="dev.galasa.inttests.core.${GALASA_TEST_NAME_SHORT}"
    export GALASA_TEST_RUN_GET_EXPECTED_SUMMARY_LINE_COUNT="4"
    export GALASA_TEST_RUN_GET_EXPECTED_DETAILS_LINE_COUNT="13"
    export GALASA_TEST_RUN_GET_EXPECTED_RAW_PIPE_COUNT="10"
    export GALASA_TEST_RUN_GET_EXPECTED_NUMBER_ARTIFACT_RUNNING_COUNT="10"

fi

# generate a random number to append to test names to avoid multiple running at once overriding each other
function get_random_property_name_number {
    minimum=100
    maximum=999
    PROP_NUM=$(($minimum + $RANDOM % $maximum))
    echo $PROP_NUM
}

#-----------------------------------------------------------------------------------------
# Tests
#-----------------------------------------------------------------------------------------

function auth_tokens_get {

    h2 "Performing auth tokens get without loginId: get..."

    set -o pipefail # Fail everything if anything in the pipeline fails. Else we are just checking the 'tee' return code.

    cmd="${BINARY_LOCATION} auth tokens get \
    --bootstrap $bootstrap \
    --log -
    "

    info "Command is: $cmd"

    $cmd
    rc=$?

    # We expect a return code of 0 because this is a properly formed auth tokens get command.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get access tokens."
        exit 1
    fi

    output_file="$ORIGINAL_DIR/temp/auth-get-output.txt"
    $cmd | tee $output_file

    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get access tokens."
        exit 1
    fi

    # Check that the previous properties set created a property
    cat $output_file | grep "Total:" -q

    success "All access tokens fetched from database successfully."

}

function auth_tokens_get_with_missing_loginId_throws_error {

    h2 "Performing auth tokens get with loginId: get..."
    loginId=""

    cmd="${BINARY_LOCATION} auth tokens get \
    --user $loginId \
    --bootstrap $bootstrap \
    --log -
    "

    info "Command is: $cmd"

    output_file="$ORIGINAL_DIR/temp/auth-get-output.txt"
    set -o pipefail # Fail everything if anything in the pipeline fails. Else we are just checking the 'tee' return code.
    $cmd | tee $output_file

    rc=$?
    if [[ "${rc}" != "1" ]]; then 
        error "Failed to get access tokens."
        exit 1
    fi

    success "galasactl auth tokens get command correctly threw an error due to missing loginId"

}

function auth_tokens_get_by_loginId {

    h2 "Performing auth tokens get with loginId: get..."
    loginId="Aashir.Siddiqui@ibm.com"

    set -o pipefail # Fail everything if anything in the pipeline fails. Else we are just checking the 'tee' return code.

    cmd="${BINARY_LOCATION} auth tokens get \
    --user $loginId \
    --bootstrap $bootstrap \
    --log -
    "

    info "Command is: $cmd"

    $cmd
    rc=$?

    # We expect a return code of 0 because this is a properly formed auth tokens get command.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to create property with name and value used."
        exit 1
    fi

    output_file="$ORIGINAL_DIR/temp/auth-get-output.txt"
    $cmd | tee $output_file
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get property with name used: command failed."
        exit 1
    fi

    # Check that the previous properties set created a property
    cat $output_file | grep "Total: 1" -q

    success "All access tokens by loginId fetched from database successfully."

}

#--------------------------------------------------------------------------

function auth_tests {
    auth_tokens_get
    auth_tokens_get_by_loginId
    auth_tokens_get_with_missing_loginId_throws_error
}

# checks if it's been called by main, set this variable if it is
if [[ "$CALLED_BY_MAIN" == "" ]]; then
    source $BASEDIR/calculate-galasactl-executables.sh
    calculate_galasactl_executable
    auth_tests
fi