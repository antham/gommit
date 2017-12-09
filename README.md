Gommit [![CircleCI](https://circleci.com/gh/antham/gommit/tree/master.svg?style=svg)](https://circleci.com/gh/antham/gommit/tree/master) [![codecov](https://codecov.io/gh/antham/gommit/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/gommit) [![codebeat badge](https://codebeat.co/badges/cc515300-053e-4b62-8184-645be6e6aa2f)](https://codebeat.co/projects/github-com-antham-gommit) [![Go Report Card](https://goreportcard.com/badge/github.com/antham/gommit)](https://goreportcard.com/report/github.com/antham/gommit) [![GoDoc](https://godoc.org/github.com/antham/gommit?status.svg)](http://godoc.org/github.com/antham/gommit)
======

Gommit analyze commits messages to ensure they follow defined pattern.

[![asciicast](https://asciinema.org/a/0j12qm7yay1kku7o3vrs67pv2.png)](https://asciinema.org/a/0j12qm7yay1kku7o3vrs67pv2)

## Summary

* [Setup](#setup)
* [Usage](#usage)
* [Practical Usage](#practical-usage)
* [Third Part Libraries](#third-part-libraries)

## Setup

Download from release page according to your architecture gommit binary : https://github.com/antham/gommit/releases

### Define a file .gommit.toml

Create a file ```.gommit.toml``` at the root of your project, for instance :

```toml
[config]
exclude-merge-commits=true
check-summary-length=true
summary-length=50

[matchers]
all="(?:ref|feat|test|fix|style)\\(.*?\\) : .*?\n(?:\n?(?:\\* |  ).*?\n)*"

[examples]
a_simple_commit="""
[feat|test|ref|fix|style](module) : A commit message
"""
an_extended_commit="""
[feat|test|ref|fix|style](module) : A commit message

* first line
* second line
* and so on...
"""
```

#### Config

* ```exclude-merge-commits``` : if set to true, will not check commit mesage for merge commit
* ```check-summary-length``` : if set to true, check commit summary length, default is 50 characters
* ```summary-length``` : you can override the default value summary length, which is 50 characters, this config is used only if check-summary-length is true

#### Matchers

You can define as many matchers you want using regexp, naming is up to you, they will all be compared against a commit message till one match. Regexps used support comments, possessive match, positive lookahead, negative lookahead, positive lookbehind, negative lookbehind, back reference, named back referenc and conditionals.

#### Examples

Provided to help user to understand where is the problem, like matchers you can define as many examples as you want, they all will be displayed to the user if an error occured.

If you defined for instance  :

```
a_simple_commit="""
[feat|test|ref|fix|style](module) : A commit message
"""
```

this example will be displayed to the user like that :

```
A simple commit :

[feat|test|ref|fix|style](module) : A commit message
```

key is used as a title, underscore are replaced with withespaces.

## Usage

```bash
Ensure your commit messages are consistent

Usage:
  gommit [command]

Available Commands:
  check       Check ensure a message follows defined patterns
  version     App version

Flags:
      --config string    (default ".gommit.toml")
  -h, --help            help for gommit

Use "gommit [command] --help" for more information about a command.
```

### check

```bash
Check ensure a message follows defined patterns

Usage:
  gommit check [flags]
  gommit check [command]

Available Commands:
  commit      Check commit message
  message     Check message
  range       Check messages in commit range

Flags:
  -h, --help   help for check

Global Flags:
      --config string    (default ".gommit.toml")

Use "gommit check [command] --help" for more information about a command.


You need to provide two commit references to run matching for instance :
```

#### check commit

Check one comit from its commit ID, doesn't support short ID currently :

```gommit check commit aeb603ba83614fae682337bdce9ee1bad1da6d6e```

#### check message

Check a message, useful for script for instance when you want to use it with git hooks :

```gommit check message "Hello"```

#### check range

Check a commit range, useful if you want to use it with a CI to ensure all commits in branch are following your conventions :

* with relative references                             : ```gommit check range master~2^ master```
* with absolute references                             : ```gommit check range dev test```
* with commit ids (doesn't support short ID currently) : ```gommit check range 7bbb37ade3ff36e362d7e20bf34a1325a15b 09f25db7971c100a8c0cfc2b22ab7f872ff0c18d```

## Practical usage

If your system isn't described here and you find a way to have gommit working on it, please improve this documentation by doing a PR for the next who would like to do the same.

### Git hook

It's possible to use gommit to validate each commit when you are creating them. To do so, you need to use the ```commit-msg``` hook, you can replace default script with this one :

```
#!/bin/sh

gommit check message "$(cat "$1")";
```

### Travis

In travis, all history isn't cloned, default depth is 50 commits, you can change it : https://docs.travis-ci.com/user/customizing-the-build#Git-Clone-Depth.

First, we download the binary from the release page according to the version we want and we add in ```.travis.yml``` :

```yaml
before_install:
  - wget -O /tmp/gommit https://github.com/antham/gommit/releases/download/v2.0.0/gommit_linux_386 && chmod 777 /tmp/gommit
```

We can add a perl script in our repository to analyze the commit range against master for instance (master reference needs to be part of cloned history):

```perl
#!/bin/perl

`git ls-remote origin master` =~ /([a-f0-9]{40})/;

my $refHead = `git rev-parse HEAD`;
my $refTail = $1;

chomp($refHead);
chomp($refTail);

if ($refHead eq $refTail) {
    exit 0;
}

system "gommit check range $refTail $refHead";

if ($? > 0) {
    exit 1;
}
```

And finally in ```.travis.yml```, make it crashs when an error occured :

```yaml
script: perl test-branch-commit-messages-in-travis.pl
```

### CircleCI

In CircleCI (2.0), there is an environment variable that describe current branch : ```CIRCLE_BRANCH``` (https://circleci.com/docs/2.0/env-vars/#circleci-environment-variable-descriptions).

First, we download the binary from the release page according to the version we want and we add in ```.circleci/config.yml``` :

```yaml
- run:
    name: Get gommit binary
    command: |
      mkdir /home/circleci/bin
      wget -O ~/bin/gommit https://github.com/antham/gommit/releases/download/v2.0.0/gommit_linux_386 && chmod 777 ~/bin/gommit
```

And we can run gommit against master for instance :

```
- run:
    name: Run gommit
    command: |
      ~/bin/gommit check range $(git rev-parse origin/master) $(git rev-parse ${CIRCLE_BRANCH})
```

## Third Part Libraries

### Nodejs

* [gommitjs](https://github.com/dschnare/gommitjs) : A Nodejs wrapper for gommit
