// Package timeseries provides a unified interface for time-series databases.
//
// Supported backends:
//   - AWS Timestream
//   - InfluxDB
//   - Memory (for testing)
//
// Usage:
//
//	ts, err := timestream.New(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer ts.Close()
//
//	point := &timeseries.Point{
//	    Measurement: "cpu_usage",
//	    Tags:        map[string]string{"host": "server-1"},
//	    Fields:      map[string]interface{}{"value": 85.5},
//	    Time:        time.Now(),
//	}
//
//	if err := ts.Write(ctx, point); err != nil {
//	    log.Error("failed to write point", err)
//	}
package timeseries
