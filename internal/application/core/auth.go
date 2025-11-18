package core

import (
	"context"
	"errors"
	"example-otelgo-miniservice/internal/application/service"
	"example-otelgo-miniservice/internal/domain"
	"example-otelgo-miniservice/internal/port"
	"example-otelgo-miniservice/internal/utils"
	"strings"
	"time"

	"github.com/Maximumsoft-Co-LTD/otelgo/eto"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	repoUser port.Repository[domain.UserRepo]
	repoRole port.Repository[domain.RoleRepo]
	cache    port.Redis
	jwe      *service.JWE
}

func NewAuth(
	userRepo port.Repository[domain.UserRepo],
	roleRepo port.Repository[domain.RoleRepo],
	redis port.Redis,
	jwe *service.JWE,
) (*Auth, error) {
	return &Auth{
		repoUser: userRepo,
		repoRole: roleRepo,
		cache:    redis,
		jwe:      jwe,
	}, nil
}

// =========================
// Login
// =========================

func (a *Auth) Login(ctx context.Context, username, password string) (string, int, error) {
	var (
		user = strings.TrimSpace(strings.ToLower(username))
		pass = strings.TrimSpace(password)
	)

	// Parent span สำหรับ service layer
	rootCtx, rootSpan := eto.Trace().
		Name("Auth.Login").
		FromContext(ctx).
		Start()
	ctx = rootCtx
	defer rootSpan.End()

	// Service-level metrics
	start := time.Now()
	result := "success"
	defer func() {
		latencyMs := float64(time.Since(start).Milliseconds())

		eto.MetricCounter("auth_login_service_total").
			Attr("result", result).
			Add(ctx, 1)

		eto.MetricHistogram("auth_login_service_duration_ms").
			Attr("result", result).
			Record(ctx, latencyMs)
	}()

	// Child span 1: Validate input
	vCtx, vSpan := eto.Trace().
		Name("Auth.Login.ValidateInput").
		FromContext(ctx).
		Start()
	if user == "" || pass == "" {
		err := errors.New("invalid username or password")
		result = "invalid_input"

		vSpan.RecordError(err)
		vSpan.SetStatus(codes.Error, err.Error())
		vSpan.End()

		eto.Log().
			FromContext(vCtx).
			Warn().
			Msg("username/password empty").
			Send()

		return "", 0, err
	}
	vSpan.End()

	// Child span 2: DB Find User
	dbCtx, dbSpan := eto.Trace().
		Name("Auth.Login.FindUser").
		FromContext(ctx).
		Attr("db.system", "mongodb").
		Attr("db.operation", "findOne").
		Attr("db.collection", "users").
		Attr("username", user).
		Start()

	dbStart := time.Now()

	filter := bson.M{"username": user}
	userOne, ok, err := a.repoUser.FindOne(dbCtx, filter)
	if err != nil {
		result = "db_error"

		dbSpan.RecordError(err)
		dbSpan.SetStatus(codes.Error, err.Error())
		dbSpan.End()

		eto.Log().
			FromContext(dbCtx).
			Error().
			Msg("db error: find user failed").
			Field("error", err.Error()).
			Send()

		return "", 0, err
	}
	if !ok {
		err := errors.New("invalid username or password")
		result = "user_not_found"

		dbSpan.RecordError(err)
		dbSpan.SetStatus(codes.Error, err.Error())
		dbSpan.End()

		eto.Log().
			FromContext(dbCtx).
			Warn().
			Msg("user not found").
			Field("username", user).
			Send()

		return "", 0, err
	}

	dbLatency := float64(time.Since(dbStart).Milliseconds())
	eto.MetricHistogram("auth_login_db_lookup_latency_ms").
		Attr("collection", "users").
		Record(dbCtx, dbLatency)

	dbSpan.End()

	// Child span 3: Password Check
	pwCtx, pwSpan := eto.Trace().
		Name("Auth.Login.CheckPassword").
		FromContext(ctx).
		Start()

	if err := bcrypt.CompareHashAndPassword([]byte(userOne.Password), []byte(pass)); err != nil {
		err := errors.New("invalid username or password")
		result = "password_mismatch"

		pwSpan.RecordError(err)
		pwSpan.SetStatus(codes.Error, err.Error())
		pwSpan.End()

		eto.MetricCounter("auth_login_password_check_fail_total").
			Attr("username", user).
			Add(pwCtx, 1)

		eto.Log().
			FromContext(pwCtx).
			Warn().
			Msg("password mismatch").
			Field("username", user).
			Send()

		return "", 0, err
	}
	pwSpan.End()

	// Child span 4: Generate Token
	tkCtx, tkSpan := eto.Trace().
		Name("Auth.Login.GenerateToken").
		FromContext(ctx).
		Attr("user_id", userOne.ID.Hex()).
		Start()

	payload := domain.TokenPayload{
		Roles: utils.ObjectIDSliceToStringSlice(userOne.Roles),
		Claims: jwt.Claims{
			Subject: userOne.ID.Hex(), // sub
			ID:      uuid.NewString(), // jti
		},
	}

	accessToken, expires, err := a.jwe.GenerateToken(tkCtx, payload)
	if err != nil {
		result = "generate_token_error"

		tkSpan.RecordError(err)
		tkSpan.SetStatus(codes.Error, err.Error())
		tkSpan.End()

		eto.Log().
			FromContext(tkCtx).
			Error().
			Msg("generate token failed").
			Field("error", err.Error()).
			Send()

		return "", 0, err
	}

	expiresIn := int(expires.Sub(time.Now()).Seconds())
	if expiresIn < 0 {
		expiresIn = 0
	}

	eto.MetricCounter("auth_login_success_total").
		Attr("username", user).
		Add(tkCtx, 1)

	tkSpan.End()

	result = "success"
	return accessToken, expiresIn, nil
}

