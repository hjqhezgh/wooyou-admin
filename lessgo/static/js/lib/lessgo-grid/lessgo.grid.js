jQuery.fn.grid = function (opts) {

    var componentId = this.attr("id");

    var _this = this;

    var lastId = 0;

    var myform = this.find('form');

    var textFieldTemp = '${fieldDesc}：<input name="${fieldName}" type="text" data-field data-desc="${fieldDesc}" data-searchType="${fieldSearchType}" />';
    var selectFieldTemp = '${fieldDesc}：<select name="${fieldName}" data-field  data-desc="${fieldDesc}" data-searchType="${fieldSearchType}"></select>';
    var optionTemp = '<option value="${value}">${desc}</option>';

    opts = jQuery.extend({
        width : 600,
        height : 600
    }, opts || {});

    window[componentId] = this;
    if(opts.componentId){
        window[opts.componentId] = this;
    }

    this.getGrid = function(){
        return this.grid;
    }

    this.render = function () {

        this.grid = $('#' + opts.tableId).jqGrid({
            url: opts.dataurl,
            datatype: 'json',
            colNames: opts.colNames,
            width: opts.width,
            height: opts.height,
            colModel: opts.colModel,
            rowNum: opts.rowNum,
            multiselect: opts.mutiSelect,
            rowList:[opts.rowNum,50,100,200,500,1000],
            pager: '#' + opts.pageId,
            viewrecords: true,
            onSelectRow: function (id) {
                lastId = id;
                //todo 表格多选框效果，有空了再做吧
                /*
                $('#' + opts.tableId).find('[data-action=gird-checkbox]').prop('checked',false);
                $('#' + opts.tableId).find('[data-action=gird-checkbox][value='+id+']').prop('checked',true);
                */
            },
            loadComplete: function () {
                lastId = 0;
            }
        }).navGrid('#' + opts.pageId, {edit: false, add: false, del: false, search: false});

        //渲染多选框区域
        if(_this.find('.checkbox-box').size()>0){
            this.loadCheckboxContainer();
        }

        this.renderSearchForm();

        this.bindButtonEvent();
    }

    this.loadCheckboxContainer = function(){
        var checkBoxLoadUrl = _this.find('.checkbox-box').attr('data-loadUrl');

        if(document.URL.lastIndexOf('?')>-1){

            if(checkBoxLoadUrl.lastIndexOf('?')>-1){
                checkBoxLoadUrl += "&"
            }else{
                checkBoxLoadUrl += "?"
            }

            checkBoxLoadUrl += document.URL.substring(document.URL.lastIndexOf('?')+1,document.URL.length);
        }

        _this.find('.checkbox-container').html('');

        $.get(checkBoxLoadUrl,{},function(data){
            if(data.success){
                for(var i=0;i<data.datas.length;i++){
                    _this.find('.checkbox-container').append('<input type="checkbox" checked="checked" value="'+data.datas[i].value+'"/>'+data.datas[i].desc+'&nbsp;');
                }
            }else{
                alert(data.msg);
            }
        },'json');
    }

    this.renderSearchForm = function(){

        //文本域
        myform.find('[field-inputType=text]').each(function(index,element){
            $(element).append(juicer(textFieldTemp,{
                fieldName:$(element).attr('field-name'),
                fieldDesc:$(element).attr('field-desc'),
                fieldSearchType:$(element).attr('field-searchType')
            }));
        });

        //日期域
        myform.find('[field-inputType=date]').each(function(index,element){
            $(element).append(juicer(textFieldTemp,{
                fieldName:$(element).attr('field-name'),
                fieldDesc:$(element).attr('field-desc'),
                fieldSearchType:$(element).attr('field-searchType')
            }));

            $(element).find('input').datepicker({
                dateFormat : 'yy-mm-dd'
            });
        });

        //时间域
        myform.find('[field-inputType=datetime]').each(function(index,element){
            $(element).append(juicer(textFieldTemp,{
                fieldName:$(element).attr('field-name'),
                fieldDesc:$(element).attr('field-desc'),
                fieldSearchType:$(element).attr('field-searchType')
            }));

            $(element).find('input').datetimepicker({
                dateFormat : 'yy-mm-dd',
                timeFormat: 'HH:mm:ss'
            });
        });

        //本地下拉框
        myform.find('[field-inputType=localSelect]').each(function(index,element){
            $(element).append(juicer(selectFieldTemp,{
                fieldName:$(element).attr('field-name'),
                fieldDesc:$(element).attr('field-desc'),
                fieldSearchType:$(element).attr('field-searchType')
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
        });

        //远程下拉框
        myform.find('[field-inputType=remoteSelect]').each(function(index,element){
            $(element).append(juicer(selectFieldTemp,{
                fieldName:$(element).attr('field-name'),
                fieldDesc:$(element).attr('field-desc'),
                fieldSearchType:$(element).attr('field-searchType')
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

            if($(element).attr('field-parentSelect')){
                myform.find('[name='+$(element).attr('field-parentSelect')+']').change(function(){
                    var parentValue = $(this).val();
                    if(!parentValue){//父级下拉框没有值，则清空子级下拉框
                        select.html('<option value="">全部</option>');
                    }else{
                        select.html('<option value="">全部</option>');
                        $.get(url,{id :parentValue },function(data){
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
                    }
                });
            }else{
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
            }

        });

        //this.removeConfigAttr();

        //美化按钮
        this.find('a[data-action=reset]').button();

        this.find('a[data-action=search]').button();

        this.bindFormEvent();
    }

    this.removeConfigAttr = function(){
        this.find('div.form-field')
            .removeAttr("field-type")
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

    this.bindFormEvent = function(){

        this.find('a[data-action=search]').click(function(){

            var url = opts.dataurl;

            _this.find('[data-field]').each(function(index,search){
                if($(search).parent().attr("field-char14")=="true"){//清除时间格式为14位char入20130101121212
                    url+='&'+$(search).attr('name')+'-'+$(search).attr('data-searchType')+'='+$(search).val().replace(" ","").replace(new RegExp("-","gm"),"").replace(new RegExp(":","gm"),"");
                }else if($(search).parent().attr("field-char8")=="true"){//清除时间格式为8位char入20130101
                    url+='&'+$(search).attr('name')+'-'+$(search).attr('data-searchType')+'='+$(search).val().replace(" ","").replace(new RegExp("-","gm"),"").replace(new RegExp(":","gm"),"");
                }else{
                    url+='&'+$(search).attr('name')+'-'+$(search).attr('data-searchType')+'='+$(search).val();
                }
            });

            $('#' + opts.tableId).jqGrid('setGridParam',{url:url,page:1}).trigger("reloadGrid");

        });

        this.find('a[data-action=reset]').click(function(){
            _this.find('form').find('[field-parentSelect!=""]').find('select').html('<option value="">全部</option>');
            _this.find('form')[0].reset();
        });

    }

    this.bindButtonEvent = function(){

        _this.on('click','[data-action=checkboxSave]',function(event){
            thisButton = $(this);
            event.preventDefault();
            var checkBoxSaveUrl = _this.find('.checkbox-box').attr('data-saveUrl');

            if(document.URL.lastIndexOf('?')>-1){

                if(checkBoxSaveUrl.lastIndexOf('?')>-1){
                    checkBoxSaveUrl += "&"
                }else{
                    checkBoxSaveUrl += "?"
                }

                checkBoxSaveUrl += document.URL.substring(document.URL.lastIndexOf('?')+1,document.URL.length);
            }

            var cids = '';
            var noids = '';

            _this.find('.checkbox-box').find(':checkbox').each(function(){
                if($(this).prop('checked')){
                    cids += $(this).val()+",";
                }else{
                    noids += $(this).val()+",";
                }
            });

            if(cids){
                cids = cids.substr(0,cids.length-1);
            }

            if(noids){
                noids = noids.substr(0,noids.length-1);
            }

            $.post(checkBoxSaveUrl,{cids : cids,noids : noids},function(data){
                if(data.success){
                    alert(data.msg);
                    _this.loadCheckboxContainer();
                }else{
                    alert(data.msg);
                }
            },'json');
        });


        _this.on('click','[data-action=openIframeWindow]',function(event){
            thisButton = $(this);
            event.preventDefault();

            var url = thisButton.attr('href');

            if(document.URL.lastIndexOf('?')>-1){

                if(url.lastIndexOf('?')>-1){
                    url += "&"
                }else{
                    url += "?"
                }

                url += document.URL.substring(document.URL.lastIndexOf('?')+1,document.URL.length);
            }

            url = url.replace(new RegExp("parentComponent","gm"),"pastComponent");
            url = url.replace(new RegExp("parentWindowName","gm"),"pastWindowName");

            $.openIframeWindow({
                url : url,
                parentComponent : componentId,
                parentWindowName : window.name
            },event);
        });


        _this.on('click','[data-action=mutiSelect]',function(event){
            thisButton = $(this);
            event.preventDefault();
            var ids;
            ids = _this.grid.jqGrid('getGridParam','selarrrow');

            if(ids.length==0){
                alert('请至少选择一项');
                return;
            }

            var params = {};

            if(thisButton.attr('params')!=""){
                params = eval('('+thisButton.attr('params')+')');
            }

            params["ids"] = ids.toString();

            var url = thisButton.attr('href');

            if(document.URL.lastIndexOf('?')>-1){

                if(url.lastIndexOf('?')>-1){
                    url += "&"
                }else{
                    url += "?"
                }

                url += document.URL.substring(document.URL.lastIndexOf('?')+1,document.URL.length);
            }

            var flag = false

            if(thisButton.attr('confirmMsg')){
                flag = confirm(thisButton.attr('confirmMsg'));
            }else{
                flag = true;
            }

            if(flag){
                $.get(url,params,function(data){
                    if($.trim(thisButton.next('callback').html())){
                        eval(thisButton.next('callback').html());
                    }else{
                        if(data.success){
                            _this.grid.trigger("reloadGrid");
                            alert(data.msg);
                        }else{
                            alert(data.msg);
                        }
                    }

                },'json');
            }else{
                return;
            }
        });

        _this.on('click','[data-action=mutiSelectIframeWindow]',function(event){

            thisButton = $(this);
            event.preventDefault();

            var ids;
            ids = _this.grid.jqGrid('getGridParam','selarrrow');

            if(ids.length==0){
                alert('请至少选择一项');
                return;
            }

            var url = thisButton.attr('href');

            if(url.lastIndexOf('?')>-1){
                url += "&ids="+ids.toString();
            }else{
                url += "?ids="+ids.toString();
            }

            var params = eval('('+thisButton.attr('params')+')');

            for(var prop in params){
                url += "&" + prop + "=" + params[prop];
            }

            $.openIframeWindow({
                url : url,
                parentComponent : componentId,
                parentWindowName : window.name
            },event);
        });


        _this.on('click','[data-action=addToCheckBox]',function(event){
            thisButton = $(this);
            event.preventDefault();

            var ids;
            ids = _this.grid.jqGrid('getGridParam','selarrrow');

            if(ids.length==0){
                alert('请至少选择一项');
                return;
            }

            for(var i=0;i<ids.length;i++){
                if(_this.find('.checkbox-container').find(':checkbox[value='+ids[i]+']').size()<=0){
                    _this.find('.checkbox-container').append('<input type="checkbox" checked="checked" value="'+ids[i]+'"/>'+_this.grid.jqGrid('getRowData',ids[i])[thisButton.attr('checkboxDesc')]+'&nbsp;');
                }
            }

        });
    }

    this.render();

    return this;
}

/*****
 * 表格操作列渲染器
 **/
function gridActionRender(cellvalue, options, rowObject,action,componentId) {

    var params = "";

    if(action.actionParams!=""){
        params = "?";
        action.actionParams = eval('('+action.actionParams+')');
        for(var i=0;i<action.actionParams.length;i++){
            if(i!=0){
                params +="&";
            }
            params += action.actionParams[i].name + "=";

            if(action.actionParams[i].value=="id"){
                params += options.rowId;
            }else{
                var reg = new RegExp("^[0-9]*$");
                if(!reg.test(action.actionParams[i].value)){//支持列名形式
                    //todo 发现行不通。控件还没加载到相应数据
                    /*
                    for (var j in rowObject){
                        alert(j)
                    }
                    params += window[componentId].grid.jqGrid('getRowData',10)['name'];*/
                }else{//支持索引形式
                    params += rowObject[action.actionParams[i].value];
                }
            }
        }

    }

    if (action.linkType == "newPage"){
        var tmp = '<a style="color:blue;text-decoration: underline;" href="${url}{@if params!=""}${params}{@/if}" target="_blank">${desc}</a>&nbsp;';

        return juicer(tmp,{
            url:action.url,
            desc:action.desc,
            params:params
        });

    }else if (action.linkType == "iframeWindow"){
        var tmp = '<a style="color:blue;text-decoration: underline;" href="${url}{@if params!=""}${params}{@/if}" data-action="openIframeWindow" window-width="${width}"  window-height="${height}"   window-title="${title}">${desc}</a>&nbsp;';
        return juicer(tmp,{
            url:action.url,
            desc:action.desc,
            params:params,
            width : action.width,
            height : action.height,
            title : action.title
        });
    }
}

String.prototype.formatChar14Time = function(){
    if(this && this!=""){
        var year = this.substr(0,4);
        var month = this.substr(4,2);
        var day = this.substr(6,2);
        var hour = this.substr(8,2);
        var minute = this.substr(10,2);
        var second = this.substr(12,4);

        return year+"-"+month+"-"+day+" "+hour+":"+minute+":"+second;
    }else{
        return "";
    }

}

String.prototype.formatChar8Time = function(){
    if(this && this!=""){
        var year = this.substr(0,4);
        var month = this.substr(4,2);
        var day = this.substr(6,2);

        return year+"-"+month+"-"+day;
    }else{
        return "";
    }
}
