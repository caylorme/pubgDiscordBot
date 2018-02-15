#!/bin/bash

docker run -v "$PWD":/usr/src/pubgDiscordBot -w /usr/src/pubgDiscordBot golang:latest go build