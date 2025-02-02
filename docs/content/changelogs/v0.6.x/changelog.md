---
title: Changes by patch
weight: 10
---


## v0.6.10

This release is mainly bugfixes and updates to some dependencies and tooling.

No notable changes to features or the API.


[v0.6.9...v0.6.10 diff on GitHub](https://github.com/hofstadter-io/hof/compare/v0.6.9...v0.6.10)

## v0.6.9

This release brings more consistency and long-term stability to hof.
There are also new features like the CUE commands and the hof TUI.


### main changes

- CUE v0.9.0 (+1 bugfix commit after)
- Added the CUE commands [def,eval,export,vet]
- Adjusted flags for consistency with CUE and internally across commands
- Added `hof tui` for real-time CUE manipulation and exploration
- Upgradee hof/flow to new runtime, this marks the point all subsystems have migrated
- Added support for bulk processing with parallelism
- Support for CUE style data placement with hof extensions
- Support for embedding user files into CUE values as strings

### other changes

- several bug fixes in mod, containers, #hof
- deal with macos woes on GHA

[v0.6.8...v0.6.9 diff on GitHub](https://github.com/hofstadter-io/hof/compare/v0.6.8...v0.6.9)


## v0.6.8

Almost every system was upgraded in this release and it marks a milestone in API stability.

### main changes

- Now using CUE `v0.6.0` with required fields
- New shared runtime for core commands to enable greater cohesion and consistency
- `#hof` metadata for core objects gen, datamodel, flow
- Refresh of the datamodel command
  - more flexible with a user defined structure
  - nested history tracking and injection into code generation
  - The schema has changed significantly, see breaking changes.
- Refresh of module command 
  - focus on CUE with a simpler implementation
  - automatic dependency inference with `hof mod tidy`
  - support for both OCI & git, public and private
- General schema cleanup and refactor
- Support for Nerdctl or Podman as alternatives to Docker
- Experimental LLM Chat features

### other changes

- More formatters, support running any version
- data formatting stability
- support for CUE's `@tag()` and `-t`, additonally `--inject-env` and `--inject-data`
- `--keep-deleted` for `hof gen`
- implement feedback command
- support data placement for most commands with `file.json@path.to.value`
- Improvements to several hof/flow tasks
- CI upgrades
- many test improvements and coverage increase
- bug fixes
- dependency updates

### breaking changes

The data model schemas have changed to enable the new features.
You can [learn how to upgrade them from this link](../upgrading-0.6.7-to-0.6.8/).

[v0.6.7...v0.6.8 diff on GitHub](https://github.com/hofstadter-io/hof/compare/v0.6.7...v0.6.8)


## v0.6.7

### hof create

Introduces `hof create` to as a "npm create-react-app" for anything.
Easy bootstrap or add files to an existing project from any remote repository.

- essentially a generator with `Create` field to get inputs from user.
- configure an interactive prompt and input schema, user can also use flags to fill prompt
- adds a new `hof create <repo>` which fetches and prompts user for starting input
- adds `Create` schema and includes in generators
- add schema and code for implementing an input prompt
- works with remote repo, and also locally for directories
- makes it easy for you to provide one-line setup instructions for your own generators

To learn more about `hof create`, see
[docs/getting-started/create](/getting-started/create/).

### breaking changes:

Hof's `indent` template helper was updated to mimic the behavior
of helm's. This required swapping the order of arguments.

You should now use `{{ indent <string|int> <content> }}`


### other changes

- remove modder name from mod cache, flattening because they are always a git repo at a tag
- use cache dir for remote repositories
- enable symlinks for local, replaced, cue dependencies
- remove old or unused code
- fix and enable more formatters
- improve docker images
- improve tests & CI
- several bug fixes and edge case handling


[v0.6.6...v0.6.7 diff on GitHub](https://github.com/hofstadter-io/hof/compare/v0.6.6...v0.6.7)


## v0.6.6


Fixes several bugs and re-adds telemetry with proper user controls.


[v0.6.5...v0.6.6 diff on GitHub](https://github.com/hofstadter-io/hof/compare/v0.6.5...v0.6.6)

## v0.6.5

### hof fmt

Introduces `hof fmt` to format code, in a beta state

- used during code gen, because we need this for diff3 correctness
- subcommand for formatting arbitrary files and managing containers
- adds `prettier` and `black` formatters to expand languages

### other changes

- update a number of version used during init of various things (mods & gens)
- `lookup` template helper now supports OpenAPI refs `#/path/to/thing`
- fix template path resolving when in gen is run in a subdir
- fixes issue when same outdir is used, with the same gen, from different locations
- shadow dir moved to be next to cue.mod, so that one can gen from any directory
- shadow dir updated to reflect path from CUE mod root output, use this as the path
- Static Files now support diff3 for user additions

[v0.6.4...v0.6.5 diff on GitHub](https://github.com/hofstadter-io/hof/compare/v0.6.4...v0.6.5)

## v0.6.4


### bugfix and cleanup

- The `and` template helper was custom, wrong, and needed to be removed. This is the primary reason for the release.
- Removed many template helpers which overrode or duplicated Go's, or did not make sense or seemed unusual.
- Add newline when writing generator outputs, ending on a curly brace seemed odd.
- Update almost all deps, remove some others, need Go 1.17 to go out of fashion before rest can be... generics and breaking changes in deps...

[v0.6.3...v0.6.4 diff on GitHub](https://github.com/hofstadter-io/hof/compare/v0.6.3...v0.6.4)

## v0.6.3

### hof gen ad-hoc mode

Adds flags to `hof gen` to support ad-hoc code gen,
so you do not need to setup a generator to use.

- `--template`/`-T` flag to specify templates, input data, schemas, and output files
- `--partial`/`-P` flag to support ad-hoc partial templates
- `--watch`/`-w` flag to support watching globs and regenerating (also works for generators)
- `--as-module` flag turns your other flags into a reusable and sharable generator module
- `--init` flag bootstraps a new modular generator in the current directory

The `-T` flag has a flexible format so you can
supply multiple templates and control the data.
It lets you specify the mapping from template
to input & schema, to output filepath.

```
hof gen data.cue -T template.txt
hof gen data.yaml schema.cue -T template.txt > output.txt
```

See the following to learn more

- `hof gen -h`
- [getting-started/code-generation](/getting-started/code-generation/)

### other changes

- added `dict` to template helpers to create maps, useful for passing more than one arg to a partial template
- added `pascal` template helper for string casing
- load data with CUE code, more inline with `cue`
- (bug) remove some shell completion hacks
- more tests, bugfixes, and dep updates
- some small changes to the datamodel schema, namely attribute change to prep for enhancemnts

[v0.6.2...v0.6.3 diff on GitHub](https://github.com/hofstadter-io/hof/compare/v0.6.2...v0.6.3)


## v0.6.2

### hof flow

hof flow is a custom cue/flow runtime with more task types.

See the following to learn more:

- `hof flow -h`
- [getting-started](/getting-started/task-engine/)
- [task-engine](/task-engine/)
- [task reference](/task-engine/tasks/)

### other changes

- CUE v0.4.3
- Go v1.18
- other dep updates
- various bugfixes


[v0.6.1...v0.6.2 diff on GitHub](https://github.com/hofstadter-io/hof/compare/v0.6.1...v0.6.2)

## v0.6.1

### hof datamodel

hof datamodel is a tool to manage your models.
Define, validate, checkpoint, diff, and migrate.

See the following to learn more:

- `hof datamodel -h`
- [getting-started/data-layer](/getting-started/data-layer/)
- [data modeling section](/data-modeling/)

### other changes

- data files from generators
- various bugfixes
- cleanup and legacy code removal

[v0.6.0...v0.6.1 diff on GitHub](https://github.com/hofstadter-io/hof/compare/v0.6.0...v0.6.1)

## v0.6.0

- general cleanup, bugfixing, refactoring
- rework `hof gen` schemas
- remove some disjunctions in schema to improve performance
- better error messages
- enable subgenerators

[v0.6.0](https://github.com/hofstadter-io/hof/compare/v0.5.17...v0.6.0)

