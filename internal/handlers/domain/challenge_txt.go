package domain

import (
	"strings"

	"clevergo.tech/clevergo"
)

func (h *Handler) challengeTXT(c *clevergo.Context) error {
	domain, err := h.findDomain(c)
	if err != nil {
		return err
	}

	return c.SendFile(domain.ChallengeTXT, strings.NewReader(domain.ChallengeTXT))
}
