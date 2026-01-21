// Package vision provides interfaces and adapters for computer vision tasks.
//
// This package supports image classification, object detection, and facial recognition
// through a unified interface.
//
// Supported capabilities:
//   - Analyze: generic image analysis
//   - DetectFaces: facial detection
//   - DetectLabels: object/scene classification
//
// Supported backends:
//   - Memory: In-memory mock for testing
//   - AWS Rekognition: (Planned)
//   - Google Cloud Vision: (Planned)
//
// Basic usage:
//
//	analyzer := memory.New()
//	result, err := analyzer.Analyze(ctx, imageData)
package vision
