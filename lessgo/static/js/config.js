define(function (require, exports, module) {
    var development = false; // 发布时候改为false
    var version = '0.0.1'; // 版本，发布时候修改
    var debug = development ? 2 : false; // 线上为false，本地为2
    var href = this.location.href; //当前访问的链接
    var getHost = function (url) {
        var host = "null";
        var regex = /.*\:\/\/([^\/]*).*/;
        var match = url.match(regex);
        if (typeof match != "undefined" && null != match)
            host = match[1];
        return host;
    }
    var domain = this.location.host;//当前访问的域名
    var productUrl = 'http://' + domain + '/lessgo/static/js'; // 生产环境js代码路径
    var devUrl = 'http://' + domain + '/lessgo/static/js'; // 开发环境js源码路径
    var baseUrl = '';
    if (development) {
        baseUrl = devUrl;
    } else {
        baseUrl = productUrl;
    }

    seajs.config({
//        plugins: ['shim', 'nocache', 'debug'],
        plugins: ['shim'/*, 'combo'*/],
        alias: {
            'jquery': {
                src: productUrl + '/lib/jquery/1.10.2/jquery-1.10.2.min.js',
                exports: 'jQuery'
            },
            'bootstrap': {
                src: productUrl + '/lib/bootstrap/2.3.2/bootstrap.min.js',
                deps: ['jquery']
            },
            'juicer':{//模板合成
                src: productUrl + '/lib/juicer/v0.6.5-stable/juicer-min.js'
            },
            'jq-ui':{//时间控件、部分动画效果
                src: productUrl + '/lib/jq-ui/1.10.3/jquery-ui-1.10.3.custom.min.js',
                deps: ['jquery','jq-ui-i18n']
            },
            'jq-ui-i18n':{//时间控件、部分动画效果
                src: productUrl + '/lib/jq-ui/1.10.3/jquery.ui.datepicker-zh-CN.js',
                deps: ['jquery']
            },
            'ajaxfileupload':{//时间控件、部分动画效果
                src: productUrl + '/lib/jq-fileupload/ajaxfileupload.js',
                deps: ['jquery']
            },
            'imageuploader':{
                src: productUrl + '/lib/lessgo-form/imageuploader.js',
                deps: ['jquery','juicer','ajaxfileupload']
            },
            'kindeditor' : {
                src: productUrl + '/lib/kindeditor/4.1.7/lang/zh_CN.js',
                deps: ['jquery',productUrl + '/lib/kindeditor/4.1.7/kindeditor-min.js',productUrl + '/lib/kindeditor/4.1.7/themes/default/default.css']
            },
            'timepicker' : {
                src: productUrl + '/lib/jq-ui/timepicker.js',
                deps: ['jq-ui']
            },
            'lessgo-form':{
                src: productUrl + '/lib/lessgo-form/lessgo.form.js',
                deps: ['jquery',
                       'jq-ui',
                       'juicer',
                       'imageuploader',
                       'kindeditor',
                       'timepicker',
                       productUrl + '/lib/lessgo-form/area.js',
                       productUrl + '/lib/lessgo-form/timedim-week.js'
                ]
            },
            'jq-grid' : {
                src: productUrl + '/lib/jqGrid/jquery.jqGrid.min.js',
                deps: ['jquery','jq-ui','jq-grid-i18n',productUrl + '/lib/jqGrid/css/ui.jqgrid.css']
            },
            'jq-grid-i18n' : {
                src: productUrl + '/lib/jqGrid/grid.locale-cn.js',
                deps: ['jquery']
            },
            'lessgo-grid' : {
                src: productUrl + '/lib/lessgo-grid/lessgo.grid.js',
                deps: ['jq-grid','juicer','timepicker']
            }
        },
        version: version,
        preload: [],
        debug: true,
        base: baseUrl,
        charset: 'utf-8'
    });
});