package response

type AlertWebhookResponse struct {
	Version           string                 `json:"version"`
	GroupKey          string                 `json:"groupKey"`
	TruncatedAlerts   int                    `json:"truncatedAlerts"`
	Status            string                 `json:"status"`
	Receiver          string                 `json:"receiver"`
	GroupLabels       map[string]interface{} `json:"groupLabels"`
	CommonLabels      map[string]interface{} `json:"commonLabels"`
	CommonAnnotations map[string]interface{} `json:"commonAnnotations"`
	ExternalURL       string                 `json:"externalURL"`
	Alerts            []Alert                `json:"alerts"`
}

type AlertWebhook struct {
	Status       string                 `json:"status"`
	Labels       map[string]interface{} `json:"labels"`
	Annotations  map[string]interface{} `json:"annotations"`
	StartsAt     string                 `json:"startsAt"`
	EndsAt       string                 `json:"endsAt"`
	GeneratorURL string                 `json:"generatorURL"`
	Fingerprint  string                 `json:"fingerprint"`
}
