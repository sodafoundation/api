# OpenSDS

[![Go Report Card](https://goreportcard.com/badge/github.com/opensds/opensds?branch=master)](https://goreportcard.com/report/github.com/opensds/opensds)
[![Build Status](https://travis-ci.org/opensds/opensds.svg?branch=master)](https://travis-ci.org/opensds/opensds)
[![codecov.io](https://codecov.io/github/opensds/opensds/coverage.svg?branch=master)](https://codecov.io/github/opensds/opensds?branch=master)
[![Releases](https://img.shields.io/github/release/opensds/opensds/all.svg?style=flat-square)](https://github.com/opensds/opensds/releases)
[![LICENSE](https://img.shields.io/github/license/opensds/opensds.svg?style=flat-square)](https://github.com/opensds/opensds/blob/master/LICENSE)

<img src="https://www.opensds.io/wp-content/uploads/sites/18/2016/11/logo_opensds.png" width="100">

## Latest Release: v0.6.0 Capri

[OpenAPI doc](http://petstore.swagger.io/?url=https://raw.githubusercontent.com/opensds/opensds/v0.6.0/openapi-spec/swagger.yaml)

[Release notes](https://github.com/opensds/opensds/releases/tag/v0.6.0)

## Introduction

The [OpenSDS Project](https://opensds.io/) is a collaborative project under Linux
Foundation supported by storage users and vendors, including
Dell EMC, Intel, Huawei, Fujitsu, Western Digital, Vodafone, NTT and Oregon State University. The project
will also seek to collaborate with other upstream open source communities
such as Cloud Native Computing Foundation, Docker, OpenStack, and Open
Container Initiative.

It is a software defined storage controller that provides
unified block, file, object storage services and focuses on:

* *Simple*: well-defined API that follows the [OpenAPI](https://github.com/OAI/OpenAPI-Specification) specification.
* *Lightweight*: no external dependencies, deployed once in binary file or container.
* *Extensible*: pluggable framework available for different storage systems, identity services, capability filters, etc.

## Community

The OpenSDS community welcomes anyone who is interested in software defined
storage and shaping the future of cloud-era storage. If you are a company,
you should consider joining the [OpenSDS Project](https://opensds.io/).
If you are a developer and would like to be part of the code development
that is happening now, please refer to the Contributing sections below.

## Collaborative Testing

* [CNCF Cluster](https://github.com/cncf/cluster/issues/30)

## Contact

* Mailing list: [opensds-tech-discuss](https://lists.opensds.io/mailman/listinfo/opensds-tech-discuss)
* slack: #[opensds](https://opensds.slack.com)
* Ideas/Bugs: [issues](https://github.com/opensds/opensds/issues)

## OpenSDS Controller Work Group

See [COMMUNITY](COMMUNITY.md) for details on discussion of the OpenSDS architecture design and feature development.

## Contributing

If you're interested in being a contributor and want to get involved in
developing the OpenSDS code, please see [CONTRIBUTING](CONTRIBUTING.md) for
details on submitting patches and the contribution workflow.

## Hacking

Please refer to [HACKING](HACKING.md) for any requirements when you want to perform code
development for OpenSDS.

## Installation

Please refer to [INSTALL](INSTALL.md) for any requirements when you want to perform code
development for OpenSDS.

## Auto-generated SDK

To generate SDK (e.g. Java, C#, Ruby, etc) to access the REST API, please consider using [OpenAPI Generator](https://github.com/OpenAPITools/openapi-generator) given the [OpenAPI/Swagger spec](https://raw.githubusercontent.com/opensds/opensds/master/openapi-spec/swagger.yaml). If you need help with OpenAPI Generator, please reach out to the OpenAPI Generator community by opening an [issue](https://github.com/OpenAPITools/openapi-generator/issues/new).

## License

OpenSDS is under the Apache 2.0 license. See the [LICENSE](LICENSE) file for details.
