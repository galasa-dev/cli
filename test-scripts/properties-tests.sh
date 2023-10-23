#!/bin/bash

#
# Copyright contributors to the Galasa project
#
# SPDX-License-Identifier: EPL-2.0
#
echo "Running script properties-tests.sh"
# This script can be ran locally or executed in a pipeline to test the various built binaries of galasactl
# This script tests the 'galasactl runs submit' command against a test that is in our ecosystem's testcatalog already
# Pre-requesite: the CLI must have been built first so the binaries are present in the /bin directory

# Where is this script executing from ?
BASEDIR=$(dirname "$0");pushd $BASEDIR 2>&1 >> /dev/null ;BASEDIR=$(pwd);popd 2>&1 >> /dev/null
export ORIGINAL_DIR=$(pwd)
cd "${BASEDIR}"

source ${BASEDIR}/calculate_galasactl_executables.sh

# generate a random number to append to test names to avoid multiple running at once overriding each other
function get_random_property_name_number {
    minimum=100
    maximum=999
    PROP_NUM=$(($minimum + $RANDOM % $maximum))
    echo $PROP_NUM
}

#--------------------------------------------------------------------------
function properties_create {
    h2 "Performing properties set with name and value parameter used: create..."

    set -o pipefail # Fail everything if anything in the pipeline fails. Else we are just checking the 'tee' return code.
    cd ${BASEDIR}/temp
    prop_name="properties.test.name.value.$PROP_NUM"

    cmd="${BASEDIR}/bin/${binary} properties set --namespace ecosystemtest \
    --name $prop_name \
    --value test-value \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    # We expect a return code of 0 because this is a properly formed properties set command.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to create property with name and value used."
        exit 1
    fi

    # check that property has been created
    cmd="${BASEDIR}/bin/${binary} properties get --namespace ecosystemtest \
    --name $prop_name \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    output_file="properties-get-output.txt"
    $cmd | tee $output_file
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get property with name used: command failed."
        exit 1
    fi

    # Check that the previous properties set created a property
    cat $output_file | grep "$prop_name" -q

    rc=$?
    # We expect a return code of 0 because this is a properly formed properties get command.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to create property with name and value used."
        exit 1
    fi

    # Check that the previous properties set created a property, with the correct value
    cat $output_file | grep "$prop_name test-value" -q

    rc=$?
    # We expect a return code of 0 because this is a properly formed properties get command.
    if [[ "${rc}" != "0" ]]; then 
        error "Property successfully created, but value incorrect."
        exit 1
    fi

   

    success "Properties set with name and value used seems to have been created correctly."
}

#--------------------------------------------------------------------------
function properties_update {
    h2 "Performing properties set with name and value parameter used: update..."

    used -o pipefail # Fail everything if anything in the pipeline fails. Else we are just checking the 'tee' return code.
    cd ${BASEDIR}/temp
    prop_name="properties.test.name.value.$PROP_NUM"

    cmd="${BASEDIR}/bin/${binary} properties set --namespace ecosystemtest \
    --name $prop_name \
    --value updated-value \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    # We expect a return code of 0 because this is a properly formed properties set command.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to update property, set command failed."
        exit 1
    fi

    # check that property has been updated
    cmd="${BASEDIR}/bin/${binary} properties get --namespace ecosystemtest \
    --name $prop_name \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    output_file="properties-get-output.txt"
    $cmd | tee $output_file
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get property with name used, get command failed."
        exit 1
    fi

    # Check that the previous properties set updated the property value
    cat $output_file | grep "$prop_name updated-value" -q

    rc=$?
    # We expect a return code of 0 because this is a properly formed properties get command.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to update property, property not found in namespace."
        exit 1
    fi

    success "Properties set with name and value used seems to have been updated correctly."
}

#--------------------------------------------------------------------------
function properties_delete {
    h2 "Performing properties delete with name parameter used..."

    set -o pipefail # Fail everything if anything in the pipeline fails. Else we are just checking the 'tee' return code.
    cd ${BASEDIR}/temp
    prop_name="properties.test.name.value.$PROP_NUM"

    cmd="${BASEDIR}/bin/${binary} properties delete --namespace ecosystemtest \
    --name $prop_name \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    # We expect a return code of 0 because this is a properly formed properties delete command.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to delete property, command failed."
        exit 1
    fi

    # check that property has been updated
    cmd="${BASEDIR}/bin/${binary} properties get --namespace ecosystemtest \
    --name $prop_name \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    output_file="properties-get-output.txt"
    $cmd | tee $output_file
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get property with name used."
        exit 1
    fi

    # Check that the previous properties set updated the property value
    cat $output_file | grep "Total:0" -q

    rc=$?
    # We expect a return code of 1 because this property should not exist anymore.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to delete property, property remains in namespace."
        exit 1
    fi

    success "Properties delete with name used seems to have been deleted correctly."
}

