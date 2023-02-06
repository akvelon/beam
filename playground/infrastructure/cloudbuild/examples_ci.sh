#!/usr/bin/env bash

# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Added set -x to show output into cloud build console
set -x

export GRADLE_VERSION=7.5.1
export GO_VERSION=1.18

#Install python java8 and dependencies
apt-get update > /dev/null
apt update > /dev/null
export DEBIAN_FRONTEND=noninteractive

# Env configuration commands
apt-get install -y apt-transport-https ca-certificates software-properties-common curl unzip apt-utils > /dev/null
add-apt-repository -y ppa:deadsnakes/ppa > /dev/null && apt update > /dev/null
apt install -y python3.8 python3.8-distutils python3-pip > /dev/null
apt install --reinstall python3.8-distutils > /dev/null
pip install --upgrade google-api-python-client > /dev/null
python3.8 -m pip install pip --upgrade > /dev/null
ln -s /usr/bin/python3.8 /usr/bin/python > /dev/null
apt install python3.8-venv > /dev/null
pip install -r playground/infrastructure/requirements.txt > /dev/null

# Install jdk and gradle
apt-get install openjdk-8-jdk -y > /dev/null
curl -L https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip -o gradle-${GRADLE_VERSION}-bin.zip > /dev/null
unzip gradle-${GRADLE_VERSION}-bin.zip > /dev/null
export PATH=$PATH:gradle-${GRADLE_VERSION}/bin > /dev/null

# Install go
curl -OL https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz > /dev/null
tar -C /usr/local -xvf go$GO_VERSION.linux-amd64.tar.gz > /dev/null
export PATH=$PATH:/usr/local/go/bin > /dev/null

# Install Docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - > /dev/null
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable" > /dev/null
apt update > /dev/null && apt install -y docker-ce > /dev/null

# Assigning required values for CI_CD.py script
export \
ORIGIN=PG_EXAMPLES \
STEP=CI \
SUBDIRS="./learning/katas ./examples ./sdks" \
GOOGLE_CLOUD_PROJECT=${PROJECT_ID} \
BEAM_ROOT_DIR="." \
SDK_CONFIG="playground/sdks.yaml" \
BEAM_EXAMPLE_CATEGORIES="playground/categories.yaml" \
BEAM_CONCURRENCY=4 \
BEAM_VERSION=2.43.0 \
sdks=("java" "python" "go") \
allowlist=("playground/infrastructure" "playground/backend")

# Check whether commit is tagged
tag_name=$(git tag --points-at $commit_sha)

# Get diff
if [ -z $base_ref ] || [ $base_ref == "master" ]
then
    base_ref=origin/master
fi
diff=$(git diff --name-only $base_ref $commit_sha | tr '\n' ' ')

# Check if there are Examples
for sdk in "${sdks[@]}"
do
      python3 playground/infrastructure/checker.py \
      --verbose \
      --sdk SDK_"${sdk^^}" \
      --allowlist "${allowlist[@]}" \
      --paths ${diff}
      if [ $? -eq 0 ]
      then
          echo "Checker has found changed examples for ${sdk^^}" >> /tmp/build-log-${pr_number}-${commit_sha}-${BUILD_ID}.txt
          example_has_changed=True
      elif [ $? -eq 11 ]
      then
          echo "Checker has not found changed examples for ${sdk^^}" >> /tmp/build-log-${pr_number}-${commit_sha}-${BUILD_ID}.txt
          example_has_changed=False
      else
          echo "Error: Checker is broken" >> /tmp/build-log-${pr_number}-${commit_sha}-${BUILD_ID}.txt
          exit 1
      fi

# Run main logic if examples have been changed
      if [[ $example_has_changed == "True" ]]
      then
            if [ -z "${tag_name}" ] && [ "${commit_sha}" ]
            then
                DOCKERTAG=${commit_sha}
            elif [ "${tag_name}" ] && [ "${commit_sha}" ]
            then
                DOCKERTAG=${tag_name}
            elif [ "${tag_name}" ] && [ -z "${commit_sha}" ]
            then
                DOCKERTAG=${tag_name}
            elif [ -z "${tag_name}" ] && [ -z "${commit_sha}" ]
            then
                echo "Error: DOCKERTAG is empty"
                exit 1
            fi

            if [ "$sdk" == "python" ]
            then
                # builds apache/beam_python3.7_sdk:$DOCKERTAG image
                ./gradlew -i :sdks:python:container:py37:docker -Pdocker-tag=${DOCKERTAG}
                # and set SDK_TAG to DOCKERTAG so that the next step would find it
                SDK_TAG=${DOCKERTAG}
            else
                unset SDK_TAG
            fi

            opts=" -Pdocker-tag=${DOCKERTAG}"
            if [ -n "$SDK_TAG" ]
            then
                opts="${opts} -Psdk-tag=${SDK_TAG}"
            fi

            if [ "$sdk" == "java" ]
            then
                # Java uses a fixed BEAM_VERSION
                opts="$opts -Pbase-image=apache/beam_java8_sdk:2.43.0"
            fi

            ./gradlew -i playground:backend:containers:"${sdk}":docker ${opts}

            IMAGE_TAG=apache/beam_playground-backend-${sdk}:${DOCKERTAG}

            docker run -d -p 8080:8080 --network=cloudbuild -e PROTOCOL_TYPE=TCP --name container-${sdk} $IMAGE_TAG
            sleep 10
            export SERVER_ADDRESS=container-${sdk}:8080
            python3 playground/infrastructure/ci_cd.py \
            --step ${STEP} \
            --sdk SDK_"${sdk^^}" \
            --origin ${ORIGIN} \
            --subdirs ${SUBDIRS} >> /tmp/build-log-${pr_number}-${commit_sha}-${BUILD_ID}.txt

            docker stop container-${sdk}
            docker rm container-${sdk}
      else
            echo "Nothing changed in Examples. CI step is skipped" >> /tmp/build-log-${pr_number}-${commit_sha}-${BUILD_ID}.txt
      fi
done