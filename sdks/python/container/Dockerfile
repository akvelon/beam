###############################################################################
#  Licensed to the Apache Software Foundation (ASF) under one
#  or more contributor license agreements.  See the NOTICE file
#  distributed with this work for additional information
#  regarding copyright ownership.  The ASF licenses this file
#  to you under the Apache License, Version 2.0 (the
#  "License"); you may not use this file except in compliance
#  with the License.  You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
#  Unless required by applicable law or agreed to in writing, software
#  distributed under the License is distributed on an "AS IS" BASIS,
#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#  See the License for the specific language governing permissions and
# limitations under the License.
###############################################################################

ARG py_version
FROM python:"${py_version}"-bookworm as beam
LABEL Author "Apache Beam <dev@beam.apache.org>"
ARG TARGETOS
ARG TARGETARCH

COPY target/base_image_requirements.txt /tmp/base_image_requirements.txt
COPY target/apache-beam.tar.gz /opt/apache/beam/tars/
COPY target/launcher/${TARGETOS}_${TARGETARCH}/boot target/LICENSE target/NOTICE target/LICENSE.python /opt/apache/beam/

ENV CLOUDSDK_CORE_DISABLE_PROMPTS yes
ENV PATH $PATH:/usr/local/gcloud/google-cloud-sdk/bin

# Use one RUN command to reduce the number of layers.

# Install native bindings required for dependencies
RUN apt-get update && apt-get install -y \
    libsnappy-dev \
    libyaml-dev \
    ccache \
    libgeos-dev \
    build-essential \
    g++ \
    python3-dev

RUN export CXXFLAGS="-std=c++11"

# Clean up apt cache
RUN rm -rf /var/lib/apt/lists/*

# Upgrade pip and install essential Python packages
RUN pip install --use-pep517 --upgrade pip setuptools wheel

# Install required packages for Beam Python SDK and common dependencies
RUN pip install --use-pep517 --no-deps -r /tmp/base_image_requirements.txt \
    && rm -rf /tmp/base_image_requirements.txt

# Download NLTK stopwords
RUN python -c "import nltk; nltk.download('stopwords')" \
    && rm /root/nltk_data/corpora/stopwords.zip

# Verify that the protobuf upb is used
RUN python -c "from google.protobuf.internal import api_implementation; assert api_implementation._implementation_type == 'upb'; print('Verified fast protobuf used.')"

# Install Google Cloud SDK
RUN mkdir -p /usr/local/gcloud \
    && cd /usr/local/gcloud \
    && curl -s -O https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.tar.gz \
    && tar -xf google-cloud-sdk.tar.gz \
    && /usr/local/gcloud/google-cloud-sdk/install.sh \
    && rm -rf /usr/local/gcloud/google-cloud-sdk/.install/.backup \
    && rm google-cloud-sdk.tar.gz

# Configure ccache for faster SDK builds
RUN ln -s /usr/bin/ccache /usr/local/bin/gcc \
    && ccache --set-config=sloppiness=file_macro \
    && ccache --set-config=hash_dir=false

RUN pip install --use-pep517 protobuf==5.29.2
RUN pip install --use-pep517 numpy==1.26.4

# Install Apache Beam SDK
RUN pip install --use-pep517 --no-deps -vvv /opt/apache/beam/tars/apache-beam.tar.gz[gcp]

RUN pip check

# Validate dependencies
RUN pip check || (echo "Container does not include required Beam dependencies or has conflicting dependencies. If Beam dependencies have changed, you need to regenerate base_image_requirements.txt files. See: https://s.apache.org/beam-python-requirements-generate" && exit 1)

# Log complete list of installed packages
RUN pip freeze --all

# Remove pip cache
RUN rm -rf /root/.cache/pip

ENTRYPOINT ["/opt/apache/beam/boot"]

####
# Pull and add third party licenses to the image if pull_licenses is true.
# Use multistage build to eliminate unwanted changes to beam image due to
# extra dependencies needed to pull licenses.
####

FROM beam as third_party_licenses
ARG pull_licenses
COPY target/license_scripts /tmp/license_scripts/

# Add golang licenses.
COPY  target/go-licenses/* /opt/apache/beam/third_party_licenses/golang/

COPY target/license_scripts /tmp/license_scripts/
RUN if [ "$pull_licenses" = "true" ] ; then \
      pip install 'pip-licenses<5' pyyaml tenacity && \
      python /tmp/license_scripts/pull_licenses_py.py ; \
    fi

FROM beam
ARG pull_licenses
COPY --from=third_party_licenses /opt/apache/beam/third_party_licenses /opt/apache/beam/third_party_licenses
RUN if [ "$pull_licenses" != "true" ] ; then \
      rm -rf /opt/apache/beam/third_party_licenses ; \
    fi
