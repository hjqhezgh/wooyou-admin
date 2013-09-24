jQuery.fn.form = function (opts) {

    var _this = this;
    var myform = this.find('form');

    var hiddenFieldTemp = '<input name="${fieldName}" value="${fieldValue}" type="hidden" data-field validate="${fieldValidate}" />';
    var textFieldTemp = '<label class="control-label">${fieldDesc}：</label><div class="controls"><input name="${fieldName}" value="${fieldValue}" type="text" class="input-small" data-field validate="${fieldValidate}" data-desc="${fieldDesc}" {@if fieldReadOnly=="true"}readonly{@/if} /></div>';
    var selectFieldTemp = '<label class="control-label">${fieldDesc}：</label><div class="controls"><select name="${fieldName}" data-field validate="${fieldValidate}" data-desc="${fieldDesc}" {@if fieldReadOnly=="true"}disabled{@/if} ></select></div>';
    var optionTemp = '<option value="${value}">${desc}</option>';
    var textAreaTemp = '<label class="control-label">${fieldDesc}：</label><div class="controls"><textarea name="${fieldName}" data-field validate="${fieldValidate}" data-desc="${fieldDesc}" {@if fieldReadOnly=="true"}readonly{@/if} >${fieldValue}</textarea></div>';

    //初始化html编辑器对象，一个表单只能拥有一个编辑器
    this.editor = "";

    opts = jQuery.extend({
    }, opts || {});

    //进行各种类型表单控件的渲染和初始化操作
    this.render = function(){

        var mythis = this;

        //隐藏域
        myform.find('[field-type=hidden]').each(function(index,element){

            var defaultValue = $(element).attr('field-defaultValue');

            if($(element).attr('field-value')){
                defaultValue = $(element).attr('field-value');
            }

            $(element).append(juicer(hiddenFieldTemp,{
                fieldName:$(element).attr('field-name'),
                fieldValue:defaultValue
            }));

        });

        //文本域
        myform.find('[field-type=text]').each(function(index,element){
            $(element).append(juicer(textFieldTemp,{
                fieldName:$(element).attr('field-name'),
                fieldValue:$(element).attr('field-value'),
                fieldDesc:$(element).attr('field-desc'),
                fieldReadOnly:$(element).attr('field-readonly'),
                fieldValidate:$(element).attr('field-validate')
            }));
        });

        //大文本域
        myform.find('[field-type=textarea]').each(function(index,element){
            $(element).append(juicer(textAreaTemp,{
                fieldName:$(element).attr('field-name'),
                fieldValue:$(element).attr('field-value'),
                fieldDesc:$(element).attr('field-desc'),
                fieldReadOnly:$(element).attr('field-readonly'),
                fieldValidate:$(element).attr('field-validate')
            }));
        });

        //日期域
        myform.find('[field-type=date]').each(function(index,element){
            $(element).append(juicer(textFieldTemp,{
                fieldName:$(element).attr('field-name'),
                fieldValue:$(element).attr('field-value'),
                fieldDesc:$(element).attr('field-desc'),
                fieldReadOnly:$(element).attr('field-readonly'),
                fieldValidate:$(element).attr('field-validate')
            }));

            if($(element).attr('field-readonly')!="true"){
                $(element).find('input').datepicker({
                    dateFormat : 'yy-mm-dd'
                });
            }

        });

        //时间域
        myform.find('[field-type=datetime]').each(function(index,element){
            $(element).append(juicer(textFieldTemp,{
                fieldName:$(element).attr('field-name'),
                fieldValue:$(element).attr('field-value'),
                fieldDesc:$(element).attr('field-desc'),
                fieldReadOnly:$(element).attr('field-readonly'),
                fieldValidate:$(element).attr('field-validate')
            }));

            if($(element).attr('field-readonly')!="true"){
                $(element).find('input').datetimepicker({
                    dateFormat : 'yy-mm-dd',
                    timeFormat: 'HH:mm:ss'
                });
            }

        });

        //本地下拉框
        myform.find('[field-type=localSelect]').each(function(index,element){
            $(element).append(juicer(selectFieldTemp,{
                fieldName:$(element).attr('field-name'),
                fieldDesc:$(element).attr('field-desc'),
                fieldReadOnly:$(element).attr('field-readonly'),
                fieldValidate:$(element).attr('field-validate')
            }));

            var objects = eval($(element).attr('field-localData'));

            var select = $(element).find('select');

            select.append(juicer(optionTemp,{
                value:'',
                desc:'全部'
            }));

            for(var i=0;i<objects.length;i++){
                select.append(juicer(optionTemp,{
                    value:objects[i].value,
                    desc:objects[i].desc
                }));
            }

            if($(element).attr('field-value')){
                select.find("option[value="+$(element).attr('field-value')+"]").attr("selected","selected");
            }

        });


        //远程下拉框
        myform.find('[field-type=remoteSelect]').each(function(index,element){
            $(element).append(juicer(selectFieldTemp,{
                fieldName:$(element).attr('field-name'),
                fieldDesc:$(element).attr('field-desc'),
                fieldReadOnly:$(element).attr('field-readonly'),
                fieldValidate:$(element).attr('field-validate')
            }));

            var url = $(element).attr('field-url');
            var valueField = $(element).attr('field-valueField');
            var descField = $(element).attr('field-descField');

            var select = $(element).find('select');
            var value = $(element).attr('field-value');

            select.append(juicer(optionTemp,{
                value:'',
                desc:'全部'
            }));

            $.get(url,{},function(data){
               if(data.success){
                   var objects = data.datas;

                   for(var i=0;i<objects.length;i++){
                       select.append(juicer(optionTemp,{
                           value:objects[i][valueField],
                           desc:objects[i][descField]
                       }));
                   }

                   if(value){
                       select.find("option[value="+value+"]").attr("selected","selected");
                   }

               } else{
                   alert(data.msg);
               }
            },'json');

        });

        //area
        myform.find('[field-type=area]').each(function(index,element){
            new area($(this),$(element).attr('field-desc'),$(element).attr('field-name'),$(element).attr('field-readonly'),$(element).attr('field-value'));
        });

        //timedim-week
        myform.find('[field-type=dateweek]').each(function(index,element){
            new timeDimWeek($(this),$(element).attr('field-desc'),$(element).attr('field-name'),$(element).attr('field-readonly'));
        });

        //图片上传
        myform.find('[field-type=image]').each(function(index,element){
            new ImageUploader($(this));
        });

        //html编辑框
        myform.find('[field-type=htmleditor]').each(function(index,element){
            $(element).append(juicer(textAreaTemp,{
                fieldName:$(element).attr('field-name'),
                fieldValue:$(element).attr('field-value'),
                fieldDesc:$(element).attr('field-desc'),
                fieldReadOnly:$(element).attr('field-readonly'),
                fieldValidate:$(element).attr('field-validate')
            }));

            $(element).find('textarea').attr('htmleditor',"true");

            KindEditor.ready(function(K) {
                mythis.editor = K.create('textarea[name="'+$(element).find('textarea').attr('name')+'"]', {
                    items:[
                        'source', '|', 'undo', 'redo', '|', 'preview', 'cut', 'copy', 'paste',
                        'plainpaste', 'wordpaste', '|', 'justifyleft', 'justifycenter', 'justifyright',
                        'justifyfull', 'insertorderedlist', 'insertunorderedlist', 'indent', 'outdent', 'subscript',
                        'superscript', 'clearhtml', 'quickformat', 'selectall', '|', 'fullscreen', '/',
                        'formatblock', 'fontname', 'fontsize', '|', 'forecolor', 'hilitecolor', 'bold',
                        'italic', 'underline', 'strikethrough', 'lineheight', 'removeformat', '|', 'image',
                        'table', 'hr', 'emoticons', 'baidumap', 'pagebreak',
                        'anchor', 'link', 'unlink'
                    ],
                    uploadJson : '/kindeditorImageUpload'
                });
            });
        });


        this.removeConfigAttr();

        //美化按钮
        this.find('a[data-action=reset]').button();

        this.find('a[data-action=save]').button();

        myform.find('.nav-tabs').find('li:first').addClass('active');
    }

    this.bindEvent = function(){

        var mythis = this;

        this.find('a[data-action=reset]').click(function(e){
            myform[0].reset();
        });

        this.find('a[data-action=save]').click(function(e){
            var params = {};

            var flag = true;

            //如果存在唯一的htmleditor，需要将编辑器的内容和textarea的内容做一下同步
            if(mythis.editor){
                mythis.editor.sync();
            }

            myform.find('[data-field]').each(function(index,field){
                if(_this.validate($(field))){
                    params[$(field).attr('name')] = $(field).val();
                }else{
                    flag = false;
                    return false;
                }
            });

            if(flag){
                $.post(opts.saveUrl,params,function(data){
                    if(data.success){
                        opts.successCallback(data);
                    }else{
                        opts.failCallback(data);
                    }
                },'json');
            }

        });

        this.find('.nav-tabs a').click(function (e) {
            e.preventDefault();
            $(this).tab('show');
        });
    }

    this.removeConfigAttr = function(){
        this.find('div.form-field').removeAttr("field-type")
        .removeAttr("field-name")
        .removeAttr("field-value")
        .removeAttr("field-desc")
        .removeAttr("field-localData")
        .removeAttr("field-url")
        .removeAttr("field-valueField")
        .removeAttr("field-descField")
        .removeAttr("field-readonly")
        .removeAttr("field-defaultValue")
        .removeAttr("field-validate");
    }

    //验证对象至少需要validate方法
    this.getValidateObject = function(rule){

        if(rule=='notNull'){
            return {
                msg : '不能为空',
                validate : function(element){
                    if(!element.val()){
                        alert(element.attr('data-desc')+this.msg);
                        return false;
                    }else{
                        return true;
                    }
                }
            }
        }else if (rule.indexOf('len')>-1){
            var minLength = rule.substring(rule.lastIndexOf("(")+1,rule.lastIndexOf("-"));
            var maxLength = rule.substring(rule.lastIndexOf("-")+1,rule.lastIndexOf(")"));
            return {
                minLength : minLength,
                maxLength : maxLength,
                msg : '${desc}需要输入${min}-${max}个字符，汉字为2个字符',
                validate : function(element){
                    var b1 = element.val();
                    var valueLength = b1.match(/[^ -~]/g) == null ? b1.length : b1.length + b1.match(/[^ -~]/g).length;
                    if(valueLength>maxLength || valueLength<minLength){
                        alert(juicer(this.msg,{
                            desc : element.attr('data-desc'),
                            min : this.minLength,
                            max : this.maxLength
                        }));
                    }else{
                        return true;
                    }
                }
            }
        }else if(rule == "email"){
            return {
                msg : '应为邮箱格式',
                validate : function(element){
                    var myreg = /^([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+@([a-zA-Z0-9]+[_|\_|\.]?)*[a-zA-Z0-9]+\.[a-zA-Z]{2,3}$/;
                    if(!myreg.test(element.val())){
                        alert(element.attr('data-desc')+this.msg);
                        return false;
                    }else{
                        return true;
                    }
                }
            }
        }else if(rule == "num") {
            return {
                msg : '应为数字格式',
                validate : function(element){
                    var myreg = /^[0-9]*$/;
                    if(!myreg.test(element.val())){
                        alert(element.attr('data-desc')+this.msg);
                        return false;
                    }else{
                        return true;
                    }
                }
            }
        }else if(rule == "float") {
            return {
                msg : '应为浮点数字格式',
                validate : function(element){
                    var myreg = /^\d+(\.\d+)?$/;
                    if(!myreg.test(element.val())){
                        alert(element.attr('data-desc')+this.msg);
                        return false;
                    }else{
                        return true;
                    }
                }
            }
        }


        return

    }

    this.validate = function(element){

        if(element.attr('validate')){
            var validates = element.attr('validate').split(',');
            for(var i=0;i<validates.length;i++){
                var validateNow = _this.getValidateObject(validates[i]);
                if(validateNow){
                    if(!validateNow.validate(element)){
                        return false;
                    }
                }
            }
        }

        return true;
    }

    this.render();
    this.bindEvent();

    return this;
}


