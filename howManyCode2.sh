#!/bin/bash
find . -name "*.lucy"|xargs wc -l|grep "total"|awk '{print $1}'
