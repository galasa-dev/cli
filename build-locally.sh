#!/usr/bin/env bash


# Where is this script executing from ?
BASEDIR=$(dirname "$0");pushd $BASEDIR 2>&1 >> /dev/null ;BASEDIR=$(pwd);popd 2>&1 >> /dev/null
# echo "Running from directory ${BASEDIR}"
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
    info "Syntax: build-locally.sh [OPTIONS]"
    cat << EOF
Options are:
-c | --clean : Do a clean build. One of the --clean or --delta flags are mandatory.
-d | --delta : Do a delta build. One of the --clean or --delta flags are mandatory.
EOF
}

#--------------------------------------------------------------------------
# 
# Main script logic
#
#--------------------------------------------------------------------------

#-----------------------------------------------------------------------------------------                   
# Process parameters
#-----------------------------------------------------------------------------------------                   
build_type=""

while [ "$1" != "" ]; do
    case $1 in
        -c | --clean )          build_type="clean"
                                ;;
        -d | --delta )          build_type="delta"
                                ;;
        -h | --help )           usage
                                exit
                                ;;
        * )                     error "Unexpected argument $1"
                                usage
                                exit 1
    esac
    shift
done

if [[ "${build_type}" == "" ]]; then
    error "Need to use either the --clean or --delta parameter."
    usage
    exit 1  
fi

#--------------------------------------------------------------------------
h1 "Building the CLI component"
#--------------------------------------------------------------------------

#--------------------------------------------------------------------------
# Check that the ../framework is present.
h2 "Making sure the openapi yaml file is available..."
if [[ ! -e "../framework" ]]; then
    error "../framework is not present. Clone the framework repository."
    info "The openapi.yaml file from the framework repository is needed to generate a go client for the rest API"
    exit 1
fi

if [[ ! -e "../framework/openapi.yaml" ]]; then 
    error "File ../framework/openapi.yaml is not found."
    info "The openapi.yaml file from the framework repository is needed to generate a go client for the rest API"
    exit 1
fi
success "OK"

#--------------------------------------------------------------------------
h2 "Setting versions of things."
# Could get this bootjar from https://development.galasa.dev/main/maven-repo/obr/dev/galasa/galasa-boot/0.26.0/
export BOOT_JAR_VERSION="0.26.0"
info "BOOT_JAR_VERSION=${BOOT_JAR_VERSION}"
success "OK"