#--------------------------------------------------------------------------
function properties_delete_invalid_property {
    h2 "Performing properties delete with name parameter used..."

    set -o pipefail # Fail everything if anything in the pipeline fails. Else we are just checking the 'tee' return code.
    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${binary} properties delete --namespace ecosystemtest \
    --name this.property.shouldnt.exist \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    # We expect a return code of 0 because this is a properly formed properties delete command.
    if [[ "${rc}" != "1" ]]; then 
        error "Command should have failed due to non existent property."
        exit 1
    fi
    success "Properties delete with the name of a non existent property correctly throws an error."
}

#--------------------------------------------------------------------------
function properties_delete_without_name {
    h2 "Performing properties delete with name parameter used..."

    set -o pipefail # Fail everything if anything in the pipeline fails. Else we are just checking the 'tee' return code.
    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${binary} properties delete --namespace ecosystemtest \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    # We expect a return code of 0 because this is a properly formed properties delete command.
    if [[ "${rc}" != "1" ]]; then 
        error "Command should have failed due to name not being used."
        exit 1
    fi
    success "Properties delete without the name flag used correctly throws an error."
}

#--------------------------------------------------------------------------
function properties_set_with_name_without_value {
    h2 "Performing properties set with name parameter, but no value parameter..."

    cd ${BASEDIR}/temp
    prop_name="properties.test.name.$PROP_NUM"

    cmd="${BASEDIR}/bin/${binary} properties set --namespace ecosystemtest \
    --name $prop_name \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    # we expect a return code of 1 as properties set should not be able to run without value used.
    if [[ "${rc}" != "1" ]]; then 
        error "Failed to recognise properties set without value should error."
        exit 1
    fi
    success "Properties set with no value correctly throws an error."
}

#--------------------------------------------------------------------------
function properties_set_without_name_with_value {
    h2 "Performing properties set without name parameter and with value parameter..."

    cd ${BASEDIR}/temp

    cmd="${BASEDIR}/bin/${binary} properties set --namespace ecosystemtest \
    --value random-arbitrary-value \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    # we expect a return code of 1 as properties set should not be able to run without name used.
    if [[ "${rc}" != "1" ]]; then 
        error "Failed to recognise properties set without name used should error."
        exit 1
    fi
    success "Properties set with no name correctly throws an error."
}

#--------------------------------------------------------------------------
function properties_set_without_name_and_value {
    h2 "Performing properties set without name and value parameter used..."

    cd ${BASEDIR}/temp
    prop_name="properties.test.name.value.$PROP_NUM"

    cmd="${BASEDIR}/bin/${binary} properties set --namespace ecosystemtest \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    $cmd
    rc=$?
    # we expect a return code of 1 as properties set should not be able to run without name and value used.
    if [[ "${rc}" != "1" ]]; then 
        error "Failed to recognise properties set without name and value used should error."
        exit 1
    fi
    success "Properties set with no name and value correctly throws an error."
}

#--------------------------------------------------------------------------
function properties_get_setup {
    h2 "Performing setup for subsequent properties get commands."
    cmd="${BASEDIR}/bin/${binary} properties set --namespace ecosystemtest \
    --name get.test.property \
    --value this-shouldnt-be-deleted \
    --bootstrap $bootstrap \
    --log -"
    $cmd
}

#--------------------------------------------------------------------------
function properties_get_with_namespace {
    h2 "Performing properties get with only namespace used, expecting list of properties..."
    
    cmd="${BASEDIR}/bin/${binary} properties get --namespace ecosystemtest \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    output_file="properties-get-output.txt"
    $cmd | tee $output_file
    rc=$?
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get properties with namespace used: command failed."
        exit 1
    fi

    # Check that the previous properties set created a property
    cat $output_file | grep "Total:([1-9])+" -q -E

    rc=$?
    # We expect a return code of 0 because this is a properly formed properties get command.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get a list of properties under the namespace"
        exit 1
    fi
    success "Properties get with namespace used seems to be successful."
}

#--------------------------------------------------------------------------
function properties_get_with_name {
    h2 "Performing properties get with only name used, expecting list of properties..."
    
    cmd="${BASEDIR}/bin/${binary} properties get --namespace ecosystemtest \
    --name get.test.property \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    output_file="properties-get-output.txt"
    $cmd | tee $output_file
    rc=$?
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get property with name used: command failed."
        exit 1
    fi

    # Check that the previous properties set created a property
    cat $output_file | grep "get.test.property this-shouldnt-be-deleted" -q -E

    rc=$?
    # We expect a return code of 0 because this is a properly formed properties get command.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get property whilst name is used"
        exit 1
    fi
    success "Properties get with name used seems to be successful."
}

