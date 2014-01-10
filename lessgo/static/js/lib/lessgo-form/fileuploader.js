/*****
 * 文件上传控件
 */

var FileUploader = function(dom,inputId){

    this.container = $(dom);

    this.inputId = inputId;

    this.tip = this.container.attr('field-tip');
    this.url = this.container.attr('field-fileUploadUrl');
    this.fieldName = this.container.attr('field-name');

    this.init();
}

FileUploader.prototype = {

    fileInput : '<input id="${inputId}" name="${fieldName}" type="file" data-desc="${fieldDesc}" /><input type="hidden" data-field name="${fieldName}" /><input type="hidden" data-field name="${fieldName}_suffix" /><input type="hidden" data-field name="${fieldName}_fileName" /><input type="hidden" data-field name="${fieldName}_fileSize" /><input type="hidden" data-field name="${fieldName}_srcFileName" /><span class="loading" style="display: none;"><img src="/lessgo/static/img/loading.gif" />上传中，请勿关闭页面</span><span>${fieldTip}</span>',

    init : function(){

        this.render();

        this.bind();
    },

    render : function(){

        var mythis = this;

        if(mythis.container.attr('field-desc')!=0){
            mythis.fileInput = '<label class="control-label">${fieldDesc}：</label>' + mythis.fileInput;
        }

        mythis.container.append(juicer(mythis.fileInput,{
            fieldName:mythis.container.attr('field-name'),
            fieldDesc:mythis.container.attr('field-desc'),
            fieldTip:mythis.container.attr('field-tip'),
            inputId:mythis.inputId
        }));
    },

    bind : function(){
        var mythis = this;

        mythis.container.on('change','input[type=file]',function(){
            var thisInput = $(this);

            //获取用户输入的文件的后缀
            var suffix = thisInput.val().substr(thisInput.val().lastIndexOf('.')+1,thisInput.val().length);

            var flag = true;

            //todo 后缀验证

            if(flag){
                if ($(this).val()){
                    mythis.container.find('.loading').show();

                    $.ajaxFileUpload({
                        url:mythis.url,
                        secureuri:false,
                        fileElementId:mythis.inputId,//文件上传的id属性
                        dataType: 'json',//返回值类型为json
                        success: function (data, status){//服务器成功响应处理函数
                            mythis.container.find('.loading').hide();
                            if(data.success){
                                mythis.container.find(":hidden[name="+mythis.container.attr('field-name')+"_suffix]").val(data.suffix);
                                mythis.container.find(":hidden[name="+mythis.container.attr('field-name')+"_fileName]").val(data.fileName);
                                mythis.container.find(":hidden[name="+mythis.container.attr('field-name')+"_srcFileName]").val(data.srcFileName);
                                mythis.container.find(":hidden[name="+mythis.container.attr('field-name')+"_fileSzie]").val(data.fileSzie);
                                mythis.container.find(":hidden[name="+mythis.container.attr('field-name')+"_srcFileName]").change();//触发值改变事件
                            }else{
                                alert(data.msg);
                            }
                        },
                        error: function (data, status, e) {//服务器响应失败处理函数
                            mythis.container.find('.loading').hide();
                            alert('文件上传失败，请重新上传');
                        }
                    });
                }
            }


            thisInput.val("");
        });
    }

}