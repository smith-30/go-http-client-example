package restclient

import (
	"io"
	"io/ioutil"
	"net/http"
)

func readAndCloseResponseBody(resp *http.Response) {
	if resp == nil {
		return
	}

	// Ensure the response body is fully read and closed
	// before we reconnect, so that we reuse the same TCP
	// connection.
	const maxBodySlurpSize = 2 << 10
	defer resp.Body.Close()

	if resp.ContentLength <= maxBodySlurpSize {
		io.Copy(ioutil.Discard, &io.LimitedReader{R: resp.Body, N: maxBodySlurpSize})
	}
}
