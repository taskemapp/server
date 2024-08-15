package mapper

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/taskemapp/server/apps/server/internal/repositories/team"
	v1 "github.com/taskemapp/server/apps/server/tools/gen/grpc/v1"
	"reflect"
	"testing"
)

func TestToGetAllTeamsResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    *team.Team
		expected *v1.TeamResponse
	}{
		{
			name: "Basic case",
			input: &team.Team{
				ID:          uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name:        "Team A",
				Description: "Description A",
				Creator:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			},
			expected: &v1.TeamResponse{
				Id:          "123e4567-e89b-12d3-a456-426614174000",
				Name:        "Team A",
				Description: "Description A",
				Creator:     "123e4567-e89b-12d3-a456-426614174001",
			},
		},
		{
			name: "Empty fields",
			input: &team.Team{
				ID:          uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
				Name:        "",
				Description: "",
				Creator:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174003"),
			},
			expected: &v1.TeamResponse{
				Id:          "123e4567-e89b-12d3-a456-426614174002",
				Name:        "",
				Description: "",
				Creator:     "123e4567-e89b-12d3-a456-426614174003",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToTeamResponse(tt.input)
			expr := reflect.DeepEqual(result, tt.expected)
			assert.True(t, expr, result)
		})
	}
}

func TestToTeamResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    *[]team.Team
		expected *v1.GetAllTeamsResponse
	}{
		{
			name: "Single team",
			input: &[]team.Team{
				{
					ID:          uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Name:        "Team A",
					Description: "Description A",
					Creator:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
				},
			},
			expected: &v1.GetAllTeamsResponse{
				Teams: []*v1.TeamResponse{
					{
						Id:          "123e4567-e89b-12d3-a456-426614174000",
						Name:        "Team A",
						Description: "Description A",
						Creator:     "123e4567-e89b-12d3-a456-426614174001",
					},
				},
			},
		},
		{
			name: "Multiple teams",
			input: &[]team.Team{
				{
					ID:          uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					Name:        "Team A",
					Description: "Description A",
					Creator:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
				},
				{
					ID:          uuid.MustParse("123e4567-e89b-12d3-a456-426614174002"),
					Name:        "Team B",
					Description: "Description B",
					Creator:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174003"),
				},
				{
					ID:          uuid.MustParse("123e4567-e89b-12d3-a456-426614174004"),
					Name:        "Team C",
					Description: "Description C",
					Creator:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174005"),
				},
			},
			expected: &v1.GetAllTeamsResponse{
				Teams: []*v1.TeamResponse{
					{
						Id:          "123e4567-e89b-12d3-a456-426614174000",
						Name:        "Team A",
						Description: "Description A",
						Creator:     "123e4567-e89b-12d3-a456-426614174001",
					},
					{
						Id:          "123e4567-e89b-12d3-a456-426614174002",
						Name:        "Team B",
						Description: "Description B",
						Creator:     "123e4567-e89b-12d3-a456-426614174003",
					},
					{
						Id:          "123e4567-e89b-12d3-a456-426614174004",
						Name:        "Team C",
						Description: "Description C",
						Creator:     "123e4567-e89b-12d3-a456-426614174005",
					},
				},
			},
		},
		{
			name:  "No teams",
			input: &[]team.Team{},
			expected: &v1.GetAllTeamsResponse{
				Teams: []*v1.TeamResponse{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToGetAllTeamsResponse(tt.input)
			expr := reflect.DeepEqual(result, tt.expected)
			assert.True(t, expr, result)
		})
	}
}

func TestToCreateTeamResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    *team.Team
		expected *v1.CreateTeamResponse
	}{
		{
			name: "Single team",
			input: &team.Team{
				ID:          uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
				Name:        "Team A",
				Description: "Description A",
				Creator:     uuid.MustParse("123e4567-e89b-12d3-a456-426614174001"),
			},
			expected: &v1.CreateTeamResponse{
				TeamId:  "123e4567-e89b-12d3-a456-426614174000",
				Message: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToCreateTeamResponse(tt.input)
			assert.Equal(t, tt.input.ID.String(), result.TeamId)
		})
	}
}
