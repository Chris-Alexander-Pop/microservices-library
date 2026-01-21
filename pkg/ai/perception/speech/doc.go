// Package speech provides interfaces and adapters for Speech-to-Text (STT) and Text-to-Speech (TTS).
//
// This package unifies various speech processing services under a common interface,
// allowing for easy switching between providers.
//
// Supported capabilities:
//   - Synthesize: Convert text to audio
//   - Transcribe: Convert audio to text
//
// Supported backends:
//   - Memory: In-memory mock for testing
//   - AWS Polly/Transcribe: (Planned)
//   - Google Cloud Text-to-Speech/Speech-to-Text: (Planned)
//
// Basic usage:
//
//	synthesizer := memory.New()
//	audio, err := synthesizer.Synthesize(ctx, "Hello world")
package speech
