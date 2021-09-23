#!/bin/bash

docker build . -t metal-hammer

sudo ignite image import --runtime=docker metal-hammer:latest
sudo ignite kernel import --runtime=docker metal-kernel:latest

sudo ignite run -i --config=./ignite.yml