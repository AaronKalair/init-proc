init-proc
=========

A simple process designed to be used as PID 1 in a container and

1. Proxy signals through to another process
2. Wait on orphan processes that have reparented to it


Usage
=====

1. `go build`
2. ./init-proc <command to run> <args to this command>

You'll probably want to use it in a container, see the Dockerfile in this repo
for an example of how to do that
