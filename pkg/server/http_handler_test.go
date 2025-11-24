package server

import (
	"net/http"
	"strings"
	"testing"

	"github.com/DemoHn/Zn/pkg/value"
)

func TestHTTPRequest_GetProperty(t *testing.T) {
	rawURL := "http://example.com/path?param1=value1&param2=value2"
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	httpRequest, _ := buildIncomingRequest(req)

	tests := []struct {
		propertyName string
		expected     string
	}{
		{"路径", "/path"},
		{"方法", "GET"},
	}

	for _, tt := range tests {
		t.Run(tt.propertyName, func(t *testing.T) {
			prop, err := httpRequest.GetProperty(tt.propertyName)
			if err != nil {
				t.Fatalf("Failed to get property: %v", err)
			}

			if prop.String() != tt.expected {
				t.Errorf("Expected %s to be '%s', but return %s", tt.propertyName, tt.expected, prop.String())
			}
		})
	}
}

func TestHTTPRequest_GetHeaders(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Custom-Header", "custom-value")
	req.Header.Add("x-lowercase-header", "lowercase-value")

	httpRequest, _ := buildIncomingRequest(req)

	prop, err := httpRequest.GetProperty("头部")
	if err != nil {
		t.Fatalf("Failed to get headers: %v", err)
	}

	headers := prop.(*value.HashMap).GetValue()
	if headers["Content-Type"].String() != "application/json" {
		t.Errorf("Expected Content-Type to be 'application/json', but got '%s'", headers["Content-Type"].String())
	}
	if headers["X-Custom-Header"].String() != "custom-value" {
		t.Errorf("Expected X-Custom-Header to be 'custom-value', but got '%s'", headers["X-Custom-Header"].String())
	}

	// NOTE: golang will uppercase the header key, so we need to check the value
	if headers["X-Lowercase-Header"].String() != "lowercase-value" {
		t.Errorf("Expected x-lowercase-header to be 'lowercase-value', but got '%s'", headers["x-lowercase-header"].String())
	}
}

func TestHTTPRequest_GetQueryParams(t *testing.T) {
	rawURL := "http://example.com/path?param1=value1&param2=value2"
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	httpRequest, _ := buildIncomingRequest(req)

	prop, err := httpRequest.GetProperty("查询参数")
	if err != nil {
		t.Fatalf("Failed to get query parameters: %v", err)
	}

	queryParams := prop.(*value.HashMap).GetValue()
	if queryParams["param1"].String() != "value1" {
		t.Errorf("Expected param1 to be 'value1', but got '%s'", queryParams["param1"].String())
	}
	if queryParams["param2"].String() != "value2" {
		t.Errorf("Expected param2 to be 'value2', but got '%s'", queryParams["param2"].String())
	}
}

func TestHTTPRequest_GetQueryParams_Array(t *testing.T) {
	rawURL := "http://example.com/path?param1=value1&param2=value2&param2=value3"
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	httpRequest, _ := buildIncomingRequest(req)

	prop, err := httpRequest.GetProperty("查询参数")
	if err != nil {
		t.Fatalf("Failed to get query parameters: %v", err)
	}

	queryParams := prop.(*value.HashMap).GetValue()
	if queryParams["param1"].String() != "value1" {
		t.Errorf("Expected param1 to be 'value1', but got '%s'", queryParams["param1"].String())
	}

	param2Array := queryParams["param2"].(*value.Array).GetValue()
	if len(param2Array) != 2 {
		t.Errorf("Expected param2 to have 2 values, but got %d", len(param2Array))
	} else {
		if param2Array[0].String() != "value2" {
			t.Errorf("Expected param2[0] to be 'value2', but got '%s'", param2Array[0].String())
		}
		if param2Array[1].String() != "value3" {
			t.Errorf("Expected param2[1] to be 'value3', but got '%s'", param2Array[1].String())
		}
	}
}

func TestHTTPRequest_ReadBody(t *testing.T) {
	body := "这是请求体内容"
	req, err := http.NewRequest("POST", "http://example.com", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	httpRequest, _ := buildIncomingRequest(req)

	result, err := httpRequest.ExecMethod("读取内容", nil)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	if result.String() != body {
		t.Errorf("Expected body to be '%s', but got '%s'", body, result.String())
	}
}

func TestHTTPRequest_ReadBody_JSON(t *testing.T) {
	body := `{"key1": "value1", "key2": 2, "key3": true, "key4": {"nestedKey": "nestedValue"}}`
	req, err := http.NewRequest("POST", "http://example.com", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	httpRequest, _ := buildIncomingRequest(req)

	result, err := httpRequest.ExecMethod("读取内容", nil)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	resultMap := result.(*value.HashMap).GetValue()
	if resultMap["key1"].String() != "value1" {
		t.Errorf("Expected key1 to be 'value1', but got '%s'", resultMap["key1"].String())
	}
	if resultMap["key2"].String() != "2" {
		t.Errorf("Expected key2 to be '2', but got '%s'", resultMap["key2"].String())
	}
	if resultMap["key3"].String() != "真" {
		t.Errorf("Expected key3 to be '真', but got '%s'", resultMap["key3"].String())
	}

	// Check nested JSON
	nestedMap := resultMap["key4"].(*value.HashMap).GetValue()
	if nestedMap["nestedKey"].String() != "nestedValue" {
		t.Errorf("Expected nestedKey to be 'nestedValue', but got '%s'", nestedMap["nestedKey"].String())
	}
}
