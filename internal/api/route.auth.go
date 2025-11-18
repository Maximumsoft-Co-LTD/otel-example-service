package api

import (
	"example-otelgo-miniservice/internal/api/dto"
	"example-otelgo-miniservice/internal/application/core"
	"example-otelgo-miniservice/internal/port"
	"example-otelgo-miniservice/internal/utils"
	"net/http"
	"time"

	"github.com/Maximumsoft-Co-LTD/otelgo/eto"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opentelemetry.io/otel/codes"
)

type AuthRoute struct {
	cache    port.Redis
	auth     *core.Auth
	validate *validator.Validate
}

func NewAuthRoute(redis port.Redis, auth *core.Auth) (*AuthRoute, error) {
	return &AuthRoute{
		cache:    redis,
		auth:     auth,
		validate: validator.New(),
	}, nil
}

// =========================
// Login
// =========================

func (r *AuthRoute) Login(c *gin.Context) {
	baseCtx := c.Request.Context()

	// Parent span
	scope := eto.Trace().
		Name("AuthRoute.Login").
		FromContext(baseCtx).
		Attr("http.route", "/auth/login").
		StartScope()
	defer scope.Done()

	ctx := scope.Ctx()

	// ===== Metrics =====
	start := time.Now()
	result := "success"
	defer func() {
		latencyMs := float64(time.Since(start).Milliseconds())

		eto.MetricCounter("http_server_requests_total").
			Attr("route", "/auth/login").
			Attr("method", c.Request.Method).
			Attr("result", result).
			Add(ctx, 1)

		eto.MetricHistogram("http_server_request_duration_ms").
			Attr("route", "/auth/login").
			Attr("method", c.Request.Method).
			Attr("result", result).
			Record(ctx, latencyMs)

		eto.MetricCounter("auth_login_attempts_total").
			Attr("result", result).
			Add(ctx, 1)
	}()

	// Child 1: Validate Input (JSON → struct)
	validateCtx, validateSpan := eto.Trace().
		Name("AuthRoute.Login.Validate").
		FromContext(ctx).
		Start()
	defer validateSpan.End()

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result = "invalid_payload"

		validateSpan.RecordError(err)
		validateSpan.SetStatus(codes.Error, err.Error())

		eto.Log().
			FromContext(validateCtx).
			Warn().
			Msg("invalid login payload").
			Field("error", err.Error()).
			Send()

		RespJsonError(c, http.StatusBadRequest, dto.MsgInvalidResponse, "4010")
		return
	}

	// Child 2: Struct Validate
	vCtx, vSpan := eto.Trace().
		Name("AuthRoute.Login.StructValidate").
		FromContext(ctx).
		Attr("username", req.Username).
		Start()
	if err := r.validate.StructCtx(vCtx, req); err != nil {
		result = "validation_error"

		vSpan.RecordError(err)
		vSpan.SetStatus(codes.Error, err.Error())
		vSpan.End()

		eto.Log().
			FromContext(vCtx).
			Warn().
			Msg("login validation failed").
			Field("error", err.Error()).
			Send()

		RespJsonError(c, http.StatusBadRequest, dto.MsgInvalidResponse, "4011")
		return
	}
	vSpan.End()

	// Child 3: Business logic
	bizCtx, bizSpan := eto.Trace().
		Name("AuthRoute.Login.Authenticate").
		FromContext(ctx).
		Start()
	token, exp, err := r.auth.Login(bizCtx, req.Username, req.Password)
	if err != nil {
		result = "invalid_credential"

		bizSpan.RecordError(err)
		bizSpan.SetStatus(codes.Error, err.Error())
		bizSpan.End()

		eto.Log().
			FromContext(bizCtx).
			Warn().
			Msg("login failed").
			Field("username", req.Username).
			Field("error", err.Error()).
			Send()

		RespJsonError(c, http.StatusUnauthorized, dto.MsgInvalidRequest, "4012")
		return
	}
	bizSpan.End()

	// Child 4: Final logging
	logCtx, logSpan := eto.Trace().
		Name("AuthRoute.Login.Success").
		FromContext(ctx).
		Attr("username", req.Username).
		Start()
	eto.Log().
		FromContext(logCtx).
		Info().
		Msg("login success").
		Field("username", req.Username).
		Send()
	logSpan.End()

	// Response
	RespJsonSuccess(c, http.StatusOK, gin.H{
		"token": token,
		"exp":   exp,
	}, dto.MsgSuccess, "0")
}

