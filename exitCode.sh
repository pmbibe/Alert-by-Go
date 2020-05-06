#!/bin/bash
ssh root@$2 service $1 status > /dev/null
