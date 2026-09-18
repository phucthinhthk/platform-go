package security

import (
	"testing"
	"time"
)

func TestJWTServiceCreateAndParse(t *testing.T) {
	service, err := NewJWTService("01234567890123456789012345678901")
	if err != nil {
		t.Fatal(err)
	}
	token, err := service.Create(Claims{Subject: "1", TokenType: AccessTokenType}, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := service.Parse(token, AccessTokenType)
	if err != nil {
		t.Fatal(err)
	}
	if claims.Subject != "1" {
		t.Fatalf("subject = %q, want 1", claims.Subject)
	}
	if _, err := service.Parse(token, RefreshTokenType); err == nil {
		t.Fatal("expected token type mismatch")
	}
}

func TestJWTServiceRejectsTampering(t *testing.T) {
	service, _ := NewJWTService("01234567890123456789012345678901")
	token, _ := service.Create(Claims{Subject: "1", TokenType: AccessTokenType}, time.Hour)
	if _, err := service.Parse(token+"x", AccessTokenType); err == nil {
		t.Fatal("expected tampered token to be rejected")
	}
}
