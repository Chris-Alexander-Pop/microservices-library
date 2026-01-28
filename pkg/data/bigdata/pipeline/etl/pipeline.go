package etl

import "context"

// Pipeline represents a data pipeline.
type Pipeline interface {
	Run(ctx context.Context) error
}

// Extractor reads data from a source.
type Extractor interface {
	Extract(ctx context.Context, out chan<- interface{}) error
}

// Transformer processes data.
type Transformer interface {
	Transform(ctx context.Context, in <-chan interface{}, out chan<- interface{}) error
}

// Loader writes data to a destination.
type Loader interface {
	Load(ctx context.Context, in <-chan interface{}) error
}

// SimplePipeline is a linear E -> T -> L pipeline.
type SimplePipeline struct {
	E Extractor
	T Transformer
	L Loader
}

func (p *SimplePipeline) Run(ctx context.Context) error {
	dataCh1 := make(chan interface{})
	dataCh2 := make(chan interface{})

	errCh := make(chan error, 3)

	// Extract
	go func() {
		defer close(dataCh1)
		if err := p.E.Extract(ctx, dataCh1); err != nil {
			errCh <- err
		}
	}()

	// Transform
	go func() {
		defer close(dataCh2)
		if err := p.T.Transform(ctx, dataCh1, dataCh2); err != nil {
			errCh <- err
		}
	}()

	// Load
	if err := p.L.Load(ctx, dataCh2); err != nil {
		return err
	}

	close(errCh)
	// Check for any async errors
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}
