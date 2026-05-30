import { useEffect, useRef } from 'react'
import  { useMap } from 'react-leaflet';
import L from 'leaflet'

export function RecenterButton({ userPos }) {
    const map = useMap();
    
    // We use a ref to track the user's position so the button's 
    // click handler always knows where you are, without us having 
    // to re-render the button every time your GPS updates.
    const posRef = useRef(userPos);

    useEffect(() => {
        posRef.current = userPos;
    }, [userPos]);

    useEffect(() => {
        const customControl = L.control({ position: 'topleft' });

        customControl.onAdd = function () {
            const div = L.DomUtil.create('div', 'leaflet-bar leaflet-control');
            const a = L.DomUtil.create('a', '');
            
            a.href = '#';
            a.title = 'Recenter on my location';
            a.style.display = 'flex';
            a.style.justifyContent = 'center';
            a.style.alignItems = 'center';
            a.innerHTML = `
                <svg viewBox="0 0 24 24" width="18" height="18" stroke="currentColor" strokeWidth="2" fill="none" strokeLinecap="round" strokeLinejoin="round">
                    <circle cx="12" cy="12" r="6"></circle>
                    <line x1="12" y1="2" x2="12" y2="6"></line>
                    <line x1="12" y1="18" x2="12" y2="22"></line>
                    <line x1="2" y1="12" x2="6" y2="12"></line>
                    <line x1="18" y1="12" x2="22" y2="12"></line>
                </svg>
            `;

            // 3. Attach native Leaflet event listeners
            L.DomEvent.on(a, 'click', function (e) {
                L.DomEvent.stopPropagation(e);
                L.DomEvent.preventDefault(e);
                if (posRef.current) {
                    map.flyTo(posRef.current, 15, { duration: 1 });
                }
            });

            L.DomEvent.disableClickPropagation(div);

            div.appendChild(a);
            return div;
        };

        // Add it to the map
        customControl.addTo(map);

        // Cleanup function to remove the button if the component unmounts
        return () => {
            customControl.remove();
        };
    }, [map]);

    // We return null because Leaflet is handling the HTML rendering now, not React!
    return null; 
}