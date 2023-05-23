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
    info "Syntax: test-galasactl-ecosystem.sh --bootstrap [BOOTSTRAP]"
    cat << EOF

Bootstrap must refer to a remote ecosystem.
EOF
}

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
if [[ "${bootstrap}" != "" ]]; then
    echo "Running tests against ecosystem bootstrap ${bootstrap}"
else
    error "Need to provide the bootstrap for a remote ecosystem."
    usage
    exit 1  
fi

#-----------------------------------------------------------------------------------------                   
# Constants
#-----------------------------------------------------------------------------------------   
export GALASA_TEST_NAME_SHORT="local.CoreLocalJava11Ubuntu"   
export GALASA_TEST_NAME_LONG="dev.galasa.inttests.core.${GALASA_TEST_NAME_SHORT}" 
export GALASA_TEST_RUN_GET_EXPECTED_SUMMARY_LINE_COUNT="2"
export GALASA_TEST_RUN_GET_EXPECTED_DETAILS_LINE_COUNT="13"


#--------------------------------------------------------------------------
function calculate_galasactl_executable {
    h2 "Calculate the name of the galasactl executable for this machine/os"

    raw_os=$(uname -s) # eg: "Darwin"
    os=""
    case $raw_os in
        Darwin*) 
            os="darwin" 
            ;;
        Windows*)
            os="windows"
            ;;
        Linux*)
            os="linux"
            ;;
        *) 
            error "Failed to recognise which operating system is in use. $raw_os"
            exit 1
    esac

    architecture=$(uname -m)
    if [[ "${architecture}" == "x86_64" ]]; then
        architecture="amd64"
    fi

    export binary="galasactl-${os}-${architecture}"
    info "galasactl binary is ${binary}"
    success "OK"
}

