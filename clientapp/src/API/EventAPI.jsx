import { handleAPIResponse } from "./APIBase"

export function getPackageEvents(packageId, onSuccess, onError) {

    packageId = encodeURIComponent(packageId);

    fetch("/api/event/" + packageId, {method: "GET"}).then(response => {
        handleAPIResponse(response, onSuccess, onError);
    });
}