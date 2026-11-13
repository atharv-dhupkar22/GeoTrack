
import api from './api'
import type {
  ApiResponse,
  ViolationHistory,
} from '../types'

export interface ViolationFilters {
  vehicle_id?: string
  geofence_id?: string
  event_type?: 'entry' | 'exit'
  page?: number
  limit?: number
}

export const violationService = {
  async getHistory(
    filters?: ViolationFilters,
  ): Promise<ViolationHistory> {
    const response = await api.get<ApiResponse<ViolationHistory>>(
      '/violations/history',
      { params: filters },
    )

    return response.data.data
  },
}
