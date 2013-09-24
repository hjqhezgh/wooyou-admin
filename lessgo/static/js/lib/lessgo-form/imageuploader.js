//为fileupload插件打一个补丁
jQuery.extend({
    handleError: function (s, xhr, status, e) {
        if (s.error)
            s.error(xhr, status, e);
        if (s.global)
            jQuery.event.trigger("ajaxError", [xhr, s, e]);
    },
    active: 0,
    httpSuccess: function (xhr) {
        try {
            return!xhr.status && location.protocol == "file:" || (xhr.status >= 200 && xhr.status < 300) || xhr.status == 304 || xhr.status == 1223 || jQuery.browser.safari && xhr.status == undefined;
        } catch (e) {
        }
        return false;
    }
});

/*****
 * 图片上传预览控件、
 */

var ImageUploader = function(dom){

    this.container = $(dom);

    this.init();
}

ImageUploader.prototype = {

    fileInput : '<label class="control-label">${fieldDesc}：</label><input id="fileToUpload" name="${fieldName}" type="file" data-desc="${fieldDesc}" /><input type="hidden" data-field name="${fieldName}" /><span class="img-pre"><img src="/lessgo/static/img/default_image.gif" width="100" height="100" /></span>',

    init : function(){

        this.render();

        this.bind();
    },

    render : function(){

        var mythis = this;

        mythis.container.append(juicer(mythis.fileInput,{
            fieldName:mythis.container.attr('field-name'),
            fieldDesc:mythis.container.attr('field-desc')
        }));
    },

    bind : function(){
        var mythis = this;

        mythis.container.on('change','input[type=file]',function(){
            var thisInput = $(this);

            if ($(this).val()){
                $.ajaxFileUpload({
                    url:"/imgageuplaod",
                    secureuri:false,
                    fileElementId:'fileToUpload',//文件上传的id属性
                    dataType: 'json',//返回值类型为json
                    success: function (data, status){//服务器成功响应处理函数
                        if(data.success){
                            mythis.container.find(":hidden").val(data.tmpfile);
                            mythis.container.find("img").attr("src",data.tmpfile);
                        }else{
                            alert('图片上传视频，请重新上传');
                        }
                    },
                    error: function (data, status, e) {//服务器响应失败处理函数
                        alert('图片上传视频，请重新上传');
                    }
                });
            }

            thisInput.val("");
        });
    }

}