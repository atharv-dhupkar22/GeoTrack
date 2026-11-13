
import { useCallback, useEffect, useState } from 'react'
import {
  CircleMarker,
  MapContainer,
  Polygon,
  TileLayer,
  useMapEvents,
} from 'react-leaflet'
import {
  MapPin,
  Plus,
  RotateCcw,
  Save,
  Trash2,
} from 'lucide-react'
import toast from 'react-hot-toast'
import type { LatLngExpression } from 'leaflet'

import { geofenceService } from '../../services/geofenceService'
import type {
  CreateGeofenceRequest,
  Geofence,
} from '../../types'

import './GeofenceManager.css'
import 'leaflet/dist/leaflet.css'

type Coordinate = [number, number]

const INITIAL_CENTER: LatLngExpression = [19.0760, 72.8777]

const categories: {
  value: Geofence['category']
  label: string
}[] = [
  { value: 'delivery_zone', label: 'Delivery Zone' },
  { value: 'restricted_zone', label: 'Restricted Zone' },
  { value: 'toll_zone', label: 'Toll Zone' },
  { value: 'customer_area', label: 'Customer Area' },
]

interface DrawingMapProps {
  points: Coordinate[]
  setPoints: React.Dispatch<React.SetStateAction<Coordinate[]>>
}

function DrawingMap({ points, setPoints }: DrawingMapProps) {
  useMapEvents({
    click(event) {
      setPoints((previous) => [
        ...previous,
        [event.latlng.lat, event.latlng.lng],
      ])
    },
  })

  const polygonPositions: LatLngExpression[] =
    points.length >= 3
      ? [...points, points[0]]
      : points

  return (
    <>
      {points.map((point, index) => (
        <CircleMarker
          key={`${index}-${point[0]}-${point[1]}`}
          center={point}
          radius={6}
          pathOptions={{
            color: '#ffffff',
            weight: 2,
            fillColor: '#3979e9',
            fillOpacity: 1,
          }}
        />
      ))}

      {points.length >= 3 && (
        <Polygon
          positions={polygonPositions}
          pathOptions={{
            color: '#3979e9',
            weight: 3,
            fillColor: '#3979e9',
            fillOpacity: 0.2,
          }}
        />
      )}
    </>
  )
}

