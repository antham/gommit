Gommit [![Build Status](https://travis-ci.org/antham/gommit.svg?branch=master)](https://travis-ci.org/antham/gommit) [![codecov](https://codecov.io/gh/antham/gommit/branch/master/graph/badge.svg)](https://codecov.io/gh/antham/gommit) [![codebeat badge](https://codebeat.co/badges/cc515300-053e-4b62-8184-645be6e6aa2f)](https://codebeat.co/projects/github-com-antham-gommit)
======

Gommit analyze commits messages to ensure they follow defined pattern.

[![asciicast](https://asciinema.org/a/0j12qm7yay1kku7o3vrs67pv2.png)](https://asciinema.org/a/0j12qm7yay1kku7o3vrs67pv2)

## Setup

Download from release page according to your architecture gommit binary : https://github.com/antham/gommit/releases

### Define a file .gommit.toml

Create a file ```.gommit.toml``` at the root of your project, for instance :

```toml
[config]
exclude-merge-commits=true
check-summary-length=true

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
* ```check-summary-length``` : if set to true, check commit summary length is 50 characters

#### Matchers

You can define as many matchers you want, naming is up to you, they will all be run against a commit message till one match.

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
  check       Check commit messages
  version     App version

Flags:
      --config string    (default ".gommit.toml")
  -h, --help            help for gommit

Use "gommit [command] --help" for more information about a command.
```

### check

You need to provide two commit references to run matching for instance :

```gommit check master~2^ master```

or

```gommit check dev test```
