import { useState, useEffect } from 'react'
import { MapContainer, TileLayer, Marker, Popup, useMap, useMapEvents } from 'react-leaflet'
import './App.css'

import TimeSlider from './components/TimeSlider';
import { LocationMarker, UserMarker } from './components/LocationMarkers';
import { targetIcon } from './components/TargetIcon';
import { Route } from "./components/RouteOverlay";

import { useTargets } from './hooks/useTargets.js';
import { useWalkingRoute } from './hooks/useWalkingRoute.js';
import { defaultMapZoom, defaultMapCenter } from "./services/constants.js";


// Recenter map when user location is received
function RecenterMap({ position }) {
  const map = useMap();
  useEffect(() => {
    if (position) {
      map.setView(position, 15);
    }
  }, [position, map]);
  return null;
}

function MapClickHandler({ onMapClick }) {
  useMapEvents({
    click: (e) => {
      onMapClick([e.latlng.lat, e.latlng.lng]);
    },
  });
  return null;
}

export default function App() {
  const [userPos, setUserPos] = useState(null);
  const [geoLoading, setGeoLoading] = useState(true);
  const [geoError, setGeoError] = useState(null);
  const [endPos, setEndPos] = useState(null);
  const [maxMinutes, setMaxMinutes] = useState(30);
  const [debouncedMaxMinutes, setDebouncedMaxMinutes] = useState(30);

  const { targets, targetsLoading, targetsError } = useTargets(userPos);
  const { linePos, routeLoading, routeError } = useWalkingRoute(
    userPos,
    endPos,
    debouncedMaxMinutes,
  );

  // Get user location
  useEffect(() => {
    if (!navigator.geolocation) {
      setTimeout(() => {
        setGeoError("Geolocation is not supported by your browser");
        setGeoLoading(false);
      }, 0);
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (position) => {
        setUserPos([position.coords.latitude, position.coords.longitude]);
        setGeoLoading(false);
      },
      (err) => {
        setGeoError(`Failed to retrieve location: ${err.message}`);
        setGeoLoading(false);
      },
      { enableHighAccuracy: true },
    );
  }, []);

  // Debounce max walking time changes
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedMaxMinutes(maxMinutes);
    }, 500);

    return () => clearTimeout(timer);
  }, [maxMinutes]);

  const defaultCenter = userPos || defaultMapCenter;
  const statusMessage =
    geoError ||
    targetsError ||
    routeError ||
    (geoLoading ? "Locating you..." :
      targetsLoading ? "Loading nearby photo targets..." :
      routeLoading ? "Generating route..." :
      null);

  return (
    <div style={{ height: '100%', width: '100%', position: 'relative', overflow: 'hidden', margin: 0, padding: 0}}>
      <MapContainer center={defaultCenter} zoom={defaultMapZoom} style={{height: '100%', width: '100%'}}>
        {/* OpenStreetMap public map tile layer */}
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"/>

        {/* Loading / error overlay */}
        {statusMessage && (
          <div style={{
            position: 'absolute',
            top: 20,
            left: 20,
            zIndex: 1000,
            backgroundColor: 'rgba(255,255,255,0.95)',
            padding: '12px 16px',
            borderRadius: 8,
            color: geoError || targetsError || routeError ? '#c00' : '#333',
            boxShadow: '0 2px 8px rgba(0,0,0,0.12)',
            pointerEvents: 'none',
          }}>
            {statusMessage}
          </div>
        )}

        {!endPos && !geoLoading && !geoError && (
          <div style={{
            position: 'absolute',
            top: 20,
            right: 20,
            zIndex: 1000,
            backgroundColor: 'rgba(255,255,255,0.95)',
            padding: '12px 16px',
            borderRadius: 8,
            color: '#333',
            boxShadow: '0 2px 8px rgba(0,0,0,0.12)',
            pointerEvents: 'none',
          }}>
            Click the map to set your destination.
          </div>
        )}

        {/* Dynamic map panning hook */}
        {userPos && <RecenterMap position={userPos} />}

        {/* Map Click Handler */}
        <MapClickHandler onMapClick={(coords) => setEndPos(coords)} />
        {endPos && (
          <Marker position={endPos} icon={targetIcon}>
            <Popup>Destination</Popup>
          </Marker>
        )}

        <UserMarker userPos={userPos} />
        <LocationMarker targets={targets} />

        <Route linePos={linePos} />
      </MapContainer>

      <TimeSlider value={maxMinutes} onChange={setMaxMinutes} />
    </div>
  );
}