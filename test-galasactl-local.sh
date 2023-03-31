#!/usr/bin/env bash

# This script can be ran locally or executed in a pipeline to test the various built binaries of galasactl
# This script tests the 'galasactl project create' and 'galasactl runs submit local' commands
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
    info "Syntax: test-galasactl-local.sh --binary [OPTIONS]"
    cat << EOF
Options are:
galasactl-darwin-amd64 : Use the galasactl-darwin-amd64 binary
galasactl-darwin-arm64 : Use the galasactl-darwin-arm64 binary
galasactl-linux-amd64 : Use the galasactl-linux-amd64 binary
galasactl-linux-s390x : Use the galasactl-linux-s390x binary
galasactl-windows-amd64.exe : Use the galasactl-windows-amd64.exe binary
EOF
}

#-----------------------------------------------------------------------------------------                   
# Process parameters
#-----------------------------------------------------------------------------------------                   
binary=""

while [ "$1" != "" ]; do
    case $1 in
        --binary )                        shift
                                          binary="$1"
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

#--------------------------------------------------------------------------
# Initialise Galasa home
function galasa_home_init {
    h2 "Initialising galasa home directory"

    mkdir -p ${BASEDIR}/temp
    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${binary} local init \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to initialise galasa home"
        exit 1
    fi
    success "Galasa home initialised"
}

#--------------------------------------------------------------------------
# Invoke the galasactl command to create a project.
function generate_sample_code {
    h2 "Invoke the tool to create a sample project."

    cd ${BASEDIR}/temp

    export PACKAGE_NAME="dev.galasa.example.banking"
    ${BASEDIR}/bin/${binary} project create --package ${PACKAGE_NAME} --features payee,account --obr --maven --gradle --force
    rc=$?
    if [[ "${rc}" != "0" ]]; then
        error " Failed to create the galasa test project using galasactl command. rc=${rc}"
        exit 1
    fi
    success "OK"
}

#--------------------------------------------------------------------------
# Now build the source it created using maven
function build_generated_source_maven {
    h2 "Building the sample project we just generated."
    cd ${BASEDIR}/temp/${PACKAGE_NAME}
    mvn clean test install 
    rc=$?
    if [[ "${rc}" != "0" ]]; then
        error " Failed to build the generated source code which galasactl created."
        exit 1
    fi
    success "OK"
}

#--------------------------------------------------------------------------
# Now build the source it created using gradle
function build_generated_source_gradle {
    h2 "Building the sample project we just generated."
    cd ${BASEDIR}/temp/${PACKAGE_NAME}
    gradle build publishToMavenLocal
    rc=$?
    if [[ "${rc}" != "0" ]]; then
        error " Failed to build the generated source code which galasactl created."
        exit 1
    fi
    success "OK"
}

#--------------------------------------------------------------------------
# Run test using the galasactl locally in a JVM
function submit_local_test {

    h2 "Submitting a local test using galasactl in a local JVM"

    cd ${BASEDIR}/temp/*banking

    BUNDLE=$1
    JAVA_CLASS=$2
    OBR_GROUP_ID=$3
    OBR_ARTIFACT_ID=$4
    OBR_VERSION=$5

    # Could get this bootjar from https://development.galasa.dev/main/maven-repo/obr/dev/galasa/galasa-boot/0.26.0/
    export BOOT_JAR_VERSION="0.26.0"

    export GALASA_VERSION="0.26.0"

    export BOOT_JAR_PATH=~/.galasa/lib/${GALASA_VERSION}/galasa-boot-${BOOT_JAR_VERSION}.jar


    # Local .m2 content over-rides these anyway...
    # use development version of the OBR
    export REMOTE_MAVEN=https://development.galasa.dev/main/maven-repo/obr/
    # else go to maven central
    #export REMOTE_MAVEN=https://repo.maven.apache.org/maven2

    export GALASACTL="${BASEDIR}/bin/${binary}"

    ${GALASACTL} runs submit local \
    --obr mvn:${OBR_GROUP_ID}/${OBR_ARTIFACT_ID}/${OBR_VERSION}/obr \
    --class ${BUNDLE}/${JAVA_CLASS} \
    --throttle 1 \
    --requesttype automated-test \
    --poll 10 \
    --progress 1 \
    --log -

    # Uncomment this if testing that a test that should fail, fails
    # --noexitcodeontestfailures \

    # --remoteMaven https://development.galasa.dev/main/maven-repo/obr/ \
    # --galasaVersion 0.26.0 \

    rc=$?
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to run the test"
        exit 1
    fi
    success "Test ran OK"
}

function run_test_locally_using_galasactl {
    export LOG_FILE=$1
    
    # Run the Payee tests.
    export TEST_BUNDLE=dev.galasa.example.banking.payee
    export TEST_JAVA_CLASS=dev.galasa.example.banking.payee.TestPayee
    export TEST_OBR_GROUP_ID=dev.galasa.example.banking
    export TEST_OBR_ARTIFACT_ID=dev.galasa.example.banking.obr
    export TEST_OBR_VERSION=0.0.1-SNAPSHOT


    submit_local_test $TEST_BUNDLE $TEST_JAVA_CLASS $TEST_OBR_GROUP_ID $TEST_OBR_ARTIFACT_ID $TEST_OBR_VERSION $LOG_FILE
}

function cleanup_local_maven_repo {
    rm -fr ~/.m2/repository/dev/galasa/example
}

# Initialise Galasa home ...
galasa_home_init

# Generate sample project ...
generate_sample_code

# Maven ...
cleanup_local_maven_repo
build_generated_source_maven
run_test_locally_using_galasactl

# Gradle ...
cleanup_local_maven_repo
build_generated_source_gradle
run_test_locally_using_galasactl

