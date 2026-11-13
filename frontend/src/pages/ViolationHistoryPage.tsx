import { useCallback, useEffect, useMemo, useState } from 'react'
import {
  ArrowDownRight,
  ArrowUpRight,
  CalendarDays,
  ChevronLeft,
  ChevronRight,
  Filter,
  RefreshCw,
  ShieldAlert,
  XCircle,
} from 'lucide-react'
import toast from 'react-hot-toast'

import { geofenceService } from '../services/geofenceService'
import { vehicleService } from '../services/vehicleService'
import { violationService } from '../services/violationService'
import type { Geofence, Vehicle, Violation } from '../types'

type EventFilter = 'all' | 'entry' | 'exit'

interface ViolationResponse {
  violations: Violation[]
  total_count?: number
  total?: number
}

const PAGE_SIZE = 10

function formatDate(value: string | undefined) {
  if (!value) return '—'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return 'Unknown date'
  return date.toLocaleString('en-IN', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export default function ViolationHistoryPage() {
  const [violations, setViolations] = useState<Violation[]>([])
  const [vehicles, setVehicles] = useState<Vehicle[]>([])
  const [geofences, setGeofences] = useState<Geofence[]>([])
  const [vehicleFilter, setVehicleFilter] = useState('all')
  const [geofenceFilter, setGeofenceFilter] = useState('all')
  const [eventFilter, setEventFilter] = useState<EventFilter>('all')
  const [page, setPage] = useState(1)
  const [totalCount, setTotalCount] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const vehicleNames = useMemo(
    () => new Map(vehicles.map((vehicle) => [vehicle.id, vehicle.vehicle_number])),
    [vehicles],
  )
  const geofenceNames = useMemo(
    () => new Map(geofences.map((geofence) => [geofence.id, geofence.name])),
    [geofences],
  )

  const totalPages = Math.max(1, Math.ceil(totalCount / PAGE_SIZE))

  const loadOptions = useCallback(async () => {
    const [vehicleResult, geofenceResult] = await Promise.allSettled([
      vehicleService.getAll(),
      geofenceService.getAll(),
    ])
    if (vehicleResult.status === 'fulfilled') setVehicles(vehicleResult.value)
    if (geofenceResult.status === 'fulfilled') setGeofences(geofenceResult.value)
  }, [])

  const loadViolations = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const params: {
        page: number
        limit: number
        vehicle_id?: string
        geofence_id?: string
        event_type?: 'entry' | 'exit'
      } = {
        page,
        limit: PAGE_SIZE,
      }

      if (vehicleFilter !== 'all') params.vehicle_id = vehicleFilter
      if (geofenceFilter !== 'all') params.geofence_id = geofenceFilter
      if (eventFilter !== 'all') params.event_type = eventFilter

      const response = (await violationService.getHistory(
        params,
      )) as ViolationResponse

      setViolations(response.violations ?? [])
      setTotalCount(
        response.total_count ?? response.total ?? response.violations?.length ?? 0,
      )
    } catch {
      setError('Unable to load violation history. Check that the backend is running.')
      toast.error('Could not load violation history')
    } finally {
      setLoading(false)
    }
  }, [eventFilter, geofenceFilter, page, vehicleFilter])

  useEffect(() => {
    void loadOptions()
  }, [loadOptions])

  useEffect(() => {
    void loadViolations()
  }, [loadViolations])

  const resetFilters = () => {
    setVehicleFilter('all')
    setGeofenceFilter('all')
    setEventFilter('all')
    setPage(1)
  }

  const onVehicleChange = (value: string) => {
    setVehicleFilter(value)
    setPage(1)
  }

  const onGeofenceChange = (value: string) => {
    setGeofenceFilter(value)
    setPage(1)
  }

  const onEventChange = (value: EventFilter) => {
    setEventFilter(value)
    setPage(1)
  }

  const entries = violations.filter((item) => item.event_type === 'entry').length
  const exits = violations.filter((item) => item.event_type === 'exit').length

  return (
    <div className="min-h-screen bg-slate-50 px-4 py-6 text-slate-900 sm:px-6 lg:px-10">
      <div className="mx-auto max-w-7xl space-y-6">
        <header className="flex flex-col justify-between gap-4 sm:flex-row sm:items-center">
          <div>
            <div className="mb-2 flex items-center gap-2 text-sm font-semibold uppercase tracking-[0.16em] text-indigo-600">
              <ShieldAlert size={16} />
              GeoTrack / Monitoring
            </div>
            <h1 className="text-3xl font-bold tracking-tight">Violation History</h1>
            <p className="mt-1 text-sm text-slate-500">
              Review vehicle entry and exit events recorded by geofence monitoring.
            </p>
          </div>
          <button
            type="button"
            onClick={() => void loadViolations()}
            disabled={loading}
            className="inline-flex items-center justify-center gap-2 rounded-xl border border-slate-200 bg-white px-4 py-2.5 text-sm font-semibold shadow-sm transition hover:bg-slate-100 disabled:cursor-not-allowed disabled:opacity-60"
          >
            <RefreshCw size={16} className={loading ? 'animate-spin' : ''} />
            Refresh history
          </button>
        </header>

        <section className="grid gap-4 sm:grid-cols-3">
          <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
            <p className="text-sm font-medium text-slate-500">Matching records</p>
            <p className="mt-2 text-3xl font-bold">{totalCount}</p>
            <p className="mt-1 text-xs text-slate-400">Based on selected filters</p>
          </div>
          <div className="rounded-2xl border border-emerald-100 bg-white p-5 shadow-sm">
            <div className="flex items-center gap-2 text-sm font-medium text-emerald-700">
              <ArrowDownRight size={17} /> Entries on this page
            </div>
            <p className="mt-2 text-3xl font-bold text-slate-900">{entries}</p>
          </div>
          <div className="rounded-2xl border border-rose-100 bg-white p-5 shadow-sm">
            <div className="flex items-center gap-2 text-sm font-medium text-rose-700">
              <ArrowUpRight size={17} /> Exits on this page
            </div>
            <p className="mt-2 text-3xl font-bold text-slate-900">{exits}</p>
          </div>
        </section>

        <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <div className="mb-4 flex items-center gap-2">
            <Filter size={18} className="text-indigo-600" />
            <h2 className="font-semibold">Filter violation records</h2>
          </div>
          <div className="grid gap-3 md:grid-cols-4">
            <label className="text-sm font-medium text-slate-600">
              Vehicle
              <select
                value={vehicleFilter}
                onChange={(event) => onVehicleChange(event.target.value)}
                className="mt-1.5 w-full rounded-xl border border-slate-200 bg-white px-3 py-2.5 text-sm outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-100"
              >
                <option value="all">All vehicles</option>
                {vehicles.map((vehicle) => (
                  <option key={vehicle.id} value={vehicle.id}>
                    {vehicle.vehicle_number}
                  </option>
                ))}
              </select>
            </label>
            <label className="text-sm font-medium text-slate-600">
              Geofence
              <select
                value={geofenceFilter}
                onChange={(event) => onGeofenceChange(event.target.value)}
                className="mt-1.5 w-full rounded-xl border border-slate-200 bg-white px-3 py-2.5 text-sm outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-100"
              >
                <option value="all">All geofences</option>
                {geofences.map((geofence) => (
                  <option key={geofence.id} value={geofence.id}>
                    {geofence.name}
                  </option>
                ))}
              </select>
            </label>
            <label className="text-sm font-medium text-slate-600">
              Event type
              <select
                value={eventFilter}
                onChange={(event) => onEventChange(event.target.value as EventFilter)}
                className="mt-1.5 w-full rounded-xl border border-slate-200 bg-white px-3 py-2.5 text-sm outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-100"
              >
                <option value="all">All events</option>
                <option value="entry">Entry only</option>
                <option value="exit">Exit only</option>
              </select>
            </label>
            <div className="flex items-end">
              <button
                type="button"
                onClick={resetFilters}
                className="w-full rounded-xl border border-slate-200 px-4 py-2.5 text-sm font-semibold text-slate-600 transition hover:bg-slate-100"
              >
                Clear filters
              </button>
            </div>
          </div>
        </section>

        {error && (
          <div className="flex items-center gap-3 rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-700">
            <XCircle size={18} />
            <span>{error}</span>
            <button
              type="button"
              onClick={() => void loadViolations()}
              className="ml-auto font-semibold underline"
            >
              Retry
            </button>
          </div>
        )}

        <section className="overflow-hidden rounded-2xl border border-slate-200 bg-white shadow-sm">
          <div className="flex flex-col justify-between gap-2 border-b border-slate-100 px-5 py-4 sm:flex-row sm:items-center">
            <div>
              <h2 className="font-semibold">Recorded events</h2>
              <p className="text-xs text-slate-500">
                Showing {violations.length} record(s) on page {page}
              </p>
            </div>
            <div className="flex items-center gap-2 text-xs text-slate-500">
              <CalendarDays size={15} />
              Latest records from backend
            </div>
          </div>

          {loading ? (
            <div className="p-12 text-center text-sm text-slate-500">
              Loading violation history...
            </div>
          ) : violations.length === 0 ? (
            <div className="p-12 text-center">
              <ShieldAlert size={30} className="mx-auto text-slate-300" />
              <h3 className="mt-3 font-semibold">No violation records found</h3>
              <p className="mt-1 text-sm text-slate-500">
                Try changing the filters or record a vehicle location update.
              </p>
            </div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-[850px] text-left text-sm">
                <thead className="bg-slate-50 text-xs uppercase tracking-wide text-slate-500">
                  <tr>
                    <th className="px-5 py-3 font-semibold">Event</th>
                    <th className="px-5 py-3 font-semibold">Vehicle</th>
                    <th className="px-5 py-3 font-semibold">Geofence</th>
                    <th className="px-5 py-3 font-semibold">Coordinates</th>
                    <th className="px-5 py-3 font-semibold">Recorded at</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {violations.map((violation) => {
                    const isEntry = violation.event_type === 'entry'
                    const timestamp =
                      violation.recorded_at || violation.created_at
                    return (
                      <tr key={violation.id} className="transition hover:bg-slate-50">
                        <td className="px-5 py-4">
                          <span
                            className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-semibold ${
                              isEntry
                                ? 'bg-emerald-50 text-emerald-700'
                                : 'bg-rose-50 text-rose-700'
                            }`}
                          >
                            {isEntry ? (
                              <ArrowDownRight size={14} />
                            ) : (
                              <ArrowUpRight size={14} />
                            )}
                            {isEntry ? 'Entry' : 'Exit'}
                          </span>
                        </td>
                        <td className="px-5 py-4 font-semibold text-slate-800">
                          {vehicleNames.get(violation.vehicle_id) ?? violation.vehicle_id}
                        </td>
                        <td className="px-5 py-4 text-slate-600">
                          {geofenceNames.get(violation.geofence_id) ?? violation.geofence_id}
                        </td>
                        <td className="px-5 py-4 font-mono text-xs text-slate-600">
                          {Number(violation.latitude).toFixed(5)},{' '}
                          {Number(violation.longitude).toFixed(5)}
                        </td>
                        <td className="px-5 py-4 text-slate-600">
                          {formatDate(timestamp)}
                        </td>
                      </tr>
                    )
                  })}
                </tbody>
              </table>
            </div>
          )}

          <div className="flex flex-col justify-between gap-3 border-t border-slate-100 px-5 py-4 sm:flex-row sm:items-center">
            <p className="text-xs text-slate-500">
              Page {page} of {totalPages} · {totalCount} total record(s)
            </p>
            <div className="flex items-center gap-2">
              <button
                type="button"
                disabled={page <= 1 || loading}
                onClick={() => setPage((current) => Math.max(1, current - 1))}
                className="inline-flex items-center gap-1 rounded-lg border border-slate-200 px-3 py-2 text-sm font-medium disabled:cursor-not-allowed disabled:opacity-40"
              >
                <ChevronLeft size={16} /> Previous
              </button>
              <button
                type="button"
                disabled={page >= totalPages || loading}
                onClick={() => setPage((current) => Math.min(totalPages, current + 1))}
                className="inline-flex items-center gap-1 rounded-lg border border-slate-200 px-3 py-2 text-sm font-medium disabled:cursor-not-allowed disabled:opacity-40"
              >
                Next <ChevronRight size={16} />
              </button>
            </div>
          </div>
        </section>
      </div>
    </div>
  )
}
