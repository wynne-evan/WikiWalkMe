import { useState, useEffect } from 'react'
import { MapContainer, TileLayer, Marker, Popup, useMap, useMapEvents, Polyline } from 'react-leaflet'
import L from 'leaflet'
import './App.css'

import markerIcon2x from 'leaflet/dist/images/marker-icon-2x.png';
import markerIcon from 'leaflet/dist/images/marker-icon.png';
import markerShadow from 'leaflet/dist/images/marker-shadow.png';

delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
  iconUrl: markerIcon,
  iconRetinaUrl: markerIcon2x,
  shadowUrl: markerShadow,
});

// Blue user icon marker
const userIcon = new L.Icon({
  iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-blue.png',
  shadowUrl: markerShadow,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowSize: [41, 41]
});

// Red Photo target opportunities
const targetIcon = new L.Icon({
  iconUrl: 'https://raw.githubusercontent.com/pointhi/leaflet-color-markers/master/img/marker-icon-2x-red.png',
  shadowUrl: markerShadow,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
  popupAnchor: [1, -34],
  shadowSize: [41, 41]
});

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
  const [targets, setTargets] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [endPos, setEndPos] = useState(null);
  const [linePos, setLinePos] = useState([])

  // Get user location
  useEffect(() => {
    if (!navigator.geolocation) {
      setTimeout(() => {
      setError("Geolocation is not supported by your browser");
      setLoading(false);
      }, 0);
    }

    navigator.geolocation.getCurrentPosition( (position) => {
      setUserPos([position.coords.latitude, position.coords.longitude]);
      setLoading(false);
    }, (err) => {
      setError(`Failed to retrieve location: ${err.message}`);
      setLoading(false);
    },
  {enableHighAccuracy: true});
  }, []);

  // Fetch targets from backend when user position is resolved
  useEffect(() => {
    if (!userPos) return;

    const [lat, lon] = userPos;
    const backendUrl = `http://localhost:8080/api/targets?lat=${lat}&lon=${lon}&radius=2.0`;

    fetch(backendUrl)
      .then((res) => {
        if (!res.ok) throw new Error("Backend server error");
        return res.json();
      })
      .then((data) => {
        console.log(data.targets);
        setTargets(data.targets || []);
      })
      .catch((err) => {
        console.error("Error fetching targets", err);
      })
  }, [userPos]);

  useEffect(() => {
    if (!userPos) return;
    if (!endPos) return;

    const max_minutes = 45;
    const backendUrl = `http://localhost:8080/api/route`;
    const payload = {
      start_lat: userPos[0],
      start_lon: userPos[1],
      end_lat: endPos[0],
      end_lon: endPos[1],
      max_minutes: max_minutes
    }


    fetch(backendUrl,
      {method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
      }
    )
      .then((res) => {
        if (!res.ok) throw new Error("Backend server error");
        return res.json();
      })
      .then((data) => {
        console.log(data)
        const linePositions = [];
        linePositions.push(userPos);
        if (data.path) {
        data.path.forEach(target => {
          linePositions.push([target.lat, target.lon]);
        }); 
      }

        linePositions.push(endPos);

        setLinePos(linePositions);
        })
      .catch((err) => {
        console.error("Error fetching route", err);
    });
  }, [userPos, endPos]);

  if (loading) return <div style={{ padding: 20 }}>Locating you...</div>;
  if (error) return <div style={{ padding: 20, color: 'red'}}>{error}</div>;

  // Fallback location in downtown halifax
  const defaultCenter = userPos || [44.6488, -63.5752];



  return (
    <div style={{ height: '100vh', width: '100vw', position: 'relative'}}>
      <MapContainer center={defaultCenter} zoom={14} style={{height: '100%', width: '100%'}}>
        {/* OpenStreetMap public map tile layer */}
        <TileLayer
          attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
          url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"/>

        {/* Dynamic map panning hook */}
        {userPos && <RecenterMap position={userPos} />}

        {/* Map Click Handler */}
        <MapClickHandler onMapClick={(coords) => setEndPos(coords)} />
        {endPos && (
          <Marker position={endPos} icon={targetIcon}>
            <Popup>Destination</Popup>
          </Marker>
        )}

        {/* User current position marker (blue) */}
        {userPos && (
          <Marker position={userPos} icon={userIcon}>
            <Popup><strong>You are here</strong></Popup>
          </Marker>
        )}

        {/* Photo target markers (red) */}
        {targets.map((target, index) => (
          <Marker key={index} position={[target.lat, target.lon]} icon={targetIcon}>
            <Popup>
              <div style={{ fontSize: '14px' }}>
                <strong>{target.name}</strong>
                <br />
                <a href={target.wikidata_url} target="_blank" rel="noreferrer">
                  View on Wikidata
                </a>
              </div>
            </Popup>
          </Marker>
        ))}

        {/* DRAW THE ROUTE LINE */}
        {linePos.length > 1 && (
          <Polyline 
            positions={linePos} 
            color="#3388ff"      // Nice clean blue line
            weight={5}           // Thickness of the line in pixels
            opacity={0.7}        // Slight transparency so streets show through
            dashArray="10, 10"   // Makes it a dashed line to imply a walking trail
          />
        )}
      </MapContainer>
    </div>
  );
}