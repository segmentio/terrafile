# Terrafile

Terrafile is a binary written in Go to systematically manage external modules from Github for use in Terraform. This is a fork of the original version found here: https://github.com/coretech/terrafile

## How to install

### macOS

```sh
brew tap segmentio/packages
brew install segmentio/packages/terrafile
```

### Linux

Download your preferred flavor from the [releases](https://github.com/coretech/terrafile/releases/latest) page and install manually.

## How to use

Terrafile expects a file named `Terrafile` which will contain your terraform module dependencies in a yaml like format.

An example Terrafile:

```
git@github.com:segmentio/my-modules:
     - chamber_v2.0.0
     - constants_v1.0.1
     - iam_v1.0.0
     - service_v1.0.4
     - rds_v0.0.5
     - worker_v1.1.0
     - master
```

Terrafile config file in current directory and modules exported to .terrafile/<user>/<repo>/<ref>

```sh
$ terrafile
[*] Cloning   git@github.com:segmentio/my-modules
[*] Vendoring ref chamber_v2.0.0
[*] Vendoring ref constants_v1.0.1
[*] Vendoring ref iam_v1.0.0
[*] Vendoring ref service_v1.0.4
[*] Vendoring ref rds_v0.0.5
[*] Vendoring ref task-role_v1.0.5
[*] Vendoring ref master
```

Terrafile config file in custom directory

```sh
$ terrafile -f config/Terrafile
```

Terraform modules exported to custom directory

```sh
$ terrafile -p /path/to/custom_directory
```
