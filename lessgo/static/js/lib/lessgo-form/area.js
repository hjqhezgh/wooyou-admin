/*****
 * 地区数据展示控件
 * 支持的功能：
 * 1.三级地区的初始化，默认地市是北京
 * 2.根据默认areaCode展示数据，多用于修改功能
 * 3.数据回调改变成功后，可以有回调函数
 * 4.一个页面支持多个实例
 */
var area = function(dom,desc,name,readonly,regionCode,callback){

    this.code= regionCode;

    this.desc = '<label class="control-label">'+desc+'：</label><div class="controls"></div>';

    this.name = name;

    this.readonly = readonly;

    this.area = $(dom);

    this.callback = callback

    this.init();
}

area.prototype = {

    init : function(){
        this.area.html(this.desc);

        if(!this.code) {
            this.finishRegion("")
            this.bind();
        }else{
            this.changeRegion(this.code);
            this.bind();
        }

    },

    changeRegion : function(code){

        var firstRegionCode = code.substring(0,2)+"0000";

        var secondRegionCode = code.substring(0,4)+"00";

        for (var i=0;i<3;i++){
            var selectDom = $(document.createElement('select'))
            this.area.append(selectDom);
            if(i==2){
                $(selectDom).attr('data-field','');
                $(selectDom).attr('name',this.name);
            }

            if(this.readonly=="true"){
                $(selectDom).attr('disabled','disabled');
            }

        }

        this.finishOneRegion("",this.area.find("select:first"),firstRegionCode);
        this.finishOneRegion(firstRegionCode,this.area.find("select:eq(1)"),secondRegionCode);
        this.finishOneRegion(secondRegionCode,this.area.find("select:eq(2)"),code);


    },

    bind : function(){
        var mythis = this;

        this.area.on('change','select',function(e){
            var _this = $(e.target);

            _this.nextAll('select').remove();

            var code = _this.val();

            mythis.finishRegion(code);
        });
    },

    finishRegion :function(code){
        var mythis = this;

        $.getJSON('/region/regions',{code:code},function(data){
            if(data.success){
                if(data.regions.length>0){
                    var selectDom = $(document.createElement('select'))
                    for(var i=0;i<data.regions.length;i++){
                        selectDom.append('<option value="'+data.regions[i].code+'">'+data.regions[i].name+'</option>');
                    }
                    mythis.area.find('.controls').append(selectDom);

                    if($(selectDom).index()==2){
                        $(selectDom).attr('data-field','');
                        $(selectDom).attr('name',mythis.name);
                    }

                    if(mythis.readonly=="true"){
                        $(selectDom).attr('disabled','disabled');
                    }

                    mythis.finishRegion(data.regions[0].code);

                }

                if(mythis.callback){
                    mythis.callback();
                }
            }else{
                console.log(data.reason)
            }

        });
    },

    finishOneRegion : function(code,select,selectCode){
        var mythis = this;
        $.getJSON('/region/regions',{code:code},function(data){
            if(data.success){
                if(data.regions.length>0){
                    for(var i=0;i<data.regions.length;i++){
                        select.append('<option value="'+data.regions[i].code+'">'+data.regions[i].name+'</option>');
                    }

                    select.find("option[value="+selectCode+"]").attr("selected","selected");

                    if(mythis.callback){
                        mythis.callback();
                    }
                }
            }else{
                //console.log(data.reason)
            }

        });
    }

}