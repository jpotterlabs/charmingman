# Knowledge Base & RAG

CharmingMan enables you to ground your AI agents in specific documents and data through its **Retrieval-Augmented Generation (RAG)** pipeline. This ensures your assistant provides accurate answers based on your private data rather than just general knowledge.

## 🏗️ The RAG Pipeline

CharmingMan implements a standard RAG flow:
1. **Ingest**: You upload a document to the AI Gateway.
2. **Extract & Chunk**: The gateway extracts text and splits it into smaller, overlapping chunks.
3. **Embed**: Each chunk is converted into a high-dimensional vector using OpenAI's embedding models.
4. **Index**: Vectors and metadata are stored in a Vector Store (Pinecone or Local).
5. **Retrieve**: When you ask a question with `use_rag=true`, the gateway embeds your query and searches the Vector Store for the top matching chunks.
6. **Augment**: Relevant snippets are injected into the prompt before calling the LLM.

## 📄 Document Engine

The Document Engine in `backend/internal/document/` handles the heavy lifting of preparing your data.

### Supported Formats
- **Plain Text** (`.txt`, `.md`): Direct extraction.
- **PDF** (`.pdf`): Handled via the `ledongthuc/pdf` library.

### Intelligent Chunking
To preserve context, the `chunker.go` uses a recursive splitting strategy:
- Splits by paragraphs, then newlines, then spaces.
- **Max Chunk Size**: 1000 characters.
- **Overlap**: 200 characters (to ensure context isn't lost at the boundaries).

## 🗄️ Vector Stores

CharmingMan supports multiple backends for storing and searching embeddings:

| Store | ID | Description |
|-------|----|-------------|
| **Pinecone** | `pinecone` | Managed cloud vector database for production-grade performance. |
| **Local Store**| `local` | In-memory vector store for fast development and lightweight local use. |

## 🖼️ "The Stage" (TUI Previewer)

When an agent uses RAG to answer a question, the TUI displays the sourced document snippets in a dedicated window called **"The Stage"**.

- **Source Attribution**: The chat history includes references to the document IDs and similarity scores of the retrieved context.
- **Inspection**: You can use "The Stage" to read the full context that the AI is using to ground its answers.

## ⚙️ Configuration

Enable RAG by setting the following in your `.env`:
```env
OPENAI_API_KEY=sk-... # Required for embeddings
DOCUMENTS_ROOT=./documents # Root path for indexed files
# Optional: Pinecone settings
PINECONE_API_KEY=...
PINECONE_INDEX=...
```

## 🔒 Safety & Security

- **Path-Traversal Protection**: The AI Gateway validates all document paths to ensure they stay within the `DOCUMENTS_ROOT`.
- **Cleanup**: If document ingestion fails midway, the gateway automatically cleans up partial database records and orphaned vectors.
