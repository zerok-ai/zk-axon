package response

import "github.com/prometheus/common/model"

func ConvertMetricToPodUsage(metric model.Matrix) Usage {
	var usage Usage = make(map[ContainerName]PlotValues)
	var plotValues PlotValues

	for _, series := range metric {
		var values []float64
		var timeStamps []int64

		for _, value := range series.Values {
			values = append(values, float64(value.Value))
			timeStamps = append(timeStamps, int64(value.Timestamp))
		}

		plotValues.Values = values
		plotValues.TimeStamp = timeStamps

		// Add to container name
		containerName := ContainerName(series.Metric["container"])
		usage[containerName] = plotValues
	}

	return usage
}
