package middleware

type LoginRequireConfig struct {
	ExcludedPath map[string]struct{}
}

const (
	LoginRequiredURLFmt = "%s:%s"
	OpenAPIPrefix       = "/open-api"
)
