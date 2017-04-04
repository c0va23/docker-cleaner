# duclean
duclean (**D**ocker **u**niversal **clean**er) is command line application for cleaning docker garbage (images, containers).

## Usage

```
Usage: declean COMMAND [arg...]

Docker universal cleaner

Commands:
  images       Clean useless images
  containers   Clean containers
  volumes      Clean useless volumes
  version      Print version

Run 'declean COMMAND --help' for more information on a command.
```

## Build
### Requirements
- golang (>= v1.8)
- git
- make

### Compile
```bash
make get_deps
make build
./bin/duclean version
```

## Known issues
- *Description*: Any command return error message `Error response from daemon: client is newer than server (client API version: 1.xx, server API version: 1.24)`

  *Solution*: Run `duclean` with environment variable `DOCKER_API_VERSION=1.xx`, when `1.xx` is version from error
  message.

  *Example*:
  ```bash
  DOCKER_API_VERSION=1.24 duclean
  ```
