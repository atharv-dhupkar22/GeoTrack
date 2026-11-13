
import api from './api'
import type {
  ApiResponse,
  CreateLocationRequest,
  CreateVehicleRequest,
  Vehicle,
  VehicleLocation,
} from '../types'

export const vehicleService = {
  async getAll(): Promise<Vehicle[]> {
    const response = await api.get<ApiResponse<Vehicle[]>>('/vehicles')
    return response.data.data
  },

  async create(request: CreateVehicleRequest): Promise<Vehicle> {
    const response = await api.post<ApiResponse<Vehicle>>('/vehicles', request)
    return response.data.data
  },

  async updateLocation(request: CreateLocationRequest): Promise<VehicleLocation> {
    const response = await api.post<ApiResponse<VehicleLocation>>(
      '/vehicles/location',
      request,
    )
    return response.data.data
  },

  async getLatestLocation(vehicleId: string): Promise<VehicleLocation> {
    const response = await api.get<ApiResponse<VehicleLocation>>(
      `/vehicles/location/${encodeURIComponent(vehicleId)}`,
    )
    return response.data.data
  },
}
