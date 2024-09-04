package ext

import (
	"net/http"
	"strings"
	"testing"

	"fmt"
	r "github.com/DemoHn/Zn/pkg/runtime"
	"github.com/DemoHn/Zn/pkg/value"
)

func TestHTTPRequest_GetProperty(t *testing.T) {
	rawURL := "http://example.com/path?param1=value1&param2=value2"
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	httpRequest := NewHTTPRequest(req)
	ctx := r.NewContext(map[string]r.Element{}, r.NewMainModule(nil))

	tests := []struct {
		propertyName string
		expected     string
	}{
		{"路径", "「/path」"},
		{"方法", "「GET」"},
	}

	for _, tt := range tests {
		t.Run(tt.propertyName, func(t *testing.T) {
			prop, err := httpRequest.GetProperty(ctx, tt.propertyName)
			if err != nil {
				t.Fatalf("Failed to get property: %v", err)
			}

			if value.StringifyValue(prop) != tt.expected {
				t.Errorf("Expected %s, but got %s", tt.expected, value.StringifyValue(prop))
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

	httpRequest := NewHTTPRequest(req)
	ctx := r.NewContext(map[string]r.Element{}, r.NewMainModule(nil))

	prop, err := httpRequest.GetProperty(ctx, "头部")
	if err != nil {
		t.Fatalf("Failed to get headers: %v", err)
	}

	headers := prop.(*value.HashMap).GetValue()
	if value.StringifyValue(headers["Content-Type"]) != "「application/json」" {
		t.Errorf("Expected Content-Type to be 「application/json」, but got %s", value.StringifyValue(headers["Content-Type"]))
	}
	if value.StringifyValue(headers["X-Custom-Header"]) != "「custom-value」" {
		t.Errorf("Expected X-Custom-Header to be 「custom-value」, but got %s", value.StringifyValue(headers["X-Custom-Header"]))
	}
}

func TestHTTPRequest_GetQueryParams(t *testing.T) {
	rawURL := "http://example.com/path?param1=value1&param2=value2"
	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	httpRequest := NewHTTPRequest(req)
	ctx := r.NewContext(map[string]r.Element{}, r.NewMainModule(nil))

	prop, err := httpRequest.GetProperty(ctx, "查询参数")
	if err != nil {
		t.Fatalf("Failed to get query parameters: %v", err)
	}

	queryParams := prop.(*value.HashMap).GetValue()
	if value.StringifyValue(queryParams["param1"]) != "「value1」" {
		t.Errorf("Expected param1 to be 「value1」, but got %s", value.StringifyValue(queryParams["param1"]))
	}
	if value.StringifyValue(queryParams["param2"]) != "「value2」" {
		t.Errorf("Expected param2 to be 「value2」, but got %s", value.StringifyValue(queryParams["param2"]))
	}
}

func TestHTTPRequest_ReadBody(t *testing.T) {
	body := "这是请求体内容"
	req, err := http.NewRequest("POST", "http://example.com", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	httpRequest := NewHTTPRequest(req)
	ctx := r.NewContext(map[string]r.Element{}, r.NewMainModule(nil))

	result, err := httpRequest.ExecMethod(ctx, "读取内容", nil)
	if err != nil {
		t.Fatalf("Failed to read body: %v", err)
	}

	if value.StringifyValue(result) != fmt.Sprintf("「%s」", body) {
		t.Errorf("Expected body to be %s, but got %s", body, value.StringifyValue(result))
	}
}
