#!/bin/bash

#
# Copyright contributors to the Galasa project
#
# SPDX-License-Identifier: EPL-2.0
#
echo "Running script gherkin-runs-tests.sh"
# This script can be ran locally or executed in a pipeline to test the various built binaries of galasactl
# This script tests the 'galasactl runs submit local --gherkin' command against a test that is generated
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
        export bootstrap="https://galasa-galasa-prod.cicsk8s.hursley.ibm.com/bootstrap"
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

function generateGherkinFile {
    h1 "Creating sample gherkin feature files"

    input_file="$ORIGINAL_DIR/temp/GherkinSubmitTest.feature"

    cat << EOF > $input_file 
Feature: GherkinSubmitTest
  Scenario: Log Example Statement
    
    THEN Write to log "This is a log statement"
    
  Scenario: Log Statement Test
    
    THEN Write to log "This is a second scenario"
EOF
    success "OK"
}

function SubmittingLocalGherkinTest {
    h1 "Submitting the sample gherkin feature"

    cmd="$ORIGINAL_DIR/bin/${binary} runs submit local \
    --gherkin file://$input_file \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?

    if [[ "${rc}" != "0" ]]; then 
        error "Failed to run gherkin test"
        exit 1
    fi
}

function SubmittingBadPrefixLocalGherkinTest {
    h1 "Submitting the bad gherkin feature"

    cmd="$ORIGINAL_DIR/bin/${binary} runs submit local \
    --gherkin $input_file \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?

    if [[ "${rc}" != "1" ]]; then 
        error "Gherkin test did not fail as expected"
        exit 1
    fi
}

function SubmittingBadSuffixLocalGherkinTest {
    h1 "Submitting the bad gherkin feature"

    cmd="$ORIGINAL_DIR/bin/${binary} runs submit local \
    --gherkin file:///gherkin \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?

    if [[ "${rc}" != "1" ]]; then 
        error "Gherkin test did not fail as expected"
        exit 1
    fi
}

function test_gherkin_commands {
    generateGherkinFile
    SubmittingLocalGherkinTest
    SubmittingBadPrefixLocalGherkinTest
    SubmittingBadSuffixLocalGherkinTest
}

# checks if it's been called by main, set this variable if it is
if [[ "$CALLED_BY_MAIN" == "" ]]; then
    source $BASEDIR/calculate-galasactl-executables.sh
    calculate_galasactl_executable
    test_gherkin_commands
fi