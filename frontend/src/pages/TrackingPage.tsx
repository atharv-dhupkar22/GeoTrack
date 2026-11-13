
import { useCallback, useEffect, useState } from 'react'
import type { ChangeEvent, FormEvent } from 'react'
import { Link } from 'react-router-dom'
import {
  ArrowLeft,
  MapPin,
  RefreshCw,
  Navigation,
  Radio,
  Truck,
  AlertCircle,
  CheckCircle2,
} from 'lucide-react'
import {
  CircleMarker,
  MapContainer,
  Popup,
  TileLayer,
  useMapEvents,
} from 'react-leaflet'
import toast from 'react-hot-toast'

import 'leaflet/dist/leaflet.css'

import { vehicleService } from '../services/vehicleService'
import type { Vehicle, VehicleLocation } from '../types'

const DEFAULT_CENTER: [number, number] = [19.076, 72.8777]

type LocationPickerProps = {
  onSelect: (latitude: number, longitude: number) => void
}

function LocationPicker({ onSelect }: LocationPickerProps) {
  useMapEvents({
    click(event) {
      onSelect(event.latlng.lat, event.latlng.lng)
    },
  })

  return null
}

function TrackingPage() {
  const [vehicles, setVehicles] = useState<Vehicle[]>([])
  const [locations, setLocations] = useState<VehicleLocation[]>([])
  const [selectedVehicleId, setSelectedVehicleId] = useState('')
  const [latitude, setLatitude] = useState('19.076')
  const [longitude, setLongitude] = useState('72.8777')
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  const loadTrackingData = useCallback(async () => {
    setLoading(true)
    setError('')

    try {
      const vehicleList = await vehicleService.getAll()
      setVehicles(vehicleList)

      if (vehicleList.length === 0) {
        setLocations([])
        setSelectedVehicleId('')
        return
      }

      const results = await Promise.allSettled(
        vehicleList.map((vehicle) =>
          vehicleService.getLatestLocation(vehicle.id),
        ),
      )

      const latestLocations = results
        .filter(
          (result): result is PromiseFulfilledResult<VehicleLocation> =>
            result.status === 'fulfilled',
        )
        .map((result) => result.value)

      setLocations(latestLocations)

      setSelectedVehicleId((current) => {
        if (current && vehicleList.some((vehicle) => vehicle.id === current)) {
          return current
        }

        return vehicleList[0].id
      })
    } catch {
      setError('Unable to load tracking data. Check the backend connection.')
      toast.error('Could not load tracking data')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadTrackingData()
  }, [loadTrackingData])

  const handleVehicleChange = (event: ChangeEvent<HTMLSelectElement>) => {
    setSelectedVehicleId(event.target.value)
  }

  const handleMapSelect = (lat: number, lng: number) => {
    setLatitude(lat.toFixed(6))
    setLongitude(lng.toFixed(6))
  }

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    if (!selectedVehicleId) {
      toast.error('Please select a vehicle')
      return
    }

    const lat = Number(latitude)
    const lng = Number(longitude)

    if (
      !Number.isFinite(lat) ||
      !Number.isFinite(lng) ||
      lat < -90 ||
      lat > 90 ||
      lng < -180 ||
      lng > 180
    ) {
      toast.error('Enter valid latitude and longitude values')
      return
    }

    try {
      setSaving(true)

      await vehicleService.updateLocation({
        vehicle_id: selectedVehicleId,
        latitude: lat,
        longitude: lng,
      })

      toast.success('Vehicle location updated successfully')

      await loadTrackingData()
    } catch {
      toast.error('Unable to update vehicle location')
    } finally {
      setSaving(false)
    }
  }

  const selectedVehicle = vehicles.find(
    (vehicle) => vehicle.id === selectedVehicleId,
  )

  const selectedLocation = locations.find(
    (location) => location.vehicle_id === selectedVehicleId,
  )

  const formatDate = (value: string) => {
    const date = new Date(value)

    if (Number.isNaN(date.getTime())) {
      return 'Unknown time'
    }

    return date.toLocaleString('en-IN')
  }

  return (
    <main className="min-h-screen bg-slate-50 px-4 py-6 text-slate-800 sm:px-8">
      <div className="mx-auto max-w-7xl">
        <header className="mb-7 flex flex-wrap items-center justify-between gap-4">
          <div>
            <Link
              to="/"
              className="mb-4 inline-flex items-center gap-2 text-sm font-medium text-blue-600 hover:text-blue-800"
            >
              <ArrowLeft size={16} />
              Back to Dashboard
            </Link>

            <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">
              Live GPS Tracking
            </h1>

            <p className="mt-1 text-sm text-slate-500">
              Monitor vehicle locations and submit GPS updates.
            </p>
          </div>

          <button
            type="button"
            onClick={() => void loadTrackingData()}
            disabled={loading}
            className="inline-flex items-center gap-2 rounded-lg border border-slate-200 bg-white px-4 py-2.5 text-sm font-semibold shadow-sm hover:bg-slate-100 disabled:opacity-50"
          >
            <RefreshCw
              size={16}
              className={loading ? 'animate-spin' : ''}
            />
            Refresh Locations
          </button>
        </header>

        {error && (
          <div className="mb-5 flex items-center gap-2 rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700">
            <AlertCircle size={18} />
            {error}
          </div>
        )}

        <div className="mb-6 grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
            <div className="flex items-center gap-3">
              <div className="rounded-lg bg-blue-50 p-3 text-blue-600">
                <Truck size={21} />
              </div>
              <div>
                <p className="text-sm text-slate-500">Registered Vehicles</p>
                <p className="text-2xl font-bold">
                  {loading ? '—' : vehicles.length}
                </p>
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
            <div className="flex items-center gap-3">
              <div className="rounded-lg bg-green-50 p-3 text-green-600">
                <Navigation size={21} />
              </div>
              <div>
                <p className="text-sm text-slate-500">Vehicles With Locations</p>
                <p className="text-2xl font-bold">
                  {loading ? '—' : locations.length}
                </p>
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
            <div className="flex items-center gap-3">
              <div className="rounded-lg bg-purple-50 p-3 text-purple-600">
                <Radio size={21} />
              </div>
              <div>
                <p className="text-sm text-slate-500">Tracking Service</p>
                <p className="text-lg font-bold text-green-600">
                  {loading ? 'Checking...' : 'Connected'}
                </p>
              </div>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
          <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm lg:col-span-1">
            <div className="mb-5 flex items-center gap-3">
              <div className="rounded-xl bg-blue-50 p-3 text-blue-600">
                <MapPin size={21} />
              </div>

              <div>
                <h2 className="font-bold">Update Vehicle Location</h2>
                <p className="text-xs text-slate-500">
                  Select a vehicle and submit GPS coordinates
                </p>
              </div>
            </div>

            {vehicles.length === 0 && !loading ? (
              <div className="rounded-xl bg-amber-50 p-4 text-sm text-amber-800">
                Register a vehicle before submitting GPS updates.
              </div>
            ) : (
              <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                  <label
                    htmlFor="tracking_vehicle"
                    className="mb-1.5 block text-sm font-medium"
                  >
                    Select Vehicle
                  </label>

                  <select
                    id="tracking_vehicle"
                    value={selectedVehicleId}
                    onChange={handleVehicleChange}
                    required
                    disabled={loading}
                    className="w-full rounded-lg border border-slate-200 bg-white px-3 py-2.5 text-sm outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-100"
                  >
                    {vehicles.map((vehicle) => (
                      <option key={vehicle.id} value={vehicle.id}>
                        {vehicle.vehicle_number} — {vehicle.driver_name}
                      </option>
                    ))}
                  </select>
                </div>

                {selectedVehicle && (
                  <div className="rounded-lg bg-slate-50 p-3 text-sm">
                    <p className="font-semibold">
                      {selectedVehicle.vehicle_number}
                    </p>
                    <p className="mt-1 text-slate-500">
                      {selectedVehicle.driver_name} · {selectedVehicle.vehicle_type}
                    </p>
                  </div>
                )}

                <div>
                  <label
                    htmlFor="tracking_latitude"
                    className="mb-1.5 block text-sm font-medium"
                  >
                    Latitude
                  </label>

                  <input
                    id="tracking_latitude"
                    type="number"
                    step="any"
                    min="-90"
                    max="90"
                    value={latitude}
                    onChange={(event) => setLatitude(event.target.value)}
                    required
                    className="w-full rounded-lg border border-slate-200 px-3 py-2.5 text-sm outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-100"
                  />
                </div>

                <div>
                  <label
                    htmlFor="tracking_longitude"
                    className="mb-1.5 block text-sm font-medium"
                  >
                    Longitude
                  </label>

                  <input
                    id="tracking_longitude"
                    type="number"
                    step="any"
                    min="-180"
                    max="180"
                    value={longitude}
                    onChange={(event) => setLongitude(event.target.value)}
                    required
                    className="w-full rounded-lg border border-slate-200 px-3 py-2.5 text-sm outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-100"
                  />
                </div>

                <p className="rounded-lg bg-blue-50 p-3 text-xs leading-5 text-blue-700">
                  You can also click any point on the map to automatically fill
                  in its latitude and longitude.
                </p>

                <button
                  type="submit"
                  disabled={saving || loading || !selectedVehicleId}
                  className="flex w-full items-center justify-center gap-2 rounded-lg bg-blue-600 px-4 py-3 text-sm font-semibold text-white shadow-sm hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-60"
                >
                  {saving ? (
                    <RefreshCw size={17} className="animate-spin" />
                  ) : (
                    <Navigation size={17} />
                  )}
                  {saving ? 'Updating Location...' : 'Update GPS Location'}
                </button>
              </form>
            )}

            {selectedLocation && (
              <div className="mt-6 border-t border-slate-100 pt-5">
                <div className="mb-3 flex items-center gap-2 text-sm font-bold">
                  <CheckCircle2 size={17} className="text-green-600" />
                  Latest Recorded Location
                </div>

                <div className="space-y-2 rounded-xl bg-green-50 p-4 text-sm">
                  <p>
                    <span className="text-slate-500">Latitude: </span>
                    <strong>{selectedLocation.latitude}</strong>
                  </p>
                  <p>
                    <span className="text-slate-500">Longitude: </span>
                    <strong>{selectedLocation.longitude}</strong>
                  </p>
                  <p className="text-xs text-slate-500">
                    Recorded: {formatDate(selectedLocation.recorded_at)}
                  </p>
                </div>
              </div>
            )}
          </section>

          <section className="overflow-hidden rounded-2xl border border-slate-200 bg-white p-5 shadow-sm lg:col-span-2">
            <div className="mb-4 flex flex-wrap items-center justify-between gap-2">
              <div>
                <h2 className="font-bold">Vehicle Tracking Map</h2>
                <p className="mt-1 text-xs text-slate-500">
                  Click the map to choose a GPS location
                </p>
              </div>

              <span className="inline-flex items-center gap-2 rounded-full bg-green-50 px-3 py-1.5 text-xs font-semibold text-green-700">
                <span className="h-2 w-2 rounded-full bg-green-500" />
                Map available
              </span>
            </div>

            <div
              className="overflow-hidden rounded-xl border border-slate-200"
              style={{ height: '520px', width: '100%' }}
            >
              <MapContainer
                center={DEFAULT_CENTER}
                zoom={12}
                scrollWheelZoom
                style={{ height: '100%', width: '100%' }}
              >
                <TileLayer
                  attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                  url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />

                <LocationPicker onSelect={handleMapSelect} />

                {locations.map((location) => {
                  const vehicle = vehicles.find(
                    (item) => item.id === location.vehicle_id,
                  )

                  return (
                    <CircleMarker
                      key={`${location.vehicle_id}-${location.id}`}
                      center={[location.latitude, location.longitude]}
                      radius={9}
                      pathOptions={{
                        color: '#2563eb',
                        fillColor: '#3b82f6',
                        fillOpacity: 0.85,
                        weight: 2,
                      }}
                    >
                      <Popup>
                        <strong>
                          {vehicle?.vehicle_number || location.vehicle_id}
                        </strong>
                        <br />
                        {vehicle?.driver_name || 'Vehicle'}
                        <br />
                        Latitude: {location.latitude}
                        <br />
                        Longitude: {location.longitude}
                        <br />
                        {formatDate(location.recorded_at)}
                      </Popup>
                    </CircleMarker>
                  )
                })}

                <CircleMarker
                  center={[
                    Number(latitude) || DEFAULT_CENTER[0],
                    Number(longitude) || DEFAULT_CENTER[1],
                  ]}
                  radius={8}
                  pathOptions={{
                    color: '#dc2626',
                    fillColor: '#ef4444',
                    fillOpacity: 0.9,
                    weight: 2,
                  }}
                >
                  <Popup>Selected GPS location</Popup>
                </CircleMarker>
              </MapContainer>
            </div>

            <div className="mt-4 flex flex-wrap items-center gap-5 text-xs text-slate-600">
              <span className="inline-flex items-center gap-2">
                <span className="h-3 w-3 rounded-full bg-blue-500" />
                Recorded vehicle locations
              </span>
              <span className="inline-flex items-center gap-2">
                <span className="h-3 w-3 rounded-full bg-red-500" />
                Selected GPS location
              </span>
            </div>
          </section>
        </div>

        <section className="mt-6 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <h2 className="mb-4 font-bold">Recorded Vehicle Locations</h2>

          {loading ? (
            <p className="py-6 text-center text-sm text-slate-500">
              Loading vehicle locations...
            </p>
          ) : locations.length === 0 ? (
            <p className="py-6 text-center text-sm text-slate-500">
              No GPS locations have been recorded yet. Select a vehicle and
              submit its first location.
            </p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-[600px] text-left text-sm">
                <thead>
                  <tr className="border-b border-slate-100 text-xs uppercase text-slate-500">
                    <th className="px-3 py-3">Vehicle</th>
                    <th className="px-3 py-3">Latitude</th>
                    <th className="px-3 py-3">Longitude</th>
                    <th className="px-3 py-3">Recorded At</th>
                  </tr>
                </thead>

                <tbody>
                  {locations.map((location) => {
                    const vehicle = vehicles.find(
                      (item) => item.id === location.vehicle_id,
                    )

                    return (
                      <tr
                        key={`${location.vehicle_id}-${location.id}`}
                        className="border-b border-slate-50"
                      >
                        <td className="px-3 py-3 font-semibold">
                          {vehicle?.vehicle_number || location.vehicle_id}
                        </td>
                        <td className="px-3 py-3">{location.latitude}</td>
                        <td className="px-3 py-3">{location.longitude}</td>
                        <td className="px-3 py-3 text-slate-500">
                          {formatDate(location.recorded_at)}
                        </td>
                      </tr>
                    )
                  })}
                </tbody>
              </table>
            </div>
          )}
        </section>

        <footer className="mt-6 text-center text-xs text-slate-400">
          GeoTrack · Live GPS Tracking
        </footer>
      </div>
    </main>
  )
}

export default TrackingPage
