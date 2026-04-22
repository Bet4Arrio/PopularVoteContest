package apihandles

import (
	"github.com/gofiber/fiber/v3"
)

type ContestHandler struct {
}

func NewContestHandler() *ContestHandler {
	return &ContestHandler{}
}

func (h *ContestHandler) CreateContest(c *fiber.Ctx) {}
func (h *ContestHandler) GetContest(c *fiber.Ctx)    {}
func (h *ContestHandler) ListContests(c *fiber.Ctx)  {}
func (h *ContestHandler) UpdateContest(c *fiber.Ctx) {}

// TODO: Add DeleteContest when we implement soft deletes.
func (h *ContestHandler) AddParticipant(c *fiber.Ctx)    {}
func (h *ContestHandler) RemoveParticipant(c *fiber.Ctx) {}

func (h *ContestHandler) ListContestParticipants(c *fiber.Ctx) {}

func (h *ContestHandler) ListPairtoVote(c *fiber.Ctx)  {}
func (h *ContestHandler) VoteParticipant(c *fiber.Ctx) {}
