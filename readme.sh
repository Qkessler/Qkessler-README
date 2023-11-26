#!/bin/bash

export GH_ACCESS_TOKEN=$(pass show gh-access-token)

go run main.go

cd Qkessler
git add .

date=$(date '+%Y-%m-%d %H:%M:%S')
git commit -am "Add dynamic README $date"
git push
