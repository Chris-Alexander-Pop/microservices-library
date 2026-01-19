package mapreduce

import (
	"context"
	"sync"
)

// Mapper function: Key, Value -> []KeyValue
type Mapper func(key string, value interface{}, out chan<- KeyValue)

// Reducer function: Key, []Value -> []Value
type Reducer func(key string, values []interface{}, out chan<- interface{})

type KeyValue struct {
	Key   string
	Value interface{}
}

// Job represents a MR job.
type Job struct {
	Mapper  Mapper
	Reducer Reducer
	Inputs  map[string]interface{}
	// Shards/Workers
	NumWorkers int
}

func NewJob(m Mapper, r Reducer, inputs map[string]interface{}, workers int) *Job {
	return &Job{
		Mapper:     m,
		Reducer:    r,
		Inputs:     inputs,
		NumWorkers: workers,
	}
}

func (j *Job) Run(ctx context.Context) (map[string][]interface{}, error) {
	// Map Phase
	mapResults := make(chan KeyValue, len(j.Inputs)*2) // Buffer estimation

	var wg sync.WaitGroup
	inputCh := make(chan KeyValue, len(j.Inputs))

	// Input feeder
	go func() {
		for k, v := range j.Inputs {
			inputCh <- KeyValue{Key: k, Value: v}
		}
		close(inputCh)
	}()

	// Workers
	for i := 0; i < j.NumWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for kv := range inputCh {
				j.Mapper(kv.Key, kv.Value, mapResults)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(mapResults)
	}()

	// Shuffle / Group
	shuffled := make(map[string][]interface{})
	for res := range mapResults {
		shuffled[res.Key] = append(shuffled[res.Key], res.Value)
	}

	// Reduce Phase (Simplification: Sequential Reduce for now)
	results := make(map[string][]interface{})
	for k, vals := range shuffled {
		outCh := make(chan interface{})
		var reduceOutput []interface{}

		go func() {
			defer close(outCh)
			j.Reducer(k, vals, outCh)
		}()

		for out := range outCh {
			reduceOutput = append(reduceOutput, out)
		}
		results[k] = reduceOutput
	}

	return results, nil
}
