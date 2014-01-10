function getParamsMap() {
    var map = {}
    if (document.URL.indexOf('?') > -1) {
        var paramsStr = document.URL.substring(document.URL.indexOf('?') + 1);
        var paramsArr = paramsStr.split('&');
        for (var index in paramsArr) {
            var param = paramsArr[index].split('=');
            var key = param[0];
            var value = param[1];
            var urlParamValue = map[key];
            if (urlParamValue == null) {
                urlParamValue = value;
                map[key] = urlParamValue;
            }
        }
    }
    return map;
}