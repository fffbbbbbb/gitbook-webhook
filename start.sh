#!/bin/sh

cd ./script

for f in *.sh; do
  sh "$f"
done