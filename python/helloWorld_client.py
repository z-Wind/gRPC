"""The Python implementation of the gRPC client."""

import logging

import grpc

import helloWorld_pb2
import helloWorld_pb2_grpc


def getDoubleSync(stub, value):
    req = helloWorld_pb2.Int(value=value)
    rsp = stub.Double(req)

    print(f"double {value} => {rsp.value}")


def getDoubleAsync(stub, value):
    req = helloWorld_pb2.Int(value=value)
    rspFuture = stub.Double.future(req)

    rsp = rspFuture.result()
    print("Async")

    print(f"double {value} => {rsp.value}")


def getRange(stub, value):
    req = helloWorld_pb2.Int(value=value)
    rsp = stub.Range(req)

    for i in rsp:
        print(f"range {value} => {i.value}")


def getSumSync(stub, value):
    req = (helloWorld_pb2.Int(value=i) for i in range(value))
    rsp = stub.Sum(req)

    print(f"Sum range({value}) => {rsp.value}")


def getSumAsync(stub, value):
    req = (helloWorld_pb2.Int(value=i) for i in range(value))
    rspFuture = stub.Sum.future(req)

    rsp = rspFuture.result()
    print("Async")

    print(f"Sum range({value}) => {rsp.value}")


def getDoubleIter(stub, value):
    req = (helloWorld_pb2.Int(value=i) for i in range(value))
    rsp = stub.DoubleIter(req)
    for i in rsp:
        print(f"double range({value}) => {i.value * 2}")


def run():
    # NOTE(gRPC Python Team): .close() is possible on a channel and should be
    # used in circumstances in which the with statement does not fit the needs
    # of the code.
    with grpc.insecure_channel("localhost:50051") as channel:
        stub = helloWorld_pb2_grpc.HelloWorldStub(channel)
        print("-------------- Double Sync--------------")
        getDoubleSync(stub, 5)
        print("-------------- Double Async-------------")
        getDoubleAsync(stub, 5)
        print("-------------- Range -------------------")
        getRange(stub, 5)
        print("-------------- Sum Sync-----------------")
        getSumSync(stub, 5)
        print("-------------- Sum Async----------------")
        getSumAsync(stub, 5)
        print("-------------- DoubleIter --------------")
        getDoubleIter(stub, 5)


if __name__ == "__main__":
    logging.basicConfig()
    run()

# Output
# -------------- Double Sync--------------
# double 5 => 10
# -------------- Double Async-------------
# Async
# double 5 => 10
# -------------- Range -------------------
# range 5 => 0
# range 5 => 1
# range 5 => 2
# range 5 => 3
# range 5 => 4
# -------------- Sum Sync-----------------
# Sum range(5) => 10
# -------------- Sum Async----------------
# Async
# Sum range(5) => 10
# -------------- DoubleIter --------------
# double range(5) => 0
# double range(5) => 4
# double range(5) => 8
# double range(5) => 12
# double range(5) => 16
