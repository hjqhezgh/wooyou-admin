jQuery.fn.grid = function (opts) {

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


    this.render = function () {

        $('#' + opts.tableId).jqGrid({
            url: opts.dataurl,
            datatype: 'json',
            colNames: opts.colNames,
            width: opts.width,
            height: opts.height,
            colModel: opts.colModel,
            rowNum: opts.rowNum,
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

        for (var i = 0; i < opts.actions.length; i++) {
            if (opts.actions[i].action == "add") {
                $('#' + opts.tableId).navButtonAdd('#' + opts.pageId, {
                    caption: opts.actions[i].desc,
                    buttonicon: "ui-icon-plus",
                    onClickButton: function () {
                        window.open(opts.addurl);
                    },
                    position: "last"
                });
            } else if (opts.actions[i].action == "modify") {
                $('#' + opts.tableId).navButtonAdd('#' + opts.pageId, {
                    caption: opts.actions[i].desc,
                    buttonicon: "ui-icon-pencil",
                    onClickButton: function () {

                        if (!lastId) {
                            alert('请选择一行');
                            return;
                        }

                        window.open(opts.modifyurl + lastId);
                    },
                    position: "last"
                });
            } else if (opts.actions[i].action == "delete") {
                $('#' + opts.tableId).navButtonAdd('#' + opts.pageId, {
                    caption: opts.actions[i].desc,
                    buttonicon: "ui-icon-minus",
                    onClickButton: function () {
                        if (!lastId) {
                            alert('请选择一行');
                            return;
                        }

                        if (confirm('确认删除？')) {
                            $.post(opts.deleteurl + lastId, {}, function (data) {
                                if (data.success) {
                                    alert('删除成功');
                                    $('#' + opts.tableId).trigger('reloadGrid');
                                } else {
                                    alert(data.msg)
                                }
                            }, 'json');
                        }
                    },
                    position: "last"
                });
            }else if (opts.actions[i].action == "custom") {
                var thisurl = opts.actions[i].url;
                $('#' + opts.tableId).navButtonAdd('#' + opts.pageId, {
                    caption: opts.actions[i].desc,
                    onClickButton: function () {
                        if (!lastId) {
                            alert('请选择一行');
                            return;
                        }
                        window.open(thisurl + "?id="+lastId);
                    },
                    position: "last"
                });
            }
        }

        //todo 表格多选框效果，有空了再做吧
        /*
        $('#' + opts.tableId).on('click','input[data-action=gird-checkbox]',function(){
            if($(this).prop('checked')){
                jQuery('#' + opts.tableId).jqGrid('setSelection',$(this).val(),false)
            }
        });*/

        this.renderSearchForm();
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

        //美化按钮
        this.find('a[data-action=reset]').button();

        this.find('a[data-action=search]').button();

        this.bindFormEvent();
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

    this.render();

    return this;
}