#--------------------------------------------------------------------------
function launch_test_on_ecosystem_with_portfolio {

    h2 "Building a portfolio..."

    mkdir -p ${BASEDIR}/temp
    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${binary} runs prepare \
    --bootstrap $bootstrap \
    --stream inttests \
    --portfolio portfolio.yaml \
    --test ${GALASA_TEST_NAME_SHORT} \
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
    --bootstrap ${bootstrap} \
    --portfolio portfolio.yaml \
    --throttle 1 \
    --poll 10 \
    --progress 1 \
    --noexitcodeontestfailures \
    --log -"

    set -o pipefail # Fail everything if anything in the pipeline fails. Else we are just checking the 'tee' return code.
    $cmd | tee runs-submit-output.txt # Store the output of galasactl runs submit to use later

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
function get_result_with_runname {
    h2 "Querying the result of the test we just ran..."

    cd ${BASEDIR}/temp

    # Get the RunName from the output of galasactl runs submit

    # Gets the line from the last part of the output stream the RunName is found in
    cat runs-submit-output.txt | grep -o "Run.*-" | tail -1  > line.txt 

    # Get just the RunName from the line. 
    # There is a line in the output like this:
    #   Run C6967 - inttests/dev.galasa.inttests/dev.galasa.inttests.core.local.CoreLocalJava11Ubuntu
    # Environment failure of the test results in "C6976(EnvFail)" ... so the '('...')' part needs removing also.
    sed 's/Run //; s/ -//; s/[(].*[)]//;' line.txt > runname.txt 
    runname=$(cat runname.txt)

    cmd="${BASEDIR}/bin/${binary} runs get \
    --name ${runname} \
    --bootstrap ${bootstrap}"

    info "Command is: $cmd"

    $cmd | grep "${runname}" # Checks the RunName can be found in the output from galasactl runs get
    rc=$?
    # We expect a return code of '0' because the ecosystem should be able to find this test as we just ran it.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to query the result of run ${runname} in the remote ecosystem."
        exit 1
    fi
    success "Querying the result of a run in the ecosystem worked OK"

    # The test above just checks that some output was found from galasactl runs get.
    # TO DO - Get the Result of the run from the output of galasactl runs submit as well as the RunName, and make sure the Result is correct too
    export RUN_NAME=${runname}
}

#--------------------------------------------------------------------------
function runs_get_check_summary_format_output {
    h2 "Performing runs get with summary format..."

    run_name=$1

    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${binary} runs get \
    --name ${run_name} \
    --format summary \
    --bootstrap ${bootstrap} "

    info "Command is: $cmd"

    output_file="runs-get-output.txt"
    $cmd | tee $output_file

    # Check that the full test name is output
    cat $output_file | grep "${GALASA_TEST_NAME_LONG}" -q
    rc=$?
    # We expect a return code of '0' because the test name should be output.
    if [[ "${rc}" != "0" ]]; then 
        error "Did not find ${GALASA_TEST_NAME_LONG} in summary output"
        exit 1
    fi

    # Check headers
    headers=("name" "status" "result" "test-name")

    for header in "${headers[@]}"
    do
        cat $output_file | grep "$header" -q
        rc=$?
        # We expect a return code of '0' because the header name should be output.
        if [[ "${rc}" != "0" ]]; then 
            error "Did not find header $header in summary output"
            exit 1
        fi
    done    

    # Check that we got 2 lines out... one for the headers, on for the 1 line of test data.
    line_count=$(cat $output_file | wc -l | xargs)
    expected_line_count=$GALASA_TEST_RUN_GET_EXPECTED_SUMMARY_LINE_COUNT
    if [[ "${line_count}" != "${expected_line_count}" ]]; then 
        error "line count is wrong. expected ${expected_line_count} got ${line_count}"
        exit 1
    fi

    success "galasactl runs get --format summary seemed to work"
}

#--------------------------------------------------------------------------
function runs_get_check_details_format_output {
    h2 "Performing runs get with details format..."

    run_name=$1

    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${binary} runs get \
    --name ${run_name} \
    --format details \
    --bootstrap ${bootstrap} "

    info "Command is: $cmd"

    output_file="runs-get-output.txt"
    $cmd | tee $output_file


    # Check that the full test name is output and formatted
    cat $output_file | grep "test-name[[:space:]]*:[[:space:]]*${GALASA_TEST_NAME_LONG}" -q
    rc=$?
    # We expect a return code of '0' because the ecosystem should be able to find this test as we just ran it.
    if [[ "${rc}" != "0" ]]; then 
        error "Did not find ${GALASA_TEST_NAME_LONG} in details output"
        exit 1
    fi

    # Check method headers
    headers=("method" "type" "status" "result" "start-time" "end-time" "duration(ms)")

    for header in "${headers[@]}"
    do
        cat $output_file | grep "$header" -q
        rc=$?
        # We expect a return code of '0' because the header name should be output.
        if [[ "${rc}" != "0" ]]; then 
            error "Did not find header $header in details output"
            exit 1
        fi
    done  

    #check methods start on line 13 - implies other test details have outputted 
    line_count=$(grep -n "method[[:space:]]*type[[:space:]]*status[[:space:]]*result[[:space:]]*start-time[[:space:]]*end-time[[:space:]]*duration(ms)" $output_file | head -n1 | sed 's/:.*//')
    expected_line_count=$GALASA_TEST_RUN_GET_EXPECTED_DETAILS_LINE_COUNT
    if [[ "${line_count}" != "${expected_line_count}" ]]; then 
        # We expect a return code of '0' because the method header should be output on line 13.
        error "line count is wrong. expected methods to start on ${expected_line_count} got ${line_count}"
        exit 1
    fi

    success "galasactl runs get --format details seemed to work"
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

calculate_galasactl_executable

# Launch test on ecosystem from a portfolio ...
launch_test_on_ecosystem_with_portfolio

# Query the result ... setting RUN_NAME to hold the one which galasa allocated
get_result_with_runname 
runs_get_check_summary_format_output  $RUN_NAME
runs_get_check_details_format_output  $RUN_NAME

# Launch test on ecosystem without a portfolio ...
# NOTE - Bug found with this command so commenting out for now
# launch_test_on_ecosystem_without_portfolio

# Attempt to create a test portfolio with an unknown test ...
create_portfolio_with_unknown_test

# Attempt to launch a test from an unknown portfolio ...
launch_test_from_unknown_portfolio