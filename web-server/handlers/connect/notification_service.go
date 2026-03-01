package connect

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"

	"rootstock/web-server/auth"
	notificationflows "rootstock/web-server/flows/notification"
	userflows "rootstock/web-server/flows/user"
	rootstockv1 "rootstock/web-server/proto/rootstock/v1"
)

// NotificationServiceHandler implements the NotificationService Connect RPC interface.
type NotificationServiceHandler struct {
	getUser           *userflows.GetUserFlow
	listNotifications *notificationflows.ListNotificationsFlow
	markRead          *notificationflows.MarkReadFlow
	getPreferences    *notificationflows.GetPreferencesFlow
	updatePreferences *notificationflows.UpdatePreferencesFlow
}

// NewNotificationServiceHandler creates the handler with all required flows.
func NewNotificationServiceHandler(
	getUser *userflows.GetUserFlow,
	listNotifications *notificationflows.ListNotificationsFlow,
	markRead *notificationflows.MarkReadFlow,
	getPreferences *notificationflows.GetPreferencesFlow,
	updatePreferences *notificationflows.UpdatePreferencesFlow,
) *NotificationServiceHandler {
	return &NotificationServiceHandler{
		getUser:           getUser,
		listNotifications: listNotifications,
		markRead:          markRead,
		getPreferences:    getPreferences,
		updatePreferences: updatePreferences,
	}
}

// resolveUserID extracts the IdP user ID from context and resolves the app user ID.
func (h *NotificationServiceHandler) resolveUserID(ctx context.Context) (string, error) {
	idpID, ok := auth.SubjectFromContext(ctx)
	if !ok || idpID == "" {
		return "", connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("no authenticated subject"))
	}
	user, err := h.getUser.Run(ctx, idpID)
	if err != nil {
		return "", connect.NewError(connect.CodeInternal, fmt.Errorf("resolve user: %w", err))
	}
	if user == nil {
		return "", connect.NewError(connect.CodeNotFound, fmt.Errorf("user not found"))
	}
	return user.ID, nil
}

func (h *NotificationServiceHandler) ListNotifications(
	ctx context.Context,
	req *connect.Request[rootstockv1.ListNotificationsRequest],
) (*connect.Response[rootstockv1.ListNotificationsResponse], error) {
	userID, err := h.resolveUserID(ctx)
	if err != nil {
		return nil, err
	}

	msg := req.Msg
	var typeFilter *string
	if tf := msg.GetTypeFilter(); tf != "" {
		typeFilter = &tf
	}

	result, err := h.listNotifications.Run(ctx, notificationflows.ListInput{
		UserID:     userID,
		TypeFilter: typeFilter,
		Limit:      int(msg.GetLimit()),
		Offset:     int(msg.GetOffset()),
	})
	if err != nil {
		return nil, err
	}

	protos := make([]*rootstockv1.NotificationProto, len(result.Notifications))
	for i, n := range result.Notifications {
		protos[i] = &rootstockv1.NotificationProto{
			Id:           n.ID,
			Type:         n.Type,
			Message:      n.Message,
			Read:         n.Read,
			ResourceLink: n.ResourceLink,
			CreatedAt:    n.CreatedAt.Format(time.RFC3339),
		}
	}

	return connect.NewResponse(&rootstockv1.ListNotificationsResponse{
		Notifications: protos,
		UnreadCount:   int32(result.UnreadCount),
		Total:         int32(result.Total),
	}), nil
}

func (h *NotificationServiceHandler) MarkRead(
	ctx context.Context,
	req *connect.Request[rootstockv1.MarkReadRequest],
) (*connect.Response[rootstockv1.MarkReadResponse], error) {
	userID, err := h.resolveUserID(ctx)
	if err != nil {
		return nil, err
	}

	result, err := h.markRead.Run(ctx, notificationflows.MarkReadInput{
		UserID:          userID,
		NotificationIDs: req.Msg.GetNotificationIds(),
	})
	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&rootstockv1.MarkReadResponse{
		MarkedCount: int32(result.MarkedCount),
	}), nil
}

func (h *NotificationServiceHandler) GetPreferences(
	ctx context.Context,
	req *connect.Request[rootstockv1.GetPreferencesRequest],
) (*connect.Response[rootstockv1.GetPreferencesResponse], error) {
	userID, err := h.resolveUserID(ctx)
	if err != nil {
		return nil, err
	}

	prefs, err := h.getPreferences.Run(ctx, userID)
	if err != nil {
		return nil, err
	}

	protos := make([]*rootstockv1.NotificationPreferenceProto, len(prefs))
	for i, p := range prefs {
		protos[i] = &rootstockv1.NotificationPreferenceProto{
			Type:  p.Type,
			InApp: p.InApp,
			Email: p.Email,
		}
	}

	return connect.NewResponse(&rootstockv1.GetPreferencesResponse{
		Preferences: protos,
	}), nil
}

func (h *NotificationServiceHandler) UpdatePreferences(
	ctx context.Context,
	req *connect.Request[rootstockv1.UpdatePreferencesRequest],
) (*connect.Response[rootstockv1.UpdatePreferencesResponse], error) {
	userID, err := h.resolveUserID(ctx)
	if err != nil {
		return nil, err
	}

	prefs := make([]notificationflows.PreferenceInput, len(req.Msg.GetPreferences()))
	for i, p := range req.Msg.GetPreferences() {
		prefs[i] = notificationflows.PreferenceInput{
			Type:  p.GetType(),
			InApp: p.GetInApp(),
			Email: p.GetEmail(),
		}
	}

	if err := h.updatePreferences.Run(ctx, notificationflows.UpdatePreferencesInput{
		UserID:      userID,
		Preferences: prefs,
	}); err != nil {
		return nil, err
	}

	return connect.NewResponse(&rootstockv1.UpdatePreferencesResponse{}), nil
}
