
import { useCallback, useEffect, useState } from 'react'
import { Link, Route, Routes } from 'react-router-dom'
import {
  Activity,
  ArrowDownRight,
  ArrowUpRight,
  Bell,
  CarFront,
  CheckCircle2,
  Clock3,
  LayoutDashboard,
  Map,
  MapPin,
  Navigation,
  RefreshCw,
  Settings,
  ShieldAlert,
  Truck,
  Wifi,
  XCircle,
} from 'lucide-react'

import './App.css'

import { geofenceService } from './services/geofenceService'
import { vehicleService } from './services/vehicleService'
import { alertService } from './services/alertService'
import { violationService } from './services/violationService'

import GeofencesPage from './pages/GeofencesPage'
import VehiclesPage from './pages/VehiclesPage'
import TrackingPage from './pages/TrackingPage'
import AlertsPage from './pages/AlertsPage'
import ViolationHistoryPage from './pages/ViolationHistoryPage'

import type {
  AlertEvent,
  Geofence,
  Vehicle,
  Violation,
} from './types'

function Dashboard() {
  const [geofences, setGeofences] = useState<Geofence[]>([])
  const [vehicles, setVehicles] = useState<Vehicle[]>([])
  const [alerts, setAlerts] = useState<AlertEvent[]>([])
  const [violations, setViolations] = useState<Violation[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const loadDashboard = useCallback(async () => {
    setLoading(true)
    setError('')

    const results = await Promise.allSettled([
      geofenceService.getAll(),
      vehicleService.getAll(),
      alertService.getHistory({ page: 1, limit: 5 }),
      violationService.getHistory({ page: 1, limit: 5 }),
    ])

    const errors: string[] = []

    if (results[0].status === 'fulfilled') {
      setGeofences(results[0].value)
    } else {
      errors.push('Geofences')
    }

    if (results[1].status === 'fulfilled') {
      setVehicles(results[1].value)
    } else {
      errors.push('Vehicles')
    }

    if (results[2].status === 'fulfilled') {
      setAlerts(results[2].value.alerts)
    } else {
      errors.push('Alerts')
    }

    if (results[3].status === 'fulfilled') {
      setViolations(results[3].value.violations)
    } else {
      errors.push('Violations')
    }

    if (errors.length > 0) {
      setError(`Unable to load: ${errors.join(', ')}. Check the backend connection.`)
    }

    setLoading(false)
  }, [])

  useEffect(() => {
    void loadDashboard()
  }, [loadDashboard])

  const activeVehicles = vehicles.filter(
    (vehicle) => vehicle.status.toLowerCase() === 'active',
  ).length

  const activeGeofences = geofences.filter(
    (geofence) => geofence.status.toLowerCase() === 'active',
  ).length

  const entryCount = violations.filter(
    (violation) => violation.event_type === 'entry',
  ).length

  const exitCount = violations.filter(
    (violation) => violation.event_type === 'exit',
  ).length

  const formatDate = (value: string) => {
    const date = new Date(value)

    if (Number.isNaN(date.getTime())) {
      return 'Unknown date'
    }

    return date.toLocaleString('en-IN', {
      day: '2-digit',
      month: 'short',
      hour: '2-digit',
      minute: '2-digit',
    })
  }

  return (
    <div className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <div className="brand-icon">
            <MapPin size={23} />
          </div>
          <div>
            <h1>GeoTrack</h1>
            <span>INTELLIGENCE PLATFORM</span>
          </div>
        </div>

        <div className="sidebar-label">WORKSPACE</div>

        <nav className="sidebar-nav">
          <Link to="/" className="nav-link active">
            <LayoutDashboard size={19} />
            <span>Dashboard</span>
          </Link>

          <Link to="/geofences" className="nav-link">
            <Map size={19} />
            <span>Geofences</span>
          </Link>

          <Link to="/vehicles" className="nav-link">
            <Truck size={19} />
            <span>Vehicles</span>
          </Link>

          <Link to="/tracking" className="nav-link">
            <Navigation size={19} />
            <span>Live GPS Tracking</span>
          </Link>

          <Link to="/alerts" className="nav-link">
            <Bell size={19} />
            <span>Alerts</span>
            <span className="nav-count">{alerts.length}</span>
          </Link>

          <Link to="/violations" className="nav-link">
            <ShieldAlert size={19} />
            <span>Violation History</span>
          </Link>
        </nav>

        <div className="sidebar-bottom">
          <div className="sidebar-label">SYSTEM</div>

          <a href="#settings" className="nav-link">
            <Settings size={19} />
            <span>Settings</span>
          </a>

          <div className="sidebar-status">
            <div className="status-indicator" />
            <div>
              <strong>System online</strong>
              <span>GeoTrack backend</span>
            </div>
          </div>

          <div className="sidebar-footer">
            <span className="footer-avatar">GT</span>
            <div>
              <strong>GeoTrack Admin</strong>
              <span>Administrator</span>
            </div>
          </div>
        </div>
      </aside>

      <main className="main-content" id="dashboard">
        <header className="topbar">
          <div className="breadcrumb">
            <span>Workspace</span>
            <span className="breadcrumb-divider">/</span>
            <strong>Dashboard</strong>
          </div>

          <div className="topbar-actions">
            <div className="connection-badge">
              <Wifi size={15} />
              <span>API Connected</span>
            </div>

            <button
              className="icon-button"
              type="button"
              aria-label="Refresh dashboard"
              onClick={() => void loadDashboard()}
              disabled={loading}
            >
              <RefreshCw
                size={18}
                className={loading ? 'spin' : ''}
              />
            </button>

            <button
              className="icon-button notification-button"
              type="button"
              aria-label="Notifications"
            >
              <Bell size={18} />
              {alerts.length > 0 && <span className="notification-dot" />}
            </button>
          </div>
        </header>

        <section className="dashboard-content">
          <div className="page-heading">
            <div>
              <div className="eyebrow">OVERVIEW</div>
              <h2>Dashboard</h2>
              <p>Monitor your fleet, geofences and real-time activity.</p>
            </div>

            <div className="date-display">
              <Clock3 size={16} />
              {new Date().toLocaleDateString('en-IN', {
                weekday: 'short',
                day: '2-digit',
                month: 'long',
                year: 'numeric',
              })}
            </div>
          </div>

          {error && (
            <div className="error-banner">
              <XCircle size={18} />
              <span>{error}</span>
              <button
                type="button"
                onClick={() => void loadDashboard()}
              >
                Retry
              </button>
            </div>
          )}

          <div className="section-heading">
            <div>
              <h3>Platform overview</h3>
              <p>Live statistics from your GeoTrack backend</p>
            </div>
            <span className="live-label">
              <span className="live-dot" />
              LIVE DATA
            </span>
          </div>

          <div className="stats-grid">
            <StatCard
              title="Total Vehicles"
              value={vehicles.length}
              subtitle={`${activeVehicles} active vehicles`}
              icon={<Truck size={21} />}
              color="blue"
              loading={loading}
            />
            <StatCard
              title="Geofences"
              value={geofences.length}
              subtitle={`${activeGeofences} active zones`}
              icon={<MapPin size={21} />}
              color="purple"
              loading={loading}
            />
            <StatCard
              title="Total Alerts"
              value={alerts.length}
              subtitle="Latest 5 notifications"
              icon={<Bell size={21} />}
              color="orange"
              loading={loading}
            />
            <StatCard
              title="Violations"
              value={violations.length}
              subtitle="Latest 5 recorded events"
              icon={<ShieldAlert size={21} />}
              color="red"
              loading={loading}
            />
          </div>

          <div className="section-heading activity-heading">
            <div>
              <h3>Activity overview</h3>
              <p>Recent geofence movement events</p>
            </div>
          </div>

          <div className="overview-grid">
            <section className="panel event-summary">
              <div className="panel-heading">
                <div>
                  <h3>Geofence activity</h3>
                  <p>Recent entry and exit events</p>
                </div>
                <Activity size={20} className="panel-icon" />
              </div>

              <div className="event-stat-grid">
                <div className="event-stat entry-stat">
                  <div className="event-stat-icon">
                    <ArrowDownRight size={20} />
                  </div>
                  <span>Entries</span>
                  <strong>{loading ? '—' : entryCount}</strong>
                  <small>In latest 5 violations</small>
                </div>

                <div className="event-stat exit-stat">
                  <div className="event-stat-icon">
                    <ArrowUpRight size={20} />
                  </div>
                  <span>Exits</span>
                  <strong>{loading ? '—' : exitCount}</strong>
                  <small>In latest 5 violations</small>
                </div>
              </div>

              <div className="summary-footer">
                <CheckCircle2 size={16} />
                <span>Geofence monitoring is configured</span>
              </div>
            </section>

            <section className="panel activity-panel" id="alerts">
              <div className="panel-heading">
                <div>
                  <h3>Recent alerts</h3>
                  <p>Latest notifications from your system</p>
                </div>
                <Bell size={19} className="panel-icon" />
              </div>

              {loading ? (
                <div className="empty-state">Loading recent alerts...</div>
              ) : alerts.length === 0 ? (
                <div className="empty-state">
                  <Bell size={25} />
                  <strong>No alerts yet</strong>
                  <span>New geofence notifications will appear here.</span>
                </div>
              ) : (
                <div className="activity-list">
                  {alerts.map((alert) => (
                    <div className="activity-item" key={alert.id}>
                      <div
                        className={`activity-icon ${
                          alert.event_type === 'entry' ? 'entry' : 'exit'
                        }`}
                      >
                        {alert.event_type === 'entry' ? (
                          <ArrowDownRight size={18} />
                        ) : (
                          <ArrowUpRight size={18} />
                        )}
                      </div>
                      <div className="activity-details">
                        <strong>
                          {alert.event_type === 'entry'
                            ? 'Vehicle entered zone'
                            : 'Vehicle exited zone'}
                        </strong>
                        <span>{alert.message}</span>
                        <small>{formatDate(alert.created_at)}</small>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </section>
          </div>

          <section className="panel vehicle-panel" id="vehicles">
            <div className="panel-heading">
              <div>
                <h3>Registered vehicles</h3>
                <p>Vehicles connected to your tracking platform</p>
              </div>
              <CarFront size={20} className="panel-icon" />
            </div>

            {loading ? (
              <div className="empty-state">Loading vehicles...</div>
            ) : vehicles.length === 0 ? (
              <div className="empty-state">
                <Truck size={26} />
                <strong>No vehicles registered</strong>
                <span>Registered vehicles will appear here.</span>
              </div>
            ) : (
              <div className="table-wrapper">
                <table className="vehicle-table">
                  <thead>
                    <tr>
                      <th>Vehicle</th>
                      <th>Driver</th>
                      <th>Type</th>
                      <th>Phone</th>
                      <th>Status</th>
                    </tr>
                  </thead>
                  <tbody>
                    {vehicles.slice(0, 5).map((vehicle) => (
                      <tr key={vehicle.id}>
                        <td>
                          <div className="vehicle-name">
                            <span className="vehicle-icon">
                              <Truck size={17} />
                            </span>
                            <strong>{vehicle.vehicle_number}</strong>
                          </div>
                        </td>
                        <td>{vehicle.driver_name}</td>
                        <td>{vehicle.vehicle_type}</td>
                        <td>{vehicle.phone}</td>
                        <td>
                          <span
                            className={`vehicle-status ${
                              vehicle.status.toLowerCase() === 'active'
                                ? 'status-active'
                                : 'status-inactive'
                            }`}
                          >
                            {vehicle.status}
                          </span>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </section>

          <section className="panel geofence-panel" id="geofences">
            <div className="panel-heading">
              <div>
                <h3>Configured geofences</h3>
                <p>Geographic monitoring zones</p>
              </div>
              <Map size={20} className="panel-icon" />
            </div>

            {loading ? (
              <div className="empty-state">Loading geofences...</div>
            ) : geofences.length === 0 ? (
              <div className="empty-state">
                <MapPin size={26} />
                <strong>No geofences created</strong>
                <span>Create a geofence to start monitoring locations.</span>
              </div>
            ) : (
              <div className="geofence-list">
                {geofences.slice(0, 5).map((geofence) => (
                  <div className="geofence-item" key={geofence.id}>
                    <div className="geofence-icon">
                      <MapPin size={18} />
                    </div>
                    <div className="geofence-details">
                      <strong>{geofence.name}</strong>
                      <span>{geofence.category.replaceAll('_', ' ')}</span>
                    </div>
                    <span
                      className={`vehicle-status ${
                        geofence.status.toLowerCase() === 'active'
                          ? 'status-active'
                          : 'status-inactive'
                      }`}
                    >
                      {geofence.status}
                    </span>
                  </div>
                ))}
              </div>
            )}

            <div style={{ marginTop: 16 }}>
              <Link to="/geofences" className="open-geofence-link">
                Manage geofences
                <MapPin size={15} />
              </Link>
            </div>
          </section>

          <footer className="dashboard-footer">
            <span>GeoTrack © 2026</span>
            <span>Geofencing & Real-time Alert System</span>
            <span className="footer-version">Development build</span>
          </footer>
        </section>
      </main>
    </div>
  )
}

interface StatCardProps {
  title: string
  value: number
  subtitle: string
  icon: React.ReactNode
  color: string
  loading: boolean
}

function StatCard({
  title,
  value,
  subtitle,
  icon,
  color,
  loading,
}: StatCardProps) {
  return (
    <div className="stat-card">
      <div className={`stat-icon ${color}`}>{icon}</div>
      <div className="stat-card-content">
        <span className="stat-title">{title}</span>
        <strong className="stat-value">{loading ? '—' : value}</strong>
        <span className="stat-subtitle">{subtitle}</span>
      </div>
    </div>
  )
}

function App() {
  return (
    <Routes>
      <Route path="/" element={<Dashboard />} />
      <Route path="/geofences" element={<GeofencesPage />} />
      <Route path="/vehicles" element={<VehiclesPage />} />
      <Route path="/tracking" element={<TrackingPage />} />
      <Route path="/alerts" element={<AlertsPage />} />
      <Route path="/violations" element={<ViolationHistoryPage />} />
      <Route path="*" element={<Dashboard />} />
    </Routes>
  )
}

export default App
