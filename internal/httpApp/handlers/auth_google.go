package handlers

import "github.com/gofiber/fiber/v3"

func (h *AuthHandler) RedirectGoogleAuth(c fiber.Ctx) error {
	//state := generateState() // сохранить в redis / cookie

	url := "" // oauthConfig.AuthCodeURL(state)
	return c.Redirect().To(url)
}
