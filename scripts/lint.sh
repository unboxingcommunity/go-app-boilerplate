#!/bin/bash

#current directory
dir="$(pwd)"
#parent directory
parentdir="$(dirname "$dir")"
#make flags
makeArgs=$1

#sets proper path based on if invoked from make or directly
if [ "$makeArgs" != "TRUE" ]; then 
    cd $parentdir    
fi

# checks if golint exists or not , if not installs it
wq=$(which golint)
if [ "$wq" == "" ]; then go get -u golang.org/x/lint/golint; fi

# run linter
golint -set_exit_status $(go list ./... | grep -v /vendor/)