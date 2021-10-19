// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package environment

import (
	pb "beam.apache.org/playground/backend/internal/api"
	"github.com/tkanos/gonfig"
	"google.golang.org/grpc/grpclog"
	"os"
	"strconv"
	"strings"
)

const (
	serverIpKey          = "SERVER_IP"
	serverPortKey        = "SERVER_PORT"
	beamSdkKey           = "BEAM_SDK"
	beamRunner           = "BEAM_RUNNER"
	jdk                  = "JDK"
	defaultIp            = "localhost"
	defaultPort          = 8080
	defaultSdk           = pb.Sdk_SDK_JAVA
	executorConfigFolder = "CONFIG_FOLDER"
	defaultConfigFolder  = "./"
	defaultBeamSdk       = "/opt/apache/beam/jars/beam-sdks-java-harness.jar"
	defaultBeamRunner    = "/opt/apache/beam/jars/beam-runners-direct.jar"
	defaultJDK           = "/opt/apache/beam/jars/slf4j-jdk14.jar"
	jsonExt              = ".json"
)

// Environment operates with environment structures: ServerEnvs, LogWriters, BeamEnvs
type Environment struct {
	ServerEnvs  ServerEnvs
	BeamSdkEnvs BeamEnvs
}

// NewEnvironment is a constructor for Environment.
// Default values:
// LogWriters: by default using os.Stdout
// ServerEnvs: by default using defaultIp and defaultPort from constants
// BeamEnvs: by default using pb.Sdk_SDK_JAVA
func NewEnvironment() *Environment {
	svc := Environment{}
	svc.ServerEnvs = *getServerEnvsFromOsEnvs()
	svc.BeamSdkEnvs = *getSdkEnvsFromOsEnvs()

	return &svc
}

// getServerEnvsFromOsEnvs lookups in os environment variables and takes value for ip and port. If not exists - using default
func getServerEnvsFromOsEnvs() *ServerEnvs {
	port := defaultPort
	ip := getEnv(serverIpKey, defaultIp)
	if value, present := os.LookupEnv(serverPortKey); present {
		if converted, err := strconv.Atoi(value); err == nil {
			port = converted
		} else {
			grpclog.Errorf("couldn't convert provided port. Using default %s\n", defaultPort)
		}
	}
	return NewServerEnvs(ip, port)
}

// getServerEnvsFromOsEnvs lookups in os environment variables and takes value for Apache Beam SDK. If not exists - using default
func getSdkEnvsFromOsEnvs() *BeamEnvs {
	sdk := pb.Sdk_SDK_UNSPECIFIED
	if value, present := os.LookupEnv(beamSdkKey); present {

		switch value {
		case pb.Sdk_SDK_JAVA.String():
			sdk = pb.Sdk_SDK_JAVA
		case pb.Sdk_SDK_GO.String():
			sdk = pb.Sdk_SDK_GO
		case pb.Sdk_SDK_PYTHON.String():
			sdk = pb.Sdk_SDK_PYTHON
		case pb.Sdk_SDK_SCIO.String():
			sdk = pb.Sdk_SDK_SCIO
		}
	}
	if sdk == pb.Sdk_SDK_UNSPECIFIED {
		grpclog.Infof("couldn't get sdk from %s os env, using default: %s", beamSdkKey, defaultSdk)
		sdk = defaultSdk
	}
	executorConfig := readExecutorConfig(sdk)
	return NewBeamEnvs(sdk, executorConfig)
}

func readExecutorConfig(apacheBeamSdk pb.Sdk) *ExecutorConfig {
	executorConfig, _ := getConfigFromJson(apacheBeamSdk)
	switch apacheBeamSdk {
	case pb.Sdk_SDK_JAVA:
		executorConfig.CompileArgs = append(executorConfig.CompileArgs,
			getEnv(beamSdkKey, defaultBeamSdk))
		jars := strings.Join([]string{
			getEnv(beamSdkKey, defaultBeamSdk),
			getEnv(beamRunner, defaultBeamRunner),
			getEnv(jdk, defaultJDK),
		}, ":")
		executorConfig.RunArgs[1] += jars
	case pb.Sdk_SDK_GO:
	case pb.Sdk_SDK_PYTHON:
	case pb.Sdk_SDK_SCIO:
	}

	return &executorConfig
}

func getConfigFromJson(apacheBeamSdk pb.Sdk) (ExecutorConfig, error) {
	folder := getEnv(executorConfigFolder, defaultConfigFolder)
	CmdConfigPath := folder + apacheBeamSdk.String() + jsonExt
	executorConfig := ExecutorConfig{}
	err := gonfig.GetConf(CmdConfigPath, &executorConfig)
	return executorConfig, err
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
