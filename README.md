# Adam's Go TODO app

## Requirements

* Docker
* Go (built on 1.15.6)
* GNU Make if using the makefile
* Postgresql installation - _the makefile will attempt to spin up a container version for testing_

## Building and running

The easiest way to build the app would be via the Makefile:
```bash
make build
```

Which would give you the binary.

To run something more useful, you can use the provided kubernetes config which can be run in a test cluster:
```bash
kubectl apply -f k8s/
```
_Not prod ready!_

Adapted from
[Semaphore Tutorial](https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql)
and
[Better Programming Tutorial](https://medium.com/better-programming/build-a-simple-todolist-app-in-golang-82297ec25c7d)
