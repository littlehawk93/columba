import React from "react"

export const PopoverContext = React.createContext({
    show: false,
    anchorElem: null,
    text: "",
    vertical: "",
    horizontal: "",
    onShowPopover: (anchorElem, text, vertical, horizontal, onClose) => {}
});