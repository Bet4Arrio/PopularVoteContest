package contest

type CreateContestCommand struct {
	UserID      int    `json:"user_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	MaxVotes    int    `json:"max_votes" validate:"gt=0"`
}

type CreateParticipantDTO struct {
	ContestID string `json:"contest_id" validate:"required"`
	Name      string `json:"name" validate:"required"`
}

type CreateVoteDTO struct {
	ContestID         string `json:"contest_id" validate:"required"`
	ParticipantID     string `json:"participant_id" validate:"required"`
	VoterHash         string `json:"voter_hash" validate:"required"`
	VoterIP           string `json:"voter_ip" validate:"required,ip"`
	VoterRandomCookie string `json:"voter_random_cookie" validate:"required"`
	OptionsID         []int  `json:"options_id" validate:"required,dive,gt=0"` // IDs of other participants being voted against
}
