#!/usr/bin/env bats

@test "gluing exported file produces same output" {
    diff \
        <(cat main.go | ./psync export | ./psync glue --verify) \
        <(cat main.go)
    [ "$?" -eq 0 ]
}

@test "directory exists with blocks" {
    hashlist=$(cat main.go | ./psync export)
    for line in $(echo $hashlist)
    do
        test -f "$HOME/.psync/blocks/$line"
    done
}

@test "cat returns exact contents" {
    hashlist=$(echo "abc" | ./psync export)
    [[ $(./psync cat $hashlist) == "abc" ]]
}
