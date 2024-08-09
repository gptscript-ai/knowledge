#!/usr/bin/env bash

# check if homebrew is setup
if ! command -v brew &> /dev/null; then
    echo "Homebrew is not installed. Please install homebrew and try again."
    exit 1
fi

# install dependencies
brew install promptfoo
