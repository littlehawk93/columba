import { handleAPIResponse } from "./APIBase"

export function getActivePackages(onSuccess, onError) {

    fetch("/api/package?status=active", {method: "GET"}).then(response => {
        handleAPIResponse(response, onSuccess, onError);
    });
}

export function createPackage(pkg, onSuccess, onError) {

    var body = JSON.stringify(pkg);

    fetch("/api/package", {
        method: "POST",
        body: body,
        headers: {
            "Content-Type": "application/json"
        }
    }).then(response => {
        handleAPIResponse(response, onSuccess, onError);
    });
}

export function deletePackage(packageId, onSuccess, onError) {

    packageId = encodeURIComponent(packageId);

    fetch("/api/package/" + packageId, {method: "DELETE"}).then(response => {
        handleAPIResponse(response, onSuccess, onError);
    });
}

export function getLatestEvent(pkg) {

    if (pkg != null && pkg.events != null && Array.isArray(pkg.events) && pkg.events.length > 0) {
        return pkg.events[0];
    }
    return null;
}