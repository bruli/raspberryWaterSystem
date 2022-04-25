#! /bin/bash

require() {
    command -v "$1" >/dev/null 2>&1 || { echo >&2 "$1 program is required"; exit 1; } 
}
