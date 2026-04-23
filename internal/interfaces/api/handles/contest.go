package apihandles

import (
	"github.com/PopularVote/internal/domain/contest"
	"github.com/PopularVote/internal/domain/user"
	sqlrepo "github.com/PopularVote/internal/infrastructure/sql"
	apimiddleware "github.com/PopularVote/internal/interfaces/api/middleware"
	apipayloads "github.com/PopularVote/internal/interfaces/api/payloads"
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

	user, err := userService.GetUserByPublicID(c.Context(), user_id)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve user",
		})
		return err
	}
	var req apipayloads.CreateContestRequest
	if err := c.Bind().JSON(&req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
		return err
	}
	contestRepo := sqlrepo.NewContestRepo(db)
	contestService := contest.NewService(contestRepo)
	createdContest, err := contestService.CreateContest(c.Context(), user.ID, req.Name, req.Description, req.MaxVotes)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create contest",
		})
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(createdContest)

}
func (h *ContestHandler) GetContest(c fiber.Ctx) error {
	paramID := c.Params("id")
	db := c.Locals("db").(*sqlrepo.DB)
	user_id := c.Locals(apimiddleware.LocalUserID).(string)
	contestRepo := sqlrepo.NewContestRepo(db)
	contestService := contest.NewService(contestRepo)
	contest, err := contestService.GetContestByPublicIDandUserID(c.Context(), paramID, user_id)
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Contest not found",
		})
		return err
	}
	return c.JSON(contest)
}
func (h *ContestHandler) ListContests(c fiber.Ctx) error {
	user_id := c.Locals(apimiddleware.LocalUserID).(string)
	db := c.Locals("db").(*sqlrepo.DB)
	contestRepo := sqlrepo.NewContestRepo(db)
	contestService := contest.NewService(contestRepo)
	contests, err := contestService.ListContestUserID(c.Context(), user_id)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list contests",
		})
		return err
	}
	return c.JSON(contests)
}
func (h *ContestHandler) UpdateContest(c fiber.Ctx) error { return nil }

// TODO: Add DeleteContest when we implement soft deletes.
func (h *ContestHandler) AddParticipant(c fiber.Ctx) error    { return nil }
func (h *ContestHandler) RemoveParticipant(c fiber.Ctx) error { return nil }

func (h *ContestHandler) ListContestParticipants(c fiber.Ctx) error { return nil }

func (h *ContestHandler) ListPairtoVote(c fiber.Ctx) error  { return nil }
func (h *ContestHandler) VoteParticipant(c fiber.Ctx) error { return nil }
