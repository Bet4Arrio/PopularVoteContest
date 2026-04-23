package apipayloads

type CreateContestRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	MaxVotes    int    `json:"max_votes" validate:"gt=0"`
}

type ContestResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsUp         bool   `json:"is_up"`
	MaxVotesUser int    `json:"max_votes_user"`
}

type ListContestsResponse struct {
	Contests []ContestResponse `json:"contests"`
}

type ParticipantResponse struct {
	ID          string  `json:"id"`
	ContestID   string  `json:"contest_id"`
	Email       *string `json:"email,omitempty"`
	Nome        string  `json:"nome"`
	Telefone    *string `json:"telefone,omitempty"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	ImagePath   string  `json:"image_path"`
	Votes       int     `json:"votes"`
}

type ContestWithParticipantsResponse struct {
	ID           string                `json:"id"`
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	IsUp         bool                  `json:"is_up"`
	MaxVotesUser int                   `json:"max_votes_user"`
	Participants []ParticipantResponse `json:"participants"`
}

type ParticipantOption struct {
	ContestID     string `json:"contest_id"`
	ParticipantID string `json:"participant_id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ImagePath     string `json:"image_path"`
}
type VotesOptionsResponse struct {
	Options [][]ParticipantOption `json:"options"`
}

type VoteRequest struct {
	ContestID         string   `json:"contest_id" validate:"required"`
	VoteID            string   `json:"vote_id" validate:"required"`
	VoterHash         string   `json:"voter_hash" validate:"required"`
	VoterIP           string   `json:"voter_ip" validate:"required"`
	VoterRandomCookie string   `json:"voter_random_cookie" validate:"required"`
	Options           []string `json:"options" validate:"required"` // List of participant IDs being voted for

}
