package simple

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/osins/oauth2wall/common"
	"golang.org/x/oauth2"
)

func Authorize(ctx *fiber.Ctx) error {
	s, r := GenOAuthState(ctx)
	if r != nil {
		return ctx.JSON(common.NewResult(fmt.Sprintf("state 生成失败")).SetSuccess(false))
	}

	clientSecret := oauth2.SetAuthURLParam("client_secret", config.ClientSecret)
	return ctx.Redirect(oAuth2Config.AuthCodeURL(s, clientSecret), http.StatusTemporaryRedirect)
}

func Callback(ctx *fiber.Ctx) error {
	state := ctx.Query("state")
	code := ctx.Query("code")

	if s, r := GetOAuthState(ctx); r != nil && state != s {
		return ctx.JSON(common.NewResult("state 验证失败").SetSuccess(false))
	}

	sessionStore.Storage.Delete(SESSION_STATE_KEY)
	fmt.Printf("callback code: %s\n", code)
	token, err := oAuth2Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return ctx.JSON(common.NewResult("获取用户信息失败, 有可能是该授权已失效，请重新访问授权接口").SetSuccess(false).SetError(err))
	}

	return ctx.JSON(GetUser(token.AccessToken))
}

func Middleware(ctx *fiber.Ctx) error {
	if strings.HasPrefix(ctx.Route().Path, "/oauth") {
		return nil
	}

	token := ctx.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	result := GetUser(token)
	if result.Success {
		ctx.Locals("user", result)

		return nil
	}

	return fmt.Errorf("simple oauth get user faild.")
}
