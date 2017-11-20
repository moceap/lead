#!/bin/bash
set -euo pipefail

ctrls=$(mktemp)
function clean {
    rm -f "$ctrls"
}
trap clean EXIT

lead discover 172.16.33.0/24 > "$ctrls"

lead --file "$ctrls" color 255 192 32
sleep 1
lead --file "$ctrls" brightness 1
sleep 1
lead --file "$ctrls" on

for ((i=2; i<32; i++)) ; do
	sleep 60
	lead --file "$ctrls" brightness "$i"
done
