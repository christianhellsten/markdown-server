#!/bin/bash
if [ ! -z "$PID" ]; then
  kill $PID
fi
go run main.go &
PID=$!