// =========================
// Logout
// =========================

func (a *Auth) Logout(ctx context.Context, token string) error {
	rootCtx, rootSpan := eto.Trace().
		Name("Auth.Logout").
		FromContext(ctx).
		Start()
	ctx = rootCtx
	defer rootSpan.End()

	start := time.Now()
	result := "success"
	defer func() {
		latencyMs := float64(time.Since(start).Milliseconds())

		eto.MetricCounter("auth_logout_service_total").
			Attr("result", result).
			Add(ctx, 1)

		eto.MetricHistogram("auth_logout_service_duration_ms").
			Attr("result", result).
			Record(ctx, latencyMs)
	}()

	trimmed := strings.TrimSpace(token)

	// decrypt token
	claims, err := a.jwe.DecryptToken(ctx, trimmed)
	if err != nil {
		// behavior เดิม: logout เงียบ ๆ ถ้า decrypt ไม่ได้ → mark เป็น noop
		result = "noop"

		eto.Log().
			FromContext(ctx).
			Warn().
			Msg("logout decrypt token failed, treat as noop").
			Field("error", err.Error()).
			Send()

		return nil
	}

	// ถ้า JTI มีอยู่ เช็ค blacklist ก่อน
	if claims.ID != "" {
		black, blkErr := a.jwe.IsBlacklisted(ctx, claims.ID)
		if blkErr != nil {
			result = "blacklist_error"

			rootSpan.RecordError(blkErr)
			rootSpan.SetStatus(codes.Error, blkErr.Error())

			eto.Log().
				FromContext(ctx).
				Error().
				Msg("logout check blacklist failed").
				Field("error", blkErr.Error()).
				Send()

			return blkErr
		}
		if black {
			// เคย logout แล้ว → noop
			result = "already_blacklisted"

			eto.Log().
				FromContext(ctx).
				Info().
				Msg("logout noop: token already blacklisted").
				Field("jti", claims.ID).
				Send()

			return nil
		}
	}

	if claims.Expiry == nil || claims.ID == "" {
		err := errors.New("token missing exp or jti")
		result = "invalid_token"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		return err
	}

	ttl := claims.Expiry.Time().Unix() - time.Now().Unix()
	if ttl <= 0 {
		result = "expired"

		eto.Log().
			FromContext(ctx).
			Info().
			Msg("logout noop: token already expired").
			Send()

		return nil
	}

	if err := a.jwe.BlacklistJTI(ctx, claims.ID, ttl); err != nil {
		result = "blacklist_error"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		eto.Log().
			FromContext(ctx).
			Error().
			Msg("logout blacklist failed").
			Field("error", err.Error()).
			Send()

		return err
	}

	result = "success"

	eto.Log().
		FromContext(ctx).
		Info().
		Msg("logout success").
		Field("jti", claims.ID).
		Send()

	return nil
}