// =========================
// Logout
// =========================

func (r *AuthRoute) Logout(c *gin.Context) {
	baseCtx := c.Request.Context()

	scope := eto.Trace().
		Name("AuthRoute.Logout").
		FromContext(baseCtx).
		Attr("http.route", "/auth/logout").
		StartScope()
	defer scope.Done()

	ctx := scope.Ctx()

	start := time.Now()
	result := "success"
	defer func() {
		latencyMs := float64(time.Since(start).Milliseconds())

		eto.MetricCounter("http_server_requests_total").
			Attr("route", "/auth/logout").
			Attr("method", c.Request.Method).
			Attr("result", result).
			Add(ctx, 1)

		eto.MetricHistogram("http_server_request_duration_ms").
			Attr("route", "/auth/logout").
			Attr("method", c.Request.Method).
			Attr("result", result).
			Record(ctx, latencyMs)

		// business metric
		eto.MetricCounter("auth_logout_total").
			Attr("result", result).
			Add(ctx, 1)
	}()

	// Child 1: Extract token
	tCtx, tSpan := eto.Trace().
		Name("AuthRoute.Logout.ExtractToken").
		FromContext(ctx).
		Start()
	token, ok := utils.BearerTokenFrom(c.GetHeader("Authorization"))
	if !ok {
		result = "missing_token"

		err := http.ErrNoCookie // fake error ใช้ mark span
		tSpan.RecordError(err)
		tSpan.SetStatus(codes.Error, "missing bearer token")
		tSpan.End()

		eto.Log().
			FromContext(tCtx).
			Warn().
			Msg("logout missing bearer token").
			Send()

		RespJsonError(c, http.StatusUnauthorized, dto.MsgInvalidResponse, "4010")
		return
	}
	tSpan.End()

	// Child 2: Call Auth.Logout
	bCtx, bSpan := eto.Trace().
		Name("AuthRoute.Logout.LogoutUser").
		FromContext(ctx).
		Start()
	if err := r.auth.Logout(bCtx, token); err != nil {
		result = "failed"

		bSpan.RecordError(err)
		bSpan.SetStatus(codes.Error, err.Error())
		bSpan.End()

		eto.Log().
			FromContext(bCtx).
			Error().
			Msg("logout failed").
			Field("error", err.Error()).
			Send()

		RespJsonError(c, http.StatusUnauthorized, dto.MsgInvalidRequest, "4011")
		return
	}
	bSpan.End()

	// Child 3: Success log
	logCtx, logSpan := eto.Trace().
		Name("AuthRoute.Logout.Success").
		FromContext(ctx).
		Start()
	eto.Log().
		FromContext(logCtx).
		Info().
		Msg("logout success").
		Send()
	logSpan.End()

	result = "success"

	RespJsonSuccess(c, http.StatusOK, nil, dto.MsgSuccess, "0")
}

// =========================
// Refresh Token
// =========================

