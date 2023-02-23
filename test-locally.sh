#!/usr/bin/env bash

# Objectives: 
# Give the tooling a spin to basically make sure it still works.


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

#-------------------------------------------------------------------------
# Clean
#-------------------------------------------------------------------------
rm -fr temp
mkdir -p temp

# rm -fr ~/.galasa/*

#-------------------------------------------------------------------------
# Run tool, generate source
#-------------------------------------------------------------------------
cd temp
../bin/galasactl-darwin-arm64 project create --package dev.galasa.example.banking --features payee,account --obr --log -
rc=$?
if [[ "${rc}" != "0" ]]; then 
    error "Failed. rc=${rc}"
    exit 1
fi

#-------------------------------------------------------------------------
# Add a long-running test
#-------------------------------------------------------------------------
cat << EOF >> dev.galasa.example.banking/dev.galasa.example.banking.account/src/main/java/dev/galasa/example/banking/account/TestLongRunningAccount.java
package dev.galasa.example.banking.account;

import static org.assertj.core.api.Assertions.*;

import java.lang.Thread;

import dev.galasa.core.manager.*;
import dev.galasa.Test;

/**
 * A sample galasa test class 
 */
@Test
public class TestLongRunningAccount {

	// Galasa will inject an instance of the core manager into the following field
	@CoreManager
	public ICoreManager core;

	/**
	 * Test which demonstrates that the managers have been injected ok.
	 */
	@Test
	public void simpleSampleTest() throws Exception {
        int secondsToSleep = 20;
        Thread.sleep(secondsToSleep*1000);
		assertThat(core).isNotNull();
	}

}
EOF

#-------------------------------------------------------------------------
# Add a quick-running failing test
#-------------------------------------------------------------------------
cat << EOF >> dev.galasa.example.banking/dev.galasa.example.banking.account/src/main/java/dev/galasa/example/banking/account/TestFailingAccount.java
package dev.galasa.example.banking.account;

import static org.assertj.core.api.Assertions.*;

import java.lang.Thread;

import dev.galasa.core.manager.*;
import dev.galasa.Test;

/**
 * A sample galasa test class 
 */
@Test
public class TestFailingAccount {

	// Galasa will inject an instance of the core manager into the following field
	@CoreManager
	public ICoreManager core;

	/**
	 * Test which demonstrates that the managers have been injected ok.
	 */
	@Test
	public void simpleSampleTest() throws Exception {
        // Fail on purpose.
		assertThat(true).isFalse();
	}

}
EOF

#-------------------------------------------------------------------------
# Build the generated code
#-------------------------------------------------------------------------
cd dev.galasa.example.banking
mvn clean test install

#-------------------------------------------------------------------------
# Run the test
#-------------------------------------------------------------------------

function submit_local_test {

    TEST_BUNDLE=$1
    TEST_JAVA_CLASS=$2
    TEST_OBR_GROUP_ID=$3
    TEST_OBR_ARTIFACT_ID=$4
    TEST_OBR_VERSION=$5

    # Could get this bootjar from https://development.galasa.dev/main/maven-repo/obr/dev/galasa/galasa-boot/0.24.0/
    export BOOT_JAR_VERSION="0.24.0"

    export OBR_VERSION="0.25.0"

    export M2_PATH=$(cd ~/.m2 ; pwd)
    export BOOT_JAR_PATH=~/.galasa/lib/${OBR_VERSION}/galasa-boot-${BOOT_JAR_VERSION}.jar

    


    # Local .m2 content over-rides these anyway...
    # use development version of the OBR
    export REMOTE_MAVEN=https://development.galasa.dev/main/maven-repo/obr/
    # else go to maven central
    #export REMOTE_MAVEN=https://repo.maven.apache.org/maven2

    export GALASACTL="${BASEDIR}/bin/galasactl-darwin-arm64"

    ${GALASACTL} runs submit local \
    --obr mvn:dev.galasa.example.banking/dev.galasa.example.banking.obr/0.0.1-SNAPSHOT/obr \
    --class dev.galasa.example.banking.account/dev.galasa.example.banking.account.TestAccount \
    --class dev.galasa.example.banking.account/dev.galasa.example.banking.account.TestAccountExtended \
    --class dev.galasa.example.banking.account/dev.galasa.example.banking.account.TestFailingAccount \
    --class dev.galasa.example.banking.payee/dev.galasa.example.banking.payee.TestPayee \
    --class dev.galasa.example.banking.payee/dev.galasa.example.banking.payee.TestPayeeExtended \
    --class dev.galasa.example.banking.account/dev.galasa.example.banking.account.TestLongRunningAccount \
    --throttle 7 \
    --requesttype MikeCLI \
    --poll 10 \
    --progress 1 \
    --log - 2>&1 | tee ${BASEDIR}/temp/log.txt

    # --reportjson myreport.json \
    # --reportyaml myreport.yaml \

    
    # --noexitcodeontestfailures \

    # --remoteMaven https://development.galasa.dev/main/maven-repo/obr/ \
    # --galasaVersion 0.25.0 \

    rc=$?
    if [[ "${rc}" != "0" ]]; then 
        echo "Failed to run the test"
        exit 1
    fi
    echo "Test ran OK"
}


# Run the Payee tests.
export TEST_BUNDLE=dev.galasa.example.banking.payee
export TEST_JAVA_CLASS=dev.galasa.example.banking.payee.TestPayee
export TEST_OBR_GROUP_ID=dev.galasa.example.banking
export TEST_OBR_ARTIFACT_ID=dev.galasa.example.banking.obr
export TEST_OBR_VERSION=0.0.1-SNAPSHOT

submit_local_test $TEST_BUNDLE $TEST_JAVA_CLASS $TEST_OBR_GROUP_ID $TEST_OBR_ARTIFACT_ID $TEST_OBR_VERSION

    