// =========================
// Refresh Token
// =========================

func (a *Auth) RefreshToken(ctx context.Context, token string) (string, int, error) {
	rootCtx, rootSpan := eto.Trace().
		Name("Auth.RefreshToken").
		FromContext(ctx).
		Start()
	ctx = rootCtx
	defer rootSpan.End()

	start := time.Now()
	result := "success"
	defer func() {
		latencyMs := float64(time.Since(start).Milliseconds())

		eto.MetricCounter("auth_refresh_service_total").
			Attr("result", result).
			Add(ctx, 1)

		eto.MetricHistogram("auth_refresh_service_duration_ms").
			Attr("result", result).
			Record(ctx, latencyMs)
	}()

	token = strings.TrimSpace(token)
	if token == "" {
		err := errors.New("empty token")
		result = "invalid_input"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		return "", 0, err
	}

	// decrypt
	claims, err := a.jwe.DecryptToken(ctx, token)
	if err != nil {
		result = "decrypt_failed"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		eto.Log().
			FromContext(ctx).
			Warn().
			Msg("refresh decrypt token failed").
			Field("error", err.Error()).
			Send()

		return "", 0, err
	}

	if claims.Expiry == nil || time.Now().After(claims.Expiry.Time()) {
		err := errors.New("token expired")
		result = "expired"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		return "", 0, err
	}

	// check blacklist
	if claims.ID != "" {
		black, blkErr := a.jwe.IsBlacklisted(ctx, claims.ID)
		if blkErr != nil {
			result = "blacklist_error"

			rootSpan.RecordError(blkErr)
			rootSpan.SetStatus(codes.Error, blkErr.Error())

			return "", 0, blkErr
		}
		if black {
			err := errors.New("forbidden")
			result = "blacklisted"

			rootSpan.RecordError(err)
			rootSpan.SetStatus(codes.Error, err.Error())

			return "", 0, err
		}
	}

	if claims.Subject == "" {
		err := errors.New("forbidden")
		result = "forbidden"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		return "", 0, err
	}

	uid, err := bson.ObjectIDFromHex(claims.Subject)
	if err != nil {
		result = "invalid_subject"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		return "", 0, err
	}

	// DB lookup user
	dbStart := time.Now()
	userOne, ok, err := a.repoUser.FindOne(ctx, bson.M{"_id": uid})
	if err != nil {
		result = "db_error"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		return "", 0, err
	}
	if !ok {
		err := errors.New("forbidden")
		result = "user_not_found"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		return "", 0, err
	}
	if !userOne.Active {
		err := errors.New("user is not active")
		result = "inactive"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		return "", 0, err
	}

	dbLatency := float64(time.Since(dbStart).Milliseconds())
	eto.MetricHistogram("auth_refresh_db_lookup_latency_ms").
		Attr("collection", "users").
		Record(ctx, dbLatency)

	newPayload := domain.TokenPayload{
		Roles: utils.ObjectIDSliceToStringSlice(userOne.Roles),
		Claims: jwt.Claims{
			Subject: userOne.ID.Hex(),
			ID:      uuid.NewString(), // rotate jti
		},
	}

	newToken, newExp, err := a.jwe.GenerateToken(ctx, newPayload)
	if err != nil {
		result = "generate_token_error"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		return "", 0, err
	}

	// blacklist old jti (best effort)
	if claims.Expiry != nil && claims.ID != "" {
		if ttl := claims.Expiry.Time().Unix() - time.Now().Unix(); ttl > 0 {
			_ = a.jwe.BlacklistJTI(ctx, claims.ID, ttl)
		}
	}

	expiresIn := int(time.Until(newExp).Seconds())
	if expiresIn < 0 {
		expiresIn = 0
	}

	result = "success"

	eto.Log().
		FromContext(ctx).
		Info().
		Msg("refresh token success").
		Field("user_id", userOne.ID.Hex()).
		Send()

	return newToken, expiresIn, nil
}

