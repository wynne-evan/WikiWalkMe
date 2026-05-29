import { useEffect, useReducer } from "react";
import { getWalkingRoute } from "../services/api.js";

function routeReducer(state, action) {
  switch (action.type) {
    case "fetchStarted":
      return { ...state, loading: true, error: null };
    case "fetchSucceeded":
      return { ...state, loading: false, linePos: action.payload };
    case "fetchFailed":
      return { ...state, loading: false, error: action.payload };
    default:
      return state;
  }
}

const initialState = {
  linePos: [],
  loading: false,
  error: null,
};

export function useWalkingRoute(userPos, endPos, maxMinutes) {
  const [state, dispatch] = useReducer(routeReducer, initialState);

  useEffect(() => {
    if (!userPos || !endPos) {
      return;
    }

    const controller = new AbortController();
    dispatch({ type: "fetchStarted" });

    getWalkingRoute(userPos, endPos, maxMinutes, controller.signal)
      .then((routeLine) => dispatch({ type: "fetchSucceeded", payload: routeLine }))
      .catch((err) => {
        if (err.name === "AbortError") return;
        console.error("Error generating walking route", err);
        dispatch({ type: "fetchFailed", payload: "Failed to calculate the walking route." });
      });

    return () => controller.abort();
  }, [userPos, endPos, maxMinutes]);

  return {
    linePos: state.linePos,
    routeLoading: state.loading,
    routeError: state.error,
  };
}
