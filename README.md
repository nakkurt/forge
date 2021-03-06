# Forge

[![CircleCI](https://circleci.com/gh/dominodatalab/forge.svg?style=shield)](https://app.circleci.com/pipelines/github/dominodatalab/forge)
[![Go Report Card](https://goreportcard.com/badge/github.com/dominodatalab/forge)](https://goreportcard.com/report/github.com/dominodatalab/forge)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/dominodatalab/forge)](https://pkg.go.dev/mod/github.com/dominodatalab/forge)

Forge is a Kubernetes controller designed to securely build OCI-compliant images
inside a cluster and push them to one or more target registries. This project
was derived from the work done in the [img][img] project and extended to support
build dispatch via a [custom resource definition][crd].

### Preparer Plugins

Forge supports the inclusion of custom plugins for the "preparation" phase of a build (between the initialization of the context and the actual image build).
This functionality is built with the [go-plugin](https://github.com/hashicorp/go-plugin) framework from Hashicorp.

#### Creation

[example/preparer_plugin.go](./docs/example/preparer_plugin.go) has the necessary structure for creating a new preparer plugin.
Functionality is implemented through two primary functions:

`Prepare(contextPath string, pluginData map[string]string) error`

Prepare runs between the context creation and image build starting. `contextPath` is an absolute path to the context for the build.
`pluginData` is the key-value data passed through the [ContainerImageBuild](./config/crd/bases/forge.dominodatalab.com_containerimagebuilds.yaml#L77-L82).

`Cleanup() error`

Cleanup runs after the build has finished (successfully or otherwise).

#### Using

To add a new runtime plugin for Forge, place a file in `/usr/local/share/forge/plugins/` (by default) or specify it with `--preparer-plugins-path`.

[img]: https://github.com/genuinetools/img
[crd]: https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/
