
import { useCallback, useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { Link } from 'react-router-dom'
import {
  ArrowLeft,
  Bell,
  RefreshCw,
  ShieldCheck,
  Truck,
  MapPin,
  ArrowDownRight,
  ArrowUpRight,
  Settings2,
} from 'lucide-react'
import toast from 'react-hot-toast'

import { alertService } from '../services/alertService'
import { vehicleService } from '../services/vehicleService'
import { geofenceService } from '../services/geofenceService'

import type {
  AlertConfig,
  AlertEvent,
  Vehicle,
  Geofence,
} from '../types'

type AlertEventType = 'both' | 'entry' | 'exit'

function AlertsPage() {
  const [vehicles, setVehicles] = useState<Vehicle[]>([])
  const [geofences, setGeofences] = useState<Geofence[]>([])
  const [configurations, setConfigurations] = useState<AlertConfig[]>([])
  const [alerts, setAlerts] = useState<AlertEvent[]>([])

  const [vehicleId, setVehicleId] = useState('')
  const [geofenceId, setGeofenceId] = useState('')
  const [eventType, setEventType] = useState<AlertEventType>('both')

  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  const loadData = useCallback(async () => {
    setLoading(true)

    try {
      const results = await Promise.all([
        vehicleService.getAll(),
        geofenceService.getAll(),
        alertService.getConfigurations(),
        alertService.getHistory({ page: 1, limit: 50 }),
      ])

      setVehicles(results[0])
      setGeofences(results[1])
      setConfigurations(results[2])
      setAlerts(results[3].alerts)
    } catch {
      toast.error('Unable to load alert management data')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadData()
  }, [loadData])

  const handleConfigure = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    if (!vehicleId || !geofenceId) {
      toast.error('Select a vehicle and geofence')
      return
    }

    try {
      setSaving(true)

      await alertService.createConfiguration({
        vehicle_id: vehicleId,
        geofence_id: geofenceId,
        event_type: eventType,
      })

      toast.success('Alert configuration created successfully')

      setVehicleId('')
      setGeofenceId('')
      setEventType('both')

      await loadData()
    } catch {
      toast.error('Unable to create alert configuration')
    } finally {
      setSaving(false)
    }
  }

  const getVehicleName = (id: string | undefined) => {
    if (!id) {
      return 'All vehicles'
    }

    return (
      vehicles.find((vehicle) => vehicle.id === id)?.vehicle_number || id
    )
  }

  const getGeofenceName = (id: string | undefined) => {
    if (!id) {
      return 'All geofences'
    }

    return (
      geofences.find((geofence) => geofence.id === id)?.name || id
    )
  }

  const formatDate = (value: string) => {
    const date = new Date(value)

    return Number.isNaN(date.getTime())
      ? 'Unknown time'
      : date.toLocaleString('en-IN')
  }

  return (
    <main className="min-h-screen bg-slate-50 px-4 py-6 text-slate-800 sm:px-8">
      <div className="mx-auto max-w-7xl">

        {/* Page Header */}
        <header className="mb-8 flex flex-wrap items-center justify-between gap-4">
          <div>
            <Link
              to="/"
              className="mb-4 inline-flex items-center gap-2 text-sm font-medium text-blue-600 hover:text-blue-800"
            >
              <ArrowLeft size={16} />
              Back to Dashboard
            </Link>

            <h1 className="text-2xl font-bold sm:text-3xl">
              Alert Management
            </h1>

            <p className="mt-1 text-sm text-slate-500">
              Configure geofence notifications and monitor alert activity.
            </p>
          </div>

          <button
            type="button"
            onClick={() => void loadData()}
            disabled={loading}
            className="inline-flex items-center gap-2 rounded-lg border border-slate-200 bg-white px-4 py-2.5 text-sm font-semibold shadow-sm hover:bg-slate-100 disabled:opacity-50"
          >
            <RefreshCw
              size={16}
              className={loading ? 'animate-spin' : ''}
            />
            Refresh
          </button>
        </header>

        {/* Statistics */}
        <div className="mb-6 grid grid-cols-1 gap-4 sm:grid-cols-3">

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="flex items-center gap-3">
              <div className="rounded-xl bg-blue-50 p-3 text-blue-600">
                <Bell size={21} />
              </div>
              <div>
                <p className="text-sm text-slate-500">Notifications</p>
                <p className="text-2xl font-bold">
                  {loading ? '—' : alerts.length}
                </p>
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="flex items-center gap-3">
              <div className="rounded-xl bg-green-50 p-3 text-green-600">
                <ShieldCheck size={21} />
              </div>
              <div>
                <p className="text-sm text-slate-500">Configurations</p>
                <p className="text-2xl font-bold">
                  {loading ? '—' : configurations.length}
                </p>
              </div>
            </div>
          </div>

          <div className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
            <div className="flex items-center gap-3">
              <div className="rounded-xl bg-purple-50 p-3 text-purple-600">
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

        </div>

        {/* Configuration and Active Rules */}
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">

          {/* Configure Alerts */}
          <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm lg:col-span-1">

            <div className="mb-5 flex items-center gap-3">
              <div className="rounded-xl bg-blue-50 p-3 text-blue-600">
                <Settings2 size={21} />
              </div>

              <div>
                <h2 className="font-bold">Configure Alerts</h2>
                <p className="text-xs text-slate-500">
                  Create a notification rule
                </p>
              </div>
            </div>

            <form onSubmit={handleConfigure} className="space-y-4">

              <div>
                <label
                  htmlFor="alert-vehicle"
                  className="mb-1.5 block text-sm font-medium"
                >
                  Select Vehicle
                </label>

                <select
                  id="alert-vehicle"
                  value={vehicleId}
                  onChange={(event) => setVehicleId(event.target.value)}
                  required
                  disabled={loading}
                  className="w-full rounded-lg border border-slate-200 bg-white px-3 py-2.5 text-sm outline-none focus:border-blue-500"
                >
                  <option value="">Choose vehicle</option>

                  {vehicles.map((vehicle) => (
                    <option key={vehicle.id} value={vehicle.id}>
                      {vehicle.vehicle_number} — {vehicle.driver_name}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label
                  htmlFor="alert-geofence"
                  className="mb-1.5 block text-sm font-medium"
                >
                  Select Geofence
                </label>

                <select
                  id="alert-geofence"
                  value={geofenceId}
                  onChange={(event) => setGeofenceId(event.target.value)}
                  required
                  disabled={loading}
                  className="w-full rounded-lg border border-slate-200 bg-white px-3 py-2.5 text-sm outline-none focus:border-blue-500"
                >
                  <option value="">Choose geofence</option>

                  {geofences.map((geofence) => (
                    <option key={geofence.id} value={geofence.id}>
                      {geofence.name}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label
                  htmlFor="alert-event"
                  className="mb-1.5 block text-sm font-medium"
                >
                  Notification Event
                </label>

                <select
                  id="alert-event"
                  value={eventType}
                  onChange={(event) =>
                    setEventType(event.target.value as AlertEventType)
                  }
                  disabled={loading}
                  className="w-full rounded-lg border border-slate-200 bg-white px-3 py-2.5 text-sm outline-none focus:border-blue-500"
                >
                  <option value="both">Entry and Exit</option>
                  <option value="entry">Entry Only</option>
                  <option value="exit">Exit Only</option>
                </select>
              </div>

              <button
                type="submit"
                disabled={
                  saving ||
                  loading ||
                  vehicles.length === 0 ||
                  geofences.length === 0
                }
                className="flex w-full items-center justify-center gap-2 rounded-lg bg-blue-600 px-4 py-3 text-sm font-semibold text-white hover:bg-blue-700 disabled:opacity-50"
              >
                {saving ? (
                  <RefreshCw size={17} className="animate-spin" />
                ) : (
                  <Bell size={17} />
                )}

                {saving ? 'Saving...' : 'Create Alert Configuration'}
              </button>

            </form>
          </section>

          {/* Active Configurations */}
          <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm lg:col-span-2">

            <div className="mb-5 flex items-center gap-3">
              <div className="rounded-xl bg-green-50 p-3 text-green-600">
                <ShieldCheck size={21} />
              </div>

              <div>
                <h2 className="font-bold">Active Alert Configurations</h2>
                <p className="text-xs text-slate-500">
                  Configured vehicle and geofence rules
                </p>
              </div>
            </div>

            {loading ? (
              <p className="py-10 text-center text-sm text-slate-500">
                Loading configurations...
              </p>
            ) : configurations.length === 0 ? (
              <p className="py-10 text-center text-sm text-slate-500">
                No alert configurations found.
              </p>
            ) : (
              <div className="space-y-3">

                {configurations.map((configuration: AlertConfig) => (
                  <article
                    key={configuration.id}
                    className="rounded-xl border border-slate-100 bg-slate-50 p-4"
                  >
                    <div className="flex flex-wrap items-center justify-between gap-3">
                      <div>
                        <h3 className="font-bold">
                          {getVehicleName(configuration.vehicle_id)}
                        </h3>

                        <p className="mt-1 text-xs text-slate-500">
                          {getGeofenceName(configuration.geofence_id)}
                        </p>
                      </div>

                      <span className="rounded-full bg-green-50 px-3 py-1 text-xs font-semibold capitalize text-green-700">
                        {configuration.status}
                      </span>
                    </div>

                    <div className="mt-3 flex flex-wrap gap-2">
                      <span className="rounded-lg bg-white px-3 py-1.5 text-xs font-medium capitalize text-slate-600">
                        Event: {configuration.event_type}
                      </span>

                      <span className="rounded-lg bg-white px-3 py-1.5 text-xs font-medium text-slate-600">
                        ID: {configuration.id.slice(0, 14)}...
                      </span>
                    </div>
                  </article>
                ))}

              </div>
            )}

          </section>

        </div>

        {/* Notification History */}
        <section className="mt-6 rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">

          <div className="mb-5 flex flex-wrap items-center justify-between gap-3">
            <div>
              <h2 className="font-bold">Notification History</h2>
              <p className="mt-1 text-xs text-slate-500">
                Recent geofence entry and exit notifications
              </p>
            </div>

            <span className="rounded-full bg-blue-50 px-3 py-1.5 text-xs font-semibold text-blue-700">
              {alerts.length} Notifications
            </span>
          </div>

          {loading ? (
            <p className="py-10 text-center text-sm text-slate-500">
              Loading notifications...
            </p>
          ) : alerts.length === 0 ? (
            <div className="py-10 text-center text-sm text-slate-500">
              No notifications available.
            </div>
          ) : (
            <div className="space-y-3">

              {alerts.map((alert: AlertEvent) => (
                <article
                  key={alert.id}
                  className="flex flex-wrap items-start gap-3 rounded-xl border border-slate-100 p-4"
                >
                  <div
                    className={`rounded-xl p-3 ${
                      alert.event_type === 'entry'
                        ? 'bg-green-50 text-green-600'
                        : 'bg-orange-50 text-orange-600'
                    }`}
                  >
                    {alert.event_type === 'entry' ? (
                      <ArrowDownRight size={21} />
                    ) : (
                      <ArrowUpRight size={21} />
                    )}
                  </div>

                  <div className="min-w-0 flex-1">
                    <h3 className="text-sm font-bold">
                      Vehicle{' '}
                      {alert.event_type === 'entry' ? 'Entered' : 'Exited'}{' '}
                      Geofence
                    </h3>

                    <p className="mt-1 break-words text-sm text-slate-600">
                      {alert.message}
                    </p>

                    <div className="mt-2 flex flex-wrap gap-3 text-xs text-slate-400">
                      <span className="inline-flex items-center gap-1">
                        <Truck size={13} />
                        {getVehicleName(alert.vehicle_id)}
                      </span>

                      <span className="inline-flex items-center gap-1">
                        <MapPin size={13} />
                        {getGeofenceName(alert.geofence_id)}
                      </span>

                      <span>{formatDate(alert.created_at)}</span>
                    </div>
                  </div>

                  <span
                    className={`rounded-full px-2.5 py-1 text-xs font-semibold capitalize ${
                      alert.event_type === 'entry'
                        ? 'bg-green-50 text-green-700'
                        : 'bg-orange-50 text-orange-700'
                    }`}
                  >
                    {alert.event_type}
                  </span>
                </article>
              ))}

            </div>
          )}

        </section>

        <footer className="mt-6 text-center text-xs text-slate-400">
          GeoTrack · Alert Management System
        </footer>

      </div>
    </main>
  )
}

export default AlertsPage
