package strapi

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/charisworks/charisworks-service-go/util"
)

func requestToStrapi(method httpMethod, path string, reqBody []byte) (*http.Response, error) {
	req := &http.Request{}
	if method == http.MethodGet || method == http.MethodDelete {
		req, _ = http.NewRequest(string(method), util.STRAPI_URL+"/api"+path, nil)
	} else {
		req, _ = http.NewRequest(string(method), util.STRAPI_URL+"/api"+path, bytes.NewBuffer(reqBody))
	}
	util.Logger(
		fmt.Sprintf(
			`
			**********************************************************************************************
			Requesting to Strapi... 
			method: %s
			path: %s
			reqBody: %s
			util.STRAPI_URL %s
			**********************************************************************************************`,
			string(method), // Convert method to string
			path,
			string(reqBody),
			util.STRAPI_URL+"/api"+path,
		),
	)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+util.STRAPI_JWT)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type httpMethod string

const (
	GET    httpMethod = "GET"
	POST   httpMethod = "POST"
	PUT    httpMethod = "PUT"
	DELETE httpMethod = "DELETE"
)
