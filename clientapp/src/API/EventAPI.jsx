import { handleAPIResponse } from "./APIBase"

export function getPackageEvents(packageId, onSuccess, onError) {

    packageId = encodeURIComponent(packageId);

    fetch("/api/event/" + packageId, {method: "GET"}).then(response => {
        handleAPIResponse(response, onSuccess, onError);
    });
}

export function eventIsDelivered(event) {
    if (event) {
        const eventText = event.event_text.toLowerCase().replace(/[^a-z]+/, " ");
        return eventText.includes("delivered"); 
    }
    return false;
}