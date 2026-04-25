package contest

import (
	"context"
	"fmt"
)

// Service holds domain logic for the avaliacao aggregate.
type Service struct {
	repo       Repository
	acessToken map[string]int // simple in-memory store for dyn acess token in vote contest not suitable for production

}

func NewService(repo Repository) *Service {
	return &Service{repo: repo, acessToken: make(map[string]int)}
}

func (s *Service) CreateContest(ctx context.Context, user_id int, name string, description string, maxVotes int) (Contest, error) {
	// add validation and business logic here if needed
	if name == "" {
		return Contest{}, fmt.Errorf("name is required")
	}
	if description == "" {
		return Contest{}, fmt.Errorf("description is required")
	}
	dto := &CreateContestCommand{
		UserID:      user_id,
		Name:        name,
		Description: description,
		MaxVotes:    maxVotes,
	}
	return s.repo.SaveContest(ctx, dto)
}

func (s *Service) GetContestByPublicIDandUserID(ctx context.Context, publicID string, userID string) (*Contest, error) {
	return s.repo.FindContestByPublicIDandUserID(ctx, publicID, userID)
}
func (s *Service) GetContestByID(ctx context.Context, id string) (*Contest, error) {
	return s.repo.FindContestByID(ctx, id)
}

func (s *Service) GetAllContests(ctx context.Context) ([]*Contest, error) {
	return s.repo.FindAllContests(ctx)
}
func (s *Service) ListContestUserID(ctx context.Context, userID string) ([]*Contest, error) {
	// This function is not implemented in the repository, so we return an empty slice for now.
	return s.repo.FindAllContestsByUserID(ctx, userID)
}

func (s *Service) GetContestWithParticipants(ctx context.Context, contestID string) (*ContestWithParticipants, error) {
	return s.repo.GetContestWithParticipants(ctx, contestID)
}

func (s *Service) CreateParticipant(ctx context.Context, dto *CreateParticipantDTO) (Participant, error) {
	return s.repo.SaveParticipant(ctx, dto)
}

func (s *Service) ApplyVote(ctx context.Context, dto *CreateVoteDTO) (Vote, error) {
	contest, err := s.repo.FindContestByID(ctx, dto.ContestID)
	if err != nil {
		return Vote{}, fmt.Errorf("contest not found: %w", err)
	}

	// Check if the contest is active
	if !contest.IsUp {
		return Vote{}, fmt.Errorf("competição não está ativa")
	}

	count, err := s.repo.CountContestVotesFromVoterHash(ctx, dto.ContestID, dto.VoterHash)
	if err != nil {
		return Vote{}, fmt.Errorf("falha ao contar votos: %w", err)
	}

	if count >= contest.MaxVotesUser {
		return Vote{}, fmt.Errorf("usuário atingiu o número máximo de votos para esta competição")
	}
	// Here you can add any business logic related to voting, such as checking if the user has already voted.
	fmt.Printf("Applying vote for contest %s and participant %s by user %s\n", dto.ContestID, dto.ParticipantID, dto.VoterHash)
	return s.repo.SaveVote(ctx, dto)
}

func (s *Service) GetVotesFromVoterHash(ctx context.Context, contestID string, voterHash string) (int, error) {
	count, err := s.repo.CountContestVotesFromVoterHash(ctx, contestID, voterHash)
	if err != nil {
		return 0, err
	}

	fmt.Printf("User with hash %s has cast %d votes in contest %s\n", voterHash, count, contestID)

	// You can return the count or a list of votes depending on your needs.
	return count, nil // Replace with actual return value if needed
}

func (s *Service) GetContestAcesstoken(ctx context.Context, contestID string) (string, error) {
	contest, err := s.repo.FindContestByID(ctx, contestID)
	if err != nil {
		return "", fmt.Errorf("contest not found: %w", err)
	}
	// gen a random  number

	// Assuming the access ID is the same as the public ID for simplicity
	return contest.PublicID, nil
}

func (s *Service) UpdateContest(ctx context.Context, contestID string, name *string, description *string, isUp *bool, maxVotes *int) (*Contest, error) {
	contest, err := s.repo.FindContestByID(ctx, contestID)
	if err != nil {
		return nil, fmt.Errorf("contest not found: %w", err)
	}
	if name != nil {
		contest.Name = *name
	}
	if description != nil {
		contest.Description = *description
	}
	if isUp != nil {
		contest.IsUp = *isUp
	}
	if maxVotes != nil {
		contest.MaxVotesUser = *maxVotes
	}
	// Here you would typically call a repository method to update the contest in the database.
	// For example: return s.repo.UpdateContest(ctx, contest)
	return s.repo.UpdateContest(ctx, contest)
}
