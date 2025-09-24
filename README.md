# Building a Dynatrace OneAgent image for side-car container integrations
Dynatrace provides OneAgent code-module images on public registries (e.g. https://gallery.ecr.aws/dynatrace/dynatrace-codemodules) to be used by the Dynatrace K8s Operator with cloud-native injection. 

This image can also be used without the Dynatrace Operator for container services that support [init-container](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/).
The dynatrace-codemodules image contains the OneAgent artefacts and a small cli - the [dynatrace-bootstrapper](https://github.com/Dynatrace/dynatrace-bootstrapper) that is executed at startup to copy the artefacts into a volume that is shared with the application container. 

A tutorial integrating Dynatrace OneAgent into [Azure Container Apps](https://azure.microsoft.com/en-us/products/container-apps) using the init-container approach can be found here: https://github.com/dtPaTh/cloud-service-integrations/blob/main/azure-container-apps.md

### Side-Cars
Not every serverless container service supports init-containers, but instead side-cars. Side-cars may come with additional requirements or limitations e.g. such as startup probes or requiring to keep the side-car up and running during the service lifecycle.

For this the Dynatrace code-module image needs to be enhanced to fit the requirements of a side-car integration. 

This repostory includes an additional [cli](#Serverless-Boostrap) to be used in combination with the [dynatrace-bootstrapper](https://github.com/Dynatrace/dynatrace-bootstrapper) to enhance the required functionality

## Serverless-Boostrap CLI

```serverless-boostrap [--keepalive] [--healthprobe] [<command-to-execute> ...]```

#### --keepalive
After all commands have been executed, pauses the cli to prevent the side-car to be terminated. 

#### --healthprobe
Enables the following healthprobe endpoints on **port 8080**
* **/startup** Startup probe indicates that the cli has started.
* **/readiness** Readiness probe, indicating readiness to as soon as the ```<command-to-execute> ...``` has finished.
* **/liveness** Liveness probe, indicating the cli is running

E.g. Google CloudRun requires to configure a startup probe for side-cars: https://cloud.google.com/run/docs/deploying#sidecars

Example command adding a health-probe endpoint to the sidecar while the artefacts are copied to the shared volume: 
```
serverless-boostrap --healthprobe /opt/dynatrace/oneagent/agent/lib64/dynatrace-bootstrapper --source=/opt/dynatrace/oneagent --target=/shared/dynatrace/oneagent
```

## Building the image

The Dockerfile ```Dockerfile.native``` builds the cli and adds it into the Dynatrace code-modules image as a base. 

The Dockerfile can be configured using container build arguments:
* ```DT_BASEIMG``` Defines the source image containing the Dynatrace codemodules. 

The following docker build command creates the container image directly from the github repository 
using a specific immutable code-module image from the Dynatrace public container registry on ECR.  

To simplify updating patch versions of the dynatrace code-modules, a rolling tag using a ```major.minor``` versioning scheme is set to the built image. 
When a new version is available, one can build the new image and with the next container restart, the new (patch)-version is automatically applied.   

```
docker build -f Dockerfile.native https://github.com/dtPaTh/dt-codemodule-images.git#serverless-boostrapper --build-arg DT_BASEIMG=public.ecr.aws/dynatrace/dynatrace-codemodules:1.321.51.20250905-075429 -t oneagent-codemodules:1.321
```

## Example tutorial to test with docker-compose

To test using docker-compose, healthprobes require to have e.g. curl available in the side-car. 
This repository contains an alternative Dockerfile ```Dockerfile.test``` for this purpose.

### Step 1: Build the image

```
docker build -f Dockerfile.test https://github.com/dtPaTh/dt-codemodule-images.git#serverless-boostrapper --build-arg DT_BASEIMG=public.ecr.aws/dynatrace/dynatrace-codemodules:1.321.51.20250905-075429 -t oneagent-codemodules:1.321-test
```

### Step 2: Create the docker-compose file
The following compose file, starts the asp.net sample app, integrating Dynatrace OneAgent as a side-car.
The side-car is setup as a dependency for the application container, whereas the application container is started when the side-car is ready (after it has copied all artefacts). 

``` 
services:
  dtsidecar:
    image: localhost/oneagent-codemodules:${DT_IMAGE_TAG}
    healthcheck:
      test: "curl -f http://localhost:8080/readiness"
      interval: 5s
      timeout: 5s
      retries: 10
    user: 0:0
    volumes:
      - sharedvolume:/shared/
  appcontainer:
    image: "mcr.microsoft.com/dotnet/samples:aspnetapp"
    depends_on:
      dtsidecar: 
        condition: service_healthy 
    ports:
      - "8080"
    volumes:
      - sharedvolume:/shared/
    environment:
      - LD_PRELOAD=/shared/dynatrace/oneagent/agent/lib64/liboneagentproc.so
      - DT_AGENTACTIVE=${DT_AGENTACTIVE}   
      - DT_LOGSTREAM=${DT_LOGSTREAM}   
      - DT_LOGLEVELCON=${DT_LOGLEVELCON}   
      - DT_TENANT=${DT_TENANT}
      - DT_TENANTTOKEN=${DT_TENANTTOKEN}
      - DT_CONNECTION_POINT=${DT_CONNECTION_POINT}
volumes:
  sharedvolume:
    driver:
      local
```

### Step 3: Create a .env file for configurations
Provide the necessary configuration via runtime environment variables as defined in an .env file
``` 
DT_IMAGE_TAG=1.321-test

DT_TENANT=<YOUR-TENANT-ID>
DT_TENANTTOKEN=<YOUR-TENANT-TOKEN> 
DT_CONNECTION_POINT=<YOUR-CONNECTION-ENDPOINT>

DT_AGENTACTIVE=true
DT_LOGSTREAM=stdout 
DT_LOGLEVELCON=info
```

### Step 4: Run the project
```
docker-compose --env-file .env up
```




