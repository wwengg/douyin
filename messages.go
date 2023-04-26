package main

import (
	"douyin/regexputil"
	"fmt"
	"strconv"
	"strings"
)

// HTTPHeader represents a http header
type HTTPHeader struct {
	Name  string `bson:"name" json:"name"`
	Value string `bson:"value" json:"value"`
}

// HTTPRequest represents http request
type HTTPRequest struct {
	Method    string        `bson:"method" json:"method"`
	URI       string        `bson:"uri" json:"uri"`
	Version   string        `bson:"version" json:"version"`
	Timestamp int64         `bson:"timestamp" json:"timestamp"`
	ClientIP  string        `bson:"clientip" json:"clientip"`
	Headers   *[]HTTPHeader `bson:"headers" json:"headers"`
	Body      []byte        `bson:"body" json:"body"`
}

// HTTPResponse represents http response
type HTTPResponse struct {
	StatusCode int           `bson:"statuscode" json:"statuscode"`
	StatusText string        `bson:"statustext" json:"statustext"`
	Version    string        `bson:"version" json:"version"`
	Timestamp  int64         `bson:"timestamp" json:"timestamp"`
	Headers    *[]HTTPHeader `bson:"headers" json:"headers"`
	Body       []byte        `bson:"body" json:"body"`
}

// HTTPMessage represents a message including the request and response
type HTTPMessage struct {
	ID       string        `bson:"_id" json:"id"`
	Request  *HTTPRequest  `bson:"request" json:"request"`
	Response *HTTPResponse `bson:"response" json:"response"`
}

// HTTPMessageSummary represents minimal set of message attributes
type HTTPMessageSummary struct {
	ID         string `bson:"_id" json:"id"`
	Timestamp  int64  `bson:"timestamp" json:"timestamp"`
	Method     string `bson:"method" json:"method"`
	URI        string `bson:"uri" json:"uri"`
	StatusCode int    `bson:"statuscode" json:"statuscode"`
}

// unmarshalHTTPRequest deserializes bytes to HTTPRequest
func unmarshalHTTPRequest(data []byte) (request *HTTPRequest) {
	requestLines := strings.Split(string(data), "\r\n")

	match := regexputil.RegexMapString("^(?P<method>[^ ]+) (?P<uri>[^ ]+) (?P<version>.+)$", requestLines[0])
	if match != nil {
		result := HTTPRequest{}
		result.Method = (*match)["method"]
		result.URI = (*match)["uri"]
		result.Version = (*match)["version"]

		headers := []HTTPHeader{}

		for _, line := range requestLines[1:] {
			if line == "" {
				// result.Body = strings.Join(requestLines[i+2:], "\r\n")
				break
			}
			match = regexputil.RegexMapString("^(?P<name>[^:]+): (?P<value>.+)$", line)
			if match != nil {
				headers = append(headers, HTTPHeader{(*match)["name"], (*match)["value"]})
			}
		}

		result.Headers = &headers
		request = &result
	}

	return
}

func unmarshalHTTPResponse(data []byte) (response *HTTPResponse) {
	responseLines := strings.Split(string(data), "\r\n")

	match := regexputil.RegexMapString("^(?P<version>[^ ]+) (?P<statuscode>[^ ]+) (?P<status>.+)$", responseLines[0])
	if match != nil {
		result := HTTPResponse{}
		statusCode, _ := strconv.ParseInt((*match)["statuscode"], 10, 32)
		result.StatusCode = int(statusCode)
		result.StatusText = (*match)["status"]
		result.Version = (*match)["version"]

		headers := []HTTPHeader{}

		for _, line := range responseLines[1:] {
			if line == "" {
				// result.Body = strings.Join(responseLines[i+2:], "\r\n")
				break
			}
			match = regexputil.RegexMapString("^(?P<name>[^:]+): (?P<value>.+)$", line)
			if match != nil {
				headers = append(headers, HTTPHeader{(*match)["name"], (*match)["value"]})
			}
		}

		result.Headers = &headers
		response = &result
	}

	return
}

func getMessageSummary(messages []HTTPMessage) []HTTPMessageSummary {
	summary := make([]HTTPMessageSummary, len(messages))
	for i, message := range messages {
		summary[i] = summariseMessage(message)
	}
	return summary
}

func summariseMessage(message HTTPMessage) (summary HTTPMessageSummary) {
	summary = HTTPMessageSummary{
		message.ID,
		message.Request.Timestamp,
		message.Request.Method,
		message.Request.URI,
		0,
	}

	if message.Response != nil {
		summary.StatusCode = message.Response.StatusCode
	}

	if !strings.HasPrefix(strings.ToLower(summary.URI), "http:") {
		for _, header := range *message.Request.Headers {
			if strings.ToLower(header.Name) == "host" {
				summary.URI = fmt.Sprintf("https://%s%s", header.Value, summary.URI)
				break
			}
		}
	}

	return
}
