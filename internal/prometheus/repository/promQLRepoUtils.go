package repository

import (
	"axon/internal/prometheus/model/request"
	"context"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	logger "github.com/zerok-ai/zk-utils-go/logs"
	"os"
	"reflect"
	"strings"
	"text/template"
	"time"
)

func GetPromQueryString(templateString string, podInfoReq request.PromRequestMeta) (string, error) {
	// Create a PromQL query
	queryTemplate, err := template.New("query").Parse(templateString)
	if err != nil {
		logger.Fatal(err.Error())
	}
	query := new(strings.Builder)
	err = queryTemplate.Execute(query, podInfoReq)

	// Query Prometheus
	logger.Debug(LogTag, "query: ", query.String())
	logger.Debug(LogTag, "over: "+podInfoReq.StartTime.String()+" to "+podInfoReq.EndTime.String())
	return query.String(), nil
}

func (r promQLRepo) GetGenericQuery(query string, startTime time.Time, endTime time.Time) (model.Matrix, error) {
	// Execute the query
	ctx := context.Background()
	result, warnings, err := r.queryAPI.QueryRange(ctx, query, v1.Range{
		Start: startTime,
		End:   endTime,
		Step:  1 * time.Minute, // Adjust the step as needed
	})
	if err != nil {
		logger.Error(LogTag, os.Stderr, "Error executing query: %v\n", err)
		return nil, err
	}

	// Check for query warnings
	if len(warnings) > 0 {
		logger.Warn(LogTag, "Query warnings:\n")
		for _, warning := range warnings {
			logger.Warn(LogTag, "%s\n", warning)
		}
	}

	logger.Debug(LogTag, "Result type: ", reflect.TypeOf(result).Name())

	reflect.TypeOf(result).Name()
	// Process query result
	if matrix, ok := result.(model.Matrix); ok {
		return matrix, nil
	} else {
		logger.Debug(LogTag, "Query did not return a matrix\n")
	}

	return nil, nil
}

func (r promQLRepo) GetPromData(query string, startTime time.Time, endTime time.Time, duration time.Duration, step time.Duration) (interface{}, string, error) {
	// Execute the query
	ctx := context.Background()
	var result model.Value
	var warnings v1.Warnings
	var err error
	if duration == 0 {
		result, warnings, err = r.queryAPI.Query(ctx, query, endTime)
	} else {
		result, warnings, err = r.queryAPI.QueryRange(ctx, query, v1.Range{
			Start: startTime,
			End:   endTime,
			Step:  step,
		})
	}
	if err != nil {
		logger.Error(LogTag, os.Stderr, "Error executing query: %v\n", err)
		return nil, "", err
	}

	// Check for query warnings
	if len(warnings) > 0 {
		logger.Warn(LogTag, "Query warnings:\n")
		for _, warning := range warnings {
			logger.Warn(LogTag, "%s\n", warning)
		}
	}

	resultType := reflect.TypeOf(result).Name()
	logger.Debug(LogTag, "Result type: ", resultType)
	return result, resultType, nil
}
func (r promQLRepo) GetPromMatrixData(query string, startTime time.Time, endTime time.Time, step time.Duration) (model.Matrix, error) {
	// Execute the query
	ctx := context.Background()
	result, warnings, err := r.queryAPI.QueryRange(ctx, query, v1.Range{
		Start: startTime,
		End:   endTime,
		Step:  step,
	})
	if err != nil {
		logger.Error(LogTag, os.Stderr, "Error executing query: %v\n", err)
		return nil, err
	}

	// Check for query warnings
	if len(warnings) > 0 {
		logger.Warn(LogTag, "Query warnings:\n")
		for _, warning := range warnings {
			logger.Warn(LogTag, "%s\n", warning)
		}
	}

	logger.Debug(LogTag, "Result type: ", reflect.TypeOf(result).Name())

	// Process query result
	if matrix, ok := result.(model.Matrix); ok {
		return matrix, nil
	} else {
		logger.Debug(LogTag, "Query did not return a matrix\n")
	}

	return model.Matrix{}, nil
}

func (r promQLRepo) GetPromVectorData(query string, endTime time.Time) (model.Vector, error) {
	// Execute the query
	ctx := context.Background()
	result, warnings, err := r.queryAPI.Query(ctx, query, endTime)
	if err != nil {
		logger.Error(LogTag, os.Stderr, "Error executing query: %v\n", err)
		return nil, err
	}

	// Check for query warnings
	if len(warnings) > 0 {
		logger.Warn(LogTag, "Query warnings:\n")
		for _, warning := range warnings {
			logger.Warn(LogTag, "%s\n", warning)
		}
	}

	logger.Debug(LogTag, "Result type: ", reflect.TypeOf(result).Name())

	// Process query result
	if vector, ok := result.(model.Vector); ok {
		return vector, nil
	} else {
		logger.Debug(LogTag, "Query did not return a Vector\n")
	}

	return model.Vector{}, nil
}
