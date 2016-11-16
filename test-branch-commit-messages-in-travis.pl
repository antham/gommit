#!/bin/perl

if ($ENV{'TRAVIS_PULL_REQUEST'} == 'false') {
    exit 0;
}

`git fetch --depth=1 origin master 2>&1 >/dev/null`;

my $head = `git rev-parse HEAD`;
my $master = `git rev-parse FETCH_HEAD`;

chomp($head);
chomp($master);

system "gommit check range $master $head";

if ($? > 0) {
    exit 1;
}
