package restclient

type RESTClient struct {
	// base is the root URL for all invocations of the client
	base *url.URL

	// versionedAPIPath is a path segment connecting the base URL to the resource root
	versionedAPIPath string

	// content describes how a RESTClient encodes and decodes responses.
	content ClientContentConfig

	// Todo

	// creates BackoffManager that is passed to requests.
	// createBackoffMgr func() BackoffManager

	// Todo

	// rateLimiter is shared among all requests created by this client unless specifically
	// overridden.
	// rateLimiter flowcontrol.RateLimiter

	// Todo

	// warningHandler is shared among all requests created by this client.
	// If not set, defaultWarningHandler is used.
	// warningHandler WarningHandler
	
	// Set specific behavior of the client.  If not set http.DefaultClient will be used.
	Client *http.Client
}

type Interface interface {
	Verb(verb string) *Request
	Post() *Request
	Put() *Request
	Patch() *Request
	Get() *Request
	Delete() *Request
}

// NewRESTClient creates a new RESTClient. This client performs generic REST functions
// such as Get, Put, Post, and Delete on specified paths.
func NewRESTClient(baseURL *url.URL, versionedAPIPath string, config ClientContentConfig, rateLimiter flowcontrol.RateLimiter, client *http.Client) (*RESTClient, error) {
	if len(config.ContentType) == 0 {
		config.ContentType = "application/json"
	}



	base := *baseURL
	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}
	base.RawQuery = ""
	base.Fragment = ""

	return &RESTClient{
		base:             &base,
		versionedAPIPath: versionedAPIPath,
		content:          config,
		createBackoffMgr: readExpBackoffConfig,
		rateLimiter:      rateLimiter,

		Client: client,
	}, nil
}


func (c *RESTClient) Verb(verb string) *Request {
	return NewRequest(c).Verb(verb)
}

func (c *RESTClient) Post() *Request {
	return c.Verb(http.MethodPost))
}

func (c *RESTClient) Put() *Request {
	return c.Verb(http.MethodPut)
}

func (c *RESTClient) Patch(pt types.PatchType) *Request {
	return c.Verb(http.MethodPatch)
}

func (c *RESTClient) Get() *Request {
	return c.Verb(http.MethodGet)
}

func (c *RESTClient) Delete() *Request {
	return c.Verb(http.MethodDelete)
}