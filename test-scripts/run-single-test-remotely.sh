#!/bin/bash

#
# Copyright contributors to the Galasa project
#
# SPDX-License-Identifier: EPL-2.0
#

# Where is this script executing from ?
BASEDIR=$(dirname "$0");pushd $BASEDIR 2>&1 >> /dev/null ;BASEDIR=$(pwd);popd 2>&1 >> /dev/null
export ORIGINAL_DIR=$(pwd)


mkdir -p ${BASEDIR}/../temp/home
cd ${BASEDIR}/../temp/home
export GALASA_HOME=$(pwd)

cd "${BASEDIR}/.."

galasactl runs submit \
--class dev.galasa.inttests/dev.galasa.inttests.core.local.CoreLocalJava11Ubuntu \
--stream inttests \
--throttle 1 \
--poll 10 \
--progress 1 \
--noexitcodeontestfailures \
--log - \
--overridefile ${GALASA_HOME}/overrides.properties

# galasactl runs prepare --portfolio my.portfolio --class dev.galasa.inttests/dev.galasa.inttests.core.local.CoreLocalJava11Ubuntu --stream inttests
# galasactl runs submit --portfolio my.portfolio --log -