import  { useEffect } from 'react';
import { useMap } from "react-leaflet"
import L from 'leaflet';
import 'leaflet-polylinedecorator';

export function Route({ linePos, maxMinutes }) {
    const map = useMap();

    const minRepeat = (3/maxMinutes)*100;
    const repeat = minRepeat < 10 ? `${minRepeat}%` : '10%' // Have an arrow at least every 3 minutes

    useEffect(() => {
        if (!linePos || linePos.length < 2) {
            return;
        }
        
        const polyline = L.polyline(linePos, {
            color: "#0066ff",
            weight: 6,
            opacity: 0.8,
        }).addTo(map);

        const decorator = L.polylineDecorator(polyline, {
            patterns: [
                {
                    offset: "5%", // Start arrows 5% into the route
                    repeat: repeat, // Place an arrow every 15%
                    symbol: L.Symbol.arrowHead({
                        pixelSize: 14,
                        polygon: false, // Chevron (V) instead of filled arrow
                        pathOptions: {stroke: true, weight: 3, color: '#014ec2'}
                    })
                }
            ]
        }).addTo(map)

        // Cleanup old lines
        return () => {
            map.removeLayer(polyline);
            map.removeLayer(decorator);
        };
    }, [linePos, map, repeat]);

    return null;
}