function GeofenceManager() {
  const [geofences, setGeofences] = useState<Geofence[]>([])
  const [points, setPoints] = useState<Coordinate[]>([])
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [category, setCategory] =
    useState<Geofence['category']>('delivery_zone')

  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState('')

  const loadGeofences = useCallback(async () => {
    setLoading(true)
    setError('')

    try {
      const result = await geofenceService.getAll()
      setGeofences(result)
    } catch {
      setError('Unable to load geofences. Check that the backend is running.')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    void loadGeofences()
  }, [loadGeofences])

  const resetDrawing = () => {
    setPoints([])
    setName('')
    setDescription('')
    setCategory('delivery_zone')
  }

  const undoLastPoint = () => {
    setPoints((previous) => previous.slice(0, -1))
  }

  const saveGeofence = async () => {
    if (!name.trim()) {
      toast.error('Please enter a geofence name.')
      return
    }

    if (points.length < 3) {
      toast.error('Select at least three points on the map.')
      return
    }

    const coordinates: number[][] = [
      ...points.map(([latitude, longitude]) => [
        latitude,
        longitude,
      ]),
      [points[0][0], points[0][1]],
    ]

    const request: CreateGeofenceRequest = {
      name: name.trim(),
      description: description.trim(),
      coordinates,
      category,
    }

    setSaving(true)

    try {
      const created = await geofenceService.create(request)

      setGeofences((previous) => [created, ...previous])

      toast.success('Geofence created successfully!')
      resetDrawing()
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
          : 'Unable to create geofence. Please try again.'

      toast.error(message)
    } finally {
      setSaving(false)
    }
  }

  return (
    <section className="geofence-manager">
      <div className="geofence-manager-heading">
        <div>
          <div className="geofence-eyebrow">LOCATION MANAGEMENT</div>
          <h2>Geofence Management</h2>
          <p>Draw geographic zones and manage your monitoring boundaries.</p>
        </div>

        <div className="geofence-total">
          <MapPin size={17} />
          <strong>{geofences.length}</strong>
          <span>Configured zones</span>
        </div>
      </div>

      <div className="geofence-workspace">
        <div className="geofence-map-panel">
          <div className="geofence-panel-title">
            <div>
              <h3>Geofence map</h3>
              <p>Click on the map to add polygon boundary points.</p>
            </div>

            <span className="map-live-badge">MAP ACTIVE</span>
          </div>

          <div className="drawing-instructions">
            <MapPin size={15} />
            <span>
              Click at least three locations to draw a geofence.
            </span>
          </div>

          <div className="geofence-map">
            <MapContainer
              center={INITIAL_CENTER}
              zoom={12}
              scrollWheelZoom
              className="geofence-leaflet-map"
            >
              <TileLayer
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
              />

              <DrawingMap
                points={points}
                setPoints={setPoints}
              />
            </MapContainer>
          </div>

          <div className="map-controls">
            <span>
              <strong>{points.length}</strong> boundary points selected
            </span>

            <div>
              <button
                type="button"
                className="map-control-button"
                onClick={undoLastPoint}
                disabled={points.length === 0 || saving}
              >
                <RotateCcw size={15} />
                Undo
              </button>

              <button
                type="button"
                className="map-control-button clear"
                onClick={() => setPoints([])}
                disabled={points.length === 0 || saving}
              >
                <Trash2 size={15} />
                Clear points
              </button>
            </div>
          </div>
        </div>

        <div className="geofence-form-panel">
          <div className="geofence-panel-title">
            <div>
              <h3>Create geofence</h3>
              <p>Enter the zone details and save.</p>
            </div>

            <div className="create-icon">
              <Plus size={19} />
            </div>
          </div>

          <form
            className="geofence-form"
            onSubmit={(event) => {
              event.preventDefault()
              void saveGeofence()
            }}
          >
            <label htmlFor="geofence-name">
              Geofence name <span>*</span>
            </label>

            <input
              id="geofence-name"
              type="text"
              placeholder="e.g. Mumbai Delivery Zone"
              value={name}
              onChange={(event) => setName(event.target.value)}
              maxLength={150}
              required
            />

            <label htmlFor="geofence-description">
              Description
            </label>

            <textarea
              id="geofence-description"
              placeholder="Describe this monitoring zone..."
              value={description}
              onChange={(event) => setDescription(event.target.value)}
              rows={4}
            />

            <label htmlFor="geofence-category">
              Geofence category <span>*</span>
            </label>

            <select
              id="geofence-category"
              value={category}
              onChange={(event) =>
                setCategory(event.target.value as Geofence['category'])
              }
            >
              {categories.map((item) => (
                <option key={item.value} value={item.value}>
                  {item.label}
                </option>
              ))}
            </select>

            <div className="coordinate-summary">
              <div>
                <span>Boundary points</span>
                <strong>{points.length}</strong>
              </div>

              <div>
                <span>Polygon status</span>
                <strong className={points.length >= 3 ? 'valid' : 'pending'}>
                  {points.length >= 3 ? 'Ready' : 'Incomplete'}
                </strong>
              </div>
            </div>

            <button
              type="submit"
              className="save-geofence-button"
              disabled={saving || loading}
            >
              <Save size={17} />
              {saving ? 'Saving geofence...' : 'Save geofence'}
            </button>

            <button
              type="button"
              className="reset-geofence-button"
              onClick={resetDrawing}
              disabled={saving}
            >
              Reset form and drawing
            </button>
          </form>
        </div>
      </div>

      <div className="geofence-existing-panel">
        <div className="geofence-panel-title">
          <div>
            <h3>Existing geofences</h3>
            <p>Zones currently saved in your GeoTrack database.</p>
          </div>

          <button
            type="button"
            className="refresh-geofence-button"
            onClick={() => void loadGeofences()}
            disabled={loading}
          >
            <RotateCcw size={15} />
            Refresh
          </button>
        </div>

        {error && (
          <div className="geofence-error">
            {error}
          </div>
        )}

        {loading ? (
          <div className="geofence-empty">
            Loading geofences...
          </div>
        ) : geofences.length === 0 ? (
          <div className="geofence-empty">
            No geofences found. Create your first zone above.
          </div>
        ) : (
          <div className="existing-geofence-list">
            {geofences.map((geofence) => (
              <div className="existing-geofence-item" key={geofence.id}>
                <div className="existing-geofence-icon">
                  <MapPin size={18} />
                </div>

                <div className="existing-geofence-details">
                  <strong>{geofence.name}</strong>
                  <span>
                    {geofence.category.replaceAll('_', ' ')}
                  </span>
                </div>

                <span
                  className={`existing-geofence-status ${
                    geofence.status.toLowerCase() === 'active'
                      ? 'active'
                      : 'inactive'
                  }`}
                >
                  {geofence.status}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>
    </section>
  )
}

export default GeofenceManager
