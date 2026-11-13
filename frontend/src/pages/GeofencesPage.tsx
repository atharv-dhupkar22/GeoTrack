
import { Link } from 'react-router-dom'
import { ArrowLeft } from 'lucide-react'
import GeofenceManager from '../components/map/GeofenceManager'
import './GeofencesPage.css'

function GeofencesPage() {
  return (
    <main className="geofences-page">
      <header className="geofences-page-header">
        <Link to="/" className="back-dashboard-link">
          <ArrowLeft size={17} />
          Back to Dashboard
        </Link>
      </header>

      <GeofenceManager />
    </main>
  )
}

export default GeofencesPage
