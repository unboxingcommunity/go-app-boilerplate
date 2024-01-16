#!/bin/bash

#current directory
dir="$(pwd)"
#parent directory
parentdir="$(dirname "$dir")"
#temporary result directory
wdir=""
#final result directory
tdir=""
#make flags
makeArgs=$1

#sets proper path based on if invoked from make or directly
if [ "$makeArgs" = "TRUE" ]; then 
    wdir=$dir"/_results" 
    tdir=$dir"/_reports"
else 
    wdir=$parentdir"/_results" 
    tdir=$parentdir"/_reports"
    cd $parentdir
fi

#if directory present clean and create new one
if [ -d $wdir ]; then rm -Rf $wdir; fi

# make clean new dir
mkdir $wdir

#install gocovmerge if not present
wq=$(which gocovmerge)
if [ "$wq" == "" ]; then go get github.com/wadey/gocovmerge; fi

#install go-junit-report if not present
wq=$(which go-junit-report)
if [ "$wq" == "" ]; then go get -u github.com/jstemmer/go-junit-report; fi

# Using the "-short" flag on the command line causes the testing.Short() function to return true. 
# You can use this to either add test cases or skip them, Example :
# func TestThatIsLong(t *testing.T) {
#     if testing.Short() {
#         t.Skip()
#     }
# }
# test for data race seperately
printf "CHECKING FOR DATA RACE\n\n"
    go test -race -short $(go list ./... | grep -v /vendor/) 2> $wdir"/race.out"

# gets list of all existing packages and runs test cover for them
printf "\nGENERATING COVERAGE DATA\n\n"
    PKG_LIST=$(go list ./... | grep -v /vendor/)
    for package in ${PKG_LIST}; do
        # echo "$package";
        go test -short -covermode=count -coverprofile "$wdir/${package##*/}.cov" "$package" ;
    done

printf "\nMERGING COVERAGE AND RUNNING UNIT TESTS (BACKGROUND)\n\n"
    # merge all .cov files into master coverage.cov
    tail -q -n +2 $wdir/*.cov >> $wdir"/coverage.cov"
    echo 'mode: count' | cat - $wdir"/coverage.cov" > $wdir"/temp.cov" && mv $wdir"/temp.cov" $wdir"/coverage.cov"
    # run all unit tests with both .out and .xml(junit) reports
    go test -short -v -test.covermode=count -test.coverprofile=$wdir"/unit.out" ./... 2>&1 | go-junit-report > $wdir"/junit.xml"
    # merge all test outputs and coverage files -> all.out
    gocovmerge $wdir"/unit.out" $wdir"/coverage.cov" > $wdir"/all.out"

    # check if existing reports dir exists or not
    if [ -d $tdir ]; then rm -Rf $tdir; fi
    # make new clean dir 
    mkdir $tdir

printf "\nGENERATING REPORTS AND PERFORMING CLEANUP\n\n"
    # move all.out to tdir directory
    mv $wdir"/all.out" $tdir"/all.out"
    mv $wdir"/junit.xml" $tdir"/junit.xml"
    # move race output if not empty
    if [ -s $wdir"/race.out" ]; then mv $wdir"/race.out" $tdir"/race.out"; fi
    # generate html report from all.out
    go tool cover -html $tdir"/all.out" -o $tdir"/coverage.html"
    # display functional coverage on console
    go tool cover -func=$tdir"/all.out"

    # clean up by removing unnecessary _results folder
    if [ -d $wdir ]; then rm -Rf $wdir; fi

