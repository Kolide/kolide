# Kolide [![CircleCI](https://circleci.com/gh/kolide/kolide-ose.svg?style=svg&circle-token=2573c239b7f18967040d2dec95ca5f71cfc90693)](https://circleci.com/gh/kolide/kolide-ose)

## Building

To build the code ensure you have `node` and `npm` installed run the
following from the root of the repository:

```
make
```

This will produce a binary called `kolide` in the root of the repo.

## Testing

To run the application's tests, run the following from the root of the
repository:

```
go vet ./... && go test -v ./...
```

## Development Environment

To setup a working local development environment run perform the following tasks:

* Install the following dependencies:
* [Docker & docker-compose](https://www.docker.com/products/overview#/install_the_platform)
* [go 1.6.x](https://golang.org/dl/)
* [nodejs 0.6.x](https://nodejs.org/en/download/current/) (and npm)
* A GNU compatible version of `make`

Once those tools are installed, to set up a canonical development environment 
via docker, run the following from the root of the repository:

```
docker-compose up
```

This requires that you have docker installed. At this point in time,
automatic configuration tools are not included with this project.

Once you `docker-compose up` and are running the databases, you can build
the code and run the following command to create the database tables:

```
kolide prepare-db
```

To install all JavScript dependencies, run 

```
npm install
```

Now you are prepared to run a Kolide development environment. Run the following:

```
make serve
```

You may have to edit the example configuration file to reflect `localhost` if
you're using Docker via a native docker engine or the output of 
`docker-machine ip` if you're using Docker via `docker-toolbox`.

If you'd like to shut down the virtual infrastructure created by docker, run
the following from the root of the repository:
1. Start up all external servers with `docker-compose up`

By default, the last command will run the development proxy on
`http://localhost:8081` which allows you to make live changes to the code and
have them hot-reload.


## Docker Deployment
This repository comes with a simple Dockerfile. You can use this to easily
deploy Kolide in any infrastructure context that can consume a docker image
(heroku, kubernetes, rancher, etc).

To build the image locally, run:

```
docker build --rm -t kolide .
```

To run the image locally, simply run:
>>>>>>> Improve README to reflect new dev workflow

```
docker-compose down
```
