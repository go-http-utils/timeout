package timeout

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TimeoutSuite struct {
	suite.Suite

	server *httptest.Server
}

func (s *TimeoutSuite) SetupSuite() {
	mux := http.NewServeMux()

	mux.Handle("/hello", http.HandlerFunc(helloHandler))
	mux.Handle("/timeout", http.HandlerFunc(timeoutHandler))

	s.server = httptest.NewServer(Handler(mux, time.Second,
		DefaultTimeoutHandler))
}

func (s *TimeoutSuite) TestTimeout() {
	req, err := http.NewRequest(http.MethodGet, s.server.URL+"/timeout", nil)
	s.Nil(err)

	res, err := sendRequest(req)
	s.Nil(err)
	s.Equal(http.StatusGatewayTimeout, res.StatusCode)
	s.Equal([]byte("Service timeout"), getResRawBody(res))
}

func (s *TimeoutSuite) TestNotTimeout() {
	req, err := http.NewRequest(http.MethodGet, s.server.URL+"/hello", nil)
	s.Nil(err)

	res, err := sendRequest(req)
	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)
	s.Equal([]byte("Hello World"), getResRawBody(res))
}

func TestTimeout(t *testing.T) {
	suite.Run(t, new(TimeoutSuite))
}

func helloHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)

	res.Write([]byte("Hello World"))
}

func timeoutHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)

	res.Write([]byte("Hello World"))
	<-time.After(2 * time.Second)
}

func sendRequest(req *http.Request) (*http.Response, error) {
	cli := &http.Client{}
	return cli.Do(req)
}

func getResRawBody(res *http.Response) []byte {
	bytes, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	return bytes
}
