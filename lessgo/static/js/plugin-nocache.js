(function (d) {
    var e = "seajs-nocache=" + (new Date).getTime(), c = {};
    d.on("fetch", function (a) {
        var b = a.requestUri || a.uri;
        -1 === b.indexOf("seajs-nocache=") && (b += (-1 === b.indexOf("?") ? "?" : "&") + e, c[b] = a.uri, a.requestUri = b)
    });
    d.on("define", function (a) {
        c[a.uri] && (a.uri = c[a.uri])
    })
})(seajs);
