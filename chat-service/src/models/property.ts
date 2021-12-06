import { ITenant } from '.';

export interface IProperty {
  server_code: string;
  channels: string[];
  tenants: ITenant[];
}
