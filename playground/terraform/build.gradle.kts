/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * License); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an AS IS BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import com.pswidersk.gradle.terraform.TerraformTask
import java.io.ByteArrayOutputStream

description = "Apache Beam :: Playground :: Deploy"
val licenseText = "################################################################################\n" +
        "#  Licensed to the Apache Software Foundation (ASF) under one\n" +
        "#  or more contributor license agreements.  See the NOTICE file\n" +
        "#  distributed with this work for additional information\n" +
        "#  regarding copyright ownership.  The ASF licenses this file\n" +
        "#  to you under the Apache License, Version 2.0 (the\n" +
        "#  \"License\"); you may not use this file except in compliance\n" +
        "#  with the License.  You may obtain a copy of the License at\n" +
        "#\n" +
        "#      http://www.apache.org/licenses/LICENSE-2.0\n" +
        "#\n" +
        "#  Unless required by applicable law or agreed to in writing, software\n" +
        "#  distributed under the License is distributed on an \"AS IS\" BASIS,\n" +
        "#  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n" +
        "#  See the License for the specific language governing permissions and\n" +
        "# limitations under the License.\n" +
        "################################################################################"

plugins {
    id("com.pswidersk.terraform-plugin") version "1.0.0"
    id("org.unbroken-dome.helm") version "1.7.0"
    id("org.unbroken-dome.helm-releases") version "1.7.0"
}

terraformPlugin {
    terraformVersion.set("1.0.9")
}

tasks {
    /* init Infrastructure for migrate */
    register<TerraformTask>("terraformInit") {
        // exec args can be passed by commandline, for example
        var environment = "unknown"
        if (project.hasProperty("project_environment")) {
            environment = project.property("project_environment") as String
        }
        args(
            "init", "-migrate-state",
            "-backend-config=./environment/$environment/state.tfbackend",
            "-var=environment=$environment",
            if (file("./environment/$environment/terraform.tfvars").exists()) {
                "-var-file=./environment/$environment/terraform.tfvars"
            } else {
                "-no-color"
            }
        )
    }
    /* refresh Infrastucture for remote state */
    register<TerraformTask>("terraformRef") {
        mustRunAfter(":playground:terraform:terraformInit")
        var project_id = "unknown"
        var environment = "unknown"
        if (project.hasProperty("project_id")) {
            project_id = project.property("project_id") as String
        }
        if (project.hasProperty("project_environment")) {
            environment = project.property("project_environment") as String
        }
        args(
            "refresh",
            "-lock=false",
            "-var=project_id=$project_id",
            "-var=environment=$environment",
            if (file("./environment/$environment/terraform.tfvars").exists()) {
                "-var-file=./environment/$environment/terraform.tfvars"
            } else {
                "-no-color"
            }
        )
    }

    /* deploy all App */
    register<TerraformTask>("terraformApplyApp") {
        var project_id = "unknown"
        var environment = "unknown"
        if (project.hasProperty("project_id")) {
            project_id = project.property("project_id") as String
        }
        if (project.hasProperty("project_environment")) {
            environment = project.property("project_environment") as String
        }
        var docker_tag = if (project.hasProperty("docker-tag")) {
            project.property("docker-tag") as String
        } else {
            environment
        }
        args(
            "apply",
            "-auto-approve",
            "-lock=false",
            "-target=module.applications",
            "-var=project_id=$project_id",
            "-var=environment=$environment",
            "-var=docker_image_tag=$docker_tag",
            if (file("./environment/$environment/terraform.tfvars").exists()) {
                "-var-file=./environment/$environment/terraform.tfvars"
            } else {
                "-no-color"
            }
        )
    }

    /* deploy  App - Only all services for  backend */
    register<TerraformTask>("terraformApplyAppBack") {
        var environment = "unknown"
        if (project.hasProperty("project_environment")) {
            environment = project.property("project_environment") as String
        }
        args(
            "apply",
            "-auto-approve",
            "-lock=false",
            "-target=module.applications.module.backend",
            "-var=environment=$environment",
            if (file("./environment/$environment/terraform.tfvars").exists()) {
                "-var-file=./environment/$environment/terraform.tfvars"
            } else {
                "-no-color"
            }
        )
    }

    /* deploy  App - Only services for frontend */
    register<TerraformTask>("terraformApplyAppFront") {
        var project_id = "unknown"
        var environment = "unknown"
        if (project.hasProperty("project_id")) {
            project_id = project.property("project_id") as String
        }
        if (project.hasProperty("project_environment")) {
            environment = project.property("project_environment") as String
        }
        var docker_tag = if (project.hasProperty("docker-tag")) {
            project.property("docker-tag") as String
        } else {
            environment
        }
        args(
            "apply",
            "-auto-approve",
            "-lock=false",
            "-target=module.applications.module.frontend",
            "-var=project_id=$project_id",
            "-var=environment=$environment",
            "-var=docker_image_tag=$docker_tag",
            if (file("./environment/$environment/terraform.tfvars").exists()) {
                "-var-file=./environment/$environment/terraform.tfvars"
            } else {
                "-no-color"
            }
        )
    }

    /* build only Infrastructurte */
    register<TerraformTask>("terraformApplyInf") {
        var environment = "unknown"
        if (project.hasProperty("project_environment")) {
            environment = project.property("project_environment") as String
        }
        args(
            "apply",
            "-auto-approve",
            "-lock=false",
            "-target=module.infrastructure",
            "-var=environment=$environment",
            if (file("./environment/$environment/terraform.tfvars").exists()) {
                "-var-file=./environment/$environment/terraform.tfvars"
            } else {
                "-no-color"
            }
        )
    }

    /* build All */
    register<TerraformTask>("terraformApply") {
        var project_id = "unknown"
        var environment = "unknown"
        if (project.hasProperty("project_id")) {
            project_id = project.property("project_id") as String
        }
        if (project.hasProperty("project_environment")) {
        environment = project.property("project_environment") as String
        }
        var docker_tag = if (project.hasProperty("docker-tag")) {
            project.property("docker-tag") as String
        } else {
            environment
        }
        args(
            "apply",
            "-auto-approve",
            "-lock=false",
            "-var=project_id=$project_id",
            "-var=environment=$environment",
            "-var=docker_image_tag=$docker_tag",
            if (file("./environment/$environment/terraform.tfvars").exists()) {
                "-var-file=./environment/$environment/terraform.tfvars"
            } else {
                "-no-color"
            }
        )
    }

    register<TerraformTask>("terraformDestroy") {
        var project_id = "unknown"
        var environment = "unknown"
        if (project.hasProperty("project_id")) {
            project_id = project.property("project_id") as String
        }
        if (project.hasProperty("project_environment")) {
            environment = project.property("project_environment") as String
        }
        args(
            "destroy",
            "-auto-approve",
            "-lock=false",
            "-var=project_id=$project_id",
            "-var=environment=$environment",
            if (file("./environment/$environment/terraform.tfvars").exists()) {
                "-var-file=./environment/$environment/terraform.tfvars"
            } else {
                "-no-color"
            }
        )
    }
}

