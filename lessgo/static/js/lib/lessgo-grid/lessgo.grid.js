jQuery.fn.grid = function (opts) {

    var componentId = this.attr("id");

    window[componentId] = this;

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

        //todo 表格多选框效果，有空了再做吧
        /*
        $('#' + opts.tableId).on('click','input[data-action=gird-checkbox]',function(){
            if($(this).prop('checked')){
                jQuery('#' + opts.tableId).jqGrid('setSelection',$(this).val(),false)
            }
        });*/

        this.renderSearchForm();

        this.bindButtonEvent();
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
                url+='&'+$(search).attr('name')+'-'+$(search).attr('data-searchType')+'='+$(search).val();
            });

            $('#' + opts.tableId).jqGrid('setGridParam',{url:url,page:1}).trigger("reloadGrid");

        });

        this.find('a[data-action=reset]').click(function(){
            _this.find('form')[0].reset();
        });

    }

    this.bindButtonEvent = function(){
        _this.on('click','[data-action=openIframeWindow]',function(event){
            thisButton = $(this);
            event.preventDefault();
            $.openIframeWindow({
                width : thisButton.attr('window-width') ,
                height : thisButton.attr('window-height') ,
                title : thisButton.attr('window-title'),
                url : thisButton.attr('href'),
                parentComponent : componentId
            },event);
        });


        _this.on('click','[data-action=mutiSelect]',function(event){
            thisButton = $(this);
            event.preventDefault();
            var ids;
            ids = _this.grid.jqGrid('getGridParam','selarrrow');

            var params = {};

            if(thisButton.attr('params')!=""){
                params = eval('('+thisButton.attr('params')+')');
            }

            params["ids"] = ids.toString();

            if(confirm(thisButton.attr('confirmMsg'))){
                $.post(thisButton.attr('href'),params,function(data){
                    if(data.success){
                        _this.grid.trigger("reloadGrid");
                        alert(data.msg);
                    }else{
                        alert(data.msg);
                    }
                },'json');
            }
        });
    }

    this.render();

    return this;
}

/*****
 * 表格操作列渲染器
 **/
function gridActionRender(cellvalue, options, rowObject,action) {

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
                params += rowObject[action.actionParams[i].value];
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