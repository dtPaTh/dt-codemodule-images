# Building a Dynatrace OneAgent code-modules images for side-car integrations
Dynatrace provides a built-in image registry for each tenant, as well as public registries with OneAgent images. 

These images can be used to integrate images into application images using e.g. the docker file command "COPY". 

If you want a runtime integration using an initcontainer or side-car pattern, you want to copy the OneAgent artefacts during container/pod startup into a mounted volume, from where the main container can enable the Oneagent. 
To rely on the same code-module images as when integrating directly into the application image, instead of using the Dynatraec REST API to download necessary artefacts, one needs a image with the code-modules that copy the necessary artefacts at container startup. 

The Dockerfile ```Dockerfile.native``` creates such a container image, inheriting the base image and adding a native binary to copy the artefacts at startup. Due to the native binary, there is no additional OS dependency needed as it would be relying on a shell script to copy the artefacts. 
The Dockerfile can be configured using buld arguments:
* ```DT_BASEIMG``` Defines the source/baseimage for the dynatrace codemodules. 

