# JC Cloud Operations Project

## Purpose
The purpose of this project is to build a simple http service that accepts a JSON payload and depending on the field and values passed - it would perform an action

## Prerequisites
This application requires a machine running docker

## Deployment
* Copy the files or clone down the branch to the working directory
* Execute the Docker Build
* `docker build . -t jc_project`
* After build, run the image.  App is listening on port 8080 - do not change this
* `docker run -p "80:8080" jc_project`

## Payload
```
{
    "action": "<action>"
}
```
Available actions:
* *download* - Downloads the requested text file for the project
* *read* - Downloads the requested text file for the project and responds with its contents

## Example
Execute a curl command as seen below:
`curl -X POST http://localhost/manage_file -H 'Content-Type: application/json' -d '{"action": "download"}'`

Note: Substitute the hostname, ip, or port to match your environment 

