# tq

- TinyQ is a lisp-powered message server based on nanomsg (aka zeroq) written in golang.
- TinyQ does only one thing: ship around and transform json messages
- TinyQ supports sending and receiving messages encoded as json via nanomsg and REST
- TinyQ is controlled by writing small lisp scripts that define channels and message transforms

## Installing

Install TinyQ using go get:

```bash
# go get github.com/SUNET/tq
```

Alternatively use the provided Dockerfile to create a distroless docker image:

```bash
# docker build -t tq .
```


## Getting started

TinyQ has a built-in tiny lisp (based on the sabre project) interpreter which can be started as a basic lisp RELP (Read EvaL Print) loop. The lisp environment comes with lambdas (functions) that is used to create message-based services. TinyQ messages are JSON-format but TinyQ doesn't really care about the content of the messages other than that they are syntactically correct JSON.

Building a TinyQ service typically involves writing and running small lisp programs. The smallest possible example is a "ping":

```lisp
(def onesec (timer "1s"))
(run (onesec))
```

Create a file named ping.tq with the above 2 lines and run it:

```bash
# tq --loglevel=debug ping.tq
```

The output should indicate that a single JSON message is created every second. The first line of the file calls the timer primitive to create a message channel that generates a JSON message every second. The second line calls the run primitive with an instantiation of the timer instance. Some primitives in TinyQ act on message channels while other primitives create new message channels. A message channel is typically created in two steps: first one is created (or configured) and then launched. 

The run primitive can be called like in the example above or without arguments last in the file to run all channels created up to that point.

## Primitives

- *pub* <url>: returns a message channel that publishes to a specified URL
- *sub* <url>: returns a message channel that subscribes to messages published to the URL
- *listen* [<host>\*]:<port>: run the API endpoint on the specified host:port
- *merge* <channel>\*: merge a set of channels into a single channel
- *script* <cmdline>: returns a message channel that runs the specified commandline for every message
- *rest* <url>: returns a message channel that accepts JSON messages by POST/PUT to the specified url
- *kazaam* <spec>: returns a message channel that transforms JSON using kazaam

## She-Bang!

TinyQ supports #! so you can do this after installing TinyQ:

```bash
#!/usr/bin/env tq
; lisp code here...
```

## Examples

### Publish/Subscribe

In two different windows (to avoid confusing the log outputs) run the following:

```bash
# tq --loglevel=debug examples/pub.tq
```

```bash
# tq --loglevel=debug examples/sub.tq
```
