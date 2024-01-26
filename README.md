# menv - maven environment manager

menv is a simple tool to manage multiple maven environments. It allows you to switch between different maven
configurations.

# Prerequisites

* MacOS or Linux
* [homebrew](https://brew.sh/) installed
* [maven](https://maven.apache.org/) installed through homebrew

# Installation

```bash
brew tap thecheerfuldev/cli
brew install thecheefuldev/cli/menv
```

# Usage

```bash
menv --help
```

## Environment variables

* MENV_EDITOR: The editor to use for editing the maven settings.xml file. Default: vi
* MENV_DISABLE_WRAPPER: If set to true, the maven wrapper will not be used. Default: false
* MENV_VERBOSE: If set to true, menv will print the active profile with every mvn execution. Default: false

## Create and use a new profile workflow

### 1. Create a new profile

```bash
menv new <profile-name>
```

### 2. Edit/create the settings.xml file of the profile

```bash
menv edit <profile-name>
```

### 3. Edit/create the MAVEN_OPTS of the profile

```bash
menv editopts <profile-name>
```

### 4. Use the profile

```bash
menv set <profile-name>
```

# Special thanks

* [IvoNet](https://github.com/IvoNet) for creating the original version of this tool, and pushing me to rewrite it
  in Go!