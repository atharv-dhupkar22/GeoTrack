
export interface ApiResponse<T> {
  status: string
  message: string
  data: T
  time_ns: string
}

export interface Geofence {
  id: string
  name: string
  description: string
  coordinates: number[][]
  category: 'delivery_zone' | 'restricted_zone' | 'toll_zone' | 'customer_area'
  status: string
  created_at: string
}

export interface Vehicle {
  id: string
  vehicle_number: string
  driver_name: string
  vehicle_type: string
  phone: string
  status: string
  created_at: string
}

export interface VehicleLocation {
  id: number
  vehicle_id: string
  latitude: number
  longitude: number
  recorded_at: string
  created_at: string
}

export interface AlertEvent {
  id: string
  violation_id: string
  alert_config_id: string
  vehicle_id: string
  geofence_id: string
  event_type: 'entry' | 'exit'
  message: string
  created_at: string
}

export interface AlertConfig {
  id: string
  geofence_id: string
  vehicle_id?: string
  event_type: 'entry' | 'exit' | 'both'
  status: string
  created_at: string
}

export interface Violation {
  id: string
  vehicle_id: string
  geofence_id: string
  event_type: 'entry' | 'exit'
  latitude: number
  longitude: number
  recorded_at: string
  created_at: string
}

export interface ViolationHistory {
  violations: Violation[]
  page: number
  limit: number
  total: number
  total_pages: number
}

export interface AlertHistory {
  alerts: AlertEvent[]
  page: number
  limit: number
  total: number
  total_pages: number
}

export interface CreateGeofenceRequest {
  name: string
  description: string
  coordinates: number[][]
  category: Geofence['category']
}

export interface CreateVehicleRequest {
  vehicle_number: string
  driver_name: string
  vehicle_type: string
  phone: string
}

export interface CreateLocationRequest {
  vehicle_id: string
  latitude: number
  longitude: number
  recorded_at?: string
}

export interface CreateAlertConfigRequest {
  geofence_id: string
  vehicle_id?: string
  event_type: 'entry' | 'exit' | 'both'
}
