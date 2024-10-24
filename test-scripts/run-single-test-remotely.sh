#!/bin/bash

#
# Copyright contributors to the Galasa project
#
# SPDX-License-Identifier: EPL-2.0
#

# Where is this script executing from ?
BASEDIR=$(dirname "$0");pushd $BASEDIR 2>&1 >> /dev/null ;BASEDIR=$(pwd);popd 2>&1 >> /dev/null
export ORIGINAL_DIR=$(pwd)
cd "${BASEDIR}"

# export GALASA_HOME=${BASEDIR}/../temp/home

# galasactl runs submit \
# --class dev.galasa.inttests/dev.galasa.inttests.core.local.CoreLocalJava11Ubuntu \
# --stream inttests \
# --throttle 1 \
# --poll 10 \
# --progress 1 \
# --noexitcodeontestfailures \
# --log - \
# --overridefile /Users/mcobbett/builds/galasa/code/src/github.com/galasa-dev/cli/temp/home/overrides.properties


galasactl runs submit \
--class dev.galasa.inttests/dev.galasa.inttests.core.local.CoreLocalJava11Ubuntu \
--stream inttests \
--throttle 1 \
--poll 10 \
--progress 1 \
--noexitcodeontestfailures \
--log - \
--overridefile /Users/mcobbett/builds/galasa/code/src/github.com/galasa-dev/cli/temp/home/overrides.properties

# galasactl runs prepare --portfolio my.portfolio --class dev.galasa.inttests/dev.galasa.inttests.core.local.CoreLocalJava11Ubuntu --stream inttests
# galasactl runs submit --portfolio my.portfolio --log -