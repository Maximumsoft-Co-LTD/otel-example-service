package service

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"example-otelgo-miniservice/internal/domain"
	"example-otelgo-miniservice/internal/port"
	"example-otelgo-miniservice/internal/utils"
	"fmt"
	"strings"
	"time"

	"github.com/Maximumsoft-Co-LTD/otelgo/eto"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/codes"
)

type JWE struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	encrypter  jose.Encrypter
	Issuer     string
	Audience   string
	cache      port.Redis
}

func NewJWE(redis port.Redis) (*JWE, error) {
	// span: generate keypair + init jwe (รันตอน bootstrap)
	ctx, span := eto.Trace().
		Name("JWE.NewJWE").
		Start()
	defer span.End()

	start := time.Now()
	result := "success"
	defer func() {
		latency := float64(time.Since(start).Milliseconds())
		eto.MetricHistogram("jwe_init_latency_ms").
			Record(ctx, latency)
		eto.MetricCounter("jwe_init_total").
			Attr("result", result).
			Add(ctx, 1)
	}()

	privateKey, publicKey, err := utils.GenerateRSAKeyPair(4096)
	if err != nil {
		result = "failed"
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		eto.Log().
			FromContext(ctx).
			Error().
			Msg("failed to generate RSA key pair").
			Field("error", err.Error()).
			Send()

		return nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
	}

	encrypter, err := jose.NewEncrypter(
		jose.A256GCM,
		jose.Recipient{
			Algorithm: jose.RSA_OAEP,
			Key:       publicKey,
		},
		(&jose.EncrypterOptions{}).WithContentType("JWT"),
	)
	if err != nil {
		result = "failed"
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		eto.Log().
			FromContext(ctx).
			Error().
			Msg("failed to create JWE encrypter").
			Field("error", err.Error()).
			Send()

		return nil, fmt.Errorf("failed to create encrypter: %w", err)
	}

	jwe := &JWE{
		publicKey:  publicKey,
		privateKey: privateKey,
		encrypter:  encrypter,
		Issuer:     "example-otelgo-miniservice",
		Audience:   "auth",
		cache:      redis,
	}

	eto.Log().
		FromContext(ctx).
		Info().
		Msg("JWE service initialized").
		Send()

	return jwe, nil
}

func (a *JWE) GenerateToken(ctx context.Context, payload domain.TokenPayload) (string, time.Time, error) {
	spanCtx, span := eto.Trace().
		Name("JWE.GenerateToken").
		FromContext(ctx).
		Attr("subject", payload.Subject).
		Start()
	defer span.End()

	start := time.Now()
	result := "success"
	defer func() {
		latency := float64(time.Since(start).Milliseconds())
		eto.MetricHistogram("jwe_generate_token_latency_ms").
			Record(spanCtx, latency)
		eto.MetricCounter("jwe_generate_token_total").
			Attr("result", result).
			Add(spanCtx, 1)
	}()

	now := time.Now()
	expiresIn := now.Add(24 * time.Hour)

	payload.Issuer = a.Issuer
	payload.Audience = jwt.Audience{a.Audience}
	payload.IssuedAt = jwt.NewNumericDate(now)
	payload.Expiry = jwt.NewNumericDate(expiresIn)
	payload.NotBefore = jwt.NewNumericDate(now)

	raw, err := jwt.Encrypted(a.encrypter).Claims(payload).Serialize()
	if err != nil {
		result = "failed"
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		eto.Log().
			FromContext(spanCtx).
			Error().
			Msg("failed to encrypt token").
			Field("error", err.Error()).
			Send()

		return "", expiresIn.UTC(), fmt.Errorf("failed to encrypt token: %w", err)
	}

	eto.Log().
		FromContext(spanCtx).
		Info().
		Msg("JWE token generated").
		Field("subject", payload.Subject).
		Send()

	return raw, expiresIn.UTC(), nil
}

