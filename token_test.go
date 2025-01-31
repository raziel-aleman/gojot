package gojot

import (
	"reflect"
	"testing"
	"time"
)

func assertEqual(t *testing.T, actual interface{}, expected interface{}) {
	t.Helper()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Assertion failed: %v != %v", actual, expected)
	}
}

func TestGenerateToken(t *testing.T) {
	user := User{ID: uint64(123456), Name: "Test", Email: "test@email.com"}
	_, err := GenerateToken(user, []byte("your-secret-key"), time.Hour*24)
	if err != nil {
		t.Fatalf("Error occurred while processing token: %v", err)
	}
}

func TestValidateToken(t *testing.T) {

	// Assuming you have a user ID from authentication
	user := User{ID: uint64(123456), Name: "Test", Email: "test@email.com"}

	// Generate the JWT token
	token, err := GenerateToken(user, []byte("your-secret-key"), time.Hour*24)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	claims, err := ValidateToken(token, []byte("your-secret-key"))
	if err != nil {
		t.Fatalf("Error occurred while processing token: %v", err)
	}

	got := claims.User
	expected := User{ID: uint64(123456), Name: "Test", Email: "test@email.com"}
	assertEqual(t, got, expected)
}
