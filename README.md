# Building a Dynatrace OneAgent image for side-car container integrations
Dynatrace provides OneAgent code-module images on public registries (e.g. https://gallery.ecr.aws/dynatrace/dynatrace-codemodules) to be used by the Dynatrace K8s Operator with cloud-native injection. 

This image can also be used without the Dynatrace Operator for container services that support [init-container](https://kubernetes.io/docs/concepts/workloads/pods/init-containers/).
The dynatrace-codemodules image contains the OneAgent artefacts and a small cli - the [dynatrace-bootstrapper](https://github.com/Dynatrace/dynatrace-bootstrapper) that is executed at startup to copy the artefacts into a volume that is shared with the application container. 

A tutorial integrating Dynatrace OneAgent into [Azure Container Apps](https://azure.microsoft.com/en-us/products/container-apps) using the init-container approach can be found here: https://github.com/dtPaTh/cloud-service-integrations/blob/main/azure-container-apps.md

### Side-Cars
Not every serverless container service supports init-containers, but instead side-cars. Side-cars may come with additional requirements or limitations e.g. such as startup probes, requring to keep the side-car up and running along with the applicaiton container, ..

For this the Dynatrace code-module image needs to be enhanced to fit the requirements of a side-car integration. 

This repostory includes an additional [cli](#Serverless-Boostrap) to be used in combination with the [dynatrace-bootstrapper](https://github.com/Dynatrace/dynatrace-bootstrapper) to enhance the required functionality

## Serverless-Boostrap CLI

```serverless-boostrap [--keepalive] [--healthprobe] [<command-to-execute> ...]```

#### --keepalive
After all commands have been executed, pauses the cli to prevent the side-car to be terminated. 

#### --healthprobe
Enables a health-probe endpoint - **:8080/health**. 

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

## Example tutorial using docker-compose
### Step 1: Create the docker-compose file
The following docker-compose file, starts the asp.net sample app, integrating Dynatrace OneAgent as a side-car.

``` 
services:
  dtsidecar:
    image: oneagent-codemodules:${DT_IMAGE_TAG}
    volumes:
      - sharedvolume:/shared/
  appcontainer:
    image: "mcr.microsoft.com/dotnet/samples:aspnetapp"
    ports:
      - "80"
      - "443"
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

### Step 2: Create a .env file for configurations
To control agent options, we use environment variables via a .env file
``` 
DT_IMAGE_TAG=1.321
DT_TENANT=<YOUR-TENANT-ID>
DT_TENANTTOKEN=<YOUR-TENANT-TOKEN> 
DT_CONNECTION_POINT=<YOUR-CONNECTION-ENDPOINT>

DT_AGENTACTIVE=true
DT_LOGSTREAM=stdout 
DT_LOGLEVELCON=info
```

### Step 3: Run the project
```
docker-compose --env-file .env up
```




