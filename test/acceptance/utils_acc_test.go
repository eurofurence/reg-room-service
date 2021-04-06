package acceptance

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func tstPerformGet(relativeUrlWithLeadingSlash string, token string) *http.Response {
	request, err := http.NewRequest(http.MethodGet, ts.URL+relativeUrlWithLeadingSlash, nil)
	if err != nil {
		log.Fatal(err)
	}
	if token != "" {
		request.Header.Set("authorization", "Bearer "+token)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	return response
}

func tstBodyToString(response *http.Response) string {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = response.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

// tip: dto := &attendee.AttendeeDto{}
func tstParseJson(body string, dto interface{}) {
	err := json.Unmarshal([]byte(body), dto)
	if err != nil {
		log.Fatal(err)
	}
}
