
import api from './api'
import type {
  ApiResponse,
  CreateGeofenceRequest,
  Geofence,
} from '../types'

export const geofenceService = {
  async getAll(category?: string): Promise<Geofence[]> {
    const response = await api.get<ApiResponse<Geofence[]>>('/geofences', {
      params: category ? { category } : {},
    })
    return response.data.data
  },

  async create(request: CreateGeofenceRequest): Promise<Geofence> {
    const response = await api.post<ApiResponse<Geofence>>(
      '/geofences',
      request,
    )
    return response.data.data
  },
}
