"""The Python implementation of the gRPC server."""

from concurrent import futures
import logging
import time

import grpc

import helloWorld_pb2
import helloWorld_pb2_grpc

_ONE_DAY_IN_SECONDS = 60 * 60 * 24


class HelloWorldServicer(helloWorld_pb2_grpc.HelloWorldServicer):
    """Provides methods that implement functionality of hello world server."""

    def __init__(self):
        pass

    def Double(self, Int, context):
        rsp = helloWorld_pb2.Int(value=Int.value * 2)

        return rsp

    def Range(self, Int, context):
        for i in range(Int.value):
            rsp = helloWorld_pb2.Int(value=i)
            yield rsp

    def Sum(self, Int_iterator, context):
        result = 0
        for i in Int_iterator:
            result += i.value

        rsp = helloWorld_pb2.Int(value=result)
        return rsp

    def DoubleIter(self, Int_iterator, context):
        for i in Int_iterator:
            rsp = helloWorld_pb2.Int(value=i.value * 2)
            yield rsp


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    helloWorld_pb2_grpc.add_HelloWorldServicer_to_server(HelloWorldServicer(), server)
    server.add_insecure_port("[::]:50051")
    server.start()
    print("Listen to", "[::]:50051")
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == "__main__":
    logging.basicConfig()
    serve()
