package orgs

import (
	"base-core/helper"
	"base-core/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type HTTPTransport interface {
	Add(c *fiber.Ctx) error
	FindMyOrgs(c *fiber.Ctx) error
}

type HttpTransport struct {
	orgApi OrgAPI
	logger log.AllLogger
}

func NewHTTPTransport(oapi OrgAPI, logger log.AllLogger) HTTPTransport {
	return &HttpTransport{
		orgApi: oapi,
		logger: logger,
	}
}

func (s *HttpTransport) Add(c *fiber.Ctx) error {
	req := &OrgRequest{}
	userId, err := middleware.CtxUserID(c)
	if err != nil {
		return helper.HTTPError(c, err, "OrgHTTPTransport.CtxUserID")
	}
	req.UserID = userId
	if err := c.BodyParser(req); err != nil {
		return helper.HTTPError(c, err, "OrgHTTPTransport.Add.BodyParser")
	}
	resp, err := s.orgApi.Add(req)
	if err != nil {
		return helper.HTTPError(c, err, "OrgHTTPTransport.Add")
	}
	return c.JSON(resp)
}

func (s *HttpTransport) FindMyOrgs(c *fiber.Ctx) error {
	req := &FindOrgRequest{}
	userId, err := middleware.CtxUserID(c)
	if err != nil {
		return helper.HTTPError(c, err, "OrgHTTPTransport.CtxUserID")
	}
	req.UserID = userId
	resp, err := s.orgApi.FindMyOrgs(req)
	if err != nil {
		return helper.HTTPError(c, err, "OrgHTTPTransport.Add")
	}
	return c.JSON(resp)
}