#--------------------------------------------------------------------------
function properties_get_with_prefix {
    h2 "Performing properties get with prefix used, expecting list of properties..."
    
    cmd="${BASEDIR}/bin/${binary} properties get --namespace ecosystemtest \
    --prefix get \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    output_file="properties-get-output.txt"
    $cmd | tee $output_file
    rc=$?
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get properties with prefix used: command failed."
        exit 1
    fi

    # Check that the previous properties set created a property
    cat $output_file | grep "get.test.property this-shouldnt-be-deleted" -q -E

    rc=$?
    # We expect a return code of 0 because this is a properly formed properties get command.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get properties whilst prefix is used"
        exit 1
    fi
    success "Properties get with prefix used seems to be successful."
}

#--------------------------------------------------------------------------
function properties_get_with_suffix {
    h2 "Performing properties get with suffix used, expecting list of properties..."
    
    cmd="${BASEDIR}/bin/${binary} properties get --namespace ecosystemtest \
    --suffix property \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    output_file="properties-get-output.txt"
    $cmd | tee $output_file
    rc=$?
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get properties with suffix used: command failed."
        exit 1
    fi

    # Check that the previous properties set created a property
    cat $output_file | grep "get.test.property this-shouldnt-be-deleted" -q -E

    rc=$?
    # We expect a return code of 0 because this is a properly formed properties get command.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get properties whilst suffix is used"
        exit 1
    fi
    success "Properties get with suffix used seems to be successful."
}

#--------------------------------------------------------------------------
function properties_get_with_infix {
    h2 "Performing properties get with infix used, expecting list of properties..."
    
    cmd="${BASEDIR}/bin/${binary} properties get --namespace ecosystemtest \
    --infix test \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    output_file="properties-get-output.txt"
    $cmd | tee $output_file
    rc=$?
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get properties with infix used: command failed."
        exit 1
    fi

    # Check that the previous properties set created a property
    cat $output_file | grep "get.test.property this-shouldnt-be-deleted" -q -E

    rc=$?
    # We expect a return code of 0 because this is a properly formed properties get command.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get properties whilst infix is used"
        exit 1
    fi
    success "Properties get with infix used seems to be successful."
}

#--------------------------------------------------------------------------
function properties_get_with_prefix_infix_and_suffix {
    h2 "Performing properties get with prefix, infix, and suffix used, expecting list of properties..."
    
    cmd="${BASEDIR}/bin/${binary} properties get --namespace ecosystemtest \
    --prefix get \
    --suffix property \
    --infix test \
    --bootstrap $bootstrap \
    --log -"

    info "Command is: $cmd"

    output_file="properties-get-output.txt"
    $cmd | tee $output_file
    rc=$?
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get properties with prefix, infix, and suffix used: command failed."
        exit 1
    fi

    # Check that the previous properties set created a property
    cat $output_file | grep "get.test.property this-shouldnt-be-deleted" -q -E

    rc=$?
    # We expect a return code of 0 because this is a properly formed properties get command.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get properties whilst prefix, infix, and suffix is used"
        exit 1
    fi
    success "Properties get with prefix, infix, and suffix used seems to be successful."
}

#--------------------------------------------------------------------------
function properties_get_with_namespace_raw_format {
    h2 "Performing properties get with only namespace used, expecting list of properties..."
    
    cmd="${BASEDIR}/bin/${binary} properties get --namespace ecosystemtest \
    --bootstrap $bootstrap \
    --format raw \
    --log -"

    info "Command is: $cmd"

    output_file="properties-get-output.txt"
    $cmd | tee $output_file
    rc=$?
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get properties with namespace used: command failed."
        exit 1
    fi

    # Check that the previous properties set created a property
    cat $output_file | grep "ecosystemtest|get.test.property|this-shouldn't-be-deleted" -q

    rc=$?
    # We expect a return code of 0 because this is a properly formed properties get command.
    if [[ "${rc}" != "0" ]]; then 
        error "Failed to get a list of properties under the namespace"
        exit 1
    fi
    success "Properties get with namespace used seems to be successful."
}

function properties_tests {
    get_random_property_name_number
    properties_create
    properties_update
    properties_delete
    properties_delete_invalid_property
    properties_delete_without_name
    properties_set_with_name_without_value
    properties_set_without_name_with_value
    properties_set_without_name_and_value
    properties_get_setup
    properties_get_with_namespace
    properties_get_with_name
    properties_get_with_prefix
    properties_get_with_suffix
    properties_get_with_infix
    properties_get_with_prefix_infix_and_suffix
}

calculate_galasactl_executable
properties_tests