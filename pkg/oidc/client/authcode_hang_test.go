package client

import (
	"context"
	"net"
	"testing"

	"github.com/int128/kubelogin/pkg/testing/logger"
)

func TestValidateBindAddresses_AllPortsBlocked(t *testing.T) {
	// Create test client
	c := &client{
		logger: logger.New(t),
	}

	// Block all test ports
	listener1, err := net.Listen("tcp", "127.0.0.1:18001")
	if err != nil {
		t.Skipf("Cannot bind to port 18001: %v", err)
	}
	defer listener1.Close()

	listener2, err := net.Listen("tcp", "127.0.0.1:18002")
	if err != nil {
		t.Skipf("Cannot bind to port 18002: %v", err)
	}
	defer listener2.Close()

	// Test validation with all ports blocked
	addresses := []string{"127.0.0.1:18001", "127.0.0.1:18002"}
	err = c.validateBindAddresses(addresses)
	
	if err == nil {
		t.Error("Expected validation to fail when all ports are blocked")
	}
	
	if !contains(err.Error(), "all bind addresses are unavailable") {
		t.Errorf("Expected error message about unavailable bind addresses, got: %v", err)
	}
}

func TestValidateBindAddresses_SomePortsAvailable(t *testing.T) {
	// Create test client
	c := &client{
		logger: logger.New(t),
	}

	// Block only one port
	listener1, err := net.Listen("tcp", "127.0.0.1:18003")
	if err != nil {
		t.Skipf("Cannot bind to port 18003: %v", err)
	}
	defer listener1.Close()

	// Test validation with some ports available
	addresses := []string{"127.0.0.1:18003", "127.0.0.1:18004"}
	err = c.validateBindAddresses(addresses)
	
	if err != nil {
		t.Errorf("Expected validation to succeed when some ports are available, got: %v", err)
	}
}

func TestGetTokenByAuthCode_PortValidationFailure(t *testing.T) {
	// Create test client with minimal setup
	c := &client{
		logger: logger.New(t),
	}

	// Block all test ports to trigger validation failure
	listener1, err := net.Listen("tcp", "127.0.0.1:18005")
	if err != nil {
		t.Skipf("Cannot bind to port 18005: %v", err)
	}
	defer listener1.Close()

	listener2, err := net.Listen("tcp", "127.0.0.1:18006")
	if err != nil {
		t.Skipf("Cannot bind to port 18006: %v", err)
	}
	defer listener2.Close()

	// Create context
	ctx := context.Background()

	// Prepare input with blocked addresses
	input := GetTokenByAuthCodeInput{
		BindAddress: []string{"127.0.0.1:18005", "127.0.0.1:18006"},
		State:       "test-state",
		Nonce:       "test-nonce",
	}

	readyChan := make(chan string, 1)

	// This should fail fast due to port validation
	_, err = c.GetTokenByAuthCode(ctx, input, readyChan)

	// Verify it failed
	if err == nil {
		t.Error("Expected GetTokenByAuthCode to fail when all ports are blocked")
	}

	// Verify the error message indicates port binding validation failure
	if !contains(err.Error(), "port binding validation failed") {
		t.Errorf("Expected error message about port binding validation, got: %v", err)
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || s[0:len(substr)] == substr || contains(s[1:], substr))
}