#--------------------------------------------------------------------------
# Create a temporary folder which is never checked in.
function download_dependencies {
    h2 "Making sure the tools folder is present."

    if [[ "${build_type}" == "clean" ]]; then
        h2 "Cleaning the dependencies out..."
        rm -fr build/dependencies
    fi

    info "Making sure the boot jar we embed is a fresh one from maven."
    rm -fr pkg/embedded/templates/galasahome/lib/*.jar

    mkdir -p build/dependencies
    rc=$? ; if [[ "${rc}" != "0" ]]; then error "Failed to ensure the tools folder is present. rc=${rc}" ; exit 1 ; fi
    success "OK"

    #--------------------------------------------------------------------------
    # Download the dependencies we define in gradle into a local folder
    h2 "Downloading dependencies"
    gradle installJarsIntoTemplates --warning-mode all
    rc=$? ; if [[ "${rc}" != "0" ]]; then  error "Failed to run the gradle build to get our dependencies. rc=${rc}" ; exit 1 ; fi
    success "OK"
}


#--------------------------------------------------------------------------
# Invoke the generator
function generate_rest_client {
    h2 "Generate the openapi client go code..."
    
    # Pick up and use the openapi generator we just downloaded.
    # We don't know which version it is (dictated by the gradle build), but as there
    # is only one we can just pick the filename up..
    # Should end up being something like: ${BASEDIR}/build/dependencies/openapi-generator-cli-6.2.0.jar
    export OPENAPI_GENERATOR_CLI_JAR=$(ls ${BASEDIR}/build/dependencies/openapi-generator-cli*)
    

    if [[ "${build_type}" == "clean" ]]; then
        h2 "Cleaning the generated code out..."
        rm -fr ${BASEDIR}/pkg/galasaapi/*
    fi

    mkdir -p build
    ./genapi.sh 2>&1 > build/generate-log.txt
    rc=$? ; if [[ "${rc}" != "0" ]]; then cat build/generate-log.txt ; error "Failed to generate the code from the yaml file. rc=${rc}" ; exit 1 ; fi
    rm -f build/generate-log.txt
    success "Code generation OK"

    #--------------------------------------------------------------------------
    # Invoke the generator again with different parameters
    h2 "Generate the openapi client go code... part II"
    ./generate.sh 2>&1 > build/generate-log.txt
    rc=$? ; if [[ "${rc}" != "0" ]]; then cat build/generate-log.txt ; error "Failed to generate II the code from the yaml file. rc=${rc}" ; exit 1 ; fi
    rm -f build/generate-log.txt
    success "Code generation part II - OK"
}


#--------------------------------------------------------------------------
#
# Build the executables
#
#--------------------------------------------------------------------------
function build_executables {
    if [[ "${build_type}" == "clean" ]]; then
        h2 "Cleaning the binaries out..."
        make clean
        rc=$? ; if [[ "${rc}" != "0" ]]; then error "Failed to build binary executable galasactl programs. rc=${rc}" ; exit 1 ; fi
        success "Binaries cleaned up - OK"
    fi

    h2 "Building new binaries..."
    make all
    rc=$? ; if [[ "${rc}" != "0" ]]; then error "Failed to build binary executable galasactl programs. rc=${rc}" ; exit 1 ; fi
    success "New binaries built - OK"
}



#--------------------------------------------------------------------------
#
# Testing what was built...
#
#--------------------------------------------------------------------------

#--------------------------------------------------------------------------
function calculate_galasactl_executable {
    h2 "Calculate the name of the galasactl executable for this machine/os"

    rm -fr ${BASEDIR}/temp
    mkdir -p ${BASEDIR}/temp
    cd ${BASEDIR}/temp

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

    export galasactl_command="galasactl-${os}-${architecture}"
    info "galasactl command is ${galasactl_command}"
    success "OK"
}

#--------------------------------------------------------------------------
# Invoke the galasactl command to create a project.
function generate_sample_code {
    h2 "Invoke the tool to create a sample project."

    cd $BASEDIR/temp

    export PACKAGE_NAME="dev.galasa.example.banking"
    ${BASEDIR}/bin/${galasactl_command} project create --package ${PACKAGE_NAME} --features payee,account --obr 
    rc=$?
    if [[ "${rc}" != "0" ]]; then
        error " Failed to create the galasa test project using galasactl command. rc=${rc}"
        exit 1
    fi
    success "OK"
}

#--------------------------------------------------------------------------
# Now build the source it created.
function build_generated_source {
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
# Execute the tests in Galasa

function run_test_java_minus_jar_method {

    cd $BASEDIR/temp

    TEST_BUNDLE=$1
    TEST_JAVA_CLASS=$2
    TEST_OBR_GROUP_ID=$3
    TEST_OBR_ARTIFACT_ID=$4
    TEST_OBR_VERSION=$5

    export M2_PATH=$(cd ~/.m2 ; pwd)
    export BOOT_JAR_PATH=${BASEDIR}/pkg/embedded/templates/galasahome/lib/galasa-boot-${BOOT_JAR_VERSION}.jar

    export OBR_VERSION="0.26.0"


    # Local .m2 content over-rides these anyway...
    # use development version of the OBR
    export REMOTE_MAVEN=https://development.galasa.dev/main/maven-repo/obr/
    # else go to maven central
    #export REMOTE_MAVEN=https://repo.maven.apache.org/maven2



    info "Running the following command..."
    cat << EOF 

    java -jar ${BOOT_JAR_PATH} \\
    --localmaven file:${M2_PATH}/repository/ \\
    --remotemaven $REMOTE_MAVEN \\
    --bootstrap file:${HOME}/.galasa/bootstrap.properties \\
    --overrides file:${HOME}/.galasa/overrides.properties \\
    --obr mvn:dev.galasa/dev.galasa.uber.obr/${OBR_VERSION}/obr \\
    --obr mvn:${TEST_OBR_GROUP_ID}/${TEST_OBR_ARTIFACT_ID}/${TEST_OBR_VERSION}/obr \\
    --test ${TEST_BUNDLE}/${TEST_JAVA_CLASS} 

EOF


    java -jar ${BOOT_JAR_PATH} \
    --localmaven file:${M2_PATH}/repository/ \
    --remotemaven $REMOTE_MAVEN \
    --bootstrap file:${HOME}/.galasa/bootstrap.properties \
    --overrides file:${HOME}/.galasa/overrides.properties \
    --obr mvn:dev.galasa/dev.galasa.uber.obr/${OBR_VERSION}/obr \
    --obr mvn:${TEST_OBR_GROUP_ID}/${TEST_OBR_ARTIFACT_ID}/${TEST_OBR_VERSION}/obr \
    --test ${TEST_BUNDLE}/${TEST_JAVA_CLASS} | tee jvm-log.txt | grep "[*][*][*]" | grep -v "[*][*][*][*]" | sed -e "s/[--]*//g"
    cat jvm-log.txt | grep "Passed - Test class ${TEST_JAVA_CLASS}"
    rc=$?
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to run the test"
        exit 1
    fi
    success "Test ran OK"
}






