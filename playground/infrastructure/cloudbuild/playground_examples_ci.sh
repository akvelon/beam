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

export GRADLE_VERSION=7.5.1
export GO_VERSION=1.18

#Install python java8 and dependencies
apt-get update > /dev/null
export DEBIAN_FRONTEND=noninteractive
apt-get install -y apt-transport-https ca-certificates software-properties-common curl unzip apt-utils > /dev/null
apt install net-tools
add-apt-repository -y ppa:deadsnakes/ppa > /dev/null && apt update > /dev/null
apt install -y python3.8 python3.8-distutils python3-pip > /dev/null
apt install --reinstall python3.8-distutils > /dev/null
pip install --upgrade google-api-python-client > /dev/null
python3.8 -m pip install pip --upgrade > /dev/null
ln -s /usr/bin/python3.8 /usr/bin/python > /dev/null
apt install python3.8-venv > /dev/null
pip install -r playground/infrastructure/requirements.txt

# Install jdk and gradle
apt-get install openjdk-8-jdk -y > /dev/null
curl -L https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip -o gradle-${GRADLE_VERSION}-bin.zip > /dev/null
unzip gradle-${GRADLE_VERSION}-bin.zip > /dev/null
export PATH=$PATH:gradle-${GRADLE_VERSION}/bin

# Install go
curl -OL https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz > /dev/null
tar -C /usr/local -xvf go$GO_VERSION.linux-amd64.tar.gz > /dev/null
export PATH=$PATH:/usr/local/go/bin

# Install Docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable" > /dev/null
apt update > /dev/null && apt install -y docker-ce > /dev/null

export \
ORIGIN=PG_EXAMPLES \
STEP=CI \
SUBDIRS="././learning/katas ././examples ././sdks" \
GOOGLE_CLOUD_PROJECT=$PROJECT_ID \
BEAM_ROOT_DIR="../../" \
SDK_CONFIG="../../playground/sdks.yaml" \
BEAM_EXAMPLE_CATEGORIES="../categories.yaml" \
BEAM_CONCURRENCY=4 \
BEAM_VERSION=2.43.0 \
sdks=("java" "python" "go") \
allowlist=("playground/infrastructure" "playground/backend")

echo "Environment variables exported"
ls -la playground

# Get Difference
# define the base ref
base_ref=refs/heads/master
if [[ -z "$base_ref" ]] || [[ "$base_ref" == "refs/heads/master" ]]
then
  base_ref=refs/heads/master
fi
diff=$(git diff --name-only $base_ref | tr '\n' ' ')

cd playground/infrastructure
# Check if there are Examples
for sdk in "${sdks[@]}"
do
      python3 checker.py \
      --verbose \
      --sdk SDK_"${sdk^^}" \
      --allowlist "${allowlist}" \
      --paths "${diff}"
if [[ $? -eq 0 ]]
  then
      example_has_changed=True
  else
      example_has_changed=False
fi
done

cd ../..

if [[ ${example_has_changed} == True ]]
then
    if [[ -z ${TAG_NAME} ]]
    then
        DOCKERTAG=${COMMIT_SHA}
    elif [[ -z ${COMMIT_SHA} ]]
    then
        DOCKERTAG=${TAG_NAME}
    fi
    for sdk in "${sdks[@]}"
    do
      if [[ "$sdk" == "python" ]]
      then
          # builds apache/beam_python3.7_sdk:$DOCKERTAG image
          ./gradlew -i :sdks:python:container:py37:docker -Pdocker-tag=$DOCKERTAG
          # and set SDK_TAG to DOCKERTAG so that the next step would find it
          SDK_TAG=${DOCKERTAG}
      fi

      opts=" -Pdocker-tag=${DOCKERTAG}"
      if [[ -n "$SDK_TAG" ]]
      then
          opts="${opts} -Psdk-tag=${SDK_TAG}"
      fi

      if [[ "$sdk" == "java" ]]
      then
          # Java uses a fixed BEAM_VERSION
          opts="$opts -Pbase-image=apache/beam_java8_sdk:${BEAM_VERSION}"
      fi

      ./gradlew -i playground:backend:containers:${sdk}:docker ${opts} -Pdocker-repository-root="us-central1-docker.pkg.dev/sandbox-playground-008/playground-repository"
      docker push us-central1-docker.pkg.dev/sandbox-playground-008/playground-repository/beam_playground-backend-${sdk}:${DOCKERTAG}

      docker run -d -p 8080:8080 --network=cloudbuild -e PROTOCOL_TYPE=TCP us-central1-docker.pkg.dev/sandbox-playground-008/playground-repository/beam_playground-backend-${sdk}:${DOCKERTAG}
      docker ps -a
      docker inspect --format='{{.Name}}' $(docker ps -aq --no-trunc) | cut -c2- >> /workspace/container_names.txt

      cat /workspace/container_names.txt
      cd playground/infrastructure
      for container_name in $(cat /workspace/container_names.txt)
      do
          export SERVER_ADDRESS=${container_name}:8080
          python3 ci_cd.py \
          --step ${STEP} \
          --sdk SDK_"${sdk^^}" \
          --origin ${ORIGIN} \
          --subdirs ${SUBDIRS}
      done
    done
else
      echo "Example has NOT been changed"
fi