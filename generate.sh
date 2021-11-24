#!/bin/bash

protoc calculatorpb/calculator.proto --go_out=. --go-grpc_out=.