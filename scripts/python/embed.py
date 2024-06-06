# scripts/python/embed.py

import warnings
import torch
import torch.nn.functional as F
from sentence_transformers import SentenceTransformer
import argparse

def split_texts(texts, chunk_size, chunk_overlap):
    chunks = []
    for text in texts:
        for i in range(0, len(text), chunk_size - chunk_overlap):
            chunk = text[i:i + chunk_size]
            chunks.append(chunk)
    return chunks

def main(model_name, revision, texts, chunk_size, chunk_overlap):
    # Disable FutureWarnings
    warnings.filterwarnings("ignore", category=FutureWarning, message="`resume_download` is deprecated and will be removed in version 1.0.0.")

    # Check if MPS is available
    device = torch.device("mps" if torch.backends.mps.is_built() else "cpu")

    # Print the device
    print(f"Device: {device}")

    # Load the model and move it to the MPS device
    model = SentenceTransformer(model_name, revision=revision)
    model.to(device)

    # Split texts into chunks
    chunks = split_texts(texts, chunk_size, chunk_overlap)

    # Compute embeddings and move them to the MPS device
    embeddings = model.encode(chunks, convert_to_tensor=True, device=device)

    # foreach chunk, print the chunk and its embedding in nice format
    for chunk, embedding in zip(chunks, embeddings):
        print(f"Chunk: {chunk}")
        print(f"Embedding: {embedding}")
        print()

    # Compute cosine-similarity for each pair of chunks
    scores = F.cosine_similarity(embeddings.unsqueeze(1), embeddings.unsqueeze(0), dim=-1)

    print(scores.cpu().numpy())

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Sentence Embedding Script")

    parser.add_argument("--model_name", type=str, default="avsolatorio/GIST-small-Embedding-v0", help="Name of the model to use")
    parser.add_argument("--revision", type=str, default=None, help="Model revision to use")
    parser.add_argument("--texts", nargs='+', default=[
        "Illustration of the REaLTabFormer model. The left block shows the non-relational tabular data model using GPT-2 with a causal LM head. In contrast, the right block shows how a relational dataset's child table is modeled using a sequence-to-sequence (Seq2Seq) model. The Seq2Seq model uses the observations in the parent table to condition the generation of the observations in the child table. The trained GPT-2 model on the parent table, with weights frozen, is also used as the encoder in the Seq2Seq model.",
        "Predicting human mobility holds significant practical value, with applications ranging from enhancing disaster risk planning to simulating epidemic spread. In this paper, we present the GeoFormer, a decoder-only transformer model adapted from the GPT architecture to forecast human mobility.",
        "As the economies of Southeast Asia continue adopting digital technologies, policy makers increasingly ask how to prepare the workforce for emerging labor demands. However, little is known about the skills that workers need to adapt to these changes"
    ], help="List of texts to encode")
    parser.add_argument("--chunk_size", type=int, default=128, help="Size of each chunk")
    parser.add_argument("--chunk_overlap", type=int, default=32, help="Overlap between chunks")

    args = parser.parse_args()

    main(args.model_name, args.revision, args.texts, args.chunk_size, args.chunk_overlap)
