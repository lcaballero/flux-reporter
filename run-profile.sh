#!/bin/bash


files=$(find . -type f -regex ".*\.go$" | tr "\n" " ")

echo $files

go tool pprof profile.bin "$files"
