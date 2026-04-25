package sqlrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/PopularVote/internal/domain/contest"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ContestRepo struct {
	db *DB
}

func NewContestRepo(db *DB) *ContestRepo {
	return &ContestRepo{db: db}
}

// SaveContest inserts a new contest and returns it.
func (r *ContestRepo) SaveContest(ctx context.Context, dto *contest.CreateContestCommand) (contest.Contest, error) {
	const q = `
		INSERT INTO contests (public_id, user_id, name, description, is_up, max_votes_user)
		VALUES ($1, $2, $3, $4, false, $5)
		RETURNING id, public_id, user_id, name, description, is_up, max_votes_user
	`
	log.Printf("Saving contest: %+v", dto)
	c := contest.Contest{PublicID: uuid.New().String()}
	err := r.db.Pool.QueryRow(ctx, q, c.PublicID, dto.UserID, dto.Name, dto.Description, dto.MaxVotes).
		Scan(&c.ID, &c.PublicID, &c.UserId, &c.Name, &c.Description, &c.IsUp, &c.MaxVotesUser)
	if err != nil {
		return contest.Contest{}, err
	}
	return c, nil
}

// FindContestByID returns a contest by its public UUID.
func (r *ContestRepo) FindContestByID(ctx context.Context, id string) (*contest.Contest, error) {
	const q = `
		SELECT id, public_id, user_id, name, description, is_up, max_votes_user
		FROM contests
		WHERE public_id = $1
	`
	c := &contest.Contest{}
	err := r.db.Pool.QueryRow(ctx, q, id).
		Scan(&c.ID, &c.PublicID, &c.UserId, &c.Name, &c.Description, &c.IsUp, &c.MaxVotesUser)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return c, nil
}

func (r *ContestRepo) FindContestByPublicIDandUserID(ctx context.Context, publicID string, userID string) (*contest.Contest, error) {
	const q = `
		SELECT c.id, c.public_id, c.name, c.description, c.is_up, c.max_votes_user
		FROM contests c
		JOIN "user" u ON u.id = c.user_id
		WHERE c.public_id = $1 AND u.public_id = $2
	`
	c := &contest.Contest{}
	err := r.db.Pool.QueryRow(ctx, q, publicID, userID).
		Scan(&c.ID, &c.PublicID, &c.Name, &c.Description, &c.IsUp, &c.MaxVotesUser)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return c, nil
}

