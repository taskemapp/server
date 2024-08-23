package mapper

import (
	"github.com/samber/lo"
	"github.com/taskemapp/server/apps/server/internal/repositories/team"
	v1 "github.com/taskemapp/server/apps/server/tools/gen/grpc/v1"
)

func ToTeamResponse(t *team.Team) *v1.TeamResponse {
	return &v1.TeamResponse{
		Id:          t.ID.String(),
		Name:        t.Name,
		Description: t.Description,
		Creator:     t.Creator.String(),
	}
}

func ToGetAllTeamsResponse(t *[]team.Team) *v1.GetAllTeamsResponse {
	teams := lo.Map(*t, func(item team.Team, index int) *v1.TeamResponse {
		return ToTeamResponse(&item)
	})
	return &v1.GetAllTeamsResponse{
		Teams: teams,
	}
}

func ToCreateTeamResponse(t *team.Team) *v1.CreateTeamResponse {
	return &v1.CreateTeamResponse{
		TeamId:  t.ID.String(),
		Message: "Successfully created a team",
	}
}
