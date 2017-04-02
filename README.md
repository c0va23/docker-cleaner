# duclean
duclean (**D**ocker **u**niversal **clean**er) is command line application for cleaning docker garbage (images, containers).

## Usage

```
Usage: declean [OPTIONS] COMMAND [arg...]

Docker universal cleaner

Options:
  --safe-period=0   Save period (seconds)
  --dry-run=false   Dry run

Commands:
  images       Clean useless images
  containers   Clean containers
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
