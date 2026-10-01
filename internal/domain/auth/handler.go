// Package auth HTTP 入口暂留 domain（与 Session/OAuth 同包，避免跨包导出洪泛）。
// Trigger 负责 module 注册并组装 Case。
//
// Author: Charlie

package auth


import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"voxel-gin-admin/internal/infrastructure/core/bind"
	contextx "voxel-gin-admin/internal/infrastructure/core/context"
	"voxel-gin-admin/internal/types/response"
	"voxel-gin-admin/internal/infrastructure/core/security"
	"voxel-gin-admin/internal/infrastructure/middleware"
)


func (s *Service) RegisterHTTPRoutes(api *gin.RouterGroup) {
	rdb := s.Repo.rdb

	api.GET("/v1/admin/auth/session/captcha", middleware.RateLimit(rdb, "admin:captcha", 30, 60), s.captcha)
	api.GET("/v1/admin/auth/session/password-key", middleware.RateLimit(rdb, "admin:password-key", 30, 60), s.passwordKey)

	api.POST("/v1/admin/auth/session/login", middleware.RateLimit(rdb, "admin:login", 20, 60), middleware.OperationAudit(s.Audit, "auth", "login"), s.login(security.AccountAdmin))
	api.POST("/v1/admin/auth/session/send-login-code", middleware.RateLimit(rdb, "admin:send-login-code", 10, 60), middleware.OperationAudit(s.Audit, "auth", "send_login_code"), s.sendLoginCode(security.AccountAdmin))

	api.POST("/v1/admin/auth/session/forgot-password", middleware.RateLimit(rdb, "admin:forgot-password", 5, 60), middleware.OperationAudit(s.Audit, "auth", "forgot_password"), s.forgotPassword(security.AccountAdmin))
	api.POST("/v1/admin/auth/session/forgot-password/phone", middleware.RateLimit(rdb, "admin:forgot-password-phone", 5, 60), middleware.OperationAudit(s.Audit, "auth", "forgot_password_phone"), s.forgotPasswordByPhone(security.AccountAdmin))
	api.POST("/v1/admin/auth/session/reset-password", middleware.RateLimit(rdb, "admin:reset-password", 10, 60), middleware.OperationAudit(s.Audit, "auth", "reset_password"), s.resetPassword(security.AccountAdmin))
	api.POST("/v1/admin/auth/session/reset-password/phone", middleware.RateLimit(rdb, "admin:reset-password-phone", 10, 60), middleware.OperationAudit(s.Audit, "auth", "reset_password_phone"), s.resetPasswordByPhone(security.AccountAdmin))

	api.GET("/v1/public/site-footer", s.siteFooter)

	api.POST("/v1/admin/auth/session/logout", middleware.RequireAccountType(security.AccountAdmin), middleware.OperationAudit(s.Audit, "auth", "logout"), s.logout)

	api.GET("/v1/admin/auth/public/auth-options", s.authOptions(security.AccountAdmin))

	api.POST("/v1/admin/auth/session/refresh", middleware.RequireAccountType(security.AccountAdmin), middleware.OperationAudit(s.Audit, "auth", "refresh"), s.refresh(security.AccountAdmin))

	api.POST("/v1/admin/auth/session/cancel", middleware.RequireAccountType(security.AccountAdmin), middleware.OperationAudit(s.Audit, "auth", "cancel"), s.cancel(security.AccountAdmin))

	if s.Oauth != nil {
		s.Oauth.RegisterRoutes(api)
		s.Oauth.RegisterBindingRoutes(api, s.Perms)
	}
	s.registerSessionRoutes(api)
}

