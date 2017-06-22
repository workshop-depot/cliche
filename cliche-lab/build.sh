#!/bin/sh

clear

TRG_PKG='main'
BUILD_TIME=$(date +"%Y%m%d.%H%M%S")
CommitHash=N/A
GoVersion=N/A
GitTag=N/A

if [[ $(go version) =~ [0-9]+\.[0-9]+\.[0-9]+ ]];
then
    GoVersion=${BASH_REMATCH[0]}
fi

if [[ $(git tag) =~ [^[:space:]]+ ]];
then
    GitTag=${BASH_REMATCH[0]}
fi

GH=$(git log -1 --pretty=format:%h)
if [[ GH =~ 'fatal' ]];
then
    CommitHash=N/A
else
    CommitHash=$GH
fi

FLAG="-X $TRG_PKG.BuildTime=$BUILD_TIME"
FLAG="$FLAG -X $TRG_PKG.CommitHash=$CommitHash"
FLAG="$FLAG -X $TRG_PKG.GoVersion=$GoVersion"
FLAG="$FLAG -X $TRG_PKG.GitTag=$GitTag"

go build -v -ldflags "$FLAG"
