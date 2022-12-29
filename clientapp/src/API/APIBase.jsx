
export function handleAPIResponse(response, onSuccess, onError) {

    if (response.ok) {
        if (response.status === 204) {
            onSuccess(null);
        } else {
            response.json().then(data => {
                onSuccess(data);
            });
        }
    } else {
        response.text().then(error => {
            onError(error);
        });
    }
}