// =========================
// Profile
// =========================

func (a *Auth) Profile(ctx context.Context, token string) (*domain.UserProfile, error) {
	rootCtx, rootSpan := eto.Trace().
		Name("Auth.Profile").
		FromContext(ctx).
		Start()
	ctx = rootCtx
	defer rootSpan.End()

	start := time.Now()
	result := "success"
	defer func() {
		latencyMs := float64(time.Since(start).Milliseconds())

		eto.MetricCounter("auth_profile_service_total").
			Attr("result", result).
			Add(ctx, 1)

		eto.MetricHistogram("auth_profile_service_duration_ms").
			Attr("result", result).
			Record(ctx, latencyMs)
	}()

	claims, err := a.jwe.DecryptToken(ctx, strings.TrimSpace(token))
	if err != nil {
		result = "decrypt_failed"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	if claims.ID != "" {
		black, blkErr := a.jwe.IsBlacklisted(ctx, claims.ID)
		if blkErr != nil {
			result = "blacklist_error"

			rootSpan.RecordError(blkErr)
			rootSpan.SetStatus(codes.Error, blkErr.Error())

			return nil, blkErr
		}
		if black {
			err := errors.New("forbidden")
			result = "blacklisted"

			rootSpan.RecordError(err)
			rootSpan.SetStatus(codes.Error, err.Error())

			return nil, err
		}
	}

	if claims.Subject == "" {
		err := errors.New("forbidden")
		result = "forbidden"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	uid, err := bson.ObjectIDFromHex(claims.Subject)
	if err != nil {
		result = "invalid_subject"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	dbStart := time.Now()
	userOne, ok, err := a.repoUser.FindOne(ctx, bson.M{"_id": uid})
	if err != nil {
		result = "db_error"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		return nil, err
	}
	if !ok {
		err := errors.New("forbidden")
		result = "user_not_found"

		rootSpan.RecordError(err)
		rootSpan.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	dbLatency := float64(time.Since(dbStart).Milliseconds())
	eto.MetricHistogram("auth_profile_db_lookup_latency_ms").
		Attr("collection", "users").
		Record(ctx, dbLatency)

	var rolesMany []domain.RoleRepo
	if len(userOne.Roles) > 0 {
		rolesMany, err = a.repoRole.Find(ctx, bson.M{"_id": bson.M{"$in": userOne.Roles}}, nil)
		if err != nil {
			result = "db_error_roles"

			rootSpan.RecordError(err)
			rootSpan.SetStatus(codes.Error, err.Error())

			return nil, err
		}
	} else {
		rolesMany = nil
	}

	rolesSlug := make([]string, 0, len(rolesMany))
	for _, r := range rolesMany {
		rolesSlug = append(rolesSlug, r.Slug)
	}

	result = "success"

	eto.Log().
		FromContext(ctx).
		Info().
		Msg("get profile success").
		Field("user_id", userOne.ID.Hex()).
		Send()

	return &domain.UserProfile{
		Name:     userOne.Name,
		Username: userOne.Username,
		Email:    userOne.Email,
		Roles:    rolesSlug,
	}, nil
}
