package crypto

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/atomreforge/daizy-night-server/internal/model"
)

func (p *ProviderCrypto) SignRegistercode(payload model.RegistercodePayload) (
	RegistercodeRaw string,
	e error) {

	payloadBytes, _ := json.Marshal(payload)
	payloadHex := hex.EncodeToString(payloadBytes)

	sig := ed25519.Sign(p.RegistercodeEnckey, payloadBytes)
	sigHex := hex.EncodeToString(sig)

	slog.Info("Signed Registercode for atomid", "atomid", payload.AtomID, "magicword", payload.Magicword, "before", payload.Before)
	return payloadHex + "." + sigHex, nil
}

// tranform regcode string to payload struct.
func (p *ProviderCrypto) AnalyzeRegistercode(codeStr string) (*model.RegistercodePayload, error) {
	// decode registercode string to model.RegistercodePacket

	// spilt
	parts := strings.SplitN(codeStr, ".", 2)
	if len(parts) != 2 {
		slog.Error("Invalid registercode format")
		return nil, &ErrRegister{
			StatusCode: 400,
			Message:    "invalid registercode format",
			Field:      "registercode",
			Type:       RegistercodeFormat,
		}
	}

	// transform
	payloadHex, sigHex := parts[0], parts[1]
	payloadStr, _ := hex.DecodeString(payloadHex)
	sig, _ := hex.DecodeString(sigHex)

	// verify sig
	if !ed25519.Verify(p.RegistercodeDeckey, []byte(payloadStr), sig) {
		slog.Error("Registercode signature verification failed")
		return nil, &ErrRegister{
			StatusCode: 400,
			Message:    "registercode signature verification failed",
			Field:      "registercode",
			Type:       RegistercodeAuthenFailed,
		}
	}

	// decode payload to struct
	var payload model.RegistercodePayload
	err := json.Unmarshal([]byte(payloadStr), &payload)
	if err != nil {
		slog.Error("Failed to unmarshal registercode", "error", err)
		return nil, &ErrRegister{
			StatusCode: 400,
			Message:    "failed to unmarshal registercode",
			Field:      "registercode",
			Type:       RegistercodeFormat,
		}
	}
	return &payload, nil
}

// 1.sig; 2.expire
func (p *ProviderCrypto) VerifyRegistercodePayload(payload model.RegistercodePayload) bool {

	return payload.Before.Before(time.Now())

}
