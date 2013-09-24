/**
 * Created with IntelliJ IDEA.
 * User: Chenm
 * Date: 13-4-27
 * Time: 下午11:17
 * To change this template use File | Settings | File Templates.
 */
define(function () {

    var rules = [];
    var host = this.location.host
    rules.push(function (url) {
        if (url.indexOf('public/js/src') > 0) {
            url = url.replace('public/js', 'el/js');
            url = url.replace(host, 'localhost:8082');
        }
        return url;
    });
    // set map rules
    seajs.config({'map': rules});
});
