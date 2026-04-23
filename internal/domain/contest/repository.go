package contest

import "context"

// Repository defines the persistence contract for the Contest aggregate.
type Repository interface {
	SaveContest(ctx context.Context, a *CreateContestCommand) (Contest, error)
	SaveVote(ctx context.Context, vote *CreateVoteDTO) (Vote, error)
	SaveParticipant(ctx context.Context, participant *CreateParticipantDTO) (Participant, error)
	FindContestByID(ctx context.Context, id string) (*Contest, error)
	FindContestByPublicIDandUserID(ctx context.Context, publicID string, userID string) (*Contest, error)
	FindAllContests(ctx context.Context) ([]*Contest, error)
	FindAllContestsByUserID(ctx context.Context, userPublicID string) ([]*Contest, error)
	FindAllParticipantsByContestID(ctx context.Context, contestID string) ([]*Participant, error)
	FindAllParticipantsWithVotesByContestID(ctx context.Context, contestID string) ([]*ParticipantWithVotes, error)
	GetContestWithParticipants(ctx context.Context, contestID string) (*ContestWithParticipants, error)
	CountContestVotesFromVoterHash(ctx context.Context, contestID string, voterHash string) (int, error)
}