func (r *AuthRoute) RefreshToken(c *gin.Context) {
	baseCtx := c.Request.Context()

	scope := eto.Trace().
		Name("AuthRoute.RefreshToken").
		FromContext(baseCtx).
		Attr("http.route", "/auth/refresh").
		StartScope()
	defer scope.Done()

	ctx := scope.Ctx()

	start := time.Now()
	result := "success"
	defer func() {
		latencyMs := float64(time.Since(start).Milliseconds())

		eto.MetricCounter("http_server_requests_total").
			Attr("route", "/auth/refresh").
			Attr("method", c.Request.Method).
			Attr("result", result).
			Add(ctx, 1)

		eto.MetricHistogram("http_server_request_duration_ms").
			Attr("route", "/auth/refresh").
			Attr("method", c.Request.Method).
			Attr("result", result).
			Record(ctx, latencyMs)

		eto.MetricCounter("auth_refresh_token_total").
			Attr("result", result).
			Add(ctx, 1)
	}()

	// Child 1: Extract token
	tCtx, tSpan := eto.Trace().
		Name("AuthRoute.RefreshToken.ExtractToken").
		FromContext(ctx).
		Start()
	oldToken, ok := utils.BearerTokenFrom(c.GetHeader("Authorization"))
	if !ok {
		result = "missing_token"

		err := http.ErrNoCookie
		tSpan.RecordError(err)
		tSpan.SetStatus(codes.Error, "missing bearer token")
		tSpan.End()

		eto.Log().
			FromContext(tCtx).
			Warn().
			Msg("refresh token missing bearer token").
			Send()

		RespJsonError(c, http.StatusUnauthorized, dto.MsgInvalidResponse, "4010")
		return
	}
	tSpan.End()

	// Child 2: Call Auth.RefreshToken
	bCtx, bSpan := eto.Trace().
		Name("AuthRoute.RefreshToken.Refresh").
		FromContext(ctx).
		Start()
	newToken, exp, err := r.auth.RefreshToken(bCtx, oldToken)
	if err != nil {
		result = "failed"

		bSpan.RecordError(err)
		bSpan.SetStatus(codes.Error, err.Error())
		bSpan.End()

		eto.Log().
			FromContext(bCtx).
			Warn().
			Msg("refresh token failed").
			Field("error", err.Error()).
			Send()

		RespJsonError(c, http.StatusUnauthorized, dto.MsgInvalidRequest, "4011")
		return
	}
	bSpan.End()

	// Child 3: Success log
	logCtx, logSpan := eto.Trace().
		Name("AuthRoute.RefreshToken.Success").
		FromContext(ctx).
		Start()
	eto.Log().
		FromContext(logCtx).
		Info().
		Msg("refresh token success").
		Send()
	logSpan.End()

	result = "success"

	RespJsonSuccess(c, http.StatusOK, gin.H{
		"token": newToken,
		"exp":   exp,
	}, dto.MsgSuccess, "0")
}

// =========================
// Profile
// =========================

func (r *AuthRoute) Profile(c *gin.Context) {
	baseCtx := c.Request.Context()

	scope := eto.Trace().
		Name("AuthRoute.Profile").
		FromContext(baseCtx).
		Attr("http.route", "/auth/profile").
		StartScope()
	defer scope.Done()

	ctx := scope.Ctx()

	start := time.Now()
	result := "success"
	defer func() {
		latencyMs := float64(time.Since(start).Milliseconds())

		eto.MetricCounter("http_server_requests_total").
			Attr("route", "/auth/profile").
			Attr("method", c.Request.Method).
			Attr("result", result).
			Add(ctx, 1)

		eto.MetricHistogram("http_server_request_duration_ms").
			Attr("route", "/auth/profile").
			Attr("method", c.Request.Method).
			Attr("result", result).
			Record(ctx, latencyMs)

		eto.MetricCounter("auth_profile_total").
			Attr("result", result).
			Add(ctx, 1)
	}()

	// Child 1: Extract token
	tCtx, tSpan := eto.Trace().
		Name("AuthRoute.Profile.ExtractToken").
		FromContext(ctx).
		Start()
	token, ok := utils.BearerTokenFrom(c.GetHeader("Authorization"))
	if !ok {
		result = "missing_token"

		err := http.ErrNoCookie
		tSpan.RecordError(err)
		tSpan.SetStatus(codes.Error, "missing bearer token")
		tSpan.End()

		eto.Log().
			FromContext(tCtx).
			Warn().
			Msg("profile missing bearer token").
			Send()

		RespJsonError(c, http.StatusUnauthorized, dto.MsgInvalidResponse, "4010")
		return
	}
	tSpan.End()

	// Child 2: Call Auth.Profile
	bCtx, bSpan := eto.Trace().
		Name("AuthRoute.Profile.GetProfile").
		FromContext(ctx).
		Start()
	profile, err := r.auth.Profile(bCtx, token)
	if err != nil {
		result = "failed"

		bSpan.RecordError(err)
		bSpan.SetStatus(codes.Error, err.Error())
		bSpan.End()

		eto.Log().
			FromContext(bCtx).
			Warn().
			Msg("get profile failed").
			Field("error", err.Error()).
			Send()

		RespJsonError(c, http.StatusUnauthorized, dto.MsgInvalidRequest, "4011")
		return
	}
	bSpan.End()

	// Child 3: Success log
	logCtx, logSpan := eto.Trace().
		Name("AuthRoute.Profile.Success").
		FromContext(ctx).
		Start()
	eto.Log().
		FromContext(logCtx).
		Info().
		Msg("get profile success").
		Send()
	logSpan.End()

	result = "success"

	RespJsonSuccess(c, http.StatusOK, profile, dto.MsgSuccess, "0")
}
