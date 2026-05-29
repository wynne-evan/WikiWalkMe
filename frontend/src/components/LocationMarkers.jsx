import { Marker, Popup } from "react-leaflet"

import { targetIcon, userIcon, endpointIcon } from "./TargetIcon"

export function UserMarker({ userPos }) {
    if (!userPos) return <></>

    return (
        <Marker position={userPos} icon={userIcon}>
            <Popup><strong>You are here</strong></Popup>
        </Marker>
    )
}

export function LocationMarker({ targets }) {
    return (
        <>
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
        </>
    )
}

export function EndpointMarker({endPos}) {
    if (!endPos) { return <></> }

    return (
        <Marker position={endPos} icon={endpointIcon}>
        </Marker>
        )
    }