import React from "react"
import { createRoot } from "react-dom/client"
import ColumbaApp from "./ColumbaApp"

const container = document.getElementById("root");
const root = createRoot(container);
root.render(<ColumbaApp />);
