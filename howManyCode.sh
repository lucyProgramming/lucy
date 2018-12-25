#!/bin/bash
find . -name "*.go" | xargs wc -l|grep "total"|awk '{print $1}'