import { createClient } from '@connectrpc/connect';
import { transport } from './transport';
import { UserService, HealthService, CampaignService, ScitizenService, NotificationService } from './gen/rootstock/v1/rootstock_connect';

export const userService = createClient(UserService, transport);
export const healthService = createClient(HealthService, transport);
export const campaignService = createClient(CampaignService, transport);
export const scitizenService = createClient(ScitizenService, transport);
export const notificationService = createClient(NotificationService, transport);
