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

    var windowWidth = document.documentElement.clientWidth - 230;
    var windowHeight = document.documentElement.clientHeight - 85;

    top.jQuery.layer({
        type : 2,
        fix : true,
        moveOut : false,
        move : false,
        shade : false,
        closeBtn : [0 , true],
        title: '',
        border : [5 , 0.7 , '#272822', true],
        offset : ['75px','220px'],
        shade : [0.5 , '#000' , true],
        area : [windowWidth+'px',windowHeight+'px'],
        iframe : {src : url},
        close : function(index){
            top.windowNum--;
            top.layer.close(index)
        }
    });

    top.windowNum++;
}

// tofix 这里应该是要获取最近的那个窗口。而不一定是parent
jQuery.getParentComponentObject = function (id) {
    return parent[id];
}