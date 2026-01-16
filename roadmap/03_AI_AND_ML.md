# AI & Machine Learning Roadmap

## LLM Integration (pkg/ai)
- [ ] **Multi-Provider Interface**: Unified `GenerateText`, `GenerateEmbedding` interfaces.
- [ ] **OpenAI**: GPT-4, Embeddings, Functions (Tools) calling support.
- [ ] **Anthropic**: Claude 3 API integration.
- [ ] **Google Gemini**: Vertex AI / AI Studio client integration.
- [ ] **Mistral**: API client.
- [ ] **Local Inference**:
    - [ ] **Ollama**: API wrapper for running Llama 3/Mistral locally.
    - [ ] **LocalAI**: Drop-in OpenAI replacement.

## Agents & Chains
- [ ] **ChainBuilder**: Middleware-style chain construction (Prompt -> LLM -> Parser).
- [ ] **Memory**: Conversation history management (Redis/Postgres backed).
- [ ] **Tools**: Interface for defining Go functions callable by LLMs.
    - [ ] Web Search Tool
    - [ ] Database Query Tool
    - [ ] API Caller Tool
- [ ] **RAG Engine**:
    - [ ] Document Ingestion pipeline (PDF/HTML -> Text).
    - [ ] Chunking strategies (Sentence, Paragraph, Semantic).
    - [ ] Embedding generation + Vector Store upsert.

## Computer Vision
- [ ] **OpenCV**: Go bindings (`gocv`) integration for image processing.
- [ ] **AWS Rekognition**: Image analysis wrapper.
- [ ] **Google Vision API**: OCR and labeling wrapper.

## Audio / Speech
- [ ] **Speech-to-Text (STT)**: Whisper (OpenAI/Local) integration.
- [ ] **Text-to-Speech (TTS)**: ElevenLabs and OpenAI TTS integration.
