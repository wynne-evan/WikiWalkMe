import { Polyline } from "react-leaflet"

export function Route({ linePos }) {
    if (!linePos || linePos.length < 1) {
        return null
    }

    return (
            <Polyline 
              positions={linePos}
              color="#0066ff"
              weight={6}
              opacity={0.8}
              dashArray="8, 6"
            />
    )
};