from concurrent import futures
import grpc
import TextGeneration_pb2
import TextGeneration_pb2_grpc
from mlx_lm import load, generate

# Load the model
model, tokenizer = load("mlx-community/Phi-3-mini-128k-instruct-8bit")

class TextGeneratorServicer(TextGeneration_pb2_grpc.TextGeneratorServicer):
    def GenerateTextStream(self, request, context):
        for response_part in generate(model, tokenizer, prompt=request.prompt, max_tokens=128000):
            yield TextGeneration_pb2.TextResponse(response=response_part)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    TextGeneration_pb2_grpc.add_TextGeneratorServicer_to_server(TextGeneratorServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()