package connect

import (
	"context"
	"time"

	"connectrpc.com/connect"

	scoreflows "rootstock/web-server/flows/score"
	rootstockv1 "rootstock/web-server/proto/rootstock/v1"
)

// ScoreServiceHandler implements the ScoreService Connect RPC interface.
type ScoreServiceHandler struct {
	getContribution *scoreflows.GetContributionFlow
}

// NewScoreServiceHandler creates the handler with all required flows.
func NewScoreServiceHandler(
	getContribution *scoreflows.GetContributionFlow,
) *ScoreServiceHandler {
	return &ScoreServiceHandler{
		getContribution: getContribution,
	}
}

func (h *ScoreServiceHandler) GetContribution(
	ctx context.Context,
	req *connect.Request[rootstockv1.GetContributionRequest],
) (*connect.Response[rootstockv1.GetContributionResponse], error) {
	contribution, err := h.getContribution.Run(ctx, req.Msg.GetScitizenId())
	if err != nil {
		return nil, err
	}

	badges := make([]*rootstockv1.BadgeProto, len(contribution.Badges))
	for i, b := range contribution.Badges {
		badges[i] = &rootstockv1.BadgeProto{
			Id:        b.ID,
			BadgeType: b.BadgeType,
			AwardedAt: b.AwardedAt.Format(time.RFC3339),
		}
	}

	return connect.NewResponse(&rootstockv1.GetContributionResponse{
		ScitizenId:  contribution.ScitizenID,
		Volume:      int32(contribution.Volume),
		QualityRate: contribution.QualityRate,
		Consistency: contribution.Consistency,
		Diversity:   int32(contribution.Diversity),
		Total:       contribution.Total,
		UpdatedAt:   contribution.UpdatedAt.Format(time.RFC3339),
		Badges:      badges,
	}), nil
}
