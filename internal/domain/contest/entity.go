package contest

// --- DB models for the contest domain
type Contest struct {
	ID           int    ` json:"id"`
	PublicID     string `json:"public_id"` // UUID for external reference
	Name         string `json:"name"`
	Description  string `json:"description"`
	UserId       string `json:"user_id"`   // Reference to the user who created the contest
	IsUp         bool   `json:"is_up"`     // Indicates if the contest is active or not
	MaxVotesUser int    `json:"max_votes"` // Maximum number of votes per user for this contest
	// Add more fields as needed, e.g., StartDate, EndDate, etc.
}

type Participant struct {
	ID        int    `json:"id"`
	PublicID  string `json:"public_id"`  // UUID for external reference
	ContestID int    `json:"contest_id"` // Reference to the contest this participant belongs to
	// contatos
	Email    *string `json:"email"`
	Nome     string  `json:"nome"`
	Telefone *string `json:"telefone"`
	// cosplay dado
	Title       string  `json:"title"`       // Title of the cosplay
	Description *string `json:"description"` // Description of the cosplay optional
	ImagePath   string  `json:"image_path"`  // Path to the participant's image
}

type Vote struct {
	ID                string `json:"id"`
	ContestID         string `json:"contest_id"`          // Reference to the contest being voted on
	ParticipantID     string `json:"participant_id"`      // Reference to the participant being voted for
	VoterHash         string `json:"voter_id"`            // Reference to the user who cast the vote
	VoterIP           string `json:"voter_ip"`            // IP address of the voter (for anti-fraud measures)
	VoterRandomCookie string `json:"voter_random_cookie"` // Random cookie to help identify unique voters
	VotedAt           string `json:"voted_at"`            // Timestamp of when the vote was cast
}

type VoteOptions struct {
	ContestID      string `json:"contest_id"`       // Reference to the contest being voted on
	VoteID         string `json:"vote_id"`          // Reference to the vote being cast
	ConterOptionID string `json:"conter_option_id"` // Reference of other participant to be voted against (pair options selected by the system)
}

// entities agrregate

type ParticipantWithVotes struct {
	Participant        Participant `json:"participant"`
	Votes              int         `json:"votes"`
	VoteOptionsCounter int         `json:"vote_options_counter"`
}

type ContestWithParticipants struct {
	Contest      Contest                `json:"contest"`
	Participants []ParticipantWithVotes `json:"participants"`
}
