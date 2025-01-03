#! /usr/bin/env bash

# Where is this script executing from ?
BASEDIR=$(dirname "$0");pushd $BASEDIR 2>&1 >> /dev/null ;BASEDIR=$(pwd);popd 2>&1 >> /dev/null
# echo "Running from directory ${BASEDIR}"
export ORIGINAL_DIR=$(pwd)
cd "${BASEDIR}"

#-----------------------------------------------------------------------------------------                   
#
# Set Colors
#
#-----------------------------------------------------------------------------------------                   
bold=$(tput bold)
underline=$(tput sgr 0 1)
reset=$(tput sgr0)
red=$(tput setaf 1)
green=$(tput setaf 76)
white=$(tput setaf 7)
tan=$(tput setaf 202)
blue=$(tput setaf 25)

#-----------------------------------------------------------------------------------------                   
#
# Headers and Logging
#
#-----------------------------------------------------------------------------------------                   
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

h2 "Querying the RAS for old records we don't want any longer."
mkdir -p $BASEDIR/temp
cd $BASEDIR/temp >> /dev/null

set -o pipefail

galasactl runs get --age 120d:90d --format raw | cut -f1 -d'|' | sort | uniq > files-to-delete.txt
rc=$? ; if [[ "${rc}" != "0" ]]; then error "Failed to query the test runs to delete." ; exit 1 ; fi
success "Queried the run records we want to delete."

h2 "Deleting the run records we don't care about."
while IFS="" read -r line 
do
    testRunName="$line"
    info "Deleting test run $testRunName"
    galasactl runs delete --name $testRunName
    rc=$? ; if [[ "${rc}" != "0" ]]; then error "Failed to delete the test runs $testRunName" ; exit 1 ; fi
    success "Deleted test run $testRunName"
done < files-to-delete.txt

cd - >> /dev/null 