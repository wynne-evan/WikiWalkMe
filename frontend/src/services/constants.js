export const defaultMapZoom = 14;
export const defaultMapCenter = [44.6488, -63.5752]; // Downtown Halifax

export const defaultPointSearchRadius = 10;

export const targetsUrl = `http://localhost:8080/api/targets`;
export const routeUrl = `http://localhost:8080/api/route`;

export function osrmUrl(coordPairs) {
  return `https://router.project-osrm.org/route/v1/walking/${coordPairs.join(";")}?overview=full&geometries=geojson`;
}
