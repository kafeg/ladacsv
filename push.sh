#!/bin/bash

git status

git commit -am "Updated models `date --iso-8601`"

git tag "v`date --iso-8601`"

git push

git push --tags