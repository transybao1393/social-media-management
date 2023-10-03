package http

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type CustomRequest struct {
	MethodName    string
	PathURL       string
	CustomHeaders map[string]string
	AccessToken   string //- optional
	Body          []byte
	Authorization string //- optional
}

type Response struct {
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    uint   `json:"expires_in"`
}

// - Define your Error struct
// - TODO: Need to move CustomError to somewhere else
type CustomError struct {
	CustomMessage string
	DefaultError  error
}

// - Create a function Error() string and associate it to the struct.
func (err *CustomError) Error() string {
	customErrorMessage := fmt.Sprintf("Message %[1]s with error %[2]s", err.CustomMessage, err.DefaultError)
	return customErrorMessage
}

func (ri *CustomRequest) Exec() ([]byte, error) {
	//- request header set
	var req *http.Request
	var err error
	if ri.MethodName == "POST" || ri.MethodName == "PUT" || ri.MethodName == "PATCH" {
		req, err = http.NewRequest(ri.MethodName, ri.PathURL, bytes.NewBuffer([]byte(ri.Body)))
	} else {
		req, err = http.NewRequest(ri.MethodName, ri.PathURL, nil)
	}

	if err != nil {
		return nil, &CustomError{
			CustomMessage: "Error when creating new request from data input",
			DefaultError:  err,
		}
	}

	if ri.AccessToken != "" {
		//- default header
		finalAccessToken := fmt.Sprintf("Bearer %[1]s", ri.AccessToken)
		req.Header.Add("Authorization", finalAccessToken)
	}

	//- Custom response header
	for key, value := range ri.CustomHeaders {
		req.Header.Add(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &CustomError{
			CustomMessage: "Error when creating response from client request",
			DefaultError:  err,
		}
	}
	defer resp.Body.Close()
	fmt.Printf("status code %d %s \n", resp.StatusCode, resp.Status)

	/**
	Informational responses (100 – 199)
	Successful responses (200 – 299)
	Redirection messages (300 – 399)
	Client error responses (400 – 499)
	Server error responses (500 – 599)
	*/
	//- Success
	//- write test cases to check response status code
	fmt.Printf("Status code %d, response message %s\n", resp.StatusCode, resp.Status)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		//- Body decode to string
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, &CustomError{
				CustomMessage: "Error when convert Body response to bytes",
				DefaultError:  err,
			}
		}
		return bodyBytes, nil
	}

	//- Error
	errorResponseBodyByteData, err := io.ReadAll(resp.Body)
	fmt.Printf("errorResponseBodyByteData %v, error %v \n", string(errorResponseBodyByteData), err)
	return nil, NewError(string(errorResponseBodyByteData))
}

type PathURL struct {
	QueryParams       map[string]string
	PathParams        map[string]string
	queryParamsString string
	pathParamsString  string
	APIDomain         string
	APIURI            string
}

func (pu PathURL) New() *PathURL {
	//- TODO: validate map[string] with support array
	//- convert map[string]string to string
	pu.pathParamsString = createPathParamsKeyValuePairs(pu.PathParams, "/")
	pu.queryParamsString = createQueryParamsKeyValuePairs(pu.QueryParams, "&")
	return &pu
}

func (pu *PathURL) Build() string {
	apiDomain := pu.APIDomain
	apiURI := pu.APIURI
	// var finalURL string
	finalURL := apiDomain + apiURI
	//- handle path paramters and query parameters
	if len(pu.PathParams) > 0 {
		finalURL = apiDomain + apiURI + "/" + pu.pathParamsString
	}
	if len(pu.QueryParams) > 0 {
		finalURL = apiDomain + apiURI + "?" + pu.queryParamsString
	}
	if len(pu.PathParams) > 0 && len(pu.QueryParams) > 0 {
		finalURL = apiDomain + apiURI + "/" + pu.pathParamsString + "?" + pu.queryParamsString
	}
	fmt.Printf("%s \n", finalURL)
	return finalURL
}

func createPathParamsKeyValuePairs(m map[string]string, symbol string) string {
	b := new(bytes.Buffer)
	for _, value := range m {
		fmt.Fprintf(b, "%[1]s\"%[2]s", value, symbol)
		//- replace string
		//- TODO: Handling use case /crm/v3/objects/calls/{callId}/something/else
	}
	return urlSpecialCharacterCleanUp(trimSuffix(b.String(), symbol))
}

func createQueryParamsKeyValuePairs(m map[string]string, symbol string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%[1]s=\"%[2]s\"%[3]s", key, value, symbol)
	}
	return urlSpecialCharacterCleanUp(trimSuffix(b.String(), symbol))
}

func trimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func urlSpecialCharacterCleanUp(url string) string {
	re, err := regexp.Compile(`"`)
	if err != nil {
		log.Fatal(err)
	}
	newStr := re.ReplaceAllString(url, "")
	return newStr
}
