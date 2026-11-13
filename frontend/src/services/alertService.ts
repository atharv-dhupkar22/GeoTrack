
import api from './api'
import type {
  AlertConfig,
  AlertHistory,
  ApiResponse,
  CreateAlertConfigRequest,
} from '../types'

export const alertService = {
  async createConfiguration(
    request: CreateAlertConfigRequest,
  ): Promise<AlertConfig> {
    const response = await api.post<ApiResponse<AlertConfig>>(
      '/alerts/configure',
      request,
    )
    return response.data.data
  },

  async getConfigurations(): Promise<AlertConfig[]> {
    const response = await api.get<ApiResponse<AlertConfig[]>>(
      '/alerts/configurations',
    )
    return response.data.data
  },

  async getHistory(params?: {
    vehicle_id?: string
    geofence_id?: string
    event_type?: 'entry' | 'exit'
    page?: number
    limit?: number
  }): Promise<AlertHistory> {
    const response = await api.get<ApiResponse<AlertHistory>>('/alerts', {
      params,
    })
    return response.data.data
  },
}
