#!/bin/bash

export GH_ACCESS_TOKEN=$(pass show gh-access-token)

./generate-readme