func (s *Service) captcha(c *gin.Context) {
	out, err := s.Captcha(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.OK(c, out)
}

func (s *Service) passwordKey(c *gin.Context) {
	out, err := s.PasswordKey(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.OK(c, out)
}

func (s *Service) login(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginParam
		if err := bind.JSON(c, &req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		out, err := s.Login(c.Request.Context(), accountType, req, c.ClientIP(), c.Request.UserAgent())
		if err != nil {
			switch err {
			case errInvalidCredentials, errInvalidOTP, errAccountLocked, errIPLocked:
				response.Fail(c, http.StatusUnauthorized, 401, err.Error())
			case errAccountFinder:
				response.Fail(c, http.StatusInternalServerError, 500, err.Error())
			default:
				response.Fail(c, http.StatusBadRequest, 400, err.Error())
			}
			return
		}
		ttlSec := out.ExpiresIn
		if ttlSec <= 0 {
			ttlSec = s.cfg.Auth.TokenTTLSeconds
		}
		if ttlSec <= 0 {
			ttlSec = 14400
		}
		remember := rememberMeOrDefault(req.RememberMe)
		s.SetSessionCookie(c, out.Token, accountType, time.Duration(ttlSec)*time.Second, remember)
		maxAge := ttlSec
		if !remember {
			maxAge = 0
		}
		middleware.IssueCSRFCookie(c, s.cfg.Auth, maxAge)
		response.OK(c, out)
	}
}

func (s *Service) logout(c *gin.Context) {
	token := s.ResolveLogoutToken(c)
	sess := contextx.Session(c.Request.Context())
	accountID, accountType := "", ""
	if sess != nil {
		accountID = sess.AccountID
		accountType = string(sess.AccountType)
	}
	_ = s.Logout(c.Request.Context(), token, accountID, accountType, c.ClientIP(), c.Request.UserAgent())
	s.ClearSessionCookie(c, contextx.AccountType(c.Request.Context()))
	middleware.ClearCSRFCookie(c)
	response.OK(c, nil)
}

func (s *Service) register(c *gin.Context) {
	var req RegisterParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	out, err := s.Register(c.Request.Context(), req)
	if err != nil {
		switch err {
		case errRegisterDisabled:
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
		case errPortalRegistrar:
			response.Fail(c, http.StatusNotImplemented, 501, err.Error())
		default:
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
		}
		return
	}
	response.OK(c, out)
}

func (s *Service) sendLoginCode(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req SendLoginCodeParam
		if err := bind.JSON(c, &req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		if err := s.SendLoginCode(c.Request.Context(), accountType, req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		response.OK(c, nil)
	}
}

func (s *Service) forgotPassword(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ForgotPasswordParam
		if err := bind.JSON(c, &req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		if err := s.ForgotPassword(c.Request.Context(), accountType, req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		response.OK(c, nil)
	}
}

func (s *Service) resetPassword(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ResetPasswordParam
		if err := bind.JSON(c, &req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		if err := s.ResetPassword(c.Request.Context(), accountType, req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		response.OK(c, nil)
	}
}

func (s *Service) forgotPasswordByPhone(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ForgotPasswordByPhoneParam
		if err := bind.JSON(c, &req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		if err := s.ForgotPasswordByPhone(c.Request.Context(), accountType, req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		response.OK(c, nil)
	}
}

func (s *Service) resetPasswordByPhone(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req ResetPasswordByPhoneParam
		if err := bind.JSON(c, &req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		if err := s.ResetPasswordByPhone(c.Request.Context(), accountType, req); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		response.OK(c, nil)
	}
}

func (s *Service) siteFooter(c *gin.Context) {
	response.OK(c, s.SiteFooter(c.Request.Context()))
}

func (s *Service) authOptions(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		response.OK(c, s.AuthOptions(c.Request.Context(), accountType))
	}
}

func (s *Service) refresh(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		out, err := s.RefreshSession(c.Request.Context(), accountType, c.ClientIP(), c.Request.UserAgent())
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, 401, err.Error())
			return
		}
		response.OK(c, out)
	}
}

func (s *Service) cancel(accountType security.AccountType) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CancelParam
		_ = bind.JSON(c, &req)
		if err := s.CancelAccount(c.Request.Context(), accountType, c.ClientIP(), c.Request.UserAgent(), req.CancelReason); err != nil {
			response.Fail(c, http.StatusBadRequest, 400, err.Error())
			return
		}
		s.ClearSessionCookie(c, contextx.AccountType(c.Request.Context()))
		middleware.ClearCSRFCookie(c)
		response.OK(c, nil)
	}
}

func (s *Service) registerSendCode(c *gin.Context) {
	var req SendLoginCodeParam
	if err := bind.JSON(c, &req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	if err := s.sendRegisterCode(c.Request.Context(), req); err != nil {
		response.Fail(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	response.OK(c, nil)
}
