import { defaultPointSearchRadius, targetsUrl, routeUrl } from "./constants.js";

export async function getTargets(userPos, signal) {
  if (!userPos) return [];

  const payload = {
    lat: userPos[0],
    lon: userPos[1],
    radius: defaultPointSearchRadius,
  };

  const res = await fetch(targetsUrl, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
    signal,
  });

  if (!res.ok) {
    throw new Error("Backend server error fetching targets");
  }

  const data = await res.json();
  return data.targets || [];
}

export async function getWalkingRoute(userPos, endPos, maxMinutes, signal) {
  if (!userPos || !endPos) return [];

  const payload = {
    start_lat: userPos[0],
    start_lon: userPos[1],
    end_lat: endPos[0],
    end_lon: endPos[1],
    max_minutes: maxMinutes,
  };

  const res = await fetch(routeUrl, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
    signal,
  });

  if (!res.ok) {
    throw new Error("Backend server error generating route");
  }

  const data = await res.json();
  return data.route_points || [];
}
