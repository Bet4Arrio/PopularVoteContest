package apihandles

import (
	"github.com/PopularVote/internal/domain/user"
	sqlrepo "github.com/PopularVote/internal/infrastructure/sql"
	apimiddleware "github.com/PopularVote/internal/interfaces/api/middleware"
	"github.com/gofiber/fiber/v3"
)

type ContestHandler struct {
}

func NewContestHandler() *ContestHandler {
	return &ContestHandler{}
}

func (h *ContestHandler) CreateContest(c fiber.Ctx) error {
	user_id := c.Locals(apimiddleware.LocalUserID).(string)
	db := c.Locals("db").(*sqlrepo.DB)
	userRepo := sqlrepo.NewUserRepo(db)
	userService := user.NewService(userRepo)
	// contestRepo := sqlrepo.NewContestRepo(db)
	// _ := contest.NewService(contestRepo)

	_, err := userService.GetUserByPublicID(c.Context(), user_id)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user",
		})
		return err
	}

	return nil

}
func (h *ContestHandler) GetContest(c *fiber.Ctx)    {}
func (h *ContestHandler) ListContests(c *fiber.Ctx)  {}
func (h *ContestHandler) UpdateContest(c *fiber.Ctx) {}

// TODO: Add DeleteContest when we implement soft deletes.
func (h *ContestHandler) AddParticipant(c *fiber.Ctx)    {}
func (h *ContestHandler) RemoveParticipant(c *fiber.Ctx) {}

func (h *ContestHandler) ListContestParticipants(c *fiber.Ctx) {}

func (h *ContestHandler) ListPairtoVote(c *fiber.Ctx)  {}
func (h *ContestHandler) VoteParticipant(c *fiber.Ctx) {}
