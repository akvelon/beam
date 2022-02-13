#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

variable "project_id" {
  description = "The GCP Project ID where Playground Applications will be created"
}

#GCS

variable "examples_bucket_name" {
  description = "Name of Bucket to Store Playground Examples"
  default     = "playground-examples"
}

variable "examples_bucket_location" {
  description = "Location of Playground Examples Bucket"
  default     = "US"
}

variable "examples_storage_class" {
  description = "Examples Bucket Storage Class"
  default     = "STANDARD"
}

# Artifact Registry

variable "repository_id" {
  description = "ID of Artifact Registry"
  default     = "playground-repository"
}

variable "repository_location" {
  description = "Location of Artifact Registry"
  default     = "us-central1"
}

#Redis

variable "redis_version" {
  description = "The GCP Project ID where Playground Applications will be created"
  default     = "REDIS_6_X"
}

variable "terraform_state_bucket_name" {
  description = "Bucket name for terraform state"
  default     = "beam_playground_terraform"
}

variable "redis_region" {
  description = "Region of Redis"
  default     = "us-central1"
}

variable "redis_name" {
  description = "Name of Redis"
  default     = "playground-backend-cache"
}

variable "redis_tier" {
  description = "Tier of Redis"
  default     = "STANDARD_HA"
}

variable "redis_replica_count" {
  description = "Redis's replica count"
  default     = 1
}

variable "redis_memory_size_gb" {
  description = "Size of Redis memory ,  if set 'read replica' it must be from 5GB to 100GB."
  default     = 5
}

#VPC

variable "vpc_name" {
  description = "Name of VPC to be created"
  default     = "playground-vpc"
}

variable "create_subnets" {
  description = "Auto Create Subnets Inside VPC"
  default     = true
}

variable "mtu" {
  description = "MTU Inside VPC"
  default     = 1460
}

# GKE

variable "gke_machine_type" {
  description = "Node pool machine types"
  default     = "e2-standard-4"
}

variable "gke_node_count" {
  description = "Node pool size"
  default     = 1
}

variable "gke_name" {
  description = "Name of GKE cluster"
  default     = "playground-examples"
}

variable "gke_location" {
  description = "Location of GKE cluster"
  default     = "us-central1-a"
}

variable "service_account" {
  description = "Service account email (id) for example service-account-playground@friendly-tower-340607.iam.gserviceaccount.com"
  default     = "service-account-playground@friendly-tower-340607.iam.gserviceaccount.com"
}

# Over

variable "environment" {
  description = "prod,dev"
}
