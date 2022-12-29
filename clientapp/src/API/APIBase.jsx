
export function handleAPIResponse(response, onSuccess, onError) {

    if (response.ok) {
        response.json().then(data => {
            onSuccess(data);
        })
    } else {
        response.text().then(error => {
            onError(error);
        })
    }
}