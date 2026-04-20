package contest

type CreateContestDTO struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	StartDate   string `json:"start_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndDate     string `json:"end_date" validate:"required,datetime=2006-01-02T15:04:05Z07:00,gtfield=StartDate"`
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
