import { createClient } from '@connectrpc/connect';
import { transport } from './transport';
import { UserService } from './gen/rootstock/v1/rootstock_connect';
import { HealthService } from './gen/rootstock/v1/rootstock_connect';

export const userService = createClient(UserService, transport);
export const healthService = createClient(HealthService, transport);
