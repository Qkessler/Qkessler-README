#!/bin/bash

export GH_ACCESS_TOKEN=$(pass show gh-access-token)

./generate-readme

cd Qkessler
git add .

date=$(date '+%Y-%m-%d %H:%M:%S')
git commit -am "Add dynamic README $date"
git push
