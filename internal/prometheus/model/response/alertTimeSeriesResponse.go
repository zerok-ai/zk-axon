package response

import "sort"

type AlertTimeSeriesResponse struct {
	Data AlertRangeDataResponse `json:"data"`
}

type AlertRangeDataResponse struct {
	ResultType string           `json:"resultType"`
	Result     []ResultResponse `json:"result"`
}

type ResultResponse struct {
	Metric    Metric     `json:"metric"`
	Durations []Duration `json:"durations"`
}

type Duration struct {
	From int `json:"from"`
	To   int `json:"to"`
}

type AlertRangePrometheusResponse struct {
	Status string         `json:"status"`
	Data   AlertRangeData `json:"data"`
}

type AlertRangeData struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

type Result struct {
	Metric Metric  `json:"metric"`
	Values []Value `json:"values"`
}

type Metric struct {
	Name       string `json:"__name__"`
	AlertName  string `json:"alertname"`
	AlertState string `json:"alertstate"`
	Pod        string `json:"pod"`
	Severity   string `json:"severity"`
}

type Value []interface{}

func findSeries(arr []Value, step int) []Duration {
	var result []Duration

	if len(arr) == 0 {
		return result
	}

	// Sort the array based on the first element of each sub-array
	sort.Slice(arr, func(i, j int) bool {
		return int(arr[i][0].(float64)) < int(arr[j][0].(float64))
	})

	var start, end int

	for i := 0; i < len(arr); i++ {
		timestamp := int(arr[i][0].(float64))
		if i == 0 {
			start = timestamp
			end = timestamp
		} else {
			if timestamp-end == step {
				end = timestamp
			} else {
				result = append(result, Duration{From: start, To: end})
				start = timestamp
				end = timestamp
			}
		}
	}

	result = append(result, Duration{From: start, To: end})

	return result
}

func ConvertAlertRangePrometheusResponseToAlertTimeSeriesResponse(alertRangePrometheusResponse AlertRangePrometheusResponse, step int64) AlertTimeSeriesResponse {
	alertTimeSeriesResponse := AlertTimeSeriesResponse{}
	alertTimeSeriesResponse.Data.ResultType = alertRangePrometheusResponse.Data.ResultType
	alertTimeSeriesResponse.Data.Result = make([]ResultResponse, 0)
	for _, result := range alertRangePrometheusResponse.Data.Result {
		resultResponse := ResultResponse{}
		resultResponse.Metric = result.Metric
		resultResponse.Durations = make([]Duration, 0)
		resultResponse.Durations = findSeries(result.Values, int(step))
		alertTimeSeriesResponse.Data.Result = append(alertTimeSeriesResponse.Data.Result, resultResponse)
	}
	return alertTimeSeriesResponse
}
