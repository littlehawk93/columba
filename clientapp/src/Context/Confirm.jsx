import React from "react"

export const ConfirmContext = React.createContext({
    show: false,
    title: null,
    message: null,
    onConfirm: null,
    onCancel: null,
    onShowConfirm: (title, message, onConfirm, onCancel) => {}
});