#--------------------------------------------------------------------------
# Build a portfolio
function build_portfolio {
    h2 "Building a portfolio file"

    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${galasactl_command} runs prepare \
    --portfolio my.portfolio \
    --bootstrap file:~/.galasa/bootstrap.properties \
    --class dev.galasa.example.banking.account/dev.galasa.example.banking.account.TestAccountExtended \
    --class dev.galasa.example.banking.account/dev.galasa.example.banking.account.TestAccount \
    --log -"

    TEST_JAVA_CLASS=dev.galasa.example.banking.payee.TestPayeeExtended
    run_test_java_minus_jar_method $TEST_BUNDLE $TEST_JAVA_CLASS $TEST_OBR_GROUP_ID $TEST_OBR_ARTIFACT_ID $TEST_OBR_VERSION

    # Run the Account tests
    TEST_BUNDLE=dev.galasa.example.banking.account
    TEST_JAVA_CLASS=dev.galasa.example.banking.account.TestAccount
    TEST_OBR_GROUP_ID=dev.galasa.example.banking
    TEST_OBR_ARTIFACT_ID=dev.galasa.example.banking.obr
    TEST_OBR_VERSION=0.0.1-SNAPSHOT
    run_test_java_minus_jar_method $TEST_BUNDLE $TEST_JAVA_CLASS $TEST_OBR_GROUP_ID $TEST_OBR_ARTIFACT_ID $TEST_OBR_VERSION

    TEST_JAVA_CLASS=dev.galasa.example.banking.account.TestAccountExtended
    run_test_java_minus_jar_method $TEST_BUNDLE $TEST_JAVA_CLASS $TEST_OBR_GROUP_ID $TEST_OBR_ARTIFACT_ID $TEST_OBR_VERSION   
 
    success "OK"
}


#--------------------------------------------------------------------------
# Build a portfolio
function build_portfolio {
    h2 "Building a portfolio file"

    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${galasactl_command} runs prepare \
    --portfolio my.portfolio \
    --bootstrap file:~/.galasa/bootstrap.properties \
    --class dev.galasa.example.banking.account/dev.galasa.example.banking.account.TestAccountExtended \
    --class dev.galasa.example.banking.account/dev.galasa.example.banking.account.TestAccount \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to build a portfolio file"
        exit 1
    fi
    success "Built portfolio OK"
}


