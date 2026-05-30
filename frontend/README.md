# WikiWalkMe Frontend

The frontend is a React application built with `react-leaflet` to visualize geographic data.

## Tech Stack

- **React**: Component-based UI.
- **Leaflet**: Interactive map rendering.
- **Hooks**: Custom hooks (`useTargets`, `useWalkingRoute`) for managing API state and side effects.

## Setup

1. Run `npm install` to install dependencies (Leaflet, React-Leaflet).
2. Run `npm start` to launch the dev server.

## Features

- **Real-time Map**: Interactive map with marker clustering.
- **Time Slider**: Dynamically adjusts the target route search area based on the user's available time.
- **Visual Feedback**: Automatic marker placement and route polyline drawing.
