import React from "react"

export const ErrorContext = React.createContext({
    error: null,
    level: "error",
    onError: (error) => {}
});