#--------------------------------------------------------------------------
# Initialise Galasa home
function galasa_home_init {
    h2 "Initialising galasa home directory"

    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${galasactl_command} local init \
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
function launch_test_on_ecosystem {
    h2 "Launching test on an ecosystem..."

    if [[ "${GALASA_BOOTSTRAP}" == "" ]]; then 
        error "GALASA_BOOTSTRAP environment variable is not set. It should refer to a remote ecosystem"
        exit 1
    fi

    # hostname=$(echo -n "${GALASA_BOOTSTRAP}" | sed -e "s/http:\/\///g" | sed -e "s/https:\/\///g" | sed -e "s/.bootstrap//g")
    # info "Host name for boostrap is ${hostname}"

    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${galasactl_command} runs submit \
    --bootstrap $GALASA_BOOTSTRAP \
    --class dev.galasa.example.banking.account/dev.galasa.example.banking.account.TestAccount \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    # We expect a return code of '2' because the ecosystem doesn't know about this testcase.
    if [[ "${rc}" != "2" ]]; then 
        error "Failed to submit a test to a remote ecosystem, and get Unknown back."
        exit 1
    fi
    success "Submitting test to ecosystem worked OK"
}



#--------------------------------------------------------------------------

#--------------------------------------------------------------------------
# Return to the top folder so we can do other things.
cd ${BASEDIR}




#--------------------------------------------------------------------------
# Initialise Galasa home
function galasa_home_init {
    h2 "Initialising galasa home directory"

    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${galasactl_command} local init \
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
function launch_test_on_ecosystem {
    h2 "Launching test on an ecosystem..."

    if [[ "${GALASA_BOOTSTRAP}" == "" ]]; then 
        error "GALASA_BOOTSTRAP environment variable is not set. It should refer to a remote ecosystem"
        exit 1
    fi

    rm -fr ~/.galasa-old
    mv ~/.galasa ~/.galasa-old

    # hostname=$(echo -n "${GALASA_BOOTSTRAP}" | sed -e "s/http:\/\///g" | sed -e "s/https:\/\///g" | sed -e "s/.bootstrap//g")
    # info "Host name for boostrap is ${hostname}"

    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${galasactl_command} runs submit \
    --bootstrap $GALASA_BOOTSTRAP \
    --class dev.galasa.example.banking.account/dev.galasa.example.banking.account.TestAccount \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    # We expect a return code of '2' because the ecosystem doesn't know about this testcase.
    if [[ "${rc}" != "2" ]]; then 
        error "Failed to submit a test to a remote ecosystem, and get Unknown back."
        exit 1
    fi
    success "Submitting test to ecosystem worked OK"

    mv ~/.galasa-old ~/.galasa
}

#--------------------------------------------------------------------------
# Build the documentation
function generate_galasactl_documentation {
    generated_docs_folder=${BASEDIR}/docs/generated
    h2 "Generating documentation"
    info "Documentation will be placed in ${generated_docs_folder}"
    mkdir -p ${generated_docs_folder}

    # Figure out which type of machine this script is currently running on.
    unameOut="$(uname -s)"
    case "${unameOut}" in
        Linux)      machine=linux;;
        Darwin)     machine=darwin;;
        *)          error "Unknown machine type ${unameOut}"
                    exit 1
    esac
    architecture="$(uname -m)"
    if [[ "${architecture}" == "x86_64" ]]; then
        architecture="amd64"
    fi

    # Call the documentation generator, which builds .md files
    info "Using program ${BASEDIR}/bin/gendocs-galasactl-${machine}-${architecture} to generate the documentation..."
    ${BASEDIR}/bin/gendocs-galasactl-${machine}-${architecture} ${generated_docs_folder}
    rc=$? ; if [[ "${rc}" != "0" ]]; then error "Failed to generate documentation. rc=${rc}" ; exit 1 ; fi

    # The files have a line "###### Auto generated by cobra at 17/12/2022"
    # As we are (currently) checking-in these .md files, we don't want them to show as 
    # changed in git (which compares the content, not timestamps).
    # So lets remove these lines from all the .md files.
    info "Removing lines with date/time in, to limit delta changes in git..."
    mkdir -p ${BASEDIR}/build
    temp_file="${BASEDIR}/build/temp.md"
    for FILE in ${generated_docs_folder}/*; do 
        mv -f ${FILE} ${temp_file}
        cat ${temp_file} | grep -v "###### Auto generated by" > ${FILE}
        rm ${temp_file}
        success "Processed file ${FILE}"
    done
    success "Documentation generated - OK"
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

    export M2_PATH=$(cd ~/.m2 ; pwd)
    export BOOT_JAR_PATH=~/.galasa/lib/${GALASA_VERSION}/galasa-boot-${BOOT_JAR_VERSION}.jar


    # Local .m2 content over-rides these anyway...
    # use development version of the OBR
    export REMOTE_MAVEN=https://development.galasa.dev/main/maven-repo/obr/
    # else go to maven central
    #export REMOTE_MAVEN=https://repo.maven.apache.org/maven2

    export GALASACTL="${BASEDIR}/bin/galasactl-darwin-arm64"

    ${GALASACTL} runs submit local \
    --obr mvn:${OBR_GROUP_ID}/${OBR_ARTIFACT_ID}/${OBR_VERSION}/obr \
    --class ${BUNDLE}/${JAVA_CLASS} \
    --throttle 1 \
    --requesttype MikeCLI \
    --poll 10 \
    --progress 1 \
    --log - 2>&1 | tee ${BASEDIR}/temp/local-run-log.txt

    # --reportjson myreport.json \
    # --reportyaml myreport.yaml \

    
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

    # Run the Payee tests.
    export TEST_BUNDLE=dev.galasa.example.banking.payee
    export TEST_JAVA_CLASS=dev.galasa.example.banking.payee.TestPayee
    export TEST_OBR_GROUP_ID=dev.galasa.example.banking
    export TEST_OBR_ARTIFACT_ID=dev.galasa.example.banking.obr
    export TEST_OBR_VERSION=0.0.1-SNAPSHOT

    submit_local_test $TEST_BUNDLE $TEST_JAVA_CLASS $TEST_OBR_GROUP_ID $TEST_OBR_ARTIFACT_ID $TEST_OBR_VERSION
}

download_dependencies
generate_rest_client
build_executables
calculate_galasactl_executable
galasa_home_init
generate_sample_code
build_generated_source
run_tests_java_minus_jar_method
run_test_locally_using_galasactl
# build_portfolio
generate_galasactl_documentation
launch_test_on_ecosystem

#--------------------------------------------------------------------------
h2 "Use the results.."
info "Binary executable programs are found in the 'bin' folder."
ls ${BASEDIR}/bin | grep -v "gendocs"