/* set Docker Registry to params from Inf */
task("setDockerRegistry") {
    group = "deploy"
    //get Docker Registry
    dependsOn(":playground:terraform:terraformInit")
    dependsOn(":playground:terraform:terraformRef")
    try {
        var stdout = ByteArrayOutputStream()
        //set Docker Registry
        exec {
            commandLine = listOf("terraform", "output", "docker-repository-root")
            standardOutput = stdout
        }
        project.rootProject.extra["docker-repository-root"] = stdout.toString().trim().replace("\"", "")
    } catch (e: Exception) {
    }
}

task("readState") {
    group = "deploy"
    dependsOn(":playground:terraform:terraformInit")
    dependsOn(":playground:terraform:terraformRef")
}

task("pushBack") {
    group = "deploy"
    dependsOn(":playground:backend:containers:go:dockerTagsPush")
    dependsOn(":playground:backend:containers:java:dockerTagsPush")
    dependsOn(":playground:backend:containers:python:dockerTagsPush")
    dependsOn(":playground:backend:containers:scio:dockerTagsPush")
    dependsOn(":playground:backend:containers:router:dockerTagsPush")
}

task("pushFront") {
    group = "deploy"
    dependsOn(":playground:frontend:dockerTagsPush")
}

task("prepareConfig") {
    group = "deploy"
    doLast {
        var extip = ""
        var stdout = ByteArrayOutputStream()
        exec {
            commandLine = listOf("terraform", "output", "playground_static_ip_address")
            standardOutput = stdout
        }
        extip = stdout.toString().trim().replace("\"", "")
        stdout = ByteArrayOutputStream()
        val configFileName = "config.g.dart"
        val modulePath = project(":playground:frontend").projectDir.absolutePath
        var file = File("$modulePath/lib/$configFileName")

        file.writeText(
            """${licenseText}
const String kAnalyticsUA = 'UA-73650088-2';
const String kApiClientURL =
      'https://router.${extip}.nip.io';
const String kApiJavaClientURL =
      'https://java.${extip}.nip.io';
const String kApiGoClientURL =
      'https://go.${extip}.nip.io';
const String kApiPythonClientURL =
      'https://python.${extip}.nip.io';
const String kApiScioClientURL =
      'https://scio.${extip}.nip.io';
"""
        )
        try {
            var stdout = ByteArrayOutputStream()
            //set Docker Registry
            exec {
                commandLine = listOf("terraform", "output", "docker-repository-root")
                standardOutput = stdout
            }
            project.rootProject.extra["docker-repository-root"] = stdout.toString().trim().replace("\"", "")
        } catch (e: Exception) {
        }
    }
}
/* initialization infrastructure */
task("InitInfrastructure") {
    group = "deploy"
    description = "initialization infrastructure"
    val init = tasks.getByName("terraformInit")
    val apply = tasks.getByName("terraformApplyInf")
    dependsOn(init)
    dependsOn(apply)
    apply.mustRunAfter(init)
}

/* build, push, deploy Frontend app */
task("deployFrontend") {
    group = "deploy"
    description = "deploy Frontend app"
    val read = tasks.getByName("readState")
    val push = tasks.getByName("pushFront")
    val deploy = tasks.getByName("terraformApplyAppFront")
    dependsOn(read)
    Thread.sleep(10)
    push.mustRunAfter(read)
    deploy.mustRunAfter(push)
    dependsOn(push)
    dependsOn(deploy)
}

/* build, push, deploy Backend app */
task("deployBackend") {
    group = "deploy"
    description = "deploy Backend app"
    //TODO please add default tag from project_environment property
    //if !(project.hasProperty("docker-tag")) {
    //    project.extra.set("docker-tag", project.property("project_environment") as String)
    //}
    val config = tasks.getByName("setDockerRegistry")
    val push = tasks.getByName("pushBack")
    val deploy = tasks.getByName("terraformApplyAppBack")
    dependsOn(config)
    Thread.sleep(10)
    push.mustRunAfter(config)
    deploy.mustRunAfter(push)
    dependsOn(push)
    dependsOn(deploy)
}

task("takeConfig") {
  group = "deploy"
  doLast {
   var ipaddr = ""
   var redis = ""
   var project = ""
   var registry = ""
   var ipaddrname = ""
   var stdout = ByteArrayOutputStream()
   var environment = "unknown"
   var docker_tag = if (project.hasProperty("docker-tag")) {
       project.property("docker-tag") as String
   }   else {
         environment
   }
   exec {
       commandLine = listOf("terraform", "output", "playground_static_ip_address")
       standardOutput = stdout
   }
   ipaddr = stdout.toString().trim().replace("\"", "")
   stdout = ByteArrayOutputStream()

   exec {
       commandLine = listOf("terraform", "output", "playground_redis_ip")
       standardOutput = stdout
   }
   redis = stdout.toString().trim().replace("\"", "")
   stdout = ByteArrayOutputStream()
   exec {
       commandLine = listOf("terraform", "output", "playground_gke_project")
       standardOutput = stdout
   }
   project = stdout.toString().trim().replace("\"", "")
   stdout = ByteArrayOutputStream()
   exec {
       commandLine = listOf("terraform", "output", "docker-repository-root")
       standardOutput = stdout
   }
   registry = stdout.toString().trim().replace("\"", "")
   stdout = ByteArrayOutputStream()
   exec {
       commandLine = listOf("terraform", "output", "playground_static_ip_address_name")
       standardOutput = stdout
   }
   ipaddrname = stdout.toString().trim().replace("\"", "")
   stdout = ByteArrayOutputStream()

   val configFileName = "values.yaml"
   val modulePath = project(":playground").projectDir.absolutePath
   var file = File("$modulePath/infrastructure/helm-backend/$configFileName")
   file.appendText("""
static_ip: ${ipaddr}
redis_ip: ${redis}:6379
project_id: ${project}
registry: ${registry}
static_ip_name: ${ipaddrname}
tag: ${docker_tag}
    """)
 }
}
helm {
    val backend by charts.creating {
        chartName.set("backend")
        sourceDir.set(file("../infrastructure/helm-backend"))
    }
    releases {
        create("backend") {
            from(backend)
        }
    }
}
task ("gkebackend") {
  group = "deploy"
  val init = tasks.getByName("terraformInit")
  val apply = tasks.getByName("terraformApplyInf")
  val takeConfig = tasks.getByName("takeConfig")
  val push = tasks.getByName("pushBack")
  val helm = tasks.getByName("helmInstallBackend")
  val prepare = tasks.getByName("prepareConfig")
  dependsOn(init)
  dependsOn(apply)
  dependsOn(takeConfig)
  dependsOn(push)
  dependsOn(helm)
  dependsOn(prepare)
  apply.mustRunAfter(init)
  takeConfig.mustRunAfter(apply)
  push.mustRunAfter(takeConfig)
  helm.mustRunAfter(push)
  prepare.mustRunAfter(helm)
}