func (a *JWE) DecryptToken(ctx context.Context, token string) (*domain.TokenPayload, error) {
	spanCtx, span := eto.Trace().
		Name("JWE.DecryptToken").
		FromContext(ctx).
		Start()
	defer span.End()

	start := time.Now()
	result := "success"
	defer func() {
		latency := float64(time.Since(start).Milliseconds())
		eto.MetricHistogram("jwe_decrypt_token_latency_ms").
			Record(spanCtx, latency)
		eto.MetricCounter("jwe_decrypt_token_total").
			Attr("result", result).
			Add(spanCtx, 1)
	}()

	token = strings.TrimSpace(token)
	if token == "" {
		result = "empty_token"
		err := fmt.Errorf("empty token")

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		eto.Log().
			FromContext(spanCtx).
			Warn().
			Msg("decrypt token: empty token").
			Send()

		return nil, err
	}

	jwe, err := jose.ParseEncrypted(token, []jose.KeyAlgorithm{jose.RSA_OAEP}, []jose.ContentEncryption{jose.A256GCM})
	if err != nil {
		result = "parse_failed"
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		eto.MetricCounter("jwe_decrypt_token_parse_error_total").
			Add(spanCtx, 1)

		eto.Log().
			FromContext(spanCtx).
			Error().
			Msg("decrypt token: parse JWE failed").
			Field("error", err.Error()).
			Send()

		return nil, fmt.Errorf("parse JWE failed: %w", err)
	}

	plaintext, err := jwe.Decrypt(a.privateKey)
	if err != nil {
		result = "decrypt_failed"
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		eto.MetricCounter("jwe_decrypt_token_decrypt_error_total").
			Add(spanCtx, 1)

		eto.Log().
			FromContext(spanCtx).
			Error().
			Msg("decrypt token: decrypt JWE failed").
			Field("error", err.Error()).
			Send()

		return nil, fmt.Errorf("decrypt JWE failed: %w", err)
	}

	var payload domain.TokenPayload
	if err := json.Unmarshal(plaintext, &payload); err != nil {
		result = "unmarshal_failed"
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		eto.MetricCounter("jwe_decrypt_token_unmarshal_error_total").
			Add(spanCtx, 1)

		eto.Log().
			FromContext(spanCtx).
			Error().
			Msg("decrypt token: unmarshal claims failed").
			Field("error", err.Error()).
			Send()

		return nil, fmt.Errorf("unmarshal claims failed: %w", err)
	}

	if err := payload.Validate(jwt.Expected{
		Time:        time.Now(),
		Issuer:      a.Issuer,
		AnyAudience: jwt.Audience{a.Audience},
	}); err != nil {
		result = "validate_failed"
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		eto.MetricCounter("jwe_decrypt_token_validate_error_total").
			Add(spanCtx, 1)

		eto.Log().
			FromContext(spanCtx).
			Warn().
			Msg("decrypt token: validation failed").
			Field("error", err.Error()).
			Send()

		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	span.SetAttributes(
	// ถ้ามี subject ใส่ลงไปเป็น attr
	)

	eto.Log().
		FromContext(spanCtx).
		Info().
		Msg("JWE token decrypted & validated").
		Send()

	return &payload, nil
}

func (a *JWE) key(jti string) string {
	return fmt.Sprintf("bl:jti:%s", jti)
}

func (a *JWE) BlacklistJTI(ctx context.Context, jti string, ttlSec int64) error {
	spanCtx, span := eto.Trace().
		Name("JWE.BlacklistJTI").
		FromContext(ctx).
		Attr("jti", jti).
		Start()
	defer span.End()

	start := time.Now()
	result := "success"
	defer func() {
		latency := float64(time.Since(start).Milliseconds())
		eto.MetricHistogram("jwe_blacklist_jti_latency_ms").
			Record(spanCtx, latency)
		eto.MetricCounter("jwe_blacklist_jti_total").
			Attr("result", result).
			Add(spanCtx, 1)
	}()

	if jti == "" {
		result = "empty_jti"
		err := errors.New("empty jti")

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		eto.Log().
			FromContext(spanCtx).
			Warn().
			Msg("blacklist JTI: empty jti").
			Send()

		return err
	}

	err := a.cache.Set(spanCtx, a.key(jti), true, time.Duration(ttlSec)*time.Second)
	if err != nil {
		result = "redis_error"

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		eto.Log().
			FromContext(spanCtx).
			Error().
			Msg("blacklist JTI: redis set failed").
			Field("error", err.Error()).
			Send()

		return err
	}

	eto.Log().
		FromContext(spanCtx).
		Info().
		Msg("blacklist JTI success").
		Field("jti", jti).
		Field("ttl_sec", ttlSec).
		Send()

	return nil
}

func (a *JWE) IsBlacklisted(ctx context.Context, jti string) (bool, error) {
	spanCtx, span := eto.Trace().
		Name("JWE.IsBlacklisted").
		FromContext(ctx).
		Attr("jti", jti).
		Start()
	defer span.End()

	start := time.Now()
	result := "not_blacklisted"
	defer func() {
		latency := float64(time.Since(start).Milliseconds())
		eto.MetricHistogram("jwe_is_blacklisted_latency_ms").
			Record(spanCtx, latency)
		eto.MetricCounter("jwe_is_blacklisted_total").
			Attr("result", result).
			Add(spanCtx, 1)
	}()

	if jti == "" {
		err := errors.New("empty jti")

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		eto.Log().
			FromContext(spanCtx).
			Warn().
			Msg("isBlacklisted: empty jti").
			Send()

		return false, err
	}

	_, err := a.cache.Get(spanCtx, a.key(jti))
	if err == nil {
		result = "blacklisted"

		eto.Log().
			FromContext(spanCtx).
			Info().
			Msg("isBlacklisted: jti is blacklisted").
			Field("jti", jti).
			Send()

		return true, nil
	}

	if errors.Is(err, redis.Nil) {
		result = "not_blacklisted"
		return false, nil
	}

	result = "redis_error"

	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())

	eto.Log().
		FromContext(spanCtx).
		Error().
		Msg("isBlacklisted: redis error").
		Field("error", err.Error()).
		Send()

	return false, err
}
