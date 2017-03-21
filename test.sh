#!/bin/bash

peg cpppeg/parser.peg
go build main.go
./main -enum test.h

#
