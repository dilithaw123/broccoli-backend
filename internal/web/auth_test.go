package web

import "testing"

func TestBadToken(t *testing.T) {
	tokenStr := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE3MjY3OTg2MDl9.E6_L9orwKMskkm7f9ikYaF0xNXOYIi9iTTXZy9peL4w"
	secretKey := "secretkey"
	_, val := ParseAndValidateToken(tokenStr, secretKey)
	if val {
		t.Error("Token should be invalid")
	}
}
