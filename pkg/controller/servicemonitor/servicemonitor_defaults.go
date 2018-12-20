package servicemonitor

type endpointValues struct {
	port            string
	path            string
	scheme          string
	scrapeInterval  string
	bearerTokenFile string
}

func (e *endpointValues) Port(definedValue string) string {

	if definedValue != "" {
		return definedValue
	}
	return "8080"
}

func (e *endpointValues) Path(definedValue string) string {

	if definedValue != "" {
		return definedValue
	}
	return "/metrics"
}

func (e *endpointValues) Scheme(definedValue string) string {

	if definedValue != "" {
		return definedValue
	}
	return "http"
}

func (e *endpointValues) ScrapeInterval(definedValue string) string {

	if definedValue != "" {
		return definedValue
	}
	return "30s"
}

func (e *endpointValues) BearerTokenFile(definedValue string) string {

	if definedValue != "" {
		return definedValue
	}
	return "/var/run/secrets/kubernetes.io/serviceaccount/token"
}
