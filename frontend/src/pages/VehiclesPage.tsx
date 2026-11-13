
import { useCallback, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import {
  ArrowLeft,
  CarFront,
  RefreshCw,
  Plus,
  User,
  Phone,
  Truck,
  CheckCircle2,
  AlertCircle,
} from 'lucide-react'
import toast from 'react-hot-toast'

import { vehicleService } from '../services/vehicleService'
import type { Vehicle } from '../types'

type VehicleForm = {
  vehicle_number: string
  driver_name: string
  vehicle_type: string
  phone: string
}

const initialForm: VehicleForm = {
  vehicle_number: '',
  driver_name: '',
  vehicle_type: 'Delivery Van',
  phone: '',
}

function VehiclesPage() {
  const [vehicles, setVehicles] = useState<Vehicle[]>([])
  const [form, setForm] = useState<VehicleForm>(initialForm)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  const loadVehicles = useCallback(async () => {
    try {
      setLoading(true)
      setError('')
      const result = await vehicleService.getAll()
      setVehicles(result)
    } catch {
      setError('Unable to load vehicles. Check that the backend is running.')
      toast.error('Could not load vehicles')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadVehicles()
  }, [loadVehicles])

  const handleChange = (
    event: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const { name, value } = event.target
    setForm((previous) => ({ ...previous, [name]: value }))
  }

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault()

    if (!form.vehicle_number.trim() || !form.driver_name.trim()) {
      toast.error('Vehicle number and driver name are required')
      return
    }

    try {
      setSaving(true)

      await vehicleService.create({
        vehicle_number: form.vehicle_number.trim().toUpperCase(),
        driver_name: form.driver_name.trim(),
        vehicle_type: form.vehicle_type,
        phone: form.phone.trim(),
      })

      toast.success('Vehicle registered successfully')
      setForm(initialForm)
      await loadVehicles()
    } catch (err: unknown) {
      const message =
        typeof err === 'object' &&
        err !== null &&
        'response' in err &&
        typeof err.response === 'object' &&
        err.response !== null &&
        'data' in err.response &&
        typeof err.response.data === 'object' &&
        err.response.data !== null &&
        'message' in err.response.data &&
        typeof err.response.data.message === 'string'
          ? err.response.data.message
          : 'Could not register vehicle. Please check the details.'

      toast.error(message)
    } finally {
      setSaving(false)
    }
  }

  return (
    <main className="min-h-screen bg-slate-50 px-4 py-6 text-slate-800 sm:px-8">
      <div className="mx-auto max-w-7xl">
        <header className="mb-8 flex flex-wrap items-center justify-between gap-4">
          <div>
            <Link
              to="/"
              className="mb-4 inline-flex items-center gap-2 text-sm font-medium text-blue-600 hover:text-blue-800"
            >
              <ArrowLeft size={16} />
              Back to Dashboard
            </Link>

            <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">
              Vehicle Management
            </h1>
            <p className="mt-1 text-sm text-slate-500">
              Register vehicles and manage your GeoTrack fleet.
            </p>
          </div>

          <button
            type="button"
            onClick={() => void loadVehicles()}
            disabled={loading}
            className="inline-flex items-center gap-2 rounded-lg border border-slate-200 bg-white px-4 py-2.5 text-sm font-semibold shadow-sm hover:bg-slate-100 disabled:opacity-50"
          >
            <RefreshCw size={16} className={loading ? 'animate-spin' : ''} />
            Refresh
          </button>
        </header>

        <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
          <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm lg:col-span-1">
            <div className="mb-5 flex items-center gap-3">
              <div className="rounded-xl bg-blue-50 p-3 text-blue-600">
                <Plus size={21} />
              </div>
              <div>
                <h2 className="font-bold">Register Vehicle</h2>
                <p className="text-xs text-slate-500">
                  Enter the vehicle details
                </p>
              </div>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label
                  htmlFor="vehicle_number"
                  className="mb-1.5 block text-sm font-medium"
                >
                  Vehicle Number *
                </label>
                <input
                  id="vehicle_number"
                  name="vehicle_number"
                  value={form.vehicle_number}
                  onChange={handleChange}
                  placeholder="MH46AB1234"
                  required
                  maxLength={20}
                  className="w-full rounded-lg border border-slate-200 px-3 py-2.5 text-sm outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-100"
                />
              </div>

              <div>
                <label
                  htmlFor="driver_name"
                  className="mb-1.5 block text-sm font-medium"
                >
                  Driver Name *
                </label>
                <input
                  id="driver_name"
                  name="driver_name"
                  value={form.driver_name}
                  onChange={handleChange}
                  placeholder="Enter driver name"
                  required
                  className="w-full rounded-lg border border-slate-200 px-3 py-2.5 text-sm outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-100"
                />
              </div>

              <div>
                <label
                  htmlFor="vehicle_type"
                  className="mb-1.5 block text-sm font-medium"
                >
                  Vehicle Type
                </label>
                <select
                  id="vehicle_type"
                  name="vehicle_type"
                  value={form.vehicle_type}
                  onChange={handleChange}
                  className="w-full rounded-lg border border-slate-200 bg-white px-3 py-2.5 text-sm outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-100"
                >
                  <option>Delivery Van</option>
                  <option>Truck</option>
                  <option>Car</option>
                  <option>Bike</option>
                  <option>Bus</option>
                  <option>Other</option>
                </select>
              </div>

              <div>
                <label
                  htmlFor="phone"
                  className="mb-1.5 block text-sm font-medium"
                >
                  Driver Phone
                </label>
                <input
                  id="phone"
                  name="phone"
                  type="tel"
                  value={form.phone}
                  onChange={handleChange}
                  placeholder="9876543210"
                  maxLength={15}
                  className="w-full rounded-lg border border-slate-200 px-3 py-2.5 text-sm outline-none focus:border-blue-500 focus:ring-2 focus:ring-blue-100"
                />
              </div>

              <button
                type="submit"
                disabled={saving}
                className="flex w-full items-center justify-center gap-2 rounded-lg bg-blue-600 px-4 py-3 text-sm font-semibold text-white shadow-sm hover:bg-blue-700 disabled:cursor-not-allowed disabled:opacity-60"
              >
                {saving ? (
                  <RefreshCw size={17} className="animate-spin" />
                ) : (
                  <Plus size={17} />
                )}
                {saving ? 'Registering...' : 'Register Vehicle'}
              </button>
            </form>
          </section>

          <section className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm lg:col-span-2">
            <div className="mb-5 flex items-center justify-between">
              <div>
                <h2 className="font-bold">Registered Vehicles</h2>
                <p className="mt-1 text-xs text-slate-500">
                  Vehicles stored in the GeoTrack database
                </p>
              </div>

              <span className="rounded-full bg-blue-50 px-3 py-1.5 text-sm font-bold text-blue-700">
                {vehicles.length} Total
              </span>
            </div>

            {loading ? (
              <div className="flex min-h-48 items-center justify-center gap-2 text-sm text-slate-500">
                <RefreshCw size={18} className="animate-spin" />
                Loading vehicles...
              </div>
            ) : error ? (
              <div className="flex min-h-48 flex-col items-center justify-center gap-3 text-center">
                <AlertCircle size={28} className="text-red-500" />
                <p className="text-sm text-red-600">{error}</p>
                <button
                  type="button"
                  onClick={() => void loadVehicles()}
                  className="rounded-lg bg-blue-600 px-4 py-2 text-sm font-semibold text-white"
                >
                  Try Again
                </button>
              </div>
            ) : vehicles.length === 0 ? (
              <div className="flex min-h-48 flex-col items-center justify-center text-center">
                <CarFront size={36} className="mb-3 text-slate-300" />
                <p className="font-semibold">No vehicles registered yet</p>
                <p className="mt-1 text-sm text-slate-500">
                  Register your first vehicle using the form.
                </p>
              </div>
            ) : (
              <div className="space-y-3">
                {vehicles.map((vehicle) => (
                  <article
                    key={vehicle.id}
                    className="rounded-xl border border-slate-100 bg-slate-50 p-4"
                  >
                    <div className="flex flex-wrap items-start justify-between gap-3">
                      <div className="flex items-center gap-3">
                        <div className="rounded-xl bg-white p-3 text-blue-600 shadow-sm">
                          <Truck size={22} />
                        </div>
                        <div>
                          <h3 className="font-bold">
                            {vehicle.vehicle_number}
                          </h3>
                          <p className="mt-0.5 text-xs text-slate-500">
                            {vehicle.vehicle_type}
                          </p>
                        </div>
                      </div>

                      <span className="inline-flex items-center gap-1 rounded-full bg-green-50 px-2.5 py-1 text-xs font-semibold text-green-700">
                        <CheckCircle2 size={13} />
                        {vehicle.status || 'Active'}
                      </span>
                    </div>

                    <div className="mt-4 grid grid-cols-1 gap-3 border-t border-slate-200 pt-3 sm:grid-cols-2">
                      <div className="flex items-center gap-2 text-sm text-slate-600">
                        <User size={15} className="text-slate-400" />
                        {vehicle.driver_name}
                      </div>
                      <div className="flex items-center gap-2 text-sm text-slate-600">
                        <Phone size={15} className="text-slate-400" />
                        {vehicle.phone || 'No phone provided'}
                      </div>
                    </div>
                  </article>
                ))}
              </div>
            )}
          </section>
        </div>
      </div>
    </main>
  )
}

export default VehiclesPage
