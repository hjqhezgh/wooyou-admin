/*****
 * 时间维度，周数据的三级联动，目前不支持修改操作 todo
 * 支持的功能：
 * 1.三级数据的初始化，默认是第一个年份
 * 2.一个页面支持多个实例
 */
var timeDimWeek = function(dom,desc,name,readonly,callback){

    this.desc = desc+ "：";

    this.name = name;

    this.readonly = readonly;

    this.timedim = $(dom);

    this.callback = callback

    this.init();
}

timeDimWeek.prototype = {

    init : function(){

        this.timedim.html('<label class="control-label">'+this.desc+'</label><div class="controls"></div>');

        this.timedim = this.timedim.find('.controls');

        this.render()

        this.bind();
    },

    render : function(){

        var mythis = this;

        $.get('/timedim/years',function(data){
            if(data.success){
                if(data.datas.length>0){
                    var selectDom = $(document.createElement('select'))
                    for(var i=0;i<data.datas.length;i++){
                        selectDom.append('<option value="'+data.datas[i]+'">'+data.datas[i]+'年</option>');
                    }

                    mythis.timedim.append(selectDom);

                    if(mythis.readonly=="true"){
                        $(selectDom).attr('disabled','disabled');
                    }

                    mythis.changeYear(data.datas[0]);

                }
            }else{
                alert(data.msg);
            }

        },'json');


    },

    changeYear : function(year){

        var mythis = this;

        mythis.timedim.find('select:gt(0)').remove();

        $.get('/timedim/months',{year:year},function(data){
            if(data.success){
                if(data.datas.length>0){
                    var selectDom = $(document.createElement('select'))
                    for(var i=0;i<data.datas.length;i++){
                        selectDom.append('<option value="'+data.datas[i]+'">'+data.datas[i]+'月</option>');
                    }

                    mythis.timedim.append(selectDom);

                    if(mythis.readonly=="true"){
                        $(selectDom).attr('disabled','disabled');
                    }

                    mythis.changeMonth(year,data.datas[0]);

                }
            }else{
                alert(data.msg);
            }

        },'json');

    },

    changeMonth : function(year,month){

        var mythis = this;

        mythis.timedim.find('select:gt(1)').remove();

        $.get('/timedim/weeks',{year:year,month:month},function(data){
            if(data.success){
                if(data.datas.length>0){
                    var selectDom = $(document.createElement('select'));
                    selectDom.attr('data-field','data-field');
                    selectDom.attr('name',mythis.name);
                    for(var i=0;i<data.datas.length;i++){
                        selectDom.append('<option value="'+data.datas[i].weekKey+'">第'+data.datas[i].currentWeek+'周</option>');
                    }

                    mythis.timedim.append(selectDom);

                    if(mythis.readonly=="true"){
                        $(selectDom).attr('disabled','disabled');
                    }

                }
            }else{
                alert(data.msg);
            }

        },'json');

    },

    bind : function(){
        var mythis = this;

        this.timedim.on('change','select:first',function(e){
            mythis.changeYear($(this).val());
        });

        this.timedim.on('change','select:eq(1)',function(e){
            mythis.changeMonth(mythis.timedim.find('select:first').val(),$(this).val());
        });
    }

}