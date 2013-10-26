jQuery.openIframeWindow = function (opts,event) {

    opts = jQuery.extend({
        width : 500,
        height : 500
    }, opts || {});

    var url = opts.url;

    if(url.lastIndexOf('?')>-1){
        url = url + "&parentComponentId="+ opts.parentComponent;
    }else{
        url = url + "?parentComponentId="+ opts.parentComponent;
    }

    top.jQuery.layer({
        type : 2,
        fix : true,
        moveOut : false,
        shade : [0.8 , '#E3E3E3' , true],
        shadeClose : false,
        border : [10 , 0.7 , '#272822', true],
        title : opts.title,
        offset : ['50px',''],
        area : [opts.width+'px',opts.height+'px'],
        iframe : {src : url}
    });
}

// tofix 这里应该是要获取最近的那个窗口。而不一定是parent
jQuery.getParentComponentObject = function (id) {
    return parent[id];
}