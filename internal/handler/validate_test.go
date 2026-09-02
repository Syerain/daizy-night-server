package handler

import (
	"strings"
	"testing"

	"github.com/atomreforge/daizy-night-server/internal/errs"
	"github.com/atomreforge/daizy-night-server/internal/model"
)

// validSigHex returns 128 hex chars — the exact size of an Ed25519 signature.
func validSigHex() string { return strings.Repeat("ab", 64) }

func TestValidateFormatRegistercode(t *testing.T) {
	cases := []struct {
		name     string
		code     model.RegistercodeRawHex
		wantErr  bool
		wantType errs.ErrType
	}{
		{"valid", model.RegistercodeRawHex("abcd." + validSigHex()), false, 0},
		{"empty", "", true, errs.ValidationKeyNull},
		{"no dot", model.RegistercodeRawHex("abcdef"), true, errs.ValidationKeyBadFormat},
		{"leading dot", model.RegistercodeRawHex("." + validSigHex()), true, errs.ValidationKeyBadFormat},
		{"trailing dot", model.RegistercodeRawHex("abcd."), true, errs.ValidationKeyBadFormat},
		{"non-hex payload", model.RegistercodeRawHex("zzzz." + validSigHex()), true, errs.ValidationKeyBadFormat},
		{"non-hex signature", model.RegistercodeRawHex("abcd." + strings.Repeat("g", 128)), true, errs.ValidationKeyBadFormat},
		{"short signature", model.RegistercodeRawHex("abcd." + validSigHex()[:127]), true, errs.ValidationKeyBadFormat},
		{"oversized payload", model.RegistercodeRawHex(strings.Repeat("a", 2049) + "." + validSigHex()), true, errs.ValidationKeyBadFormat},
		{"max payload", model.RegistercodeRawHex(strings.Repeat("a", 2048) + "." + validSigHex()), false, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := ValidateFormatRegistercode(tc.code)
			if !tc.wantErr {
				if !ok || err != nil {
					t.Fatalf("expected accept, got ok=%v err=%v", ok, err)
				}
				return
			}
			if ok || err == nil {
				t.Fatalf("expected rejection, got ok=%v err=%v", ok, err)
			}
			v, isVal := errs.Easx[*errs.ErrValidation](err)
			if !isVal {
				t.Fatalf("expected *errs.ErrValidation, got %T: %v", err, err)
			}
			if v.Type != tc.wantType {
				t.Fatalf("type = %v, want %v", v.Type, tc.wantType)
			}
		})
	}
}

func TestValidateRefreshParams(t *testing.T) {
	if err := ValidateRefreshParams(""); err == nil {
		t.Fatal("empty refresh token must be rejected")
	} else if v, ok := errs.Easx[*errs.ErrValidation](err); !ok ||
		v.Type != errs.ValidationKeyNull || v.StatusCode() != 400 {
		t.Fatalf("wrong error for empty token: %v", err)
	}

	if err := ValidateRefreshParams("some-token"); err != nil {
		t.Fatalf("valid token rejected: %v", err)
	}
}