// FindAllContests returns every contest.
func (r *ContestRepo) FindAllContests(ctx context.Context) ([]*contest.Contest, error) {
	const q = `SELECT id, public_id, name, description, is_up, max_votes_user FROM contests`
	rows, err := r.db.Pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*contest.Contest
	for rows.Next() {
		c := &contest.Contest{}
		if err := rows.Scan(&c.ID, &c.PublicID, &c.Name, &c.Description, &c.IsUp, &c.MaxVotesUser); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}
func (r *ContestRepo) FindAllContestsByUserID(ctx context.Context, userPublicID string) ([]*contest.Contest, error) {
	const q = `
		SELECT c.id, c.public_id, c.name, c.description, c.is_up, c.max_votes_user
		FROM contests c
		JOIN "user" u ON u.id = c.user_id
		WHERE u.public_id = $1
	`
	rows, err := r.db.Pool.Query(ctx, q, userPublicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*contest.Contest
	for rows.Next() {
		c := &contest.Contest{}
		if err := rows.Scan(&c.ID, &c.PublicID, &c.Name, &c.Description, &c.IsUp, &c.MaxVotesUser); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

// SaveParticipant inserts a new participant into a contest.
func (r *ContestRepo) SaveParticipant(ctx context.Context, dto *contest.CreateParticipantDTO) (contest.Participant, error) {
	const q = `
		INSERT INTO participants (public_id, contest_id, nome, title, image_path)
		SELECT $1, c.id, $3, $4, ''
		FROM contests c WHERE c.public_id = $2
		RETURNING id, public_id, contest_id, email, nome, telefone, title, description, image_path
	`
	p := contest.Participant{PublicID: uuid.New().String()}
	err := r.db.Pool.QueryRow(ctx, q, p.PublicID, dto.ContestID, dto.Name, dto.Name).
		Scan(&p.ID, &p.PublicID, &p.ContestID, &p.Email, &p.Nome, &p.Telefone, &p.Title, &p.Description, &p.ImagePath)
	if err != nil {
		return contest.Participant{}, err
	}
	return p, nil
}

// FindAllParticipantsByContestID returns all participants of a contest.
func (r *ContestRepo) FindAllParticipantsByContestID(ctx context.Context, contestID string) ([]*contest.Participant, error) {
	const q = `
		SELECT p.id, p.public_id, p.contest_id, p.email, p.nome, p.telefone, p.title, p.description, p.image_path
		FROM participants p
		JOIN contests c ON c.id = p.contest_id
		WHERE c.public_id = $1
	`
	rows, err := r.db.Pool.Query(ctx, q, contestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*contest.Participant
	for rows.Next() {
		p := &contest.Participant{}
		if err := rows.Scan(&p.ID, &p.PublicID, &p.ContestID, &p.Email, &p.Nome, &p.Telefone, &p.Title, &p.Description, &p.ImagePath); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

// FindAllParticipantsWithVotesByContestID returns participants with their vote counts.
func (r *ContestRepo) FindAllParticipantsWithVotesByContestID(ctx context.Context, contestID string) ([]*contest.ParticipantWithVotes, error) {
	const q = `
		SELECT p.id, p.public_id, p.contest_id, p.email, p.nome, p.telefone, p.title, p.description, p.image_path,
		       COUNT(v.id) AS votes,
		       COUNT(vo.id) AS vote_options_counter
		FROM participants p
		JOIN contests c ON c.id = p.contest_id
		LEFT JOIN votes v ON v.participant_id = p.id
		LEFT JOIN vote_options vo ON vo.contest_id = c.id
		WHERE c.public_id = $1
		GROUP BY p.id
	`
	rows, err := r.db.Pool.Query(ctx, q, contestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*contest.ParticipantWithVotes
	for rows.Next() {
		pw := &contest.ParticipantWithVotes{}
		if err := rows.Scan(
			&pw.Participant.ID, &pw.Participant.PublicID, &pw.Participant.ContestID,
			&pw.Participant.Email, &pw.Participant.Nome, &pw.Participant.Telefone,
			&pw.Participant.Title, &pw.Participant.Description, &pw.Participant.ImagePath,
			&pw.Votes, &pw.VoteOptionsCounter,
		); err != nil {
			return nil, err
		}
		list = append(list, pw)
	}
	return list, rows.Err()
}

// GetContestWithParticipants returns a contest together with its participants and vote counts.
func (r *ContestRepo) GetContestWithParticipants(ctx context.Context, contestID string) (*contest.ContestWithParticipants, error) {
	c, err := r.FindContestByID(ctx, contestID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, nil
	}
	participants, err := r.FindAllParticipantsWithVotesByContestID(ctx, contestID)
	if err != nil {
		return nil, err
	}
	result := &contest.ContestWithParticipants{Contest: *c}
	for _, pw := range participants {
		result.Participants = append(result.Participants, *pw)
	}
	return result, nil
}

// SaveVote records a vote and its counter-option references.
func (r *ContestRepo) SaveVote(ctx context.Context, dto *contest.CreateVoteDTO) (contest.Vote, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return contest.Vote{}, err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	const voteQ = `
		INSERT INTO votes (contest_id, participant_id, voter_hash, voter_ip, voter_random_cookie, voted_at)
		SELECT c.id, p.id, $3, $4, $5, $6
		FROM contests c, participants p
		WHERE c.public_id = $1 AND p.public_id = $2
		RETURNING id
	`
	var voteID int
	err = tx.QueryRow(ctx, voteQ,
		dto.ContestID, dto.ParticipantID,
		dto.VoterHash, dto.VoterIP, dto.VoterRandomCookie,
		time.Now(),
	).Scan(&voteID)
	if err != nil {
		return contest.Vote{}, err
	}

	// Insert vote_options for each counter-option participant id
	for _, optID := range dto.OptionsID {
		const optQ = `
			INSERT INTO vote_options (contest_id, vote_id, counter_option_id)
			SELECT c.id, $2, p.public_id
			FROM contests c, participants p
			WHERE c.public_id = $1 AND p.id = $3
		`
		if _, err := tx.Exec(ctx, optQ, dto.ContestID, voteID, optID); err != nil {
			return contest.Vote{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return contest.Vote{}, err
	}

	return contest.Vote{
		ID:                string(rune(voteID)),
		ContestID:         dto.ContestID,
		ParticipantID:     dto.ParticipantID,
		VoterHash:         dto.VoterHash,
		VoterIP:           dto.VoterIP,
		VoterRandomCookie: dto.VoterRandomCookie,
		VotedAt:           time.Now().String(),
	}, nil
}

// CountContestVotesFromVoterHash counts how many votes a given voter has cast in a contest.
func (r *ContestRepo) CountContestVotesFromVoterHash(ctx context.Context, contestID string, voterHash string) (int, error) {
	const q = `
		SELECT COUNT(v.id)
		FROM votes v
		JOIN contests c ON c.id = v.contest_id
		WHERE c.public_id = $1 AND v.voter_hash = $2
	`
	var count int
	err := r.db.Pool.QueryRow(ctx, q, contestID, voterHash).Scan(&count)
	return count, err
}

func (r *ContestRepo) UpdateContest(ctx context.Context, c *contest.Contest) (*contest.Contest, error) {
	const q = `
		UPDATE contests
		SET name = $2, description = $3, is_up = $4, max_votes_user = $5
		WHERE public_id = $1
		RETURNING id, public_id, user_id, name, description, is_up, max_votes_user
	`
	err := r.db.Pool.QueryRow(ctx, q,
		c.PublicID, c.Name, c.Description, c.IsUp, c.MaxVotesUser,
	).Scan(&c.ID, &c.PublicID, &c.UserId, &c.Name, &c.Description, &c.IsUp, &c.MaxVotesUser)
	if err != nil {
		return nil, err
	}
	return c, nil
}
