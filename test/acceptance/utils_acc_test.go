package acceptance

import (
	"log"
	"net/http"
)

func tstPerformGet(relativeUrlWithLeadingSlash string) *http.Response {
	request, err := http.NewRequest(http.MethodGet, ts.URL+relativeUrlWithLeadingSlash, nil)
	if err != nil {
		log.Fatal(err)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return response
}
