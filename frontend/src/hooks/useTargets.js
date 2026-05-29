import { useEffect, useReducer } from "react";
import { getTargets } from "../services/api.js";

function targetsReducer(state, action) {
  switch (action.type) {
    case "fetchStarted":
      return { ...state, loading: true, error: null };
    case "fetchSucceeded":
      return { ...state, loading: false, targets: action.payload };
    case "fetchFailed":
      return { ...state, loading: false, error: action.payload };
    default:
      return state;
  }
}

const initialState = {
  targets: [],
  loading: false,
  error: null,
};

export function useTargets(userPos) {
  const [state, dispatch] = useReducer(targetsReducer, initialState);

  useEffect(() => {
    if (!userPos) {
      return;
    }

    const controller = new AbortController();
    dispatch({ type: "fetchStarted" });

    getTargets(userPos, controller.signal)
      .then((results) => dispatch({ type: "fetchSucceeded", payload: results }))
      .catch((err) => {
        if (err.name === "AbortError") return;
        console.error("Error fetching targets", err);
        dispatch({ type: "fetchFailed", payload: "Failed to load nearby photo targets." });
      });

    return () => controller.abort();
  }, [userPos]);

  return {
    targets: state.targets,
    targetsLoading: state.loading,
    targetsError: state.error,
  };
}
