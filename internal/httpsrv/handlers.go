package httpsrv

import "github.com/gofiber/fiber/v2"

type Handler interface {
	Mount(*fiber.App)
}
