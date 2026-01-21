// Package ocr provides interfaces and adapters for Optical Character Recognition.
//
// This package enables extracting text from images and documents using various
// backend providers. It supports structured data extraction and raw text extraction.
//
// Supported backends:
//   - Memory: In-memory mock for testing
//   - AWS Textract: (Planned)
//   - Google Cloud Vision: (Planned)
//   - Azure Computer Vision: (Planned)
//
// Basic usage:
//
//	reader := memory.New()
//	text, err := reader.ReadImage(ctx, imageData)
package ocr
