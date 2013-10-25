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

var ImageUploader = function(dom,inputId){

    this.container = $(dom);

    this.inputId = inputId;

    this.imageType = this.container.attr('field-imageType');
    this.widths = this.container.attr('field-widths');
    this.maxWidth = this.container.attr('field-maxWidth');
    this.maxHeight = this.container.attr('field-maxHeight');
    this.minWidth = this.container.attr('field-minWidth');
    this.minHeight = this.container.attr('field-minHeight');
    this.maxSize = this.container.attr('field-maxSize');
    this.resolution = this.container.attr('field-resolution');
    this.tip = this.container.attr('field-tip');

    this.init();
}

ImageUploader.prototype = {

    fileInput : '{@if fieldDesc===0}<label class="control-label">${fieldDesc}：</label> {@/if}<input id="${inputId}" name="${fieldName}" type="file" data-desc="${fieldDesc}" /><input type="hidden" data-field name="${fieldName}" /><input type="hidden" data-field name="${fieldName}_s" /><span class="img-pre"><img src="{@if fieldValue==""}/lessgo/static/img/default_image.gif{@else}${fieldValue}{@/if}" width="100" height="100" /></span><span>${fieldTip}</span>',

    init : function(){

        this.render();

        this.bind();
    },

    render : function(){

        var mythis = this;
        mythis.container.append(juicer(mythis.fileInput,{
            fieldName:mythis.container.attr('field-name'),
            fieldValue:mythis.container.attr('field-value'),
            fieldDesc:mythis.container.attr('field-desc'),
            fieldTip:mythis.container.attr('field-tip'),
            inputId:mythis.inputId
        }));

        mythis.fieldName = mythis.container.attr('field-name')
    },

    bind : function(){
        var mythis = this;

        mythis.container.on('change','input[type=file]',function(){
            var thisInput = $(this);

            //获取用户输入的文件的后缀
            var suffix = thisInput.val().substr(thisInput.val().lastIndexOf('.')+1,thisInput.val().length);

            var flag = true;

            if(mythis.imageType){
                if(mythis.imageType.lastIndexOf(suffix)<0){
                    flag = false;
                    alert('图片格式应为：'+ mythis.imageType);
                }
            }

            if(flag){
                if ($(this).val()){
                    $.ajaxFileUpload({
                        url:"/imgageuplaod",
                        secureuri:false,
                        data:{
                            fileInputName : mythis.fieldName,
                            widths : mythis.widths,
                            maxWidth : mythis.maxWidth,
                            maxHeight : mythis.maxHeight,
                            minWidth : mythis.minWidth,
                            minHeight : mythis.minHeight,
                            maxSize : mythis.maxSize,
                            resolution : mythis.resolution
                        },
                        fileElementId:mythis.inputId,//文件上传的id属性
                        dataType: 'json',//返回值类型为json
                        success: function (data, status){//服务器成功响应处理函数
                            if(data.success){
                                mythis.container.find(":hidden[name="+mythis.container.attr('field-name')+"]").val(data.tmpfile);
                                mythis.container.find(":hidden[name="+mythis.container.attr('field-name')+"_s]").val(data.tmpfiles);
                                mythis.container.find("img").attr("src",data.tmpfile);
                            }else{
                                alert(data.msg);
                            }
                        },
                        error: function (data, status, e) {//服务器响应失败处理函数
                            alert('图片上传失败，请重新上传');
                        }
                    });
                }
            }


            thisInput.val("");
        });
    }

}