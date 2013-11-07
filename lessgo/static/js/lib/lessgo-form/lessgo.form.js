jQuery.fn.form = function (opts) {

    var _this = this;
    var myform = this.find('form');

    var hiddenFieldTemp = '<input name="${fieldName}" value="${fieldValue}" type="hidden" data-field validate="${fieldValidate}" />';
    var textFieldTemp = '<label class="control-label">${fieldDesc}：</label><div class="controls"><input name="${fieldName}" value="${fieldValue}" type="text" class="input-small" data-field validate="${fieldValidate}" data-desc="${fieldDesc}" {@if fieldReadOnly=="true"}readonly{@/if} /><span>${fieldTip}</span></div>';
    var selectFieldTemp = '<label class="control-label">${fieldDesc}：</label><div class="controls"><select name="${fieldName}" data-field validate="${fieldValidate}" data-desc="${fieldDesc}" {@if fieldReadOnly=="true"}disabled{@/if} ></select><span>${fieldTip}</span></div>';
    var optionTemp = '<option value="${value}">${desc}</option>';
    var textAreaTemp = '<label class="control-label">${fieldDesc}：</label><div class="controls"><textarea name="${fieldName}" rows="7" style="width:60%" data-field validate="${fieldValidate}" data-desc="${fieldDesc}" {@if fieldReadOnly=="true"}readonly{@/if} >${fieldValue}</textarea><span>${fieldTip}</span></div>';
    var checkboxContainerTemp = '<label class="control-label">${fieldDesc}：</label><div class="controls"><input data-field type="hidden" name="${fieldName}" validate="${fieldValidate}" data-desc="${fieldDesc}"  /></div>';
    var checkboxTemp = '<input validate="${fieldValidate}" value="${value}" type="checkbox" {@if fieldReadOnly=="true"}readonly{@/if} />${desc}&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;';

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
                fieldTip:$(element).attr('field-tip'),
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
                fieldTip:$(element).attr('field-tip'),
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
                fieldTip:$(element).attr('field-tip'),
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
                fieldTip:$(element).attr('field-tip'),
                fieldValidate:$(element).attr('field-validate')
            }));

            if($(element).attr('field-readonly')!="true"){
                $(element).find('input').datetimepicker({
                    dateFormat : 'yy-mm-dd',
                    timeFormat: 'HH:mm:ss'
                });
            }

        });

        //无日期时间域
        myform.find('[field-type=time]').each(function(index,element){
            $(element).append(juicer(textFieldTemp,{
                fieldName:$(element).attr('field-name'),
                fieldValue:$(element).attr('field-value'),
                fieldDesc:$(element).attr('field-desc'),
                fieldReadOnly:$(element).attr('field-readonly'),
                fieldTip:$(element).attr('field-tip'),
                fieldValidate:$(element).attr('field-validate')
            }));

            if($(element).attr('field-readonly')!="true"){
                $(element).find('input').timepicker();
            }

        });

        //本地下拉框
        myform.find('[field-type=localSelect]').each(function(index,element){
            $(element).append(juicer(selectFieldTemp,{
                fieldName:$(element).attr('field-name'),
                fieldDesc:$(element).attr('field-desc'),
                fieldReadOnly:$(element).attr('field-readonly'),
                fieldTip:$(element).attr('field-tip'),
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
                fieldTip:$(element).attr('field-tip'),
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
            new ImageUploader($(this),$(this).attr('field-name'));
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

            var uploadUrl = '/kindeditorImageUpload';

            if($(element).attr('field-uploadUrl')){
                uploadUrl = $(element).attr('field-uploadUrl');
            }

            mythis.editor = KindEditor.create('textarea[name="' + $(element).find('textarea').attr('name') + '"]', {
                items: [
                    'source', '|', 'undo', 'redo', '|', 'preview', 'cut', 'copy', 'paste',
                    'plainpaste', 'wordpaste', '|', 'justifyleft', 'justifycenter', 'justifyright',
                    'justifyfull', 'insertorderedlist', 'insertunorderedlist', 'indent', 'outdent', 'subscript',
                    'superscript', 'clearhtml', 'quickformat', 'selectall', '|', 'fullscreen', '/',
                    'formatblock', 'fontname', 'fontsize', '|', 'forecolor', 'hilitecolor', 'bold',
                    'italic', 'underline', 'strikethrough', 'lineheight', 'removeformat', '|', 'image',
                    'table', 'hr', 'emoticons', 'baidumap', 'pagebreak',
                    'anchor', 'link', 'unlink'
                ],
                uploadJson: uploadUrl
            });
        });


        //多对多控件
        myform.find('[field-type=checkbox]').each(function(index,element){
            var url = $(element).attr('field-url');
            var valueField = $(element).attr('field-valueField');
            var descField = $(element).attr('field-descField');

            $(element).append(juicer(checkboxContainerTemp,{
                fieldName:$(element).attr('field-name'),
                fieldDesc:$(element).attr('field-desc'),
                fieldValidate:$(element).attr('field-validate')
            }));

            var checkboxContainer = $(element).find('.controls');

            if($(element).attr('field-value')){
                $(element).find('[data-field]').val($(element).attr('field-value'));
            }

            $.get(url,{},function(data){
                if(data.success){
                    var objects = data.datas;

                    for(var i=0;i<objects.length;i++){
                        checkboxContainer.append(juicer(checkboxTemp,{
                            fieldReadOnly:$(element).attr('field-readonly'),
                            value:objects[i][valueField],
                            desc:objects[i][descField]
                        }));
                    }

                    $(element).on('click',':checkbox',function(){
                        $(element).find('[data-field]').val('');

                        var value = "";

                        $(element).find(':checkbox').each(function(){
                            if($(this).prop('checked')){
                                value += $(this).val()+ ","
                            }
                        });

                        if(value.length > 0 ){
                            value = value.substring(0,value.length-1);
                        }

                        $(element).find('[data-field]').val(value);
                    });

                    var value = $(element).attr('field-value');
                    var values = value.split(',')
                    for(var i=0;i<values.length;i++){
                        $(element).find('[value='+values[i]+']').prop('checked',true);
                    }

                } else{
                    alert(data.msg);
                }
            },'json');
        });


        //this.removeConfigAttr();

        //美化按钮
        this.find('a[data-action=reset]').button();

        this.find('a[data-action=save]').button();

        myform.find('.nav-tabs').find('li:first').addClass('active');
        myform.find('.tab-content').find('.tab-pane:first').addClass('active');

        //需要渲染用户自定义的按钮
        if(opts.buttons.length>0){
            for(var i=0;i<opts.buttons.length;i++){
                _this.find('.form-buttons').find('.controls').append('<a href="javascript:void(0);" class="btn '+opts.buttons[i].buttonClass+'">'+opts.buttons[i].desc+'</a>&nbsp;');
                _this.find('.form-buttons').find('a:last').bind('click',{'b':opts.buttons[i]},function(e){
                    e.data.b.handler(mythis);
                }).button();
            }
        }

        if(opts.afterRender){
            opts.afterRender(mythis);
        }
    }

    this.loadValue = function(field,value){
        var mythis = this;
        myform.find('[field-name='+field+']').attr("field-value",value);
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
                if(opts.beforeSave){
                    flag = opts.beforeSave(mythis);
                }
            }else{
                return;
            }


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


    if(opts.loadUrl){
        var url = opts.loadUrl;
        if(document.URL.indexOf('?')>-1){
            url = opts.loadUrl + document.URL.substring(document.URL.lastIndexOf('?'));
        }

        $.get(url,{},function(data){
            if(data.success){
                for(var i=0;i<data.datas.length;i++){
                    _this.loadValue(data.datas[i].field,data.datas[i].value);
                }
                _this.render();
                _this.bindEvent();
            }else{
                alert(data.msg);
            }
        },'json');
    }else{
        _this.render();
        _this.bindEvent();
    }

    return this;
}


