# helm-microservice

This helm chart deploys a deployement, with a series of services and ingresses for a given container, hosting a web service with the following properties:

- Configuration from environment variables
- exposes a number of ports, dependent on the container

The chart exposes port 8000 over http by default, with a load balancing service, but the ports and services are configured from a list of arbitrary length in the values file.

Similarly, the environment variables are injected from a definition list, of arbitrary length.

The goal is to provide a DRY miicroservice deployment mechanism with some flexibility, to be used with some other charts orchestrating a series of microservices via a requirement file, pointing to this chart, where each instance of this chart is differenciated using aliases.

The federating chart will provide the configuration for each microservices in its values file.

## Usage

TBD

## Configuration options

TBD
