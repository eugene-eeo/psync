#!/usr/bin/env bats

@test "gluing exported file produces same output" {
    diff \
        <(cat ./psync | ./psync export | ./psync glue --verify) \
        <(cat ./psync)
}

@test "blocks/<checksum> files exist" {
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
