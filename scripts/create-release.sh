#!/bin/sh -e

release() {
    name=$1
    shift 1
    message="$name: $@"

    echo "Creating release '${name}' with message: '${message}'... "
    read -p "Type Y/y to confirm: " -n 1 -r
    echo    # move to a new line
    if [[ $REPLY =~ ^[Yy]$ ]]
    then
        git commit --allow-empty -m "$name: $message"
        git tag $name
        git push origin $name
    fi
}

release $@
