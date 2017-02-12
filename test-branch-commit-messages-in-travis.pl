#!/bin/perl

my $branch = '';

if ($ENV{'TRAVIS_PULL_REQUEST'} eq 'false') {
    $branch = `$ENV{TRAVIS_BRANCH}`;
} else {
    $branch = `$ENV{TRAVIS_PULL_REQUEST_BRANCH}`;
}

if ($branch eq 'master') {
    exit 0;
}

`git ls-remote origin master` =~ /([a-f0-9]{40})/;

my $refHead = `git rev-parse HEAD`;
my $refTail = $1;

chomp($refHead);
chomp($refTail);

system "gommit check range $refTail $refHead";

if ($? > 0) {
    exit 1;
}
