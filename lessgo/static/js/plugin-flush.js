(function (f) {
    function j(b, c) {
        c || (c = []);
        for (var a = 0; a < b.length; a++) {
            var d = b[a];
            k(d) ? j(d, c) : c.push(d)
        }
        return c
    }

    var l = [].slice, m = f.use, e = [];
    f.use = function (b, c) {
        e.push([k(b) ? b : [b], c]);
        return f
    };
    f.flush = function () {
        for (var b = [], c = [], a = 0; a < e.length; a++) {
            var d = e[a];
            b[a] = d[0];
            c[a] = d[1]
        }
        e.length = 0;
        return m(j(b), function () {
            var a;
            a = l.call(arguments);
            for (var g = [], d, f, e = 0; e < b.length; e++) {
                g[e] = f = [];
                d = b[e];
                for (var h = 0; h < d.length; h++)f[h] = a.shift()
            }
            a = g;
            for (g = 0; g < c.length; g++)(d = c[g]) && d.apply(null, a[g])
        })
    };
    var n = {}.toString, k = Array.isArray || function (b) {
        return"[object Array]" === n.call(b)
    }
})(seajs);
