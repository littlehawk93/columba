
export function formatLocationString(location) {

    if(location) {
        var result = `${location.city}, ${location.state}, ${location.zip}`;
        return result.replace(/(^(\s|,)+)|((\s|,)+$)/, "");
    }
    return "";
}