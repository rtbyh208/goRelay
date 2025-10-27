#!/bin/bash

gitLog=$(git log -1 --pretty=format:"%H")
version=$1
nowDate=$(date +"%F")

platInfo=("linux/amd64" "linux/arm64" "windows/amd64" "android/arm64")

for i in ${platInfo[@]}
do
    os=$(echo ${i} | awk -F"/" '{print $1}')
    arch=$(echo ${i} | awk -F"/" '{print $2}')

    make GoOS=${os} GoArch=${arch} VERSION=${version} BuildAt=${nowDate} GITCOMMIT=${gitLog}
done
