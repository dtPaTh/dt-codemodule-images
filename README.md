# Building a Dynatrace OneAgent image for side-car container integrations
Dynatrace provides a built-in image registry for each tenant, as well as public registries with OneAgent images. 

These images are used by e.g. the Dynatrace K8s Operator or when integrating Dynatrace OneAgent into application images for serverless containers at image build-time.  

If you prefer a runtime integration using an initcontainer or side-car pattern, you need to copy the OneAgent artefacts during container/pod startup into a mounted volume, from where the application container can enable the Oneagent via environment variables. The benefits of such an integration without the need to modify the application image.

Instead of using a init-container that downloads the necessary artefacts via e.g. a shellscript, this project aims to build **a distroless image including the OneAgent which copies the required artefacts at container startup**, resulting in  **lower startup latency** and **higher reliability due reduced risk of network issues when downloading the agent**, align with container- and security-best practices enabling **easy integration with container pipelines** and **validation and security scans**. 

The Dockerfile ```Dockerfile.native``` creates such a container image. It combines the Dynatrace OneAgent artefacts and a statically linked binary for copying. 

The Dockerfile can be configured using container build arguments:
* ```DT_BASEIMG``` Defines the source image containing the Dynatrace codemodules. 

## Example tutorial using docker-compose
### Step 1: Build the Dynatrace image 
The following docker build command creates the container image directly from the github repository 
using a specific immutable code-module image from the Dynatrace public container registry on ECR.  

To simplify updating patch versions of the dynatrace code-modules, a rolling tag using a ```major.minor``` versioning scheme is set to the built image. 
When a new version is available, one can build the new image and with the next container restart, the new (patch)-version is automatically applied.   

```
docker build -f Dockerfile.native https://github.com/dtPaTh/dt-codemodule-images.git --build-arg DT_BASEIMG=public.ecr.aws/dynatrace/dynatrace-codemodules:1.301.54.20241017-161011 -t oneagent-codemodules:1.301
```

### Step 2: Create the docker-compose file
The following docker-compose file, starts a nginx container, integrating Dynatrace OneAgent as a side-car.
``` 
services:
  dtsidecar:
    image: oneagent-codemodules:${DT_IMAGE_TAG}
    volumes:
      - sharedvolume:/home/
  appcontainer:
    image: "nginx:latest"
    depends_on:
      dtsidecar: 
        condition: service_completed_successfully
    ports:
      - "80"
      - "443"
    volumes:
      - sharedvolume:/home/
    environment:
      - LD_PRELOAD=/home/dynatrace/oneagent/agent/lib64/liboneagentproc.so
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
To control agent options, we use environment variables via a .env file
``` 
DT_IMAGE_TAG=1.301
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

## A sample integration for Azure App Service for Linux
In 2024, we've [partnered with the Azure App Service for Linux team](https://azure.github.io/AppService/2024/11/08/Global-Availability-Sidecars.html) to provide a new observability [integration for Dynatrace using the sidecar pattern with container apps](https://azure.github.io/AppService/2024/07/26/Using-Dynatrace-with-Sidecar.html). 

Prior using the tutorial to integrate Dynatrace you will need to create the Dynatrace image to be used as a side-car.  Please see the following example of creating the image using a Azure Container Registry task:  

### Create the Dynatrace image using a Azure Container Registry task. 
#### Step 1 - Configure build parameters
```
$GITHUB_REPO_URL = "https://github.com/dtPaTh/dt-codemodule-images.git"  
$DOCKER_FILE = "Dockerfile.native"

$ACR_NAME = "<YOUR-AZURECONTAINERREGISTRY-NAME>" 
$IMAGE_NAME ="code-modules"
$IMAGE_TAG = "1.301"  
$TASK_NAME = "bootstrapped_codemodules_from_ecr" 

$BUILD_ARG_BASEIMG = "DT_BASEIMG=public.ecr.aws/dynatrace/dynatrace-codemodules:1.301.54.20241017-161011" 
```

#### Step 2 - Create the ACR task using the Azure cli
```
az acr task create --registry $ACR_NAME --name $TASK_NAME --image "$($ACR_NAME).azurecr.io/$($IMAGE_NAME):$($IMAGE_TAG)" --context $GITHUB_REPO_URL --file $DOCKER_FILE --arg $BUILD_ARG_BASEIMG --base-image-trigger-enabled false --commit-trigger-enabled false --pull-request-trigger-enabled false
```

#### Step 3 - Run the task
```
az acr task run  --registry $ACR_NAME --name $TASK_NAME
```

Now you can reference the image using [AppService for Linux sidecar integration for container apps](https://azure.github.io/AppService/2024/07/26/Using-Dynatrace-with-Sidecar.html). If you are using code-based apps, you have to use an ARM template as seen in this [example](https://github.com/Azure-Samples/sidecar-samples/tree/main/sidecar-arm-template).




