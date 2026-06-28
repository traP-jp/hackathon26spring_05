package handler

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v5"
	"github.com/traP-jp/hackathon26spring_05/Qpid/domain"
)

var groupAffiliationMap = map[string]domain.UserAffiliation{
	"280bf56d-fa22-46bc-8dcc-6367d600d873": domain.UserAffiliationAlgorithm,
	"c5670065-75d4-4851-bfba-9ff05201fc44": domain.UserAffiliationCtf,
	"af240e80-8526-4f21-925e-b20eded06284": domain.UserAffiliationGame,
	"867b3529-696f-4bd1-af53-1947eba92e77": domain.UserAffiliationGraphics,
	"ec54d385-e5e7-4554-8aa2-878ebedc9db0": domain.UserAffiliationKaggle,
	"cb977ab2-85fa-4953-ac4d-809eaef427e6": domain.UserAffiliationSound,
	"f86db5ec-dc02-4885-aa0a-732bb229a1b5": domain.UserAffiliationSysAd,
}

func getAffiliations(groups []string) []domain.UserAffiliation {
	affiliations := []domain.UserAffiliation{}
	for _, groupUUID := range groups {
		if affiliation, ok := groupAffiliationMap[groupUUID]; ok {
			affiliations = append(affiliations, affiliation)
		}
	}
	return affiliations
}

func (h *handler) fetchAffiliations(c *echo.Context, uuid string) []domain.UserAffiliation {
	groups, err := h.getUserGroups(uuid)
	if err != nil {
		c.Logger().Warn("failed to get traQ groups, affiliations will be empty",
			slog.String("uuid", uuid),
			slog.Any("error", err),
		)
		return []domain.UserAffiliation{}
	}
	return getAffiliations(groups)
}

func (h *handler) getUserUUID(username string) (string, error) {
	if cached, ok := h.uuidCache.Load(username); ok {
		return cached.(string), nil
	}
	if h.traq.client == nil {
		return "", nil
	}
	users, _, err := h.traq.client.UserAPI.GetUsers(h.traq.context).Name(username).Execute()
	if err != nil {
		return "", err
	}
	for _, u := range users {
		if u.Name == username {
			h.uuidCache.Store(username, u.Id)
			return u.Id, nil
		}
	}
	return "", fmt.Errorf("user not found: %s", username)
}

func (h *handler) getUserGroups(uuid string) ([]string, error) {
	if h.traq.client == nil {
		return nil, nil
	}
	user, _, err := h.traq.client.UserAPI.GetUser(h.traq.context, uuid).Execute()
	if err != nil {
		return nil, err
	}
	return user.Groups, nil
}
