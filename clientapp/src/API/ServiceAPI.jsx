import { handleAPIResponse } from "./APIBase"

export function getAllServiceProviders(onSuccess, onError) {

    fetch("/api/service", { method: "GET" }).then(response => {
        handleAPIResponse(response, onSuccess, onError);
    });
}
