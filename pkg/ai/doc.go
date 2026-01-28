/*
Package ai provides artificial intelligence and machine learning capabilities.

This package organizes AI functionality into the following subdomains:

  - genai: Generative AI (LLMs, image generation, agents)
  - ml: Machine Learning (training, inference, feature stores)
  - nlp: Natural Language Processing (embeddings, RAG)
  - perception: Computer vision, speech, OCR

# Embedding Generation

For text embeddings, use nlp/embedding for focused embedding tasks:

	import "github.com/chris-alexander-pop/system-design-library/pkg/ai/nlp/embedding"
	vectors, err := embedder.Embed(ctx, texts)

For LLM chat/generation with optional embedding capabilities, use genai/llm:

	import "github.com/chris-alexander-pop/system-design-library/pkg/ai/genai/llm"
	resp, err := client.Chat(ctx, messages)
*/
package ai
