#!/bin/bash
docker build -t composer22/ringoexp_build .
docker run -v /var/run/docker.sock:/var/run/docker.sock -v $(which docker):$(which docker) -ti --name ringoexp_build composer22/ringoexp_build
docker rm ringoexp_build
docker rmi composer22/ringoexp_build
