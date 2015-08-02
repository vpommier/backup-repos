## Synopsis
Tool for backup Github Repositories.

## Build
To build the docker container from sources:

```bash

make build

```

## Run
Get usage:

```bash

docker run --rm vpommier/backup-repos 

```
Example to backup Github repositories of the user octocat:

```bash

docker run --rm -v `pwd`/backup-repos/:/var/backup-repos/archives/ vpommier/backup-repos octocat 

```

## Using with Fleet
These are examples only.

```bash

fleetctl submit fleet/*
fleetctl load backup-repos@octocat.service
fleetctl start backup-repos@octocat.timer

```

## Contributors
Vincent